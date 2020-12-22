package dal

import (
	"context"
	"github.com/hjertnes/roam/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/georgysavva/scany/pgxscan"
)

// exists
// add
// update
// method to set opened_at

type Dal struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

func New(ctx context.Context, conn *pgxpool.Pool) *Dal{
	return &Dal{
		ctx: ctx,
		conn: conn,
	}
}

func (d *Dal) Exists(path string) (bool, error){
	res := d.conn.QueryRow(d.ctx, `select exists (select 1 from files where path=$1)`, path)

	var result bool
	err := res.Scan(&result)
	if err != nil{
		return false, err
	}

	return result, nil
}

func (d *Dal) Create(path, title, content string) error{
	_, err := d.conn.Exec(
		d.ctx,
		`insert into files (processed_at, opened_at, path, title, title_tokens, content, content_tokens) values(timezone('utc', now()), timezone('utc', now()), $1, $2, to_tsvector($2), $3, to_tsvector($3))`,
		path, title, content)
	if err != nil{
		return err
	}

	return nil
}

func (d *Dal) Update(path, title, content string) error{

	_, err := d.conn.Exec(
		d.ctx,
		`update files set processed_at=timezone('utc', now()), title=$2, title_tokens=to_tsvector($2), content=$3, content_tokens=to_tsvector($3) where path=$1`,
		path, title, content)
	if err != nil{
		return err
	}

	return nil
}

func (d *Dal)Find(search string) ([]models.File, error){
	result := make([]models.File, 0)
	res, err := d.conn.Query(d.ctx, `SELECT path, title FROM files WHERE title_tokens @@ to_tsquery($1);`, search)
	if err != nil{
		return result, err
	}
	err = pgxscan.ScanAll(&result, res)
	if err != nil{
		return result, err
	}

	return result, nil
}