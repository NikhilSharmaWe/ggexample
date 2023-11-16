-- DDL: Create the 'questions' table
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text VARCHAR(500),
    options TEXT[],
    answer TEXT
);

-- DML: Insert initial data to questions table 
INSERT INTO questions (text, options, answer)
VALUES
    ('Sample Question 1', ARRAY['Option A', 'Option B', 'Option C'], 'Option A'),
    ('Sample Question 2', ARRAY['Option X', 'Option Y', 'Option Z'], 'Option Z');


-- DDL: Create the 'quiz_sessions' table
CREATE TABLE IF NOT EXISTS quiz_sessions (
    id TEXT PRIMARY KEY,
);

-- DDL: Create the 'quiz_responses' table
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    quiz_session_id TEXT REFERENCES quiz_sessions(id),
    question_id INT REFERENCES questions(id),
    answer TEXT,
    is_correct BOOLEAN
);

-- DML: Insert initial data to 'quiz_responses' table
INSERT INTO quiz_responses (quiz_session_id, question_id, answer, is_correct)
VALUES
    ('Sample Session ID 1', 'Sample Question ID 1', 'ANSWER1', TRUE),
    ('Sample Session ID 2', 'Sample Question ID 2', 'ANSWER12', FALSE);