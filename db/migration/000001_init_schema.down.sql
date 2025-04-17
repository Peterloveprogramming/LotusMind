-- Drop rules first
-- DROP RULE IF EXISTS tummo_breathing_mr_soft_deletion ON tummo_breathing_mr;
-- DROP RULE IF EXISTS tibetan_singing_bowl_mr_soft_deletion ON tibetan_singing_bowl_mr_soft_deletion;
-- DROP RULE IF EXISTS session_logs_soft_deletion_for_tummo_breathing_mr ON session_logs;
-- DROP RULE IF EXISTS session_logs_soft_deletion_for_tibetan_singing_bowl_mr ON session_logs;
-- DROP RULE IF EXISTS session_logs_soft_deletion ON session_logs;
-- DROP RULE IF EXISTS users_session_logs_soft_deletion_for_session_logs ON users;
-- DROP RULE IF EXISTS users_soft_deletion ON users;


-- Drop dependent tables
DROP TABLE IF EXISTS tibetan_singing_bowl_mr;
DROP TABLE IF EXISTS tummo_breathing_mr;

-- Drop session_logs table
DROP TABLE IF EXISTS session_logs;
DROP TABLE IF EXISTS users_profile_mr;
DROP TABLE IF EXISTS users_profile_mobile;
DROP TABLE IF EXISTS email_registrations;
DROP TABLE IF EXISTS chakra_test_results;
DROP TABLE IF EXISTS chakra_bracelet;
DROP TABLE IF EXISTS chakra_test_option_answers;



-- Drop users table
DROP TABLE IF EXISTS users;
