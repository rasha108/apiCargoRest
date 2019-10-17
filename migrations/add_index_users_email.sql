-- +migrate Up
-- +migrate StatementBegin

-- object: index_users_email
CREATE UNIQUE INDEX index_users_email ON public.create_users (lower(email));
-- ddl-end --
-- +migrate StatementEnd


-- +migrate Down
-- +migrate StatementBegin
DROP INDEX index_users_email;
-- +migrate StatementEnd