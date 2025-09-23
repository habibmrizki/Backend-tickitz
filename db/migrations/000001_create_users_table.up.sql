-- public.users definition

-- Drop table

-- DROP TABLE public.users;
CREATE TYPE type_role AS ENUM ('admin', 'user');

-- CREATE TYPE IF NOT EXISTS type_role AS ENUM ('admin', 'user');

CREATE TABLE public.users (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	email varchar(255) NOT NULL,
	"password" text NOT NULL,
	"role" public."type_role" NOT NULL,
	created_at timestamp NULL,
	update_at timestamp NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);