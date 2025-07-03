\set ON_ERROR_STOP 1

SELECT
  'CREATE ROLE ' || quote_ident(:'worker') || ' LOGIN'
WHERE NOT EXISTS (
  SELECT 1 FROM pg_roles WHERE rolname = :'worker'
)\gexec

ALTER ROLE :worker WITH PASSWORD :'wpw';

SELECT
  'CREATE DATABASE ' || quote_ident(:'appdb') ||
  ' OWNER '          || quote_ident(:'worker')
WHERE NOT EXISTS (
  SELECT 1 FROM pg_database WHERE datname = :'appdb'
)\gexec

\connect :appdb
CREATE SCHEMA IF NOT EXISTS api AUTHORIZATION :"worker";
ALTER  ROLE :"worker"  SET search_path = api, public;
GRANT ALL PRIVILEGES ON DATABASE :"appdb" TO :"worker";
GRANT ALL PRIVILEGES ON SCHEMA   api      TO :"worker";
