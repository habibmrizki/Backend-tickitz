-- public.order_seats definition

-- Drop table

-- DROP TABLE public.order_seats;

CREATE TABLE public.order_seats (
	orders_id int4 NOT NULL,
	seats_id int4 NOT NULL,
	CONSTRAINT order_seats_pkey PRIMARY KEY (orders_id, seats_id)
);


-- public.order_seats foreign keys

ALTER TABLE public.order_seats ADD CONSTRAINT order_seats_orders_id_fkey FOREIGN KEY (orders_id) REFERENCES public.orders(id) ON DELETE CASCADE;
ALTER TABLE public.order_seats ADD CONSTRAINT order_seats_seats_id_fkey FOREIGN KEY (seats_id) REFERENCES public.seats(id);