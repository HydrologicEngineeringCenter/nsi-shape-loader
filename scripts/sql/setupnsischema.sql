CREATE SCHEMA IF NOT EXISTS nsiv29;
CREATE extension postgis;

UPDATE pg_extension
  SET extrelocatable = TRUE
    WHERE extname = 'postgis';

ALTER EXTENSION postgis
  SET SCHEMA nsiv29;

ALTER EXTENSION postgis
  UPDATE TO "2.5.2next";

ALTER EXTENSION postgis
  UPDATE TO "2.5.2";

ALTER ROLE dbuser SET search_path = nsiv29;

