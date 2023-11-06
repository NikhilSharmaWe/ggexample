-- DDL: Create the 'questions' table
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text VARCHAR(500),
    options TEXT[],
    answer TEXT
);

-- DML: Insert initial data
INSERT INTO questions (text, options, answer)
VALUES
    ('Sample Question 1', ARRAY['Option A', 'Option B', 'Option C'], 'Option A'),
    ('Sample Question 2', ARRAY['Option X', 'Option Y', 'Option Z'], 'Option Z');
