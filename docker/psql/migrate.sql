DROP TABLE IF EXISTS source;

DROP TABLE IF EXISTS article;

CREATE TABLE source (
    id SERIAL NOT NULL,
    name TEXT,
    description TEXT,
    home_url TEXT,
    api_name TEXT NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE article (
    id SERIAL NOT NULL,
    url TEXT,
    pub_date TEXT,
    title TEXT,
    description TEXT,
    cover_url TEXT,
    add_date TEXT,
    source_id INTEGER,
    PRIMARY KEY (id),
    CONSTRAINT fk_source FOREIGN KEY (source_id) REFERENCES source (id)
);

INSERT INTO source (name, description, home_url, api_name) VALUES 
    ('Jet Propulsion Laboratory', 'JPL is a research and development lab federally funded by NASA and managed by Caltech.', 'https://www.jpl.nasa.gov', 'jpl'),
    ('Vestirama', 'Оренбургская государственная телерадиовещательная компания.', 'https://vestirama.ru', 'vestirama'),
    ('National Geographic', 'A world leader in geography, cartography and exploration.', 'https://www.nationalgeographic.com', 'natgeo');