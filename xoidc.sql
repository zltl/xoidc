--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Debian 15.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auth_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.auth_request (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    creation_date timestamp with time zone DEFAULT now() NOT NULL,
    done boolean DEFAULT false NOT NULL,
    auth_time timestamp with time zone DEFAULT now() NOT NULL,
    content text DEFAULT ''::text NOT NULL,
    namespace_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public.auth_request OWNER TO postgres;

--
-- Name: client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    secret character varying(500) DEFAULT ''::character varying NOT NULL,
    redirect_uris text[] DEFAULT '{}'::text[] NOT NULL,
    application_type integer DEFAULT 0 NOT NULL,
    auth_method character varying(200) DEFAULT ''::character varying NOT NULL,
    response_types character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL,
    access_token_type integer DEFAULT 0 NOT NULL,
    dev_mode boolean DEFAULT true NOT NULL,
    id_token_user_info_claims_assertion boolean DEFAULT true NOT NULL,
    clock_skew interval(6) DEFAULT '00:00:00'::interval(6) NOT NULL,
    post_logout_redirect_uri_globs text[] DEFAULT '{}'::text[] NOT NULL,
    redirect_uri_globs text[] DEFAULT '{}'::text[] NOT NULL,
    user_namespace_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL,
    grant_types character varying[] DEFAULT '{}'::character varying[] NOT NULL
);


ALTER TABLE public.client OWNER TO postgres;

--
-- Name: COLUMN client.application_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client.application_type IS '0: web
1: user_agent
2: native';


--
-- Name: COLUMN client.auth_method; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client.auth_method IS 'client_secret_basic
client_secret_post
none
private_key_jwt';


--
-- Name: COLUMN client.access_token_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client.access_token_type IS '0: bearer
1:jwt';


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    username character varying(200) DEFAULT ''::character varying NOT NULL,
    password text DEFAULT ''::text NOT NULL,
    nickname character varying(200) DEFAULT ''::character varying NOT NULL,
    given_name character varying(200) DEFAULT ''::character varying NOT NULL,
    family_name character varying(200) DEFAULT ''::character varying NOT NULL,
    middle_name character varying(200) DEFAULT ''::character varying NOT NULL,
    preferred_username character varying(200) DEFAULT ''::character varying NOT NULL,
    profile text DEFAULT ''::text NOT NULL,
    picture text DEFAULT ''::text NOT NULL,
    website text DEFAULT ''::text NOT NULL,
    email character varying(200) DEFAULT ''::character varying NOT NULL,
    email_verified boolean DEFAULT false NOT NULL,
    gender character varying(40) DEFAULT ''::character varying NOT NULL,
    birthdate date DEFAULT now() NOT NULL,
    zoneinfo character varying(40) DEFAULT ''::character varying NOT NULL,
    locale character varying(60) DEFAULT ''::character varying NOT NULL,
    phone_number character varying(100) DEFAULT ''::character varying NOT NULL,
    phone_number_verified boolean DEFAULT false NOT NULL,
    address character varying(200) DEFAULT ''::character varying NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    namespace_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Data for Name: auth_request; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.auth_request (id, creation_date, done, auth_time, content, namespace_id, user_id) FROM stdin;
\.


--
-- Data for Name: client; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client (id, secret, redirect_uris, application_type, auth_method, response_types, access_token_type, dev_mode, id_token_user_info_claims_assertion, clock_skew, post_logout_redirect_uri_globs, redirect_uri_globs, user_namespace_id, grant_types) FROM stdin;
674fc25c-7772-45e3-835d-3b77b16a2937	123456	{custom://auth/callback,http://localhost:9999/auth/callback,http://localhost/auth/callback}	0	client_secret_basic	{code}	0	t	t	01:05:00	{}	{}	00000000-0000-0000-0000-000000000000	{authorization_code,refresh_token,urn:ietf:params:oauth:grant-type:token-exchange}
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."user" (username, password, nickname, given_name, family_name, middle_name, preferred_username, profile, picture, website, email, email_verified, gender, birthdate, zoneinfo, locale, phone_number, phone_number_verified, address, updated_at, namespace_id, id) FROM stdin;
test	$argon2id$v=19$m=19456,t=2,p=1$Z0CCH0FfcFXsHnxDTfvXXQ$KqH1dzTda/0Mrj63scfybiTVGCjHxjmZHTfwMpRyOSc	test	test	test	test	test				test@email.com	f		2023-08-13				f		2023-08-13 10:33:13.160209+00	00000000-0000-0000-0000-000000000000	744d9044-f29d-42e8-a65e-e6c52398fa1f
\.


--
-- Name: auth_request auth_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_request
    ADD CONSTRAINT auth_request_pkey PRIMARY KEY (id);


--
-- Name: client client_new_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_new_pkey PRIMARY KEY (id_token_user_info_claims_assertion);


--
-- PostgreSQL database dump complete
--

