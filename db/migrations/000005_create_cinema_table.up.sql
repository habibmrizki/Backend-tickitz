-- public.cinema definition

-- Drop table

-- DROP TABLE public.cinema;

CREATE TABLE public.cinema (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(50) NOT NULL,
	image_path varchar NULL,
	CONSTRAINT cinema_pkey PRIMARY KEY (id)
);