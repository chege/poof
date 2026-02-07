CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT,
    full_name TEXT,
    email TEXT,
    company TEXT,
    phone TEXT,
    ip_address TEXT,
    bio TEXT,
    status TEXT
);

INSERT INTO users (username, full_name, email, company, phone, ip_address, bio, status) VALUES
('jdoe', 'John Doe', 'john.doe@example.com', 'Acme Corp', '+1-555-1234', '192.168.1.1', 'Original bio 1', 'active'),
('asmith', 'Alice Smith', 'alice.smith@test.org', 'Globex', '+1-555-5678', '10.0.0.1', 'Original bio 2', 'inactive');
