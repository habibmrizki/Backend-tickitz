-- public.profile definition

-- Drop table

-- DROP TABLE public.profile;

CREATE TABLE public.profile (
	users_id int4 NOT NULL,
	first_name varchar(255) NULL,
	last_name varchar(255) NULL,
	phone_number varchar(50) NULL,
	profile_image varchar(255) NULL,
	point int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	modified_at timestamp NULL,
	CONSTRAINT profile_pkey PRIMARY KEY (users_id)
);


-- public.profile foreign keys

ALTER TABLE public.profile ADD CONSTRAINT profile_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id) ON DELETE CASCADE;