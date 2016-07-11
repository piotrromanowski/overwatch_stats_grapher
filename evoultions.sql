CREATE TABLE hero (
    id serial NOT NULL,
    name character varying(100),
    image character varying(255),
    CONSTRAINT pk_hero PRIMARY KEY (id)
);

CREATE TABLE hero_stats (
        id serial NOT NULL,
        hero_id integer REFERENCES hero (id),
    	eliminations integer,
    	deaths integer,
    	weapon_accuracy character varying(100),
    	eliminations_average integer,
    	time_played character varying(100),
    	games_played integer,
    	games_won integer,
    	win_percentage character varying(100),
    	update_time timestamp,
    	constraint pk_hero_stats PRIMARY KEY (id)
);

ALTER TABLE hero_stats ADD CONSTRAINT fk_hero_stats_hero FOREIGN KEY (hero_id) REFERENCES hero (id);

ALTER DATABASE overwatch SET timezone TO 'Canada/Eastern';


############################################################################################################################
# Example Data
############################################################################################################################

INSERT INTO hero values (
    1,
    'Soldier',
    ''
)

INSERT INTO hero_stats values (
    1,
    1,
    15,
    8,
    '75%',
    15,
    '22 hours',
    98,
    76,
    '78%',
    '2016-07-08 19:46:25.809296'
)