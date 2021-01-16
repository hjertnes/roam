// Package state contains a wrapper around configuration and dal to make it less verbose.
package state

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/errs"

	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// State is the exported type.
type State struct {
	Path string
	Conf *models.Configuration
	Ctx  context.Context
	Conn *pgxpool.Pool
	Dal  dal2.Dal
}


// New is the constructor.

func NewWithoutStatusCheck(path string) (*State, error) {
	return _new(false, path)
}
func New(path string) (*State, error) {
	return _new(true, path)
}
func _new(check bool, path string) (*State, error) {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return nil, eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()

	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return nil, eris.Wrap(err, "failed to connect to database")
	}

	dal := dal2.New(ctx, pxp, path)

	if check{
		s, err := dal.WasLastSyncSuccessful()
		if err != nil {
			if eris.Is(err, errs.ErrNever){
				fmt.Println("Sync have never been done")
			} else {
				return nil, eris.Wrap(err, "failed to check last sync")
			}

		} else if !s{
			fmt.Println("Last sync failed. Run diagnostic to see why")
		}
	}

	return &State{
		Conn: pxp,
		Conf: conf,
		Ctx:  ctx,
		Path: path,
		Dal:  dal,
	}, nil
}
