-- +migrate Up
-- +migrate StatementBegin

-- object: public.create_users | type: TABLE --
-- DROP TABLE IF EXISTS public.create_users CASCADE;
CREATE TABLE public.create_users (
	id bigserial NOT NULL,
	first_name varchar(64) NOT NULL,
	last_name varchar(128) NOT NULL,
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
ALTER TABLE public.create_users OWNER TO postgres;
-- ddl-end --
-- +migrate StatementEnd


