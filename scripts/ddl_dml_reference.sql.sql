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


-- DDL: Create the 'quizes' table
CREATE TABLE IF NOT EXISTS quizes (
    id TEXT PRIMARY KEY,
	questionids INT[],
	progress TEXT[]
);

-- DML: Insert initial data to quizes table
INSERT INTO quizes (id, questionids, progress)
VALUES
    ('Sample ID 1', ARRAY[5, 1, 3, 6, 2], ARRAY['ANSWER1', 'ANSWER2', 'ANSWER3']),
    ('Sample ID 2', ARRAY[2, 4, 5, 6, 1], ARRAY['ANSWER1', 'ANSWER2']);