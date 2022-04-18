CREATE SCHEMA IF NOT EXISTS nsiv29test;
CREATE EXTENSION postgis WITH SCHEMA nsiv29test;

UPDATE pg_extension
  SET extrelocatable = TRUE
    WHERE extname = 'postgis';

ALTER EXTENSION postgis
  SET SCHEMA nsiv29test;

/* ALTER EXTENSION postgis */
/*   UPDATE TO "2.5.2next"; */

/* ALTER EXTENSION postgis */
/*   UPDATE TO "2.5.2"; */

/* SET search_path = 'nsiv29test'; */

ALTER ROLE dbuser SET search_path = nsiv29test;
