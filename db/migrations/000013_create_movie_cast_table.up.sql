-- public.movie_cast definition

-- Drop table

-- DROP TABLE public.movie_cast;

CREATE TABLE public.movie_cast (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	movie_id int4 NULL,
	cast_id int4 NULL,
	CONSTRAINT movie_cast_pkey PRIMARY KEY (id)
);


-- public.movie_cast foreign keys

ALTER TABLE public.movie_cast ADD CONSTRAINT movie_cast_cast_id_fkey FOREIGN KEY (cast_id) REFERENCES public."cast"(id);
ALTER TABLE public.movie_cast ADD CONSTRAINT movie_cast_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movie(id) ON DELETE CASCADE;