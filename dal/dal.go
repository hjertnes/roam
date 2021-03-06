// Package dal contains all the database stuff
package dal

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/utils/pathutils"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hjertnes/roam/models"
	utilslib "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// Dal is the exported type.
type Dal interface {
	FileExists(path string) (bool, error)
	GetFiles() ([]models.File, error)
	GetFolders()([]models.Folder, error)
	GetFolderFiles(folderID string) ([]models.File, error)
	GetSubFolders(folderID string) ([]models.Folder, error)
	GetRootFolder() (*models.Folder, error)
	CreateFile(path, title, content string, private bool) error
	DeleteFiles() error
	UpdateFile(path, title, content string, private bool) error
	FindFileFuzzy(search string) ([]models.File, error)
	FindFileExact(search string) ([]models.File, error)
	Stats() (int, int, int, int, error)
	AddLink(fileID, linkedToFile string) error
	DeleteLink(fileID, linkedToFile string) error
	GetBacklinks(fileID string, includePrivate bool) ([]models.File, error)
	GetLinks(fileID string, includePrivate bool) ([]models.File, error)
	Clear() error
	AddLog(failure bool) error
	GetLog() ([]models.Log, error)
	WasLastSyncSuccessful() (bool, time.Time, error)
	ClearLog() error

}

type dal struct {
	path string
	ctx  context.Context
	conn *pgxpool.Pool
}

// New is the constructor.
func New(ctx context.Context, conn *pgxpool.Pool, path string) Dal {
	return &dal{
		path: path,
		ctx:  ctx,
		conn: conn,
	}
}

func (d *dal) Clear() error {
	_, err := d.conn.Exec(d.ctx, `delete from sync_history;`)
	if err != nil {
		return eris.Wrap(err, "failed to empty sync history")
	}

	_, err = d.conn.Exec(d.ctx, `delete from links;`)
	if err != nil {
		return eris.Wrap(err, "failed to empty links")
	}

	_, err = d.conn.Exec(d.ctx, `delete from files;`)
	if err != nil {
		return eris.Wrap(err, "failed to empty files")
	}

	_, err = d.conn.Exec(d.ctx, `delete from folders;`)
	if err != nil {
		return eris.Wrap(err, "failed to empty folders")
	}

	return nil
}

// Exists checks if a row with the path exists.
func (d *dal) FileExists(path string) (bool, error) {
	res := d.conn.QueryRow(d.ctx, `select exists (select 1 from files where path=$1)`, path)

	var result bool

	err := res.Scan(&result)
	if err != nil {
		return false, eris.Wrap(err, "failed to check if file exits")
	}

	return result, nil
}

// GetFiles returns all files in database/.
func (d *dal) GetFiles() ([]models.File, error) {
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

func (d *dal) getFolder(path string) (string, error) {
	var id string

	res := d.conn.QueryRow(d.ctx, `select id from folders where path=$1`, path)

	err := res.Scan(&id)
	if err != nil {
		return "", eris.Wrap(err, "failed to get folder")
	}

	return id, nil
}

func (d *dal) GetFolderFiles(folderID string) ([]models.File, error) {
	var files []models.File

	result, err := d.conn.Query(
		d.ctx,
		`select ID, title, path, private from files where folder_fk=$1 order by path`, folderID)
	if err != nil {
		return files, eris.Wrap(err, "failed to get list of files")
	}

	err = pgxscan.ScanAll(&files, result)
	if err != nil {
		return files, eris.Wrap(err, "could not parse")
	}

	return files, nil
}

func (d *dal) GetSubFolders(folderID string) ([]models.Folder, error) {
	var result []models.Folder

	q, err := d.conn.Query(d.ctx, `select id, path from folders where parent_fk=$1`, folderID)
	if err != nil {
		return result, eris.Wrap(err, "could not get sub folder")
	}

	err = pgxscan.ScanAll(&result, q)
	if err != nil {
		return result, eris.Wrap(err, "could not get sub folder")
	}

	return result, nil
}

func (d *dal) GetFolders() ([]models.Folder, error) {
	var result []models.Folder

	q, err := d.conn.Query(d.ctx, `select id, path from folders`)
	if err != nil {
		return result, eris.Wrap(err, "could not get folders")
	}

	err = pgxscan.ScanAll(&result, q)
	if err != nil {
		return result, eris.Wrap(err, "could not get folders")
	}

	return result, nil
}

func (d *dal) GetRootFolder() (*models.Folder, error) {
	q, err := d.conn.Query(d.ctx, `select id, path from folders where parent_fk is null`)
	if err != nil {
		return nil, eris.Wrap(err, "could not get root folder")
	}

	var result []models.Folder

	err = pgxscan.ScanAll(&result, q)
	if err != nil {
		return nil, eris.Wrap(err, "could not get root folder")
	}

	return &result[0], nil
}

func (d *dal) createFolder(path string) error {
	parentID := utilslib.NilStringPointer()

	if path != d.path {
		p := pathutils.New(path).GetParent()

		err := d.createFolder(p)
		if err != nil {
			return eris.Wrap(err, "failed to create parent folder")
		}

		pp, err := d.getFolder(p)
		if err != nil {
			return eris.Wrap(err, "failed to get parent folder")
		}

		parentID = &pp
	}

	_, err := d.conn.Exec(
		d.ctx,
		`insert into folders (path, parent_fk) values($1, $2) on conflict do nothing`,
		path,
		parentID)
	if err != nil {
		return eris.Wrap(err, "failed to create folder")
	}

	return nil
}

func (d *dal) AddLog(failure bool) error{
	_, err := d.conn.Exec(d.ctx, `insert into sync_history (failure) values($1)`, failure)
	if err != nil{
		return eris.Wrap(err, "failed to add log")
	}

	return nil
}

func (d *dal) GetLog() ([]models.Log, error){
	result := make([]models.Log, 0)

	res, err := d.conn.Query(d.ctx, `SELECT ID, created_at, failure from sync_history order by created_at asc`)
	if err != nil {
		return result, eris.Wrap(err, "failed to run list query")
	}

	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return result, eris.Wrap(err, "failed to scan list query")
	}

	return result, nil
}

