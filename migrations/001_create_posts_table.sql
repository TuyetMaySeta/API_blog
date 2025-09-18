-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create GIN index for tags array for fast searching
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN (tags);

-- Create index for created_at for sorting
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);

-- Create activity_logs table
CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    post_id INTEGER NOT NULL,
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Create index for activity_logs
CREATE INDEX IF NOT EXISTS idx_activity_logs_post_id ON activity_logs (post_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_logged_at ON activity_logs (logged_at DESC);