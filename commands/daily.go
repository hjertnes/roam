package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// Daily creates a daily note.
func Daily(path, date string) error {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "could not get config")
	}

	ctx := context.Background()

	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "could not connect to database")
	}

	dal := dal2.New(ctx, pxp)

	filename := fmt.Sprintf("%s/Daily Notes/%s.md", path, date)

	if !utilslib.FileExist(filename) {
		templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", path, "daily.txt"))
		if err != nil {
			return eris.Wrap(err, "failed to read template")
		}

		err = createFile(dal, filename, "", templatedata, conf)
		if err != nil {
			return eris.Wrap(err, "failed to create file")
		}
	}

	editor := utils.GetEditor()

	cmd := exec.Command(editor, filename) // #nosec G204

	err = cmd.Start()
	if err != nil {
		return eris.Wrap(err, "faield to editNote daily in editor")
	}

	return nil
}
