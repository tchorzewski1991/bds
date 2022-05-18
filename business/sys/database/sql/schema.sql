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
