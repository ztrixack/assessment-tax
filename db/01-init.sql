-- Create the allowances table
CREATE TABLE IF NOT EXISTS allowances (
    id SERIAL PRIMARY KEY,
    personal DECIMAL(10, 2) NOT NULL,
    donation DECIMAL(10, 2) NOT NULL,
    k_receipt DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- Insert a record if it does not exist
INSERT INTO allowances (personal, donation, k_receipt)
SELECT 60000, 100000, 50000
WHERE NOT EXISTS (
    SELECT 1 FROM allowances
    WHERE personal = 60000 AND donation = 100000 AND k_receipt = 50000
);