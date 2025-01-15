\connect postgres;
CREATE TABLE IF NOT EXISTS public.users
(
    userid   SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(60) NOT NULL,
    balance  INT DEFAULT 0
);

