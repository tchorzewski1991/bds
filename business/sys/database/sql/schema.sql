-- Version: 1.1
-- Description: Create table users
CREATE TABLE users (
   uuid          UUID,
   name          TEXT,
   email         TEXT UNIQUE,
   permissions   TEXT[],
   password_hash TEXT,
   date_created  TIMESTAMP,
   date_updated  TIMESTAMP,

   PRIMARY KEY (uuid)
);

-- Version: 1.2
-- Description: Remove column name from users
ALTER TABLE users DROP COLUMN name;

-- Version: 1.3
-- Description: Create table books
CREATE TABLE books (
   id               SERIAL,
   isbn             TEXT NOT NULL CHECK (isbn <> ''),
   title            TEXT NOT NULL CHECK (title <> ''),
   author           TEXT,
   publication_year TEXT,
   publisher        TEXT,
   created_at       TIMESTAMP DEFAULT now(),
   updated_at       TIMESTAMP DEFAULT now(),

   PRIMARY KEY (id)
);