-- Insert sample services
INSERT INTO services (id, name, description, created_at, updated_at)
VALUES
  (1, 'Config Service', 'Handles configuration management', NOW(), NOW()),
  (2, 'Payment Service', 'Handles payment processing', NOW(), NOW()),
  (3, 'User Service', 'Manages user accounts and profiles', NOW(), NOW()),
  (4, 'Authentication Service', 'Handles authentication and authorization', NOW(), NOW()),
  (5, 'Supabase Service', 'Integrates with Supabase backend', NOW(), NOW()),
  (6, 'Kong Service', 'API gateway and traffic control', NOW(), NOW()),
  (7, 'Infrastructure Service', 'Manages infrastructure resources', NOW(), NOW());

-- Insert sample versions for each service
INSERT INTO versions (service_id, version, description, is_active, created_at, updated_at)
VALUES
  (1, '1.0.0', 'Initial release', TRUE, NOW(), NOW()),
  (1, '1.1.0', 'Added new config endpoints', TRUE, NOW(), NOW()),
  (2, '1.0.0', 'Initial payment integration', TRUE, NOW(), NOW()),
  (2, '1.2.0', 'Support for new payment provider', TRUE, NOW(), NOW()),
  (3, '1.0.0', 'User service MVP', TRUE, NOW(), NOW()),
  (3, '1.1.0', 'Profile picture support', TRUE, NOW(), NOW()),
  (3, '2.0.0', 'Major refactor', TRUE, NOW(), NOW()),
  (4, '1.0.0', 'Basic authentication', TRUE, NOW(), NOW()),
  (4, '1.1.0', 'OAuth2 support', TRUE, NOW(), NOW()),
  (5, '1.0.0', 'Supabase integration', TRUE, NOW(), NOW()),
  (6, '1.0.0', 'Kong gateway setup', TRUE, NOW(), NOW()),
  (7, '1.0.0', 'Infrastructure bootstrap', TRUE, NOW(), NOW()),
  (7, '1.1.0', 'Terraform support', TRUE, NOW(), NOW());

SELECT setval('services_id_seq', (SELECT MAX(id) FROM services));
SELECT setval('versions_id_seq', (SELECT MAX(id) FROM versions));