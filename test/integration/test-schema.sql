CREATE USER test WITH PASSWORD 'integration';
CREATE DATABASE test-user-db;
GRANT ALL PRIVILEGES ON DATABASE test-user-db TO test;
REVOKE ALL PRIVILEGES ON DATABASE test-user-db FROM public;
