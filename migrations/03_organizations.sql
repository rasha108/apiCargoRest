-- +migrate Up
-- +migrate StatementBegin

-- object: public.users_organizations | type: TABLE --
-- DROP TABLE IF EXISTS public.users_organizations CASCADE;
CREATE TABLE public.users_organizations (
	id bigserial NOT NULL,
	user_id bigint,
	org_id bigint,
	CONSTRAINT users_organizations_user_id_org_id UNIQUE (user_id,org_id),
	CONSTRAINT users_organizations_pk PRIMARY KEY (id)

);
-- ddl-end --

-- object: public.organizations | type: TABLE --
-- DROP TABLE IF EXISTS public.organizations CASCADE;
CREATE TABLE public.organizations (
	id bigserial NOT NULL,
	org_name varchar(128),
	phone varchar(64),
	address varchar(256),
	email varchar(128),
	CONSTRAINT organizations_pk PRIMARY KEY (id)

);

-- object: users_organizations_org_id_fk | type: CONSTRAINT --
-- ALTER TABLE public.users_organizations DROP CONSTRAINT IF EXISTS users_organizations_org_id_fk CASCADE;
ALTER TABLE public.users_organizations ADD CONSTRAINT users_organizations_org_id_fk FOREIGN KEY (org_id)
REFERENCES public.organizations (id) MATCH SIMPLE
ON DELETE CASCADE ON UPDATE RESTRICT;
-- ddl-end --

-- object: users_organizations_user_id_fk | type: CONSTRAINT --
-- ALTER TABLE public.users_organizations DROP CONSTRAINT IF EXISTS users_organizations_user_id_fk CASCADE;
ALTER TABLE public.users_organizations ADD CONSTRAINT users_organizations_user_id_fk FOREIGN KEY (user_id)
REFERENCES public.users (id) MATCH FULL
ON DELETE CASCADE ON UPDATE RESTRICT;
-- ddl-end --
-- ddl-end --
-- +migrate StatementEnd


-- +migrate Down
-- +migrate StatementBegin
DROP TABLE public.organizations;
DROP TABLE public.users_organizations;
-- +migrate StatementEnd

