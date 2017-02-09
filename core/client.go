package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type Client struct {
	DB *sql.DB
}

func NewClient() (*Client, error) {
	db, err := sql.Open("postgres", Config.Postgres.Url)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(40)
	db.SetMaxOpenConns(10)

	return &Client{db}, nil
}

func (c *Client) DropTables(names ...string) error {
	smt := fmt.Sprintf("DROP TABLE IF EXISTS %s;", strings.Join(names, ", "))
	if _, err := c.DB.Exec(smt); err != nil {
		return err
	}

	return nil
}

// TODO temporal method to create cgit tables
func (c *Client) CreateCgitTables() error {
	_, err := c.DB.Exec(`CREATE TABLE IF NOT EXISTS cgit (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
        cgit_url varchar(255),
	url varchar(255),
	aliases text[],
	html text
	)`)

	if err != nil {
		return err
	}

	_, err = c.DB.Exec(`CREATE TABLE IF NOT EXISTS cgit_urls (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
        cgit_url varchar(255) UNIQUE NOT NULL
	)`)

	return err
}

// TODO temporal method to create bitbucket table
func (c *Client) CreateBitbucketTable() error {
	_, err := c.DB.Exec(`CREATE TABLE IF NOT EXISTS bitbucket (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	next varchar(255) not null,
	scm varchar(255) not null,
	website varchar(255) not null,
	name varchar(255) not null,
	links jsonb,
	fork_policy varchar(255) not null,
	uuid varchar(255) not null,
	language varchar(255) not null,
	created_on varchar(255) not null,
	parent jsonb,
	full_name varchar(255) not null,
	has_issues boolean not null,
	owner jsonb,
	updated_on varchar(255) not null,
	size int not null,
	type varchar(255) not null,
	slug varchar(255) not null,
	is_private boolean not null,
	description text not null
	)`)

	return err
}

// TODO temporal method to create github table
func (c *Client) CreateGithubTable() error {
	_, err := c.DB.Exec(`CREATE TABLE IF NOT EXISTS github (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	github_id int,
	name varchar(255),
	full_name varchar(511),
	owner jsonb,
	private boolean not null,
	htmlurl varchar(1023),
	description text,
	fork boolean not null
	)`)

	return err
}
