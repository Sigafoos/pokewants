package wants

import (
	"fmt"
	"log"
	"strings"

	"github.com/Sigafoos/pokewants/gamemaster"

	"github.com/Sigafoos/pokemongo"
	"github.com/gocraft/dbr"
)

const tableName = "wants"

var (
	ErrorPokemonNotFound = fmt.Errorf("Pokemon not found")
	ErrorDuplicate       = fmt.Errorf("Want already exists")
)

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
	session.Begin()
	session.Select("pokemon").
		From(tableName).
		Where("username = ?", user).
		Load(&results)

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
		return ErrorPokemonNotFound
	}

	session := w.db.NewSession(w.logger)
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
		Columns("username", "pokemon").
		Record(row).
		Exec()

	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return ErrorDuplicate
		}
		return err
	}

	tx.Commit()
	return nil
}

func (w *Wants) Delete(user, pokemon string) error {
	p, err := w.gm.PokemonByID(pokemon)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrorPokemonNotFound
	}

	session := w.db.NewSession(w.logger)
	tx, err := session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.DeleteFrom(tableName).
		Where("username = ?", user).
		Where("pokemon = ?", pokemon).
		Exec()

	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (w *Wants) createDB() error {
	createSQL := `
CREATE TABLE IF NOT EXISTS wants(
username varchar(50) NOT NULL,
pokemon varchar(100) NOT NULL,
PRIMARY KEY (username, pokemon)
)`
	_, err := w.db.Exec(createSQL)
	return err
}
