DROP TABLE IF EXISTS source;

DROP TABLE IF EXISTS article;

CREATE TABLE source (
    id SERIAL NOT NULL,
    name TEXT,
    home_url TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE article (
    id SERIAL NOT NULL,
    url TEXT,
    date TEXT,
    title TEXT,
    description TEXT,
    cover_url TEXT,
    source_id INTEGER,
    PRIMARY KEY (id),
    CONSTRAINT fk_source FOREIGN KEY (source_id) REFERENCES source (id)
);

INSERT INTO source (name, home_url) VALUES 
    ('Jet Propulsion Laboratory', 'https://www.jpl.nasa.gov'),
    ('Vestirama', 'https://vestirama.ru');