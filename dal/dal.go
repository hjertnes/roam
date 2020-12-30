// Package dal contains all the database stuff
package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// Dal is the exported type.
type Dal struct {
	path string
	ctx  context.Context
	conn *pgxpool.Pool
}

// New is the constructor.
func New(path string, ctx context.Context, conn *pgxpool.Pool) *Dal {
	return &Dal{
		path: path,
		ctx:  ctx,
		conn: conn,
	}
}

// Exists checks if a row with the path exists.
func (d *Dal) Exists(path string) (bool, error) {
	res := d.conn.QueryRow(d.ctx, `select exists (select 1 from files where path=$1)`, path)

	var result bool

	err := res.Scan(&result)
	if err != nil {
		return false, eris.Wrap(err, "failed to check if file exits")
	}

	return result, nil
}

// GetFiles returns all files in database/.
func (d *Dal) GetFiles() ([]models.File, error) {
	var files []models.File

	result, err := d.conn.Query(
		d.ctx,
		`select ID, title, path, private from files order by path`)
	if err != nil {
		return files, eris.Wrap(err, "failed to get list of files")
	}

	err = pgxscan.ScanAll(&files, result)
	if err != nil {
		return files, eris.Wrap(err, "could not parse")
	}

	return files, nil
}

func getPath(path string) string{
	res := make([]string, 0)

	elems := strings.Split(path, "/")

	for _, e := range elems{
		if strings.HasSuffix(e, ".md"){
			continue
		}

		res = append(res, e)
	}

	return strings.Join(res, "/")
}

func getParent(path string) string{
	res := make([]string, 0)

	elems := strings.Split(path, "/")

	for i, e := range elems{
		if i == len(elems) -1{
			continue
		}

		res = append(res, e)
	}

	return strings.Join(res, "/")
}
//TODO
// add if not exist
// get folderid
func (d *Dal) GetFolder(path string) (string, error){
	var id string
	res := d.conn.QueryRow(d.ctx, `select id from folders where path=$1`, path)
	err := res.Scan(&id)
	if err != nil {
		return "", eris.Wrap(err, "failed to get folder")
	}

	return id, nil
}

func (d *Dal) GetFolderFiles(folderId string) ([]models.File, error){
	var files []models.File

	result, err := d.conn.Query(
		d.ctx,
		`select ID, title, path, private from files where folder_fk=$1 order by path`, folderId)
	if err != nil {
		return files, eris.Wrap(err, "failed to get list of files")
	}

	err = pgxscan.ScanAll(&files, result)
	if err != nil {
		return files, eris.Wrap(err, "could not parse")
	}

	return files, nil
}

func (d *Dal) GetSubFolders(folderId string) ([]models.Folder, error){
	var result []models.Folder

	q, err := d.conn.Query(d.ctx, `select id, path from folders where parent_fk=$1`, folderId)
	if err != nil{
		return result, eris.Wrap(err, "could not get sub folder")
	}

	err = pgxscan.ScanAll(&result, q)
	if err != nil{
		return result, eris.Wrap(err, "could not get sub folder")
	}

	return result, nil
}

func (d *Dal) GetRootFolder() (*models.Folder, error){
	q, err := d.conn.Query(d.ctx, `select id, path from folders where parent_fk is null`)
	if err != nil{
		return nil, eris.Wrap(err, "could not get root folder")
	}
	var result []models.Folder

	err = pgxscan.ScanAll(&result, q)
	if err != nil{
		return nil, eris.Wrap(err, "could not get root folder")
	}

	return &result[0], nil
}

func (d *Dal) CreateFolder(path string) error{
	parentId := utils.NilStringPointer()

	if path != d.path {
		p := getParent(path)
		err := d.CreateFolder(p)
		if err != nil{
			return eris.Wrap(err, "failed to create parent folder")
		}
		pp, err := d.GetFolder(p)
		parentId = &pp
		if err != nil{
			return eris.Wrap(err, "failed to get parent folder")
		}
	}
	_, err := d.conn.Exec(d.ctx, `insert into folders (path, parent_fk) values($1, $2) on conflict do nothing`, path, parentId)
	if err != nil {
		return eris.Wrap(err, "failed to create folder")
	}

	return nil
}

// Create adds new file to database.
func (d *Dal) Create(path, title, content string, private bool) error {
	p := getPath(path)
	err := d.CreateFolder(p)
	if err != nil{
		return eris.Wrap(err, "failed to create folder")
	}
	folderId, err := d.GetFolder(p)
	if err != nil{
		return eris.Wrap(err, "failed to get folder")
	}

	_, err = d.conn.Exec(
		d.ctx,
		`
insert into files 
(processed_at, path, title, title_tokens, content, content_tokens, private, folder_fk) 
values(timezone('utc', now()), $1, $2, to_tsvector($2), $3, to_tsvector($3), $4, $5)`,
		path, title, content, private, folderId)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	return nil
}

