--
-- PostgreSQL database dump
--

-- Dumped from database version 10.10
-- Dumped by pg_dump version 10.10

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

SET default_with_oids = false;

--
-- Name: forum; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.forum (
                              title text NOT NULL,
                              "user" text NOT NULL,
                              slug text NOT NULL
);


ALTER TABLE public.forum OWNER TO postgres;

--
-- Name: post; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.post (
                             id integer NOT NULL,
                             author text NOT NULL,
                             created text NOT NULL,
                             forum text NOT NULL,
                             is_edited boolean DEFAULT false NOT NULL,
                             message text NOT NULL,
                             parent integer DEFAULT 0 NOT NULL,
                             thread integer NOT NULL
);


ALTER TABLE public.post OWNER TO postgres;

--
-- Name: post_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.post_id_seq OWNER TO postgres;

--
-- Name: post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.post_id_seq OWNED BY public.post.id;


--
-- Name: thread; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.thread (
                               id integer NOT NULL,
                               author text NOT NULL,
                               created timestamp with time zone DEFAULT now() NOT NULL,
                               forum text NOT NULL,
                               message text NOT NULL,
                               slug text,
                               title text NOT NULL,
                               votes integer
);


ALTER TABLE public.thread OWNER TO postgres;

--
-- Name: thread_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.thread_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.thread_id_seq OWNER TO postgres;

--
-- Name: thread_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.thread_id_seq OWNED BY public.thread.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
                               nick_name text NOT NULL,
                               full_name text NOT NULL,
                               email text NOT NULL,
                               about text
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: vote; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.vote (
                             nick_name text NOT NULL,
                             voice integer NOT NULL,
                             thread_id integer NOT NULL
);


ALTER TABLE public.vote OWNER TO postgres;

--
-- Name: post id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post ALTER COLUMN id SET DEFAULT nextval('public.post_id_seq'::regclass);


--
-- Name: thread id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread ALTER COLUMN id SET DEFAULT nextval('public.thread_id_seq'::regclass);


--
-- Name: post_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.post_id_seq', 389862, true);


--
-- Name: thread_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.thread_id_seq', 10148, true);


--
-- Name: post post_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_pk PRIMARY KEY (id);


--
-- Name: thread thread_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread
    ADD CONSTRAINT thread_pk PRIMARY KEY (id);


--
-- Name: vote vote_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vote
    ADD CONSTRAINT vote_pk PRIMARY KEY (nick_name, thread_id);


--
-- Name: forum_lower(slug)_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "forum_lower(slug)_index" ON public.forum USING btree (lower(slug));


--
-- Name: forum_slug_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX forum_slug_uindex ON public.forum USING btree (slug);


--
-- Name: thread_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX thread_id_uindex ON public.thread USING btree (id);


--
-- Name: thread_slug_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_slug_index ON public.thread USING btree (lower(slug));


--
-- Name: user_email_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_email_uindex ON public."user" USING btree (email);


--
-- Name: user_nick_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_nick_name_uindex ON public."user" USING btree (nick_name);


--
-- Name: vote_nick_name_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX vote_nick_name_index ON public.vote USING btree (lower(nick_name), thread_id);


--
-- PostgreSQL database dump complete
--

