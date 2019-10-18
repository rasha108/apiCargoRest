-- +migrate Up
-- +migrate StatementBegin

-- object: public.users | type: TABLE --
-- DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users (
	id bigserial NOT NULL,
	first_name varchar(64),
	last_name varchar(128),
	email varchar(128) NOT NULL,
	encrypted_password varchar NOT NULL,
	CONSTRAINT first_name_last_name_uq_idx UNIQUE (first_name, last_name),
    CONSTRAINT email_uq_idx UNIQUE (email),
	CONSTRAINT create_users_pk PRIMARY KEY (id)
);
-- ddl-end --
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE public.users;
-- +migrate StatementEnd