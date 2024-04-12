-- Check if the database exists before creating it
SELECT 'CREATE DATABASE demo_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'demo_db');

-- Use the newly created or existing database
\c demo_db

-- Create table if it does not exist
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    status TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
