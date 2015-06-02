-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
        active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY NOT NULL,
	login TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	role INT references roles(id),
	active BOOL DEFAULT TRUE,
        name TEXT NOT NULL,
	mail TEXT,
        notify_by_mail BOOL DEFAULT FALSE,
        notify_by_telegram BOOL DEFAULT FALSE,
        gitlab_token TEXT
);

CREATE TABLE IF NOT EXISTS departments (
	id SERIAL PRIMARY KEY NOT NULL,
	title TEXT UNIQUE NOT NULL,
	active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS candidates (
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	lname TEXT,
	note TEXT,
	active BOOL DEFAULT TRUE,
        img TEXT, 
        fname TEXT, 
        addr TEXT NOT NULL,
        married BOOL NOT NULL,
        ref_dep INT references departments(id),
        salary  NUMERIC(11,2) NOT NULL,
        currency TEXT,
        phone TEXT,
        email TEXT
);

CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    ref_cand INT references candidates(id),
    active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS criteria (
	id SERIAL PRIMARY KEY NOT NULL,
	title TEXT UNIQUE NOT NULL,
	ref_dep INT references departments(id),
	active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS ratings (
	id SERIAL PRIMARY KEY NOT NULL,
	ref_crit INT references criteria(id),
	ref_cand INT references candidates(id),
	value INT NOT NULL,
	active BOOL DEFAULT TRUE,
        ref_user INT references users(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id SERIAL PRIMARY KEY NOT NULL,
	ref_user INT references users(id),
	ref_cand INT references candidates(id),
	created_dt timestamp without time zone DEFAULT NOW(),
	comment TEXT,
	active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS files (
	id SERIAL PRIMARY KEY NOT NULL,
	ref_cand INT references candidates(id),
	url TEXT,
	created timestamp without time zone DEFAULT NOW(),
	active BOOL DEFAULT TRUE,
        name TEXT,
	mime TEXT,
	fsize INT, 
	ref_user INT references users(id)
);

CREATE TABLE IF NOT EXISTS mails (
	id SERIAL PRIMARY KEY NOT NULL,
	to_mail TEXT NOT NULL,
	to_name TEXT,
	to_type TEXT NOT NULL,
	from_mail TEXT NOT NULL,
	from_name TEXT,
	subject TEXT,
	html TEXT,
	text TEXT,
	try INT,
	created timestamp without time zone DEFAULT NOW(),
	updated timestamp without time zone DEFAULT NOW(),
	status BOOL DEFAULT FALSE,
	active BOOL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS langs(
	id SERIAL PRIMARY KEY NOT NULL,
        code TEXT NOT NULL,
	is_default BOOL DEFAULT FALSE,
	active BOOL DEFAULT TRUE
);

INSERT INTO roles (name) SELECT 'Admin' WHERE 'Admin' NOT IN (SELECT name FROM roles);
INSERT INTO roles (name) SELECT 'Manager' WHERE 'Manager' NOT IN (SELECT name FROM roles);
INSERT INTO roles (name) SELECT 'HR' WHERE 'HR' NOT IN (SELECT name FROM roles);
INSERT INTO langs (code, is_default) SELECT 'en-US', true WHERE 'en-US' NOT IN (SELECT code FROM langs);

-- +migrate Down
DROP TABLE IF EXISTS langs;
DROP TABLE IF EXISTS mails;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS candidates;
DROP TABLE IF EXISTS criteria;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
