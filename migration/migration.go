package migration

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)
type Migration struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

func New(ctx context.Context, conn *pgxpool.Pool) *Migration{
	return &Migration{
		ctx: ctx,
		conn: conn,
	}
}

func (m *Migration) v1() error{
	_, err := m.conn.Exec(m.ctx, `CREATE EXTENSION pgcrypto;`)
	if err != nil{
		return err
	}

	_, err = m.conn.Exec(m.ctx, `create table migration_history(migration int not null unique);`)
	if err != nil{
		return err
	}

	_, err = m.conn.Exec(m.ctx, `
create table files(
id uuid primary key not null default gen_random_uuid() unique, 
processed_at timestamp not null default timezone('utc', now()),
opened_at timestamp not null default timezone('utc', now()),
path text not null, 
title text not null, title_tokens tsvector not null, content text not null, content_tokens tsvector not null);`)
	if err != nil{
		return err
	}

	_, err = m.conn.Exec(m.ctx, `insert into migration_history (migration) values(1);`)
	if err != nil{
		return err
	}

	return nil
}

func (m *Migration) checkIfMigrationHistoryExist() (bool, error){
	q := m.conn.QueryRow(m.ctx, `select exists(select 1 from information_schema.tables where table_schema='public' and table_name='migration_history');`)
	var res bool
	err := q.Scan(&res)
	if err != nil{
		return false, err
	}

	return res, err

}

func (m *Migration) getMigrationNumber() (int, error){
	q := m.conn.QueryRow(m.ctx, `select coalesce(max(migration), 0) from migration_history;`)
	var res int
	err := q.Scan(&res)
	if err != nil{
		return 0, err
	}

	return res, err

}

func (m *Migration) Migrate() error{
	exist, err := m.checkIfMigrationHistoryExist()
	if err != nil{
		return err
	}

	if !exist{
		err = m.v1()
		if err != nil{
			return err
		}
	}

	migrationNumber, err := m.getMigrationNumber()
	if err != nil{
		return err
	}

	if migrationNumber == 0{
		err = m.v1()
		if err != nil{
			return err
		}
	}

	return nil
}
/*
create table runs (run_at timestamp not null unique);
select coalesce(max(run_at), TO_TIMESTAMP('2020-01-01 00:00:00', 'YYYY-MM-DD HH24:MI:SS')) from runs;
*/