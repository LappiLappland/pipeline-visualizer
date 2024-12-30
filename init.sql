--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8 (Debian 15.8-0+deb12u1)
-- Dumped by pg_dump version 15.8 (Debian 15.8-0+deb12u1)

-- Started on 2024-11-17 01:42:31 EST

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
-- TOC entry 231 (class 1259 OID 294418)
-- Name: branches; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.branches (
    id integer NOT NULL,
    name character varying(255)
);


ALTER TABLE public.branches OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 294417)
-- Name: branches_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.branches_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.branches_id_seq OWNER TO postgres;

--
-- TOC entry 3551 (class 0 OID 0)
-- Dependencies: 230
-- Name: branches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.branches_id_seq OWNED BY public.branches.id;


--
-- TOC entry 229 (class 1259 OID 294411)
-- Name: commits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.commits (
    id integer NOT NULL,
    user_id integer,
    branch_id integer,
    hash character varying(8)
);


ALTER TABLE public.commits OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 294410)
-- Name: commits_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.commits_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.commits_id_seq OWNER TO postgres;

--
-- TOC entry 3552 (class 0 OID 0)
-- Dependencies: 228
-- Name: commits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.commits_id_seq OWNED BY public.commits.id;


--
-- TOC entry 248 (class 1259 OID 294654)
-- Name: http_sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.http_sessions (
    id bigint NOT NULL,
    key bytea,
    data bytea,
    created_on timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified_on timestamp with time zone,
    expires_on timestamp with time zone
);


ALTER TABLE public.http_sessions OWNER TO postgres;

--
-- TOC entry 247 (class 1259 OID 294653)
-- Name: http_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.http_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.http_sessions_id_seq OWNER TO postgres;

--
-- TOC entry 3553 (class 0 OID 0)
-- Dependencies: 247
-- Name: http_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.http_sessions_id_seq OWNED BY public.http_sessions.id;


--
-- TOC entry 221 (class 1259 OID 294367)
-- Name: job_dependencies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.job_dependencies (
    id integer NOT NULL,
    parent_id integer,
    child_id integer
);


ALTER TABLE public.job_dependencies OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 294366)
-- Name: job_dependencies_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.job_dependencies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.job_dependencies_id_seq OWNER TO postgres;

--
-- TOC entry 3554 (class 0 OID 0)
-- Dependencies: 220
-- Name: job_dependencies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_dependencies_id_seq OWNED BY public.job_dependencies.id;


--
-- TOC entry 219 (class 1259 OID 294360)
-- Name: job_workers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.job_workers (
    id integer NOT NULL,
    status_id integer DEFAULT 1,
    job_id integer,
    pipeline_worker_id integer,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    user_id integer
);


ALTER TABLE public.job_workers OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 294359)
-- Name: job_workers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.job_workers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.job_workers_id_seq OWNER TO postgres;

--
-- TOC entry 3555 (class 0 OID 0)
-- Dependencies: 218
-- Name: job_workers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_workers_id_seq OWNED BY public.job_workers.id;


--
-- TOC entry 217 (class 1259 OID 294352)
-- Name: jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.jobs (
    id integer NOT NULL,
    name_id character varying(255),
    name character varying(255),
    pipeline_id integer
);


ALTER TABLE public.jobs OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 294404)
-- Name: log_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_types (
    id integer NOT NULL,
    name character varying(255)
);


ALTER TABLE public.log_types OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 294403)
-- Name: log_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.log_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.log_type_id_seq OWNER TO postgres;

--
-- TOC entry 3556 (class 0 OID 0)
-- Dependencies: 226
-- Name: log_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.log_type_id_seq OWNED BY public.log_types.id;


--
-- TOC entry 246 (class 1259 OID 294640)
-- Name: oauth_providers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.oauth_providers (
    id integer NOT NULL,
    name character varying(255)
);


ALTER TABLE public.oauth_providers OWNER TO postgres;

--
-- TOC entry 244 (class 1259 OID 294627)
-- Name: oauth_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.oauth_users (
    id integer NOT NULL,
    user_id integer,
    provider_user_id character varying(255),
    access_token character varying(255),
    refresh_token character varying(255),
    oauth_provider_id integer
);


ALTER TABLE public.oauth_users OWNER TO postgres;

--
-- TOC entry 243 (class 1259 OID 294626)
-- Name: oauth_providers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.oauth_providers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.oauth_providers_id_seq OWNER TO postgres;

--
-- TOC entry 3557 (class 0 OID 0)
-- Dependencies: 243
-- Name: oauth_providers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.oauth_providers_id_seq OWNED BY public.oauth_users.id;


--
-- TOC entry 245 (class 1259 OID 294639)
-- Name: oauth_providers_id_seq1; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.oauth_providers_id_seq1
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.oauth_providers_id_seq1 OWNER TO postgres;

--
-- TOC entry 3558 (class 0 OID 0)
-- Dependencies: 245
-- Name: oauth_providers_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.oauth_providers_id_seq1 OWNED BY public.oauth_providers.id;


--
-- TOC entry 238 (class 1259 OID 294575)
-- Name: pipeline_statistics_short; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.pipeline_statistics_short AS
 SELECT (min(job_workers.started_at) + ('00:00:01'::interval * (avg(EXTRACT(epoch FROM (job_workers.finished_at - job_workers.started_at))))::double precision)) AS average_timestamp,
    max(job_workers.finished_at) AS last_timestamp,
    count(*) FILTER (WHERE (job_workers.status_id = 3)) AS total_errors,
        CASE
            WHEN (count(*) FILTER (WHERE (job_workers.status_id = 4)) = count(*)) THEN 4
            WHEN (count(*) FILTER (WHERE (job_workers.status_id = 3)) > 0) THEN 3
            WHEN (count(*) FILTER (WHERE (job_workers.status_id = 2)) > 0) THEN 2
            ELSE 1
        END AS pipeline_status
   FROM public.job_workers
  WHERE ((job_workers.finished_at IS NOT NULL) AND (job_workers.started_at IS NOT NULL));


ALTER TABLE public.pipeline_statistics_short OWNER TO postgres;

--
-- TOC entry 242 (class 1259 OID 294603)
-- Name: pipeline_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_users (
    id integer NOT NULL,
    pipeline_id integer,
    user_id integer,
    role_id integer
);


ALTER TABLE public.pipeline_users OWNER TO postgres;

--
-- TOC entry 241 (class 1259 OID 294602)
-- Name: pipeline_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pipeline_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pipeline_users_id_seq OWNER TO postgres;

--
-- TOC entry 3559 (class 0 OID 0)
-- Dependencies: 241
-- Name: pipeline_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pipeline_users_id_seq OWNED BY public.pipeline_users.id;


--
-- TOC entry 225 (class 1259 OID 294388)
-- Name: pipeline_workers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_workers (
    id integer NOT NULL,
    pipeline_id integer,
    user_id integer,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    status_id integer DEFAULT 1
);


ALTER TABLE public.pipeline_workers OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 294387)
-- Name: pipeline_workers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pipeline_workers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pipeline_workers_id_seq OWNER TO postgres;

--
-- TOC entry 3560 (class 0 OID 0)
-- Dependencies: 224
-- Name: pipeline_workers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pipeline_workers_id_seq OWNED BY public.pipeline_workers.id;


--
-- TOC entry 216 (class 1259 OID 294347)
-- Name: pipelines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipelines (
    id integer NOT NULL,
    commit_id integer,
    name character varying(255),
    created_at timestamp with time zone,
    is_private boolean DEFAULT true,
    description character varying(500)
);


ALTER TABLE public.pipelines OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 294596)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(255),
    permissions smallint
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 294595)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO postgres;

--
-- TOC entry 3561 (class 0 OID 0)
-- Dependencies: 239
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 237 (class 1259 OID 294551)
-- Name: stage_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stage_logs (
    id integer NOT NULL,
    text character varying(1000),
    line integer,
    log_type_id integer,
    stage_worker_id integer,
    created_at timestamp with time zone
);


ALTER TABLE public.stage_logs OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 294550)
-- Name: stage_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stage_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stage_logs_id_seq OWNER TO postgres;

--
-- TOC entry 3562 (class 0 OID 0)
-- Dependencies: 236
-- Name: stage_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stage_logs_id_seq OWNED BY public.stage_logs.id;


--
-- TOC entry 235 (class 1259 OID 294529)
-- Name: stage_workers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stage_workers (
    id integer NOT NULL,
    status_id integer DEFAULT 1,
    stage_id integer,
    job_worker_id integer,
    started_at timestamp with time zone,
    finished_at timestamp with time zone
);


ALTER TABLE public.stage_workers OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 294528)
-- Name: stage_workers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stage_workers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stage_workers_id_seq OWNER TO postgres;

--
-- TOC entry 3563 (class 0 OID 0)
-- Dependencies: 234
-- Name: stage_workers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stage_workers_id_seq OWNED BY public.stage_workers.id;


--
-- TOC entry 233 (class 1259 OID 294512)
-- Name: stages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stages (
    id integer NOT NULL,
    name character varying(255),
    job_id integer,
    runs_for integer,
    will_fail boolean,
    "order" integer
);


ALTER TABLE public.stages OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 294511)
-- Name: stages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stages_id_seq OWNER TO postgres;

--
-- TOC entry 3564 (class 0 OID 0)
-- Dependencies: 232
-- Name: stages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stages_id_seq OWNED BY public.stages.id;


--
-- TOC entry 223 (class 1259 OID 294374)
-- Name: statuses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.statuses (
    id integer NOT NULL,
    name character varying(255)
);


ALTER TABLE public.statuses OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 294373)
-- Name: statuses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.statuses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.statuses_id_seq OWNER TO postgres;

--
-- TOC entry 3565 (class 0 OID 0)
-- Dependencies: 222
-- Name: statuses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.statuses_id_seq OWNED BY public.statuses.id;


--
-- TOC entry 215 (class 1259 OID 294341)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(255),
    password character varying(255)
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 294340)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3566 (class 0 OID 0)
-- Dependencies: 214
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 3296 (class 2604 OID 294421)
-- Name: branches id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches ALTER COLUMN id SET DEFAULT nextval('public.branches_id_seq'::regclass);


--
-- TOC entry 3295 (class 2604 OID 294414)
-- Name: commits id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commits ALTER COLUMN id SET DEFAULT nextval('public.commits_id_seq'::regclass);


--
-- TOC entry 3305 (class 2604 OID 294657)
-- Name: http_sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.http_sessions ALTER COLUMN id SET DEFAULT nextval('public.http_sessions_id_seq'::regclass);


--
-- TOC entry 3290 (class 2604 OID 294370)
-- Name: job_dependencies id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_dependencies ALTER COLUMN id SET DEFAULT nextval('public.job_dependencies_id_seq'::regclass);


--
-- TOC entry 3288 (class 2604 OID 294363)
-- Name: job_workers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers ALTER COLUMN id SET DEFAULT nextval('public.job_workers_id_seq'::regclass);


--
-- TOC entry 3294 (class 2604 OID 294407)
-- Name: log_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_types ALTER COLUMN id SET DEFAULT nextval('public.log_type_id_seq'::regclass);


--
-- TOC entry 3304 (class 2604 OID 294643)
-- Name: oauth_providers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_providers ALTER COLUMN id SET DEFAULT nextval('public.oauth_providers_id_seq1'::regclass);


--
-- TOC entry 3303 (class 2604 OID 294630)
-- Name: oauth_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_users ALTER COLUMN id SET DEFAULT nextval('public.oauth_providers_id_seq'::regclass);


--
-- TOC entry 3302 (class 2604 OID 294606)
-- Name: pipeline_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_users ALTER COLUMN id SET DEFAULT nextval('public.pipeline_users_id_seq'::regclass);


--
-- TOC entry 3292 (class 2604 OID 294391)
-- Name: pipeline_workers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_workers ALTER COLUMN id SET DEFAULT nextval('public.pipeline_workers_id_seq'::regclass);


--
-- TOC entry 3301 (class 2604 OID 294599)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 3300 (class 2604 OID 294554)
-- Name: stage_logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_logs ALTER COLUMN id SET DEFAULT nextval('public.stage_logs_id_seq'::regclass);


--
-- TOC entry 3298 (class 2604 OID 294532)
-- Name: stage_workers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_workers ALTER COLUMN id SET DEFAULT nextval('public.stage_workers_id_seq'::regclass);


--
-- TOC entry 3297 (class 2604 OID 294515)
-- Name: stages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stages ALTER COLUMN id SET DEFAULT nextval('public.stages_id_seq'::regclass);


--
-- TOC entry 3291 (class 2604 OID 294377)
-- Name: statuses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses ALTER COLUMN id SET DEFAULT nextval('public.statuses_id_seq'::regclass);


--
-- TOC entry 3286 (class 2604 OID 294344)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3529 (class 0 OID 294418)
-- Dependencies: 231
-- Data for Name: branches; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.branches (id, name) FROM stdin;
1	main
\.


--
-- TOC entry 3527 (class 0 OID 294411)
-- Dependencies: 229
-- Data for Name: commits; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.commits (id, user_id, branch_id, hash) FROM stdin;
2	1	1	30f5485
\.


