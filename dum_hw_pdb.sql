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
                              id integer NOT NULL,
                              slug text NOT NULL,
                              threads integer DEFAULT 0 NOT NULL,
                              posts integer DEFAULT 0 NOT NULL,
                              title text NOT NULL,
                              "user" text NOT NULL
);


ALTER TABLE public.forum OWNER TO postgres;

--
-- Name: forum_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.forum_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_id_seq OWNER TO postgres;

--
-- Name: forum_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.forum_id_seq OWNED BY public.forum.id;


--
-- Name: forum_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.forum_user (
                                   forum_id integer NOT NULL,
                                   user_id integer NOT NULL
);


ALTER TABLE public.forum_user OWNER TO postgres;

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
                             thread integer NOT NULL,
                             path bigint[] DEFAULT '{0}'::bigint[] NOT NULL
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
                               votes integer DEFAULT 0
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
                               id integer NOT NULL,
                               nick_name text NOT NULL,
                               email text NOT NULL,
                               full_name text NOT NULL,
                               about text
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id_seq OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: vote; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.vote (
                             user_id integer,
                             voice integer NOT NULL,
                             thread_id integer NOT NULL
);


ALTER TABLE public.vote OWNER TO postgres;

--
-- Name: forum id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forum ALTER COLUMN id SET DEFAULT nextval('public.forum_id_seq'::regclass);


--
-- Name: post id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post ALTER COLUMN id SET DEFAULT nextval('public.post_id_seq'::regclass);


--
-- Name: thread id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thread ALTER COLUMN id SET DEFAULT nextval('public.thread_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);

--
-- Name: forum_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.forum_id_seq', 1, false);


--
-- Name: post_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.post_id_seq', 1, false);


--
-- Name: thread_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.thread_id_seq', 1, false);


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id_seq', 1, false);


--
-- Name: forum forum_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forum
    ADD CONSTRAINT forum_pk PRIMARY KEY (id);


--
-- Name: forum_user forum_user_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forum_user
    ADD CONSTRAINT forum_user_pk PRIMARY KEY (forum_id, user_id);


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
-- Name: user user_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pk PRIMARY KEY (id);


--
-- Name: forum_slug_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX forum_slug_uindex ON public.forum USING btree (lower(slug));


--
-- Name: post_author_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_author_forum_index ON public.post USING btree (lower(author), lower(forum));


--
-- Name: post_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_forum_index ON public.post USING btree (lower(forum));


--
-- Name: post_parent_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_parent_index ON public.post USING btree (parent);


--
-- Name: post_path_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_path_index ON public.post USING gin (path);


--
-- Name: post_thread_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX post_thread_index ON public.post USING btree (thread);


--
-- Name: thread_forum_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX thread_forum_index ON public.thread USING btree (lower(forum));


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

CREATE UNIQUE INDEX user_email_uindex ON public."user" USING btree (lower(email));


--
-- Name: user_nick_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_nick_name_uindex ON public."user" USING btree (lower(nick_name));


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

