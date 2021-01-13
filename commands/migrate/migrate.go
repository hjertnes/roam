package migrate

import (
	"github.com/hjertnes/roam/migration"
	"github.com/hjertnes/roam/state"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	"github.com/rotisserie/eris"
)

func Run(path string) error {
	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "Failed to create state")
	}

	spinner, err := spinner2.Run("")
	if err != nil{
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil{
		return eris.Wrap(err, "failed to start spinner")
	}

	mig := migration.New(s.Ctx, s.Conn)

	err = mig.Migrate()
	if err != nil {
		return eris.Wrap(err, "failed to migrate database")
	}

	err = spinner.Stop()
	if err != nil{
		return eris.Wrap(err, "failed to stop spinner")
	}

	return nil
}