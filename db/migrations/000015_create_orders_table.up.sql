-- public.orders definition

-- Drop table

-- DROP TABLE public.orders;

CREATE TABLE public.orders (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	users_id int4 NULL,
	schedule_id int4 NULL,
	payment_method_id int4 NULL,
	total_price int4 NULL,
	ispaid bool NULL,
	created_at timestamp NULL,
	update_at timestamp NULL,
	full_name varchar(255) NULL,
	email varchar(255) NULL,
	phone_number varchar(20) NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id)
);


-- public.orders foreign keys

ALTER TABLE public.orders ADD CONSTRAINT orders_payment_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_method(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES public.schedule(id) ON DELETE CASCADE;
ALTER TABLE public.orders ADD CONSTRAINT orders_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id);