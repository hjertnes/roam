package migrate

import (
	"github.com/hjertnes/roam/migration"
	"github.com/hjertnes/roam/state"
	"github.com/rotisserie/eris"
)

func Run(path string) error {
	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "Failed to create state")
	}

	mig := migration.New(s.Ctx, s.Conn)

	err = mig.Migrate()
	if err != nil {
		return eris.Wrap(err, "failed to migrate database")
	}

	return nil
}