-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for performance
CREATE INDEX idx_todos_completed ON todos(completed);

-- Create trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_todos_updated_at 
    BEFORE UPDATE ON todos 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
    
-- Insert 10 mock data samples
INSERT INTO todos (title, description, completed) VALUES 
('Learn Docker Basics', 'Understand containers, images, and docker-compose.', true),
('Setup PostgreSQL', 'Install and configure Postgres via Docker Desktop.', true),
('Build Todo API', 'Create a REST API to handle todo tasks.', false),
('Grocery Shopping', 'Buy milk, eggs, bread, and coffee beans.', false),
('Read Technical Book', 'Read 20 pages of "Clean Code" by Robert C. Martin.', false),
('Morning Workout', '30 minutes of cardio and light weightlifting.', true),
('Prepare Presentation', 'Create slides for the upcoming sprint review.', false),
('Book Flight Tickets', 'Check prices for the holiday trip to London.', false),
('Write Blog Post', 'Summarize how to handle database migrations.', false),
('House Cleaning', 'Vacuum the living room and change bed sheets.', true);
