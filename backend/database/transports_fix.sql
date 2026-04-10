-- Fix for transports column NULL handling
-- The transports column in PostgreSQL was being stored as empty array or NULL
-- but Go's pq library now handles this properly with pq.StringArray

-- Check current state of credentials
SELECT id, user_id, credential_id, transports FROM credentials LIMIT 5;

-- If you need to convert NULL transports to empty arrays:
UPDATE credentials SET transports = '{}' WHERE transports IS NULL;

-- Verify the fix
SELECT COUNT(*) as total_credentials,
       COUNT(*) FILTER (WHERE transports IS NULL) as null_transports,
       COUNT(*) FILTER (WHERE transports = '{}') as empty_transports
FROM credentials;

