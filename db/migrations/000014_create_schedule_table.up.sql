-- public.schedule definition

-- Drop table

-- DROP TABLE public.schedule;

CREATE TABLE public.schedule (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	movie_id int4 NULL,
	cinema_id int4 NULL,
	location_id int4 NULL,
	time_id int4 NULL,
	"date" date NULL,
	CONSTRAINT schedule_pkey PRIMARY KEY (id)
);


-- public.schedule foreign keys

ALTER TABLE public.schedule ADD CONSTRAINT schedule_cinema_id_fkey FOREIGN KEY (cinema_id) REFERENCES public.cinema(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_location_id_fkey FOREIGN KEY (location_id) REFERENCES public."location"(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movie(id) ON DELETE CASCADE;
ALTER TABLE public.schedule ADD CONSTRAINT schedule_time_id_fkey FOREIGN KEY (time_id) REFERENCES public."time"(id);