CREATE TABLE IF NOT EXISTS public.accounts
(
    id integer NOT NULL DEFAULT nextval('accounts_id_seq'::regclass),
    user_id integer,
    account_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    balance numeric(15,2) DEFAULT 0.00,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT accounts_pkey PRIMARY KEY (id),
    CONSTRAINT accounts_account_number_key UNIQUE (account_number),
    CONSTRAINT accounts_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.accounts
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.cards
(
    id integer NOT NULL DEFAULT nextval('cards_id_seq'::regclass),
    account_id integer,
    card_number_encrypted text COLLATE pg_catalog."default" NOT NULL,
    expiry_date_encrypted text COLLATE pg_catalog."default" NOT NULL,
    cvv_hash character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT cards_pkey PRIMARY KEY (id),
    CONSTRAINT cards_account_id_fkey FOREIGN KEY (account_id)
        REFERENCES public.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.cards
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.credit_payments
(
    id integer NOT NULL DEFAULT nextval('credit_payments_id_seq'::regclass),
    credit_id integer NOT NULL,
    due_date timestamp without time zone NOT NULL,
    amount numeric(15,2) NOT NULL,
    paid boolean NOT NULL DEFAULT false,
    paid_at timestamp without time zone,
    fine numeric(15,2) NOT NULL DEFAULT 0,
    CONSTRAINT credit_payments_pkey PRIMARY KEY (id),
    CONSTRAINT credit_payments_credit_id_fkey FOREIGN KEY (credit_id)
        REFERENCES public.credits (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.credit_payments
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.credits
(
    id integer NOT NULL DEFAULT nextval('credits_id_seq'::regclass),
    user_id integer NOT NULL,
    account_id integer NOT NULL,
    amount numeric(15,2) NOT NULL,
    rate numeric(5,2) NOT NULL,
    months integer NOT NULL,
    monthly_payment numeric(15,2) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status character varying(16) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT credits_pkey PRIMARY KEY (id),
    CONSTRAINT credits_account_id_fkey FOREIGN KEY (account_id)
        REFERENCES public.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT credits_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.credits
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_username_key UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;