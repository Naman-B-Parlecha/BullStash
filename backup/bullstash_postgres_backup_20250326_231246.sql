--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Ubuntu 17.4-1.pgdg24.04+2)
-- Dumped by pg_dump version 17.4 (Ubuntu 17.4-1.pgdg24.04+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: students; Type: TABLE; Schema: public; Owner: bullstash_user_1
--

CREATE TABLE public.students (
    name character varying(255),
    age integer
);


ALTER TABLE public.students OWNER TO bullstash_user_1;

--
-- Data for Name: students; Type: TABLE DATA; Schema: public; Owner: bullstash_user_1
--

COPY public.students (name, age) FROM stdin;
A	\N
B	\N
C	\N
\.


--
-- PostgreSQL database dump complete
--

