-- public.movie definition

-- Drop table

-- DROP TABLE public.movie;

CREATE TABLE public.movie (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	director_id int4 NULL,
	title varchar(255) NULL,
	synopsis text NULL,
	popularity int4 NULL,
	backdrop_path varchar(255) NULL,
	poster_path varchar(255) NULL,
	duration int4 NULL,
	release_date date NULL,
	created_at timestamp NULL,
	update_at timestamp NULL,
	archived_at timestamp NULL,
	CONSTRAINT movie_pkey PRIMARY KEY (id)
);


-- public.movie foreign keys

ALTER TABLE public.movie ADD CONSTRAINT movie_director_id_fkey FOREIGN KEY (director_id) REFERENCES public.director(id);