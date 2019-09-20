package gamemaster

import (
	"net/http"
	"time"

	"github.com/Sigafoos/pokemongo"
)

const gamemasterURL = "https://raw.githubusercontent.com/pvpoke/pvpoke/master/src/data/gamemaster.json"

type Gamemaster struct {
	c       *http.Client
	gm      *pokemongo.Gamemaster
	updated time.Time
}

func New(c *http.Client) *Gamemaster {
	return &Gamemaster{c: c}
}

func (g *Gamemaster) PokemonByID(ID string) (*pokemongo.Pokemon, error) {
	if err := g.update(); err != nil {
		return nil, err
	}

	return g.gm.PokemonByID(ID), nil
}

func (g *Gamemaster) Pokemon() []pokemongo.Pokemon {
	if err := g.update(); err != nil {
		return []pokemongo.Pokemon{}
	}

	return g.gm.Pokemon
}

func (g *Gamemaster) update() error {
	if g.gm != nil && time.Since(g.updated) < 24*time.Hour {
		return nil
	}

	resp, err := g.c.Get(gamemasterURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	gm, err := pokemongo.NewGamemaster(resp.Body)
	if err != nil {
		return err
	}
	g.gm = gm
	g.updated = time.Now()
	return nil
}
