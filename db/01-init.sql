CREATE SEQUENCE IF NOT EXISTS public.hello_id_seq
  INCREMENT 1
  START 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  CACHE 1
  CYCLE;

CREATE TABLE IF NOT EXISTS public.hello
(
    id integer NOT NULL DEFAULT nextval('hello_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT hello_pkey PRIMARY KEY (id)
)