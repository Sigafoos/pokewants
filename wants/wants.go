package wants

import (
	"fmt"
	"log"

	"github.com/Sigafoos/pokewants/gamemaster"

	"github.com/Sigafoos/pokemongo"
	"github.com/gocraft/dbr"
)

const tableName = "wants"

type Row struct {
	User    string
	Pokemon string
}

type Wants struct {
	db     *dbr.Connection
	logger dbr.EventReceiver
	gm     *gamemaster.Gamemaster
}

func New(db *dbr.Connection, logger dbr.EventReceiver, gm *gamemaster.Gamemaster) (*Wants, error) {
	w := &Wants{
		db:     db,
		logger: logger,
		gm:     gm,
	}
	err := w.createDB()
	if err != nil {
		_ = w.db.Close()
		return nil, err
	}
	return w, nil
}

func (w *Wants) Get(user string) []*pokemongo.Pokemon {
	var results []Row

	session := w.db.NewSession(w.logger)
	defer session.Close()
	session.Begin()
	session.Select("pokemon").From(tableName).Where("user = ?", user).Load(&results)

	var pokemon []*pokemongo.Pokemon
	for _, row := range results {
		p, err := w.gm.PokemonByID(row.Pokemon)
		if err != nil {
			log.Println(err)
			continue
		}
		pokemon = append(pokemon, p)
	}
	return pokemon
}

func (w *Wants) Add(user, pokemon string) error {
	p, err := w.gm.PokemonByID(pokemon)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("Pokemon '%s' not found", pokemon)
	}

	session := w.db.NewSession(w.logger)
	defer session.Close()
	tx, err := session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	row := &Row{
		User:    user,
		Pokemon: pokemon,
	}
	_, err = tx.InsertInto(tableName).
		Record(row).
		Exec()

	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (w *Wants) createDB() error {
	// this *might* be dependent on sqlite?
	createSQL := `
CREATE TABLE IF NOT EXISTS wants(
user text NOT NULL,
pokemon text NOT NULL,
PRIMARY KEY (user, pokemon)
)`
	_, err := w.db.Exec(createSQL)
	return err
}
