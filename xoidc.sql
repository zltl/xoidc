--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Debian 15.4-1.pgdg120+1)
-- Dumped by pg_dump version 15.4 (Debian 15.4-1.pgdg120+1)

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
-- Name: client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client (
    id bigint DEFAULT 0 NOT NULL,
    secret character varying(500) DEFAULT ''::character varying NOT NULL,
    application_type integer DEFAULT 0 NOT NULL,
    auth_method character varying(200) DEFAULT ''::character varying NOT NULL,
    access_token_type integer DEFAULT 0 NOT NULL,
    dev_mode boolean DEFAULT false NOT NULL,
    id_token_userinfo_claims_assertion boolean DEFAULT false NOT NULL,
    clock_skew timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
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
1: JWT';


--
-- Name: client_grant_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client_grant_types (
    client_id bigint NOT NULL,
    grant_type character varying(200) NOT NULL
);


ALTER TABLE public.client_grant_types OWNER TO postgres;

--
-- Name: client_redirect_uris; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client_redirect_uris (
    client_id bigint NOT NULL,
    redirect_uri text NOT NULL
);


ALTER TABLE public.client_redirect_uris OWNER TO postgres;

--
-- Name: client_response_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client_response_types (
    client_id bigint NOT NULL,
    response_type character varying(200) NOT NULL
);


ALTER TABLE public.client_response_types OWNER TO postgres;

--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id bigint DEFAULT 0 NOT NULL,
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
    namespace bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Data for Name: client; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client (id, secret, application_type, auth_method, access_token_type, dev_mode, id_token_userinfo_claims_assertion, clock_skew, created_at, updated_at) FROM stdin;
0	123456	0	none	0	f	f	2023-08-13 11:16:34.000629+00	2023-08-13 11:16:34.000629+00	2023-08-13 11:16:34.000629+00
\.


--
-- Data for Name: client_grant_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client_grant_types (client_id, grant_type) FROM stdin;
\.


--
-- Data for Name: client_redirect_uris; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client_redirect_uris (client_id, redirect_uri) FROM stdin;
\.


--
-- Data for Name: client_response_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client_response_types (client_id, response_type) FROM stdin;
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."user" (id, username, password, nickname, given_name, family_name, middle_name, preferred_username, profile, picture, website, email, email_verified, gender, birthdate, zoneinfo, locale, phone_number, phone_number_verified, address, updated_at, namespace) FROM stdin;
1	test	$argon2id$v=19$m=19456,t=2,p=1$Z0CCH0FfcFXsHnxDTfvXXQ$KqH1dzTda/0Mrj63scfybiTVGCjHxjmZHTfwMpRyOSc	test	test	test	test	test				test@email.com	f		2023-08-13				f		2023-08-13 10:33:13.160209+00	0
\.


--
-- Name: client_grant_types client_grant_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_grant_types
    ADD CONSTRAINT client_grant_types_pkey PRIMARY KEY (client_id, grant_type);


--
-- Name: client client_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_pkey PRIMARY KEY (id_token_userinfo_claims_assertion);


--
-- Name: client_redirect_uris client_redirect_uris_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_redirect_uris
    ADD CONSTRAINT client_redirect_uris_pkey PRIMARY KEY (client_id, redirect_uri);


--
-- Name: client_response_types client_response_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_response_types
    ADD CONSTRAINT client_response_type_pkey PRIMARY KEY (client_id, response_type);


--
-- Name: user idid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT idid PRIMARY KEY (id, namespace);


--
-- PostgreSQL database dump complete
--