// Delete removes all files from database that don't exist on the file system.
func (d *Dal) Delete() error {
	files, err := d.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, r := range files {
		if !utils.FileExist(r.Path) {
			_, err = d.conn.Exec(d.ctx, `delete from files where path=$1`, r)
			if err != nil {
				return eris.Wrap(err, "failed to delete file")
			}
		}
	}

	return nil
}

// Update updates a file.
func (d *Dal) Update(path, title, content string, private bool) error {
	_, err := d.conn.Exec(
		d.ctx,
		`
update files set 
processed_at=timezone('utc', now()), 
title=$2, title_tokens=to_tsvector($2), 
content=$3, 
content_tokens=to_tsvector($3), 
private=$4 where path=$1`,
		path, title, content, private)
	if err != nil {
		return eris.Wrap(err, "failed to update file")
	}

	return nil
}

func buildVectorSearch(input string) string {
	if !strings.Contains(input, " ") {
		return fmt.Sprintf("%s:*", input)
	}

	output := make([]string, 0)

	for _, l := range strings.Split(input, " ") {
		output = append(output, fmt.Sprintf("%s:*", l))
	}

	return strings.Join(output, "&")
}

// Find returns a list of notes that match a search.
func (d *Dal) Find(search string) ([]models.File, error) {
	result := make([]models.File, 0)

	res, err := d.conn.Query(d.ctx, `
SELECT ID, path, title, private
FROM files 
WHERE title_tokens @@ to_tsquery($1);`, buildVectorSearch(search))
	if err != nil {
		return result, eris.Wrap(err, "failed to run search query")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan search query")
	}

	return result, nil
}

// FindExact returns a list of notes with an exact title.
func (d *Dal) FindExact(search string) ([]models.File, error) {
	result := make([]models.File, 0)

	res, err := d.conn.Query(d.ctx, `SELECT ID, path, title, private FROM files WHERE title=$1;`, search)
	if err != nil {
		return result, eris.Wrap(err, "failed to run search query")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan search query")
	}

	return result, nil
}

// Stats returns stats.
func (d *Dal) Stats() (int, int, int, int, error) {
	var all, private, public, links int

	res1 := d.conn.QueryRow(d.ctx, "select count(*) from files")
	res2 := d.conn.QueryRow(d.ctx, "select count(*) from files where private=true")
	res3 := d.conn.QueryRow(d.ctx, "select count(*) from files where private=false")
	res4 := d.conn.QueryRow(d.ctx, "select count(*) from links")

	err := res1.Scan(&all)
	if err != nil {
		return 0, 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	err = res2.Scan(&private)
	if err != nil {
		return 0, 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	err = res3.Scan(&public)
	if err != nil {
		return 0, 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	err = res4.Scan(&links)
	if err != nil {
		return 0, 0, 0, 0, eris.Wrap(err, "failed to scan")
	}

	return all, public, private, links, nil
}

// AddLink adds a link.
func (d *Dal) AddLink(fileID, linkedToFile string) error {
	_, err := d.conn.Exec(
		d.ctx,
		`insert into links (file_fk, link_fk) values($1, $2) on conflict do nothing`,
		fileID, linkedToFile)
	if err != nil {
		return eris.Wrap(err, "failed to create link")
	}

	return nil
}

// DeleteLink removes a link from a file.
func (d *Dal) DeleteLink(fileID, linkedToFile string) error {
	_, err := d.conn.Exec(
		d.ctx,
		`delete from links where file_fk=$ and link_fk=$2`, fileID, linkedToFile,
		fileID, linkedToFile)
	if err != nil {
		return eris.Wrap(err, "failed to create link")
	}

	return nil
}

// GetCurrentLinks returns the path of all links of a file.
func (d *Dal) GetCurrentLinks(fileID string) ([]string, error) {
	result := make([]string, 0)

	res, err := d.conn.Query(d.ctx, `SELECT link_fk from links WHERE file_fk=$1;`, fileID)
	if err != nil {
		return result, eris.Wrap(err, "failed to query for list of links")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan query for list of links")
	}

	return result, nil
}

// GetBacklinks returns the backlinks of a file.
func (d *Dal) GetBacklinks(fileID string) ([]models.File, error) {
	result := make([]models.File, 0)

	res, err := d.conn.Query(d.ctx, `
select ID, path, title, private from files where ID in 
(SELECT file_fk from links WHERE link_fk=$1);`, fileID)
	if err != nil {
		return result, eris.Wrap(err, "failed to query for list of links")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan query for list of links")
	}

	return result, nil
}

// GetLinks returns the links of a file.
func (d *Dal) GetLinks(fileID string) ([]models.File, error) {
	result := make([]models.File, 0)

	res, err := d.conn.Query(
		d.ctx, `
select ID, path, title, private from files where ID in 
(SELECT link_fk from links WHERE file_fk=$1);`, fileID)
	if err != nil {
		return result, eris.Wrap(err, "failed to query for list of links")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan query for list of links")
	}

	return result, nil
}