func (d *dal) WasLastSyncSuccessful() (bool, time.Time, error){
	result := make([]models.Log, 0)

	res, err := d.conn.Query(d.ctx, `SELECT ID, created_at, failure from sync_history order by created_at desc limit 1`)
	if err != nil {
		return false, time.Time{}, eris.Wrap(err, "failed to run list query")
	}



	err = pgxscan.ScanAll(&result, res)
	if err != nil {
		return false, time.Time{}, eris.Wrap(err, "failed to scan list query")
	}

	if len(result) == 0{
		return false, time.Time{}, eris.Wrap(errs.ErrNever, "never ran")
	}

	return !result[0].Failure, result[0].Timestamp, nil
}

func (d *dal) ClearLog() error{
	_, err := d.conn.Exec(d.ctx, `delete from sync_history`)
	if err != nil{
		return eris.Wrap(err, "failed to add log")
	}

	return nil
}

// Create adds new file to database.
func (d *dal) CreateFile(path, title, content string, private bool) error {
	p := pathutils.New(path).GetParent()

	err := d.createFolder(p)
	if err != nil {
		return eris.Wrap(err, "failed to create folder")
	}

	folderID, err := d.getFolder(p)
	if err != nil {
		return eris.Wrap(err, "failed to get folder")
	}

	_, err = d.conn.Exec(
		d.ctx,
		`
insert into files 
(processed_at, path, title, title_tokens, content, content_tokens, private, folder_fk) 
values(timezone('utc', now()), $1, $2, to_tsvector($2), $3, to_tsvector($3), $4, $5)`,
		path, title, content, private, folderID)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	return nil
}

// Delete removes all files from database that don't exist on the file system.
func (d *dal) DeleteFiles() error {
	files, err := d.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, r := range files {
		if !utilslib.FileExist(r.Path) {
			_, err = d.conn.Exec(d.ctx, `delete from links where file_fk=$1`, r.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete links")
			}

			_, err = d.conn.Exec(d.ctx, `delete from links where link_fk=$1`, r.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete links")
			}

			_, err = d.conn.Exec(d.ctx, `delete from files where id=$1`, r.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete file")
			}
		}
	}

	folders, err := d.GetFolders()
	if err != nil {
		return eris.Wrap(err, "failed to get list of folders")
	}

	for _, r := range folders {
		if !utilslib.FileExist(r.Path) {
			_, err = d.conn.Exec(d.ctx, `delete from folders where parent_fk=$1`, r.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete folder")
			}

			_, err = d.conn.Exec(d.ctx, `delete from folders where id=$1`, r.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete folder")
			}
		}
	}

	return nil
}

// Update updates a file.
func (d *dal) UpdateFile(path, title, content string, private bool) error {
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

// Find returns a list of notes that match a search.
func (d *dal) FindFileFuzzy(search string) ([]models.File, error) {
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
func (d *dal) FindFileExact(search string) ([]models.File, error) {
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
func (d *dal) Stats() (int, int, int, int, error) {
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
func (d *dal) AddLink(fileID, linkedToFile string) error {
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
func (d *dal) DeleteLink(fileID, linkedToFile string) error {
	_, err := d.conn.Exec(
		d.ctx,
		`delete from links where file_fk=$1 and link_fk=$2`, fileID, linkedToFile)
	if err != nil {
		return eris.Wrap(err, "failed to create link")
	}

	return nil
}

// GetBacklinks returns the backlinks of a file.
func (d *dal) GetBacklinks(fileID string, includePrivate bool) ([]models.File, error) {
	result := make([]models.File, 0)

	if includePrivate{
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
	}else{
		res, err := d.conn.Query(d.ctx, `
select ID, path, title, private from files where ID in 
(SELECT file_fk from links WHERE link_fk=$1) and private=false;`, fileID)
		if err != nil {
			return result, eris.Wrap(err, "failed to query for list of links")
		}

		err = pgxscan.ScanAll(&result, res)
		if err != nil {
			return result, eris.Wrap(err, "failed to scan query for list of links")
		}

		return result, nil
	}
}

// GetLinks returns the links of a file.
func (d *dal) GetLinks(fileID string, includePrivate bool) ([]models.File, error) {
	result := make([]models.File, 0)

	if includePrivate{
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
	}else{
		res, err := d.conn.Query(
			d.ctx, `
select ID, path, title, private from files where ID in 
(SELECT link_fk from links WHERE file_fk=$1) and private=false;`, fileID)
		if err != nil {
			return result, eris.Wrap(err, "failed to query for list of links")
		}

		err = pgxscan.ScanAll(&result, res)
		if err != nil {
			return result, eris.Wrap(err, "failed to scan query for list of links")
		}

		return result, nil
	}

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