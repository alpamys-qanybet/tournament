CREATE TABLE public.teams (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	division varchar(1) NULL,
	wins int4 NOT NULL DEFAULT 0,
	draws int4 NOT NULL DEFAULT 0,
	loses int4 NOT NULL DEFAULT 0,
	goals_scored int4 NOT NULL DEFAULT 0,
	goals_conceded int4 NOT NULL DEFAULT 0,
	goal_diff int4 NOT NULL DEFAULT 0,
	points int4 NOT NULL DEFAULT 0,
	CONSTRAINT teams_name_un UNIQUE (name),
	CONSTRAINT teams_pk PRIMARY KEY (id)
);
CREATE INDEX teams_division_idx ON public.teams USING btree (division);


CREATE TABLE public.matches (
	id serial4 NOT NULL,
	first_team_score int4 NULL,
	second_team_score int4 NULL,
	winner_id int4 NULL,
	match_type varchar(2) NOT NULL,
	played bool NOT NULL DEFAULT false,
	first_team_id int4 NOT NULL,
	second_team_id int4 NOT NULL,
	CONSTRAINT matches_pk PRIMARY KEY (id),
	CONSTRAINT matches_first_team_fk FOREIGN KEY (first_team_id) REFERENCES public.teams(id),
	CONSTRAINT matches_second_team_fk FOREIGN KEY (second_team_id) REFERENCES public.teams(id)
);
CREATE INDEX matches_match_type_idx ON public.matches USING btree (match_type);