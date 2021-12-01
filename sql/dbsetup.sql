CREATE DATABASE docker_metadata ENCODING 'UTF8';
\c docker_metadata;

CREATE TABLE metadata(id character(64) not null primary key, name varchar(100), base_image varchar(100), approved boolean not null);
ALTER TABLE metadata ALTER COLUMN approved SET DEFAULT false;