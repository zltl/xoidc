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
    grant_types character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL
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
-- Name: code_request_id; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.code_request_id (
    code character varying(256) NOT NULL,
    request_id uuid NOT NULL,
    create_time timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.code_request_id OWNER TO postgres;

--
-- Name: refresh_token; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_token (
    id character varying(200) DEFAULT ''::character varying NOT NULL,
    token text DEFAULT ''::text NOT NULL,
    auth_time timestamp(3) without time zone DEFAULT now() NOT NULL,
    amr character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL,
    audience character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL,
    user_id character varying(200) DEFAULT ''::character varying NOT NULL,
    application_id character varying(200) DEFAULT ''::character varying NOT NULL,
    expiration timestamp(3) without time zone DEFAULT now() NOT NULL,
    scopes character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL
);


ALTER TABLE public.refresh_token OWNER TO postgres;

--
-- Name: token; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.token (
    id character varying(200) DEFAULT ''::character varying NOT NULL,
    audience character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL,
    expiration timestamp(3) without time zone DEFAULT now() NOT NULL,
    scopes character varying(200)[] DEFAULT '{}'::character varying[] NOT NULL,
    application_id uuid DEFAULT gen_random_uuid() NOT NULL,
    subject uuid DEFAULT gen_random_uuid() NOT NULL,
    refresh_token_id uuid DEFAULT gen_random_uuid() NOT NULL
);


ALTER TABLE public.token OWNER TO postgres;

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
30fe0ae9-d940-4d2a-a4d8-8c539622104e	2023-11-26 07:06:05.95332+00	f	0001-01-01 00:00:00+00	{"scope":"openid profile","response_type":"code","client_id":"674fc25c-7772-45e3-835d-3b77b16a2937","redirect_uri":"http://localhost:9999/auth/callback","state":"64a85d42-e863-4407-a923-5af760bec3a2","nonce":"","response_mode":"","display":"","prompt":"Welcome back!","max_age":null,"ui_locales":null,"id_token_hint":"","login_hint":"","acr_values":"","code_challenge":"","code_challenge_method":"","RequestParam":""}	00000000-0000-0000-0000-000000000000	00000000-0000-0000-0000-000000000000
5f141e2c-4bfb-449f-b082-21752c4080f9	2023-12-02 09:36:45.410367+00	t	0001-01-01 00:00:00+00	{"scope":"openid profile","response_type":"code","client_id":"674fc25c-7772-45e3-835d-3b77b16a2937","redirect_uri":"http://localhost:9999/auth/callback","state":"e37f7a11-c7a7-47b5-85d9-adcde28bd31a","nonce":"","response_mode":"","display":"","prompt":"Welcome back!","max_age":null,"ui_locales":null,"id_token_hint":"","login_hint":"","acr_values":"","code_challenge":"","code_challenge_method":"","RequestParam":""}	00000000-0000-0000-0000-000000000000	744d9044-f29d-42e8-a65e-e6c52398fa1f
7537efeb-31d8-41f6-a92f-c9f1567cc347	2023-11-26 06:53:34.610723+00	t	0001-01-01 00:00:00+00	{"scope":"openid profile","response_type":"code","client_id":"674fc25c-7772-45e3-835d-3b77b16a2937","redirect_uri":"http://localhost:9999/auth/callback","state":"e0dc027f-7422-4ce0-94c6-482015208e83","nonce":"","response_mode":"","display":"","prompt":"Welcome back!","max_age":null,"ui_locales":null,"id_token_hint":"","login_hint":"","acr_values":"","code_challenge":"","code_challenge_method":"","RequestParam":""}	00000000-0000-0000-0000-000000000000	744d9044-f29d-42e8-a65e-e6c52398fa1f
f8e80ace-06e0-4b73-9a20-e7f873a588f7	2023-12-02 11:20:43.741457+00	t	0001-01-01 00:00:00+00	{"scope":"openid profile","response_type":"code","client_id":"674fc25c-7772-45e3-835d-3b77b16a2937","redirect_uri":"http://localhost:9999/auth/callback","state":"3324b9ad-bafc-4880-b294-f797dd706b8d","nonce":"","response_mode":"","display":"","prompt":"Welcome back!","max_age":null,"ui_locales":null,"id_token_hint":"","login_hint":"","acr_values":"","code_challenge":"","code_challenge_method":"","RequestParam":""}	00000000-0000-0000-0000-000000000000	744d9044-f29d-42e8-a65e-e6c52398fa1f
22a55ec0-0171-44f8-85b5-99bdfcfc9318	2023-12-02 09:45:47.652939+00	f	0001-01-01 00:00:00+00	{"scope":"openid profile","response_type":"code","client_id":"674fc25c-7772-45e3-835d-3b77b16a2937","redirect_uri":"http://localhost:9999/auth/callback","state":"4d678a56-c938-425f-9a50-3287e8728ee5","nonce":"","response_mode":"","display":"","prompt":"Welcome back!","max_age":null,"ui_locales":null,"id_token_hint":"","login_hint":"","acr_values":"","code_challenge":"","code_challenge_method":"","RequestParam":""}	00000000-0000-0000-0000-000000000000	00000000-0000-0000-0000-000000000000
\.


--
-- Data for Name: client; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.client (id, secret, redirect_uris, application_type, auth_method, response_types, access_token_type, dev_mode, id_token_user_info_claims_assertion, clock_skew, post_logout_redirect_uri_globs, redirect_uri_globs, user_namespace_id, grant_types, name) FROM stdin;
674fc25c-7772-45e3-835d-3b77b16a2937	123456	{custom://auth/callback,http://localhost:9999/auth/callback,http://localhost/auth/callback}	0	client_secret_basic	{code}	0	t	t	01:05:00	{}	{}	00000000-0000-0000-0000-000000000000	{authorization_code,refresh_token,urn:ietf:params:oauth:grant-type:token-exchange}	
\.


--
-- Data for Name: code_request_id; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.code_request_id (code, request_id, create_time) FROM stdin;
s2o3ld5MBb9AaB3vKlWPDo9EnfoGwgqMgHhwEnzCovgpZ9rqgQMFhr44W4Hh-lOFsiVfUw	5f141e2c-4bfb-449f-b082-21752c4080f9	2023-12-02 09:36:50.965+00
r2dx_0-YIsxVCD7uNisPnKbQtMCfzu89tqIGLLu2QsdnQaGq0qDs0WDWJCoQKpUORGcPNg	f8e80ace-06e0-4b73-9a20-e7f873a588f7	2023-12-02 11:20:52.255+00
\.


--
-- Data for Name: refresh_token; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.refresh_token (id, token, auth_time, amr, audience, user_id, application_id, expiration, scopes) FROM stdin;
\.


--
-- Data for Name: token; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.token (id, audience, expiration, scopes, application_id, subject, refresh_token_id) FROM stdin;
78488baa-54f7-464e-ac64-a953d5cb182c	{674fc25c-7772-45e3-835d-3b77b16a2937}	2023-12-02 19:28:54.796	{openid,profile}	674fc25c-7772-45e3-835d-3b77b16a2937	744d9044-f29d-42e8-a65e-e6c52398fa1f	00000000-0000-0000-0000-000000000000
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
-- Name: code_request_id code_request_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.code_request_id
    ADD CONSTRAINT code_request_id_pkey PRIMARY KEY (code, request_id);


--
-- Name: refresh_token refresh_token_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT refresh_token_pkey PRIMARY KEY (id);


--
-- Name: token token_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.token
    ADD CONSTRAINT token_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