--
-- TOC entry 3545 (class 0 OID 294654)
-- Dependencies: 248
-- Data for Name: http_sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.http_sessions (id, key, data, created_on, modified_on, expires_on) FROM stdin;
2	\\x594e4d474e55495949444e424e5641484c584d42344248364c3355513537435a485141325a474747484b51513557353734494e51	\\x4d54637a4d5455794d4445354e587845574468465156464d5832644251554a4651555652515546454c5546525831396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463933515642666445673063306c4251554642515546425156393556456c694d48564555554a6e515468504f586c7a53474e306144426c656b4e53536b523564316479624734744d6d68465244413561336c324e6c6f7a656d3530547a6732544852494e30307a646e686c4c5568795246685859574a526245787861314a504c565a68566e4e5865444631565852724e6a4e5765556c7654464e34575746684c584e6d4d6c684e52455132656d684a524552755569303564484e48626b526e5643314b5554646f654452795246564b58337056615452584f546446526a46566157704a5932567951566c6653556b784d57525a5355686f576c396f64316c6d4d6c4a4b4e307433615539705954646d5630394e637a4a504e555679575330344e6d733457454635546e5a30646c513255453535575570735a6a6c5354584d334e31425765586c3254315256656d706b636d64564d6a4d77643051354f444e6664475271625746614d336377597a46764f45644e516e565452334a4361305277556b6f31516d7447593146754e3139425155464258313834516b464252463966656b747a6344467562554642515546386c7637594c7867316b464b4c59316f4b64355a5576376572654d4c686642386437536b315270686d742d413d	2024-11-13 12:49:55.521204-05	2024-11-13 12:49:55.521204-05	2024-12-13 12:49:55.521204-05
3	\\x4a41485a483756343333475147374849343458354f474f4947434d5535524446365a55513358594e5237344e3746555358495541	\\x4d54637a4d5455794d4463304f487845574468465156464d5832644251554a4651555652515546454c5546535246396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463934515642666455673063306c4251554642515546425156393556456b775633467954554a6e5154524959306f35547a5a6a5957463056455658556a426b5a544642656c56684d55785a56304a68554868795457307a6158704b58315a6b6458706b65446c71546d5266526a6c7259565a43646e52316255704455574e6a5a454e6f576d4a5651335668626e5a7452484269646c6478526e5247556d7872526e566663555a49597a524a5344464261564e58623239735a56426f5a45394d4d454e7a4e566449656a64505a57497963487033656d52714d6e646f64454a454d484278533142464e6c464e556c56424d7931326331564c5358427456486c7963466431646c7069564530344e6d35705a7a6332587a52784d6e4579556a6335546a4e7a4e455133556a4579647a525164446379596d784965556c73646a4e6c4e476c75644556314f47644257576f344d6a4a7255566f35567a68665a474e59616a565a5a6d3435635752304d54527456473534656b6c75536c4133536d744554464631566b466c5532684a556a686664304642515642665830465251554666587a6c306445646e4d54566e5155464251543039664f683647536f7a6470646e7034535675335434647873777a357a72452d51426e707561505375312d6b2d59	2024-11-13 12:59:08.135734-05	2024-11-13 12:59:08.135734-05	2024-12-13 12:59:08.135734-05
4	\\x47494637454345514f4d3652495053494842544a4e5852584a4349445a474932464b4e4e533347544e5444413454593633504441	\\x4d54637a4d5455794d4463324d337845574468465156464d5832644251554a4651555652515546454c5546535246396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463934515642666455673063306c4251554642515546425156393556456b7a5654644454554a52515452495a484233634441304e6b70365357747a563046724c58464259564a454f476b3462454e31655539304d6a6c78546d357153316b345a444a4f4f4756684e79314d4e3070345330686a636b4a5a6130704353336830635568754e314a53533367344c555a78596e70544e3070554d6e704f59576735554452335a546c59516e52545a315a685458685653476f77593346474f48465652464135643355315231513065546455583368335a3364484f544e4a54335271596d46524e47467452314e4b5a324e6664463970526d6c486154566161544a7765464e50656e70355a4468305a584a4a56453131626e6857536e67785a6d6f77596b64325254637861336468576e4a6851564247583074316347497a626d4645523156474e7a46724e6e5a564f5855746230684a4d6b5646526a466d536b4e345158565855457869636a52744d325535646e6b304e54686d613270466555564252335a5963476443546c46725369316d5a30564251564266583046525155466658313955556e64746254566e51554642515430396643747842715f334a7269364f636d5a3538667a4c376d4c4a56746e3873427346355f7354436378664b3359	2024-11-13 12:59:23.537057-05	2024-11-13 12:59:23.537058-05	2024-12-13 12:59:23.537057-05
5	\\x564c484e4e3357364d4e33374c4343574f45414457324647565141584749354644544235504e524d51574e3541374a4b47494d41	\\x4d54637a4d5455794d4467304d337845574468465156464d5832644251554a4651555652515546454c5546524d31396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e5735455546393151564266636b673063306c4251554642515546425156393556456c694d48564555554a6e515468504f586c7a53475661587a6c615232647a56465a74633274706148707852554e5055546875583152355a48466d5a574d304d6b3132626e5246596a4d30646d5a304c555658645868514c565250536c4e4a4f44527463326778547a52484f574a684e6a5651636e52545a44424c4e573173626e4d7a5646387754454a684e445178533046336248457759325a5a576d68475330457953334e785a57527656314e5855475934597a453257484a44576e645a65474648595770345954525457545a345969307a65554a4561557432546b5661626a4e586446525459316c69596b706b646a64714f58677a4e55517863306f34636e6b35576b68744d30745a616b4e664d6c6450567a4932524842334e54645455326f774d444579554731364e6a5a7856545a4a576a4a444e5668434e55744761465245646a5a49566a4133536c7053574856576455645057454a4953303953614651325a305656615646714e533152565546425546396651564642515639664f454a72575468514e57644251554642505431383041707a766e656472655543433550744749336472663768384e583045424571435a6947746563746f32633d	2024-11-13 13:00:43.984317-05	2024-11-13 13:00:43.984318-05	2024-12-13 13:00:43.984317-05
8	\\x5535534b32564a5a5a583337464f34325844595255544c42334157415a344635594c43484f47414c5a5034444c59355443543451	\\x4d54637a4d5455794d54457a4d6e7845574468465156464d5832644251554a4651555652515546454c5546535246396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463934515642666455673063306c4251554642515546425156393556456b7a56545a3654554a6e515452496448427a636c4232613364725332745a55566c4762486c5959565261526b597a4e576c5263554d3455557857636b39666356644c59316f33546a6868564455745244564a636b68474e4670424c57744b51553170546e4646616e524f656b6848654870335a46525352574679626a4273527a463459306731556c6b784f4768736332315052576c7a5a554a3264456f355931524955476868576d35525646424f5545646b646d68786244423257456734526d394b555446564c55747261466c7863555a32656d525a535442554d316c73646b315063326c77646b397763475a785447706864446c7358323576596c424a536d6732563064314e48464d546d4a34536c647962574a6f6548526b5647787765457454554468494f5756614e4455335a545a6c55323931635459334f473147656c6c3256465a756258553352304676596e703361334e585747744d4c5774615a33684e54317047646c6c4e61306c5453475a516430464251564266583046525155466658793174526e68664d7a566e515546425154303966436d3359544146323632483441665554584a39715a6d7049384b425078352d4a4f5f6f3561636952615741	2024-11-13 13:05:32.384874-05	2024-11-13 13:05:32.384875-05	2024-12-13 13:05:32.384874-05
6	\\x564936525255324c5134414b535351523549364a5158533747444f335955334e513336485357543454494e455057534933494b51	\\x4d54637a4d5455794d544d324d487845574468465156464d5832644251554a4651555652515546454c5546535246396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463934515642666455673063306c4251554642515546425156393556456b7756545a7954554a6e515452495a48427a636e5236516b4e7265564a72535531444e553161575645336255704554573143545739664d6d3946644739594c573536646d703165473571656c685965475a4b5245784a4f5442575330467a53564a536544465a566d6c6c55573074595574785647527956556f79555778786230356a64584e59546c6c7254444e4d536b4a6e54564a4c64453949616a4a68523052425345567565544a524f544e35563035505a5574796331637a5347307751314254625846764f4568505257744c6131646662445671616c4a4253324d786357524c596c3958545578686348597862446c59526c4a544f57466d62444e494f57564e656e5a4b616a4a36626e56464d6a4e6a546a4a4e52444e5963335a4e516c5661554731744d31513459326c324e5639745579303564464243646a64795a54424c556b313261546479644764694e4859346147314f576e705262563877616b564852326f3563456878555570445247737464324e42515642665830465251554666587a673561555a725154566e5155464251543039664472765863703066774e37454c45464c45697653484c6661594753356e3058597835304c77634a695a2d69	2024-11-13 13:02:32.084918-05	2024-11-13 13:09:20.537586-05	2024-12-13 13:09:20.537586-05
9	\\x534b4832344f5a5232434147323334584b364f5841474254445343524c494548504d4c593450454e563252374350355535413651	\\x4d54637a4d5455794d6a4d324e337845574468465156464d5832644251554a4651555652515546454c5546525831396e51554643516d354f4d474e74624856616433644a51554661626d4659556d396b56306c48597a4e53655746584e57354555463933515642666445673063306c4251554642515546425156393556456b7a56545a4554554a52515452495a48427a616e4e716333646f6130704e556d646a543051344d6d565a5a306b34595556525347566e5a475245557a6c61556b31715a546c31656b63324c576b744c556834516c6f314f475a615456467a53564a536545303256476c6c555449745957467859306877566c4e6c613238796355777a55477849626d4e524d7a4e4d526d566e54564a4c5345744d4d564a4f4d57566e526d704c5a6c7076526e55335533684d5744645a4c5670364e6b7036513270725a316c7858304a7661466c3162304673656d525a535442545a57564958334e68636a6c6b4d3364365a56396b546c56755a477845656b704d5956426a5a55785558304a59635468695646633563484e47636d7374526e464c6132743363566334597a5248646e4668646d525a64544e58576b644f5a6d4a795544465751793145516e524f635442685745686b5830387a52315274536e56545333684a65554a4e596d73325a326c5261456c694f5639425155464258313834516b4642524639664f554e6f6548524d6255464251554638396e4361614c4a796e645a6e656d7136674958594a38364a436b78344437645a30696977697372596736453d	2024-11-13 13:09:57.25542-05	2024-11-13 13:26:07.683439-05	2024-12-13 13:26:07.683439-05
10	\\x4f324935484656484f5137454c4d3337505a545436454d484d3334585942494a34434a574a4451534741334e4941544657575641	\\x4d54637a4d5455794e546b314e587845574468465156464d5832644251554a465155565251554643616c383051554642643170365a45684b63474a74593031435155464459566452525752586248566b51566c445155464652324d7a556e6c68567a56755245465a51554a484e57686956315648597a4e53655746584e573545515468425246563461474e49516e425552305a335930643461474a745555646a4d314a3559566331626b52426230464453454a35596a4e6163467048566e6c43626b347759323173645670336430644251564a76596a4e4f4d487774507a6546454a3645724f35644831647830442d394764435f6a634c4a73666a576c42706d78426a6b78513d3d	2024-11-13 14:25:55.025993-05	2024-11-13 14:25:55.025993-05	2024-12-13 14:25:55.025993-05
\.


--
-- TOC entry 3519 (class 0 OID 294367)
-- Dependencies: 221
-- Data for Name: job_dependencies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.job_dependencies (id, parent_id, child_id) FROM stdin;
1	3	2
2	2	1
3	1	4
4	8	7
5	8	6
6	6	5
7	7	5
8	11	9
9	12	9
10	13	12
11	13	11
12	14	12
13	15	13
14	15	14
15	15	10
16	19	17
17	19	18
18	18	16
19	23	20
20	23	21
21	24	21
22	22	25
23	29	24
24	29	22
25	28	25
26	27	23
27	26	23
28	30	27
29	30	26
30	31	28
31	31	29
\.


--
-- TOC entry 3517 (class 0 OID 294360)
-- Dependencies: 219
-- Data for Name: job_workers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.job_workers (id, status_id, job_id, pipeline_worker_id, started_at, finished_at, user_id) FROM stdin;
8	1	8	2	\N	\N	\N
6	1	6	2	\N	\N	\N
4	1	4	1	\N	\N	\N
1	1	1	1	\N	\N	\N
3	1	3	1	\N	\N	\N
2	1	2	1	\N	\N	\N
5	1	5	2	\N	\N	\N
9	1	9	3	\N	\N	\N
7	1	7	2	\N	\N	\N
10	1	10	3	\N	\N	\N
12	1	12	3	\N	\N	\N
11	1	11	3	\N	\N	\N
13	1	13	3	\N	\N	\N
15	1	15	3	\N	\N	\N
14	1	14	3	\N	\N	\N
16	1	16	4	\N	\N	\N
17	1	17	4	\N	\N	\N
18	1	18	4	\N	\N	\N
19	1	19	4	\N	\N	\N
33	1	33	7	\N	\N	\N
32	1	32	6	\N	\N	\N
21	1	21	5	\N	\N	\N
25	1	25	5	\N	\N	\N
20	1	20	5	\N	\N	\N
23	1	23	5	\N	\N	\N
22	1	22	5	\N	\N	\N
28	1	28	5	\N	\N	\N
24	1	24	5	\N	\N	\N
27	1	27	5	\N	\N	\N
26	1	26	5	\N	\N	\N
29	1	29	5	\N	\N	\N
30	1	30	5	\N	\N	\N
31	1	31	5	\N	\N	\N
\.


--
-- TOC entry 3515 (class 0 OID 294352)
-- Dependencies: 217
-- Data for Name: jobs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.jobs (id, name_id, name, pipeline_id) FROM stdin;
1	test	\N	1
2	assert	\N	1
3	teardown	\N	1
4	setup	Setting stuff up	1
5	work1	\N	2
6	depends1	\N	2
7	depends2	\N	2
8	master	\N	2
9	SomeJob	Begin	3
10	SomeJob	Will begin too	3
13	SomeJob	Will never start	3
14	SomeJob	Will start	3
15	SomeJob	Will never start	3
11	SomeJob	Will fail	3
12	SomeJob	Will finish	3
16	prepare	\N	4
17	get_ready	\N	4
18	Almost	\N	4
19	There	\N	4
30	anotherFinal	\N	5
31	final	\N	5
20	SomeName	\N	5
21	SomeName	\N	5
22	SomeName	\N	5
23	SomeName	\N	5
24	SomeName	\N	5
25	SomeName	\N	5
26	SomeName	\N	5
27	SomeName	\N	5
28	SomeName	\N	5
29	SomeName	\N	5
32	heh	Very long job	6
33	Something	\N	7
\.


--
-- TOC entry 3525 (class 0 OID 294404)
-- Dependencies: 227
-- Data for Name: log_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_types (id, name) FROM stdin;
1	base
2	error
3	warn
4	info
\.


--
-- TOC entry 3543 (class 0 OID 294640)
-- Dependencies: 246
-- Data for Name: oauth_providers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.oauth_providers (id, name) FROM stdin;
1	github
\.


--
-- TOC entry 3541 (class 0 OID 294627)
-- Dependencies: 244
-- Data for Name: oauth_users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.oauth_users (id, user_id, provider_user_id, access_token, refresh_token, oauth_provider_id) FROM stdin;
2	3	117802025	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		1
\.


--
-- TOC entry 3539 (class 0 OID 294603)
-- Dependencies: 242
-- Data for Name: pipeline_users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pipeline_users (id, pipeline_id, user_id, role_id) FROM stdin;
2	1	1	3
3	2	1	3
4	3	1	3
5	4	1	3
6	5	1	3
7	6	1	3
8	7	1	3
9	1	4	1
10	2	4	1
11	3	4	1
12	4	4	1
13	5	4	1
14	6	4	1
15	7	4	1
\.


--
-- TOC entry 3523 (class 0 OID 294388)
-- Dependencies: 225
-- Data for Name: pipeline_workers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pipeline_workers (id, pipeline_id, user_id, started_at, finished_at, status_id) FROM stdin;
3	3	\N	\N	\N	1
4	4	\N	\N	\N	1
7	7	\N	\N	\N	1
5	5	\N	\N	\N	1
1	1	\N	\N	\N	1
6	6	\N	\N	\N	1
2	2	\N	\N	\N	1
\.


--
-- TOC entry 3514 (class 0 OID 294347)
-- Dependencies: 216
-- Data for Name: pipelines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pipelines (id, commit_id, name, created_at, is_private, description) FROM stdin;
1	2	Default case	2024-11-07 04:57:03.600831-05	t	Can check basic funcionality
2	2	Parallel tasks	2024-11-15 04:09:32.912946-05	t	Basic case with parallel execution
3	2	Paralllel failure	2024-11-15 04:49:58.132642-05	t	Basic case with parallel execution failure
4	2	Deep dependency	2024-11-15 13:31:28.214643-05	t	Showcase of horizontal placement of jobs
5	2	Huge graph	2024-11-15 13:54:04.252831-05	t	A lot of jobs. You can see how dependency connections crossins are minimal. Also can check how flexible execution of separate tasks is.
6	2	Long job	2024-11-15 14:54:14.749951-05	t	One job with long tasks. Can check progress updates
7	2	Public pipeline	2024-11-15 15:03:06.593786-05	f	The only public pipeline, which can be opened by guests
\.


--
-- TOC entry 3537 (class 0 OID 294596)
-- Dependencies: 240
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, permissions) FROM stdin;
1	Tester	1
2	Developer	7
3	Admin	15
\.


--
-- TOC entry 3535 (class 0 OID 294551)
-- Dependencies: 237
-- Data for Name: stage_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stage_logs (id, text, line, log_type_id, stage_worker_id, created_at) FROM stdin;
\.


--
-- TOC entry 3533 (class 0 OID 294529)
-- Dependencies: 235
-- Data for Name: stage_workers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stage_workers (id, status_id, stage_id, job_worker_id, started_at, finished_at) FROM stdin;
23	1	23	16	\N	\N
24	1	24	17	\N	\N
25	1	25	18	\N	\N
26	1	26	19	\N	\N
27	1	27	19	\N	\N
15	1	15	8	\N	\N
12	1	12	5	\N	\N
14	1	14	7	\N	\N
13	1	13	6	\N	\N
20	1	20	13	\N	\N
22	1	22	15	\N	\N
9	1	8	2	\N	\N
8	1	9	2	\N	\N
11	1	10	3	\N	\N
10	1	11	3	\N	\N
45	1	45	33	\N	\N
16	1	16	9	\N	\N
17	1	17	10	\N	\N
19	1	19	12	\N	\N
18	1	18	11	\N	\N
21	1	21	14	\N	\N
1	1	3	4	\N	\N
2	1	2	4	\N	\N
3	1	1	4	\N	\N
7	1	4	1	\N	\N
6	1	5	1	\N	\N
29	1	29	21	\N	\N
33	1	33	25	\N	\N
4	1	7	1	\N	\N
5	1	6	1	\N	\N
28	1	28	20	\N	\N
44	1	44	32	\N	\N
42	1	42	32	\N	\N
31	1	31	23	\N	\N
43	1	43	32	\N	\N
41	1	41	32	\N	\N
30	1	30	22	\N	\N
40	1	40	32	\N	\N
36	1	36	28	\N	\N
32	1	32	24	\N	\N
35	1	35	27	\N	\N
34	1	34	26	\N	\N
37	1	37	29	\N	\N
38	1	38	30	\N	\N
39	1	39	31	\N	\N
\.


--
-- TOC entry 3531 (class 0 OID 294512)
-- Dependencies: 233
-- Data for Name: stages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stages (id, name, job_id, runs_for, will_fail, "order") FROM stdin;
1	Almost there	4	5	f	1
2	Now let's set up second ones	4	1	f	3
3	Set up first things	4	2	f	5
4	Test feature D	1	2	f	7
5	Test feature C	1	3	f	8
7	Test feature A	1	4	f	10
8	Assert something else	2	2	f	3
9	Assert true is true	2	3	f	5
11	Fast operation	3	1	f	6
12	Prepare for the inevitable	5	3	f	3
13	Long run	6	6	f	6
14	Short run	7	3	f	3
15	The inevitable	8	7	f	7
16	SomeCoolName	9	2	f	1
17	SomeCoolName	10	3	f	1
20	SomeCoolName	13	10	f	1
21	SomeCoolName	14	8	f	1
22	Hello darkness	15	7	f	1
18	SomeCoolName	11	4	t	1
19	SomeCoolName	12	2	f	1
23	Random	16	5	f	1
24	Stuff	17	3	f	1
25	Thoings	18	3	f	1
26	Things	19	3	f	1
27	more	19	7	f	1
28	matters not	20	5	f	1
29	matters not	21	3	f	1
30	matters not	22	4	f	1
31	matters not	23	2	f	1
32	matters not	24	5	f	1
33	matters not	25	3	f	1
34	matters not	26	4	f	1
35	matters not	27	3	f	1
36	matters not	28	5	f	1
37	matters not	29	6	f	1
38	matters not	30	7	f	1
39	matters not	31	8	f	1
40	Stage 1	32	5	f	9
41	Stage 2	32	5	f	7
42	Stage 3	32	5	f	5
43	Stage 4	32	5	f	4
44	Stage 5	32	5	f	3
45	Slacking off	33	10	f	4
6	Test feature B	1	4	f	9
10	Very long operation	3	6	f	2
\.


--
-- TOC entry 3521 (class 0 OID 294374)
-- Dependencies: 223
-- Data for Name: statuses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.statuses (id, name) FROM stdin;
1	pending
2	running
3	failed
4	completed
5	planned
6	cancelled
\.


--
-- TOC entry 3513 (class 0 OID 294341)
-- Dependencies: 215
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, password) FROM stdin;
2		
1	admin	$2a$10$NRfF/GMrSYYjlYTHa.O5E.k6ziuQNbWuOoEC3ZycV2aWgsYeDTo7W
4	tester	$2a$10$CBsWTqmmHresas8pHgMqCeUr5x1w1QSXQJGReKGALgM3WGWkBFNWy
3	LappiLappland	
\.


--
-- TOC entry 3567 (class 0 OID 0)
-- Dependencies: 230
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.branches_id_seq', 1, true);


--
-- TOC entry 3568 (class 0 OID 0)
-- Dependencies: 228
-- Name: commits_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.commits_id_seq', 2, true);


--
-- TOC entry 3569 (class 0 OID 0)
-- Dependencies: 247
-- Name: http_sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.http_sessions_id_seq', 10, true);


--
-- TOC entry 3570 (class 0 OID 0)
-- Dependencies: 220
-- Name: job_dependencies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_dependencies_id_seq', 8, true);


--
-- TOC entry 3571 (class 0 OID 0)
-- Dependencies: 218
-- Name: job_workers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_workers_id_seq', 4, true);


--
-- TOC entry 3572 (class 0 OID 0)
-- Dependencies: 226
-- Name: log_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.log_type_id_seq', 4, true);


--
-- TOC entry 3573 (class 0 OID 0)
-- Dependencies: 243
-- Name: oauth_providers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.oauth_providers_id_seq', 2, true);


--
-- TOC entry 3574 (class 0 OID 0)
-- Dependencies: 245
-- Name: oauth_providers_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.oauth_providers_id_seq1', 1, true);


--
-- TOC entry 3575 (class 0 OID 0)
-- Dependencies: 241
-- Name: pipeline_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pipeline_users_id_seq', 4, true);


--
-- TOC entry 3576 (class 0 OID 0)
-- Dependencies: 224
-- Name: pipeline_workers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pipeline_workers_id_seq', 1, false);


--
-- TOC entry 3577 (class 0 OID 0)
-- Dependencies: 239
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 3, true);


--
-- TOC entry 3578 (class 0 OID 0)
-- Dependencies: 236
-- Name: stage_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stage_logs_id_seq', 22874, true);


--
-- TOC entry 3579 (class 0 OID 0)
-- Dependencies: 234
-- Name: stage_workers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stage_workers_id_seq', 11, true);


--
-- TOC entry 3580 (class 0 OID 0)
-- Dependencies: 232
-- Name: stages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stages_id_seq', 12, true);


--
-- TOC entry 3581 (class 0 OID 0)
-- Dependencies: 222
-- Name: statuses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.statuses_id_seq', 1, true);


--
-- TOC entry 3582 (class 0 OID 0)
-- Dependencies: 214
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- TOC entry 3326 (class 2606 OID 294423)
-- Name: branches branches_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_pkey PRIMARY KEY (id);


--
-- TOC entry 3324 (class 2606 OID 294416)
-- Name: commits commits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commits
    ADD CONSTRAINT commits_pkey PRIMARY KEY (id);


--
-- TOC entry 3344 (class 2606 OID 294662)
-- Name: http_sessions http_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.http_sessions
    ADD CONSTRAINT http_sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 3316 (class 2606 OID 294372)
-- Name: job_dependencies job_dependencies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_dependencies
    ADD CONSTRAINT job_dependencies_pkey PRIMARY KEY (id);


--
-- TOC entry 3314 (class 2606 OID 294365)
-- Name: job_workers job_workers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers
    ADD CONSTRAINT job_workers_pkey PRIMARY KEY (id);


--
-- TOC entry 3312 (class 2606 OID 294358)
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--
-- TOC entry 3322 (class 2606 OID 294409)
-- Name: log_types log_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_types
    ADD CONSTRAINT log_type_pkey PRIMARY KEY (id);


--
-- TOC entry 3340 (class 2606 OID 294645)
-- Name: oauth_providers oauth_providers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_pkey PRIMARY KEY (id);


--
-- TOC entry 3338 (class 2606 OID 294647)
-- Name: oauth_users oauth_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_users
    ADD CONSTRAINT oauth_users_pkey PRIMARY KEY (id);


--
-- TOC entry 3336 (class 2606 OID 294608)
-- Name: pipeline_users pipeline_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_users
    ADD CONSTRAINT pipeline_users_pkey PRIMARY KEY (id);


--
-- TOC entry 3320 (class 2606 OID 294393)
-- Name: pipeline_workers pipeline_workers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_workers
    ADD CONSTRAINT pipeline_workers_pkey PRIMARY KEY (id);


--
-- TOC entry 3310 (class 2606 OID 294351)
-- Name: pipelines pipelines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipelines
    ADD CONSTRAINT pipelines_pkey PRIMARY KEY (id);


--
-- TOC entry 3334 (class 2606 OID 294601)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 3332 (class 2606 OID 294558)
-- Name: stage_logs stage_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_logs
    ADD CONSTRAINT stage_logs_pkey PRIMARY KEY (id);


--
-- TOC entry 3330 (class 2606 OID 294534)
-- Name: stage_workers stage_workers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_workers
    ADD CONSTRAINT stage_workers_pkey PRIMARY KEY (id);


--
-- TOC entry 3328 (class 2606 OID 294517)
-- Name: stages stages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stages
    ADD CONSTRAINT stages_pkey PRIMARY KEY (id);


--
-- TOC entry 3318 (class 2606 OID 294379)
-- Name: statuses statuses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (id);


--
-- TOC entry 3308 (class 2606 OID 294346)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3341 (class 1259 OID 294663)
-- Name: http_sessions_expiry_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX http_sessions_expiry_idx ON public.http_sessions USING btree (expires_on);


--
-- TOC entry 3342 (class 1259 OID 294664)
-- Name: http_sessions_key_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX http_sessions_key_idx ON public.http_sessions USING btree (key);


--
-- TOC entry 3356 (class 2606 OID 294469)
-- Name: commits commits_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commits
    ADD CONSTRAINT commits_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id);


--
-- TOC entry 3357 (class 2606 OID 294474)
-- Name: commits commits_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.commits
    ADD CONSTRAINT commits_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3351 (class 2606 OID 294429)
-- Name: job_dependencies job_dependencies_child_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_dependencies
    ADD CONSTRAINT job_dependencies_child_id_fkey FOREIGN KEY (child_id) REFERENCES public.jobs(id);


--
-- TOC entry 3352 (class 2606 OID 294424)
-- Name: job_dependencies job_dependencies_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_dependencies
    ADD CONSTRAINT job_dependencies_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.jobs(id);


--
-- TOC entry 3346 (class 2606 OID 294494)
-- Name: jobs job_pipeline_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT job_pipeline_fk FOREIGN KEY (pipeline_id) REFERENCES public.pipelines(id) NOT VALID;


--
-- TOC entry 3347 (class 2606 OID 294454)
-- Name: job_workers job_workers_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers
    ADD CONSTRAINT job_workers_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id);


--
-- TOC entry 3348 (class 2606 OID 294489)
-- Name: job_workers job_workers_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers
    ADD CONSTRAINT job_workers_pipeline_id_fkey FOREIGN KEY (pipeline_worker_id) REFERENCES public.pipeline_workers(id);


--
-- TOC entry 3349 (class 2606 OID 294444)
-- Name: job_workers job_workers_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers
    ADD CONSTRAINT job_workers_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.statuses(id);


--
-- TOC entry 3350 (class 2606 OID 294570)
-- Name: job_workers job_workers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_workers
    ADD CONSTRAINT job_workers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) NOT VALID;


--
-- TOC entry 3367 (class 2606 OID 294648)
-- Name: oauth_users ouath_users_provider_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_users
    ADD CONSTRAINT ouath_users_provider_id_fkey FOREIGN KEY (oauth_provider_id) REFERENCES public.oauth_providers(id) NOT VALID;


--
-- TOC entry 3368 (class 2606 OID 294633)
-- Name: oauth_users ouath_users_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth_users
    ADD CONSTRAINT ouath_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3364 (class 2606 OID 294609)
-- Name: pipeline_users pipeline_users_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_users
    ADD CONSTRAINT pipeline_users_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipelines(id);


--
-- TOC entry 3365 (class 2606 OID 294619)
-- Name: pipeline_users pipeline_users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_users
    ADD CONSTRAINT pipeline_users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 3366 (class 2606 OID 294614)
-- Name: pipeline_users pipeline_users_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_users
    ADD CONSTRAINT pipeline_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3353 (class 2606 OID 294449)
-- Name: pipeline_workers pipeline_workers_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_workers
    ADD CONSTRAINT pipeline_workers_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipelines(id);


--
-- TOC entry 3354 (class 2606 OID 294691)
-- Name: pipeline_workers pipeline_workers_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_workers
    ADD CONSTRAINT pipeline_workers_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.statuses(id) NOT VALID;


--
-- TOC entry 3355 (class 2606 OID 294484)
-- Name: pipeline_workers pipeline_workers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_workers
    ADD CONSTRAINT pipeline_workers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3345 (class 2606 OID 294479)
-- Name: pipelines pipelines_commit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipelines
    ADD CONSTRAINT pipelines_commit_id_fkey FOREIGN KEY (commit_id) REFERENCES public.commits(id);


--
-- TOC entry 3362 (class 2606 OID 294564)
-- Name: stage_logs stage_logs_log_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_logs
    ADD CONSTRAINT stage_logs_log_type_id_fkey FOREIGN KEY (log_type_id) REFERENCES public.log_types(id);


--
-- TOC entry 3363 (class 2606 OID 294559)
-- Name: stage_logs stage_logs_stage_worker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_logs
    ADD CONSTRAINT stage_logs_stage_worker_id_fkey FOREIGN KEY (stage_worker_id) REFERENCES public.stage_workers(id);


--
-- TOC entry 3359 (class 2606 OID 294540)
-- Name: stage_workers stage_workers_job_worker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_workers
    ADD CONSTRAINT stage_workers_job_worker_id_fkey FOREIGN KEY (job_worker_id) REFERENCES public.job_workers(id);


--
-- TOC entry 3360 (class 2606 OID 294545)
-- Name: stage_workers stage_workers_stage_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_workers
    ADD CONSTRAINT stage_workers_stage_id_fkey FOREIGN KEY (stage_id) REFERENCES public.stages(id);


--
-- TOC entry 3361 (class 2606 OID 294535)
-- Name: stage_workers stage_workers_status_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stage_workers
    ADD CONSTRAINT stage_workers_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.statuses(id);


--
-- TOC entry 3358 (class 2606 OID 294523)
-- Name: stages stages_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stages
    ADD CONSTRAINT stages_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id);


-- Completed on 2024-11-17 01:42:31 EST

--
-- PostgreSQL database dump complete
--

