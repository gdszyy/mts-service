-- =====================================================
-- Tournament (League) Popularity Scoring Algorithm - SQL Implementation
-- Version: 1.1
-- Author: Manus AI
-- Date: 2025-11-27
-- =====================================================

-- This script calculates and stores the popularity score for tournaments.

-- Step 1: Create or ensure the existence of the scoring table
-- =====================================================
CREATE TABLE IF NOT EXISTS tournament_popularity_scores (
    tournament_id VARCHAR(50) PRIMARY KEY,
    tournament_name VARCHAR(255),
    category_id VARCHAR(50),
    category_name VARCHAR(255),
    sport_name VARCHAR(255),
    
    -- Statistical Data
    total_events INT DEFAULT 0,
    avg_market_count NUMERIC(10, 2) DEFAULT 0,
    
    -- Scoring Data
    tournament_tier_score INT DEFAULT 0,
    market_depth_score INT DEFAULT 0,
    final_popularity_score NUMERIC(5, 2) DEFAULT 0,
    
    -- Metadata
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_tps_final_score ON tournament_popularity_scores(final_popularity_score DESC);
CREATE INDEX IF NOT EXISTS idx_tps_sport_name ON tournament_popularity_scores(sport_name);

-- Step 2: Calculate and upsert tournament popularity scores
-- =====================================================

WITH 
-- 2.1 Calculate market count for each event
event_market_counts AS (
    SELECT 
        event_id,
        COUNT(id) as market_count
    FROM markets
    GROUP BY event_id
),

-- 2.2 Calculate average market count for each tournament
tournament_market_stats AS (
    SELECT 
        te.tournament_id,
        COUNT(DISTINCT te.event_id) as total_events,
        COALESCE(AVG(emc.market_count), 0) as avg_markets
    FROM tracked_events te
    LEFT JOIN event_market_counts emc ON te.event_id = emc.event_id
    WHERE te.tournament_id IS NOT NULL
    GROUP BY te.tournament_id
    HAVING COUNT(DISTINCT te.event_id) > 0
),

-- 2.3 Join with tournament and category info and assign scores
scored_tournaments AS (
    SELECT 
        tms.tournament_id,
        t.name as tournament_name,
        c.id as category_id,
        c.name as category_name,
        s.name as sport_name,
        tms.total_events,
        tms.avg_markets,
        -- Calculate Tournament Tier Score (1-10)
        CASE 
            WHEN t.name ILIKE '%World Cup%' OR c.name ILIKE '%World Cup%' OR t.name ILIKE '%Olympics%' OR c.name ILIKE '%Olympics%' THEN 10
            WHEN t.name ILIKE '%Champions League%' OR c.name ILIKE '%Champions League%' THEN 9
            WHEN t.name ILIKE '%Premier League%' OR c.name ILIKE '%Premier League%' OR t.name ILIKE '%NBA%' OR c.name ILIKE '%NBA%' THEN 8
            WHEN c.name ILIKE '%International Clubs%' OR c.name ILIKE '%International%' THEN 7
            WHEN c.country_code IS NOT NULL AND c.country_code != '' THEN 6
            ELSE 4
        END as tournament_tier_score,
        
        -- Calculate Market Depth Score (1-10)
        CASE 
            WHEN tms.avg_markets > 300 THEN 10
            WHEN tms.avg_markets > 200 THEN 9
            WHEN tms.avg_markets > 150 THEN 8
            WHEN tms.avg_markets > 100 THEN 7
            WHEN tms.avg_markets > 50 THEN 6
            WHEN tms.avg_markets > 30 THEN 5
            WHEN tms.avg_markets > 20 THEN 4
            WHEN tms.avg_markets > 10 THEN 3
            WHEN tms.avg_markets > 5 THEN 2
            ELSE 1
        END as market_depth_score
    FROM tournament_market_stats tms
    INNER JOIN tournaments t ON tms.tournament_id = t.id
    INNER JOIN categories c ON t.category_id = c.id
    INNER JOIN sports s ON c.sport_id = s.id
)

-- 2.4 Insert or update the final scores
INSERT INTO tournament_popularity_scores (
    tournament_id,
    tournament_name,
    category_id,
    category_name,
    sport_name,
    total_events,
    avg_market_count,
    tournament_tier_score,
    market_depth_score,
    final_popularity_score,
    updated_at
)
SELECT 
    tournament_id,
    tournament_name,
    category_id,
    category_name,
    sport_name,
    total_events,
    ROUND(avg_markets, 2),
    tournament_tier_score,
    market_depth_score,
    -- Calculate the final weighted score
    ROUND((tournament_tier_score * 0.5 + market_depth_score * 0.5), 2) as final_popularity_score,
    CURRENT_TIMESTAMP
FROM scored_tournaments
-- Use ON CONFLICT to perform an upsert operation
ON CONFLICT (tournament_id) 
DO UPDATE SET
    tournament_name = EXCLUDED.tournament_name,
    category_id = EXCLUDED.category_id,
    category_name = EXCLUDED.category_name,
    sport_name = EXCLUDED.sport_name,
    total_events = EXCLUDED.total_events,
    avg_market_count = EXCLUDED.avg_market_count,
    tournament_tier_score = EXCLUDED.tournament_tier_score,
    market_depth_score = EXCLUDED.market_depth_score,
    final_popularity_score = EXCLUDED.final_popularity_score,
    updated_at = CURRENT_TIMESTAMP;

-- =====================================================
-- End of Script
-- =====================================================
