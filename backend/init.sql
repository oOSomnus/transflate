CREATE TABLE IF NOT EXISTS public.user (
                                           userid SERIAL PRIMARY KEY,
                                           username VARCHAR(50) NOT NULL,
                                           password VARCHAR(50) NOT NULL,
                                           balance DECIMAL(10, 2) DEFAULT 0.00
);
