-- public.movies_genre definition

-- Drop table

-- DROP TABLE public.movies_genre;

CREATE TABLE public.movies_genre (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	movie_id int4 NULL,
	genre_id int4 NULL,
	CONSTRAINT movies_genre_pkey PRIMARY KEY (id)
);


-- public.movies_genre foreign keys

ALTER TABLE public.movies_genre ADD CONSTRAINT movies_genre_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genre(id);
ALTER TABLE public.movies_genre ADD CONSTRAINT movies_genre_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movie(id) ON DELETE CASCADE;