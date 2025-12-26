--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2
-- Dumped by pg_dump version 17.2

-- Started on 2025-12-26 18:01:56

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
-- TOC entry 226 (class 1259 OID 19393)
-- Name: apps; Type: TABLE; Schema: public; Owner: teammate_search
--

CREATE TABLE public.apps (
    id_app integer NOT NULL,
    app text NOT NULL
);


ALTER TABLE public.apps OWNER TO teammate_search;

--
-- TOC entry 225 (class 1259 OID 19392)
-- Name: apps_id_app_seq; Type: SEQUENCE; Schema: public; Owner: teammate_search
--

CREATE SEQUENCE public.apps_id_app_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.apps_id_app_seq OWNER TO teammate_search;

--
-- TOC entry 4894 (class 0 OID 0)
-- Dependencies: 225
-- Name: apps_id_app_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: teammate_search
--

ALTER SEQUENCE public.apps_id_app_seq OWNED BY public.apps.id_app;


--
-- TOC entry 222 (class 1259 OID 19289)
-- Name: games; Type: TABLE; Schema: public; Owner: teammate_search
--

CREATE TABLE public.games (
    id_game integer NOT NULL,
    game text NOT NULL
);


ALTER TABLE public.games OWNER TO teammate_search;

--
-- TOC entry 221 (class 1259 OID 19288)
-- Name: games_id_game_seq; Type: SEQUENCE; Schema: public; Owner: teammate_search
--

CREATE SEQUENCE public.games_id_game_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.games_id_game_seq OWNER TO teammate_search;

--
-- TOC entry 4895 (class 0 OID 0)
-- Dependencies: 221
-- Name: games_id_game_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: teammate_search
--

ALTER SEQUENCE public.games_id_game_seq OWNED BY public.games.id_game;


--
-- TOC entry 220 (class 1259 OID 19280)
-- Name: genres; Type: TABLE; Schema: public; Owner: teammate_search
--

CREATE TABLE public.genres (
    id_genre integer NOT NULL,
    genre text NOT NULL
);


ALTER TABLE public.genres OWNER TO teammate_search;

--
-- TOC entry 219 (class 1259 OID 19279)
-- Name: genres_id_genre_seq; Type: SEQUENCE; Schema: public; Owner: teammate_search
--

CREATE SEQUENCE public.genres_id_genre_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.genres_id_genre_seq OWNER TO teammate_search;

--
-- TOC entry 4896 (class 0 OID 0)
-- Dependencies: 219
-- Name: genres_id_genre_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: teammate_search
--

ALTER SEQUENCE public.genres_id_genre_seq OWNED BY public.genres.id_genre;


--
-- TOC entry 224 (class 1259 OID 19298)
-- Name: languages; Type: TABLE; Schema: public; Owner: teammate_search
--

CREATE TABLE public.languages (
    id_language integer NOT NULL,
    language text NOT NULL
);


ALTER TABLE public.languages OWNER TO teammate_search;

--
-- TOC entry 223 (class 1259 OID 19297)
-- Name: languages_id_language_seq; Type: SEQUENCE; Schema: public; Owner: teammate_search
--

CREATE SEQUENCE public.languages_id_language_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.languages_id_language_seq OWNER TO teammate_search;

--
-- TOC entry 4897 (class 0 OID 0)
-- Dependencies: 223
-- Name: languages_id_language_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: teammate_search
--

ALTER SEQUENCE public.languages_id_language_seq OWNED BY public.languages.id_language;


--
-- TOC entry 218 (class 1259 OID 19271)
-- Name: users; Type: TABLE; Schema: public; Owner: teammate_search
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    age integer,
    description text,
    most_like_game integer,
    most_like_genre integer,
    language integer,
    created_at date NOT NULL,
    speaking_app integer
);


ALTER TABLE public.users OWNER TO teammate_search;

--
-- TOC entry 217 (class 1259 OID 19270)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: teammate_search
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO teammate_search;

--
-- TOC entry 4898 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: teammate_search
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4720 (class 2604 OID 19396)
-- Name: apps id_app; Type: DEFAULT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.apps ALTER COLUMN id_app SET DEFAULT nextval('public.apps_id_app_seq'::regclass);


--
-- TOC entry 4718 (class 2604 OID 19292)
-- Name: games id_game; Type: DEFAULT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.games ALTER COLUMN id_game SET DEFAULT nextval('public.games_id_game_seq'::regclass);


--
-- TOC entry 4717 (class 2604 OID 19283)
-- Name: genres id_genre; Type: DEFAULT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.genres ALTER COLUMN id_genre SET DEFAULT nextval('public.genres_id_genre_seq'::regclass);


--
-- TOC entry 4719 (class 2604 OID 19301)
-- Name: languages id_language; Type: DEFAULT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.languages ALTER COLUMN id_language SET DEFAULT nextval('public.languages_id_language_seq'::regclass);


--
-- TOC entry 4716 (class 2604 OID 19274)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 4888 (class 0 OID 19393)
-- Dependencies: 226
-- Data for Name: apps; Type: TABLE DATA; Schema: public; Owner: teammate_search
--

COPY public.apps (id_app, app) FROM stdin;
1	Discord
2	TeamSpeak
3	Zoom
4	Telegram
5	VK
\.


--
-- TOC entry 4884 (class 0 OID 19289)
-- Dependencies: 222
-- Data for Name: games; Type: TABLE DATA; Schema: public; Owner: teammate_search
--

COPY public.games (id_game, game) FROM stdin;
1	League Of Legends
2	DOTA2
3	CS GO
4	Far Cry
5	R.E.P.O
6	Raft
7	RUST
8	SUPERVIVE
9	Maincraft
10	The Forest
11	PUBG
12	PEAK
13	Lethal Company
\.


--
-- TOC entry 4882 (class 0 OID 19280)
-- Dependencies: 220
-- Data for Name: genres; Type: TABLE DATA; Schema: public; Owner: teammate_search
--

COPY public.genres (id_genre, genre) FROM stdin;
1	Шутер от первого лица
2	Шутер от третьего лица
3	Слэшеры
4	Казуальные
5	Ролевые экшены
6	Стратегии
7	Японские ролевые
8	Симуляторы
9	Башенная защита
10	Спортивные симуляторы
11	Гонки
12	Хоррор
13	Научная фантастика
14	Космос
15	Аркада
16	Платформеры
17	Файтинги
18	Визуальные новеллы
19	Приключения
20	Песочницы
21	Карточные
22	Настольные
23	Аниме
24	Выживание
25	Детективы
26	Открытый мир
27	Кооператив
\.


--
-- TOC entry 4886 (class 0 OID 19298)
-- Dependencies: 224
-- Data for Name: languages; Type: TABLE DATA; Schema: public; Owner: teammate_search
--

COPY public.languages (id_language, language) FROM stdin;
1	Russian
2	English
3	German
4	French
5	Spanish
6	Italian
7	Portuguese
8	Polish
9	Belarusian
10	Bulgarian
11	Czech
12	Slovak
13	Slovenian
14	Croatian
15	Serbian
16	Bosnian
17	Macedonian
18	Greek
19	Hungarian
20	Romanian
21	Chinese
22	Japanese
23	Korean
24	Arabic
\.


--
-- TOC entry 4880 (class 0 OID 19271)
-- Dependencies: 218
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: teammate_search
--

COPY public.users (id, username, password, age, description, most_like_game, most_like_genre, language, created_at, speaking_app) FROM stdin;
2	Piryet	Piryet	33	хороший человек\r\n                                    \r\n                                    \r\n                                    \r\n                                    \r\n                                    \r\n                                    \r\n                                    	9	1	1	2025-12-21	1
4	Aroman	Aroman	23	kungfu panda	1	1	1	2025-12-24	1
3	CRIGO	CRIGO	23	crigo estriper	4	1	1	2025-12-22	1
5	NewUser	NewUser	23	good guy	\N	\N	\N	2025-12-26	\N
6	NewUser2	NewUser2	31	qwe	\N	\N	\N	2025-12-26	\N
7	NewUser3	NewUser3	32	string	\N	\N	\N	2025-12-26	\N
\.


--
-- TOC entry 4899 (class 0 OID 0)
-- Dependencies: 225
-- Name: apps_id_app_seq; Type: SEQUENCE SET; Schema: public; Owner: teammate_search
--

SELECT pg_catalog.setval('public.apps_id_app_seq', 5, true);


--
-- TOC entry 4900 (class 0 OID 0)
-- Dependencies: 221
-- Name: games_id_game_seq; Type: SEQUENCE SET; Schema: public; Owner: teammate_search
--

SELECT pg_catalog.setval('public.games_id_game_seq', 13, true);


--
-- TOC entry 4901 (class 0 OID 0)
-- Dependencies: 219
-- Name: genres_id_genre_seq; Type: SEQUENCE SET; Schema: public; Owner: teammate_search
--

SELECT pg_catalog.setval('public.genres_id_genre_seq', 27, true);


--
-- TOC entry 4902 (class 0 OID 0)
-- Dependencies: 223
-- Name: languages_id_language_seq; Type: SEQUENCE SET; Schema: public; Owner: teammate_search
--

SELECT pg_catalog.setval('public.languages_id_language_seq', 24, true);


--
-- TOC entry 4903 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: teammate_search
--

SELECT pg_catalog.setval('public.users_id_seq', 7, true);


--
-- TOC entry 4730 (class 2606 OID 19400)
-- Name: apps apps_pkey; Type: CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.apps
    ADD CONSTRAINT apps_pkey PRIMARY KEY (id_app);


--
-- TOC entry 4726 (class 2606 OID 19296)
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (id_game);


--
-- TOC entry 4724 (class 2606 OID 19287)
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id_genre);


--
-- TOC entry 4728 (class 2606 OID 19305)
-- Name: languages languages_pkey; Type: CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (id_language);


--
-- TOC entry 4722 (class 2606 OID 19278)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4731 (class 2606 OID 19386)
-- Name: users language; Type: FK CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT language FOREIGN KEY (language) REFERENCES public.languages(id_language) NOT VALID;


--
-- TOC entry 4732 (class 2606 OID 19366)
-- Name: users most_like_game1; Type: FK CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT most_like_game1 FOREIGN KEY (most_like_game) REFERENCES public.games(id_game) NOT VALID;


--
-- TOC entry 4733 (class 2606 OID 19381)
-- Name: users most_like_genre; Type: FK CONSTRAINT; Schema: public; Owner: teammate_search
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT most_like_genre FOREIGN KEY (most_like_genre) REFERENCES public.genres(id_genre) NOT VALID;


-- Completed on 2025-12-26 18:01:57

--
-- PostgreSQL database dump complete
--

