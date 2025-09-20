package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/singhJasvinder101/go_blog/internal/config"
)

type Postgres struct {
	DB   *pgxpool.Pool
	conf *config.Config
}

func NewPostgres(cfg *config.Config) (*Postgres, error) {
	db, err := pgxpool.New(context.Background(), cfg.DbConn)
	if err != nil {
		fmt.Printf("unable to create connection pool: %v\n", err)
		return nil, err
	}

	// defer db.Close()
	// dbpool return hone se pehle hi pool close ho jayega.

	return &Postgres{
		DB:   db,
		conf: cfg,
	}, nil
}

func (p *Postgres) InitSchema(ctx context.Context) error {
	_, err := p.DB.Exec(ctx, `
		create table if not exists users (
			id serial primary key,
			name text not null,
			email text unique not null,
			password_hash text not null,
 
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		);

		create table if not exists posts (
			id serial primary key,
			title text not null,
			description text not null,

			user_id integer references users(id) on delete no action,

			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		);

		create table if not exists comments (
			id serial primary key,
			content text not null,

			user_id integer references users(id) on delete set null,
			post_id integer references posts(id) on delete cascade,

			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		);

		create table if not exists likes (
			id serial primary key,

			user_id int references users(id) on delete cascade,
			post_id int references posts(id) on delete cascade,

			created_at timestamp default current_timestamp,
			unique(user_id, post_id)
		);

		create table if not exists tags (
			id serial primary key,

			name text unique not null,

			created_at timestamp default current_timestamp
		);

		create table if not exists post_tags (
			id serial primary key,

			post_id int references posts(id) on delete cascade,
			tag_id int references tags(id) on delete cascade,

			created_at timestamp default current_timestamp,
			unique(post_id, tag_id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
