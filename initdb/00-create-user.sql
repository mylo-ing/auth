CREATE USER api_worker WITH PASSWORD 'password_for_dev_only';

CREATE DATABASE my_local OWNER api_worker;

GRANT ALL PRIVILEGES ON DATABASE my_local TO api_worker;