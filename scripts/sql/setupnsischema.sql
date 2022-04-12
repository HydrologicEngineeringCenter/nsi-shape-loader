CREATE SCHEMA IF NOT EXISTS nsi;
CREATE extension postgis;

UPDATE pg_extension
  SET extrelocatable = TRUE
    WHERE extname = 'postgis';

ALTER EXTENSION postgis
  SET SCHEMA nsi;

ALTER EXTENSION postgis
  UPDATE TO "2.5.2next";

ALTER EXTENSION postgis
  UPDATE TO "2.5.2";

SET search_path TO nsi;
