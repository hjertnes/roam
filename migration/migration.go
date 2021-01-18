// Package migration migrates database changes
package migration

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

// Migration is the exported type.
type Migration struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

// New is the constructor.
func New(ctx context.Context, conn *pgxpool.Pool) *Migration {
	return &Migration{
		ctx:  ctx,
		conn: conn,
	}
}

func (m *Migration) v1() error {
	_, err := m.conn.Exec(m.ctx, `CREATE EXTENSION pgcrypto;`)
	if err != nil {
		return eris.Wrap(err, "could not enable pgcrypto")
	}

	_, err = m.conn.Exec(m.ctx, `create table migration_history(migration int not null unique);`)
	if err != nil {
		return eris.Wrap(err, "failed to create migration history")
	}

	_, err = m.conn.Exec(m.ctx, `
create table folders(
id uuid primary key not null default gen_random_uuid() unique,
parent_fk uuid null references folders(id),
path text not null unique
);`)
	if err != nil {
		return eris.Wrap(err, "failed to create folders table")
	}

	_, err = m.conn.Exec(m.ctx, `
create table files(
id uuid primary key not null default gen_random_uuid() unique, 
processed_at timestamp not null default timezone('utc', now()),
folder_fk uuid not null references folders(id),
path text not null unique, 
title text not null, 
title_tokens tsvector not null,
private bool not null,
content text not null, 
content_tokens tsvector not null);`)
	if err != nil {
		return eris.Wrap(err, "failed to create files table")
	}

	_, err = m.conn.Exec(m.ctx, `insert into migration_history (migration) values(1);`)
	if err != nil {
		return eris.Wrap(err, "failed to insert migration into migration history")
	}

	return nil
}

func (m *Migration) v2() error {
	_, err := m.conn.Exec(m.ctx, `
create table links(
id uuid primary key not null default gen_random_uuid() unique,
file_fk uuid not null references files(id),
link_fk uuid not null references files(id),
unique(file_fk, link_fk)
);
CREATE INDEX links_idx ON links (file_fk, link_fk);`)
	if err != nil {
		return eris.Wrap(err, "failed to create links table")
	}

	_, err = m.conn.Exec(m.ctx, `insert into migration_history (migration) values(2);`)
	if err != nil {
		return eris.Wrap(err, "failed to insert migration into migration history")
	}

	return nil
}

func (m *Migration) v3() error {
	_, err := m.conn.Exec(m.ctx, `
create table sync_history(
id uuid primary key not null default gen_random_uuid() unique,
created_at timestamp not null default timezone('utc', now()),
failure boolean not null
);`)
	if err != nil {
		return eris.Wrap(err, "failed to create sync history table")
	}

	_, err = m.conn.Exec(m.ctx, `insert into migration_history (migration) values(3);`)
	if err != nil {
		return eris.Wrap(err, "failed to insert migration into migration history")
	}

	return nil
}

func (m *Migration) checkIfMigrationHistoryExist() (bool, error) {
	q := m.conn.QueryRow(
		m.ctx,
		`
select exists(
select 1 
from information_schema.tables 
where table_schema='public' and table_name='migration_history'
);`)

	var res bool

	err := q.Scan(&res)
	if err != nil {
		return false, eris.Wrap(err, "failed to check if migration_history exist")
	}

	return res, nil
}

func (m *Migration) getMigrationNumber() (int, error) {
	q := m.conn.QueryRow(m.ctx, `select coalesce(max(migration), 0) from migration_history;`)

	var res int

	err := q.Scan(&res)
	if err != nil {
		return 0, eris.Wrap(err, "failed to get the latest migration done from database")
	}

	return res, nil
}

func (m *Migration) GetCurrentMigration() (int, error) {
	e, err := m.checkIfMigrationHistoryExist()
	if err != nil{
		return 0, eris.Wrap(err, "failed to check if migration history exists")
	}

	if !e {
		return 0, nil
	}

	migrationNumber, err := m.getMigrationNumber()
	if err != nil {
		return 0, eris.Wrap(err, "failed to check migration version")
	}

	return migrationNumber, nil
}

// Migrate runs migrations.
func (m *Migration) Migrate() error {
	exist, err := m.checkIfMigrationHistoryExist()
	if err != nil {
		return err
	}

	if !exist {
		err = m.v1()
		if err != nil {
			return eris.Wrap(err, "migration to v1 failed")
		}
	}

	migrationNumber, err := m.getMigrationNumber()
	if err != nil {
		return err
	}

	if migrationNumber == 0 {
		err = m.v1()
		if err != nil {
			return eris.Wrap(err, "migration to v1 failed")
		}
		migrationNumber = 1
	}

	if migrationNumber == 1 {
		err = m.v2()
		if err != nil {
			return eris.Wrap(err, "migration to v2 failed")
		}
		migrationNumber = 2
	}

	if migrationNumber == 2 {
		err = m.v3()
		if err != nil {
			return eris.Wrap(err, "migration to v3 failed")
		}
		migrationNumber = 3
	}

	return nil
}
