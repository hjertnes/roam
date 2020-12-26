package dal

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
	"strings"
)

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
		return false, eris.Wrap(err, "failed to check if file exits")
	}

	return result, nil
}

func (d *Dal) GetFiles() ([]string, error){
	var files []string

	result, err := d.conn.Query(
		d.ctx,
		`select path from files`)
	if err != nil{
		return files, eris.Wrap(err, "failed to get list of files")
	}

	err = pgxscan.ScanAll(&files, result)
	if err != nil{
		return files, eris.Wrap(err, "could not parse")
	}

	return files, nil
}

func (d *Dal) Create(path, title, content string, private bool) error{
	_, err := d.conn.Exec(
		d.ctx,
		`insert into files (processed_at, path, title, title_tokens, content, content_tokens, private) values(timezone('utc', now()), $1, $2, to_tsvector($2), $3, to_tsvector($3), $4)`,
		path, title, content, private)
	if err != nil{
		return eris.Wrap(err, "failed to create file")
	}

	return nil
}

func (d *Dal) Delete() error{
	files, err := d.GetFiles()

	if err != nil{
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, r := range files{
		if !utils.FileExist(r){
			_, err = d.conn.Exec(d.ctx, `delete from files where path=$1`, r)
			if err != nil{
				return eris.Wrap(err, "failed to delete file")
			}
		}
	}

	return nil
}

func (d *Dal) Update(path, title, content string, private bool) error{

	_, err := d.conn.Exec(
		d.ctx,
		`update files set processed_at=timezone('utc', now()), title=$2, title_tokens=to_tsvector($2), content=$3, content_tokens=to_tsvector($3), private=$4 where path=$1`,
		path, title, content, private)
	if err != nil{
		return eris.Wrap(err, "failed to update file")
	}

	return nil
}

func buildVectorSearch(input string)string{
	if !strings.Contains(input, " "){
		return fmt.Sprintf("%s:*", input)
	}
	output := make([]string, 0)

	for _,l := range strings.Split(input, " "){
		output = append(output, fmt.Sprintf("%s:*", l))
	}

	return strings.Join(output, "&")
}

func (d *Dal) Find(search string) ([]models.File, error){
	result := make([]models.File, 0)
	res, err := d.conn.Query(d.ctx, `SELECT path, title FROM files WHERE title_tokens @@ to_tsquery($1);`, buildVectorSearch(search))
	if err != nil{
		return result, eris.Wrap(err, "failed to run search query")
	}
	err = pgxscan.ScanAll(&result, res)
	if err != nil{
		return result, eris.Wrap(err, "failed to scan search query")
	}

	return result, nil
}

func (d *Dal) FindExact(search string) ([]models.File, error){
	result := make([]models.File, 0)
	res, err := d.conn.Query(d.ctx, `SELECT path, title FROM files WHERE title=$1;`, search)
	if err != nil{
		return result, eris.Wrap(err, "failed to run search query")
	}
	err = pgxscan.ScanAll(&result, res)
	if err != nil{
		return result, eris.Wrap(err, "failed to scan search query")
	}

	return result, nil
}

func (d *Dal) Stats() (int, int, int, error){
	var all int
	var public int
	var private int

	res1 := d.conn.QueryRow(d.ctx, "select count(*) from files")
	res2 := d.conn.QueryRow(d.ctx, "select count(*) from files where private=true")
	res3 := d.conn.QueryRow(d.ctx, "select count(*) from files where private=false")

	err := res1.Scan(&all)
	if err != nil{
		return 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	err = res2.Scan(&private)
	if err != nil{
		return 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	err = res3.Scan(&public)
	if err != nil{
		return 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	return all, public, private, nil
}