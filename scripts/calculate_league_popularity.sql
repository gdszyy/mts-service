-- =====================================================
-- League Popularity Scoring Algorithm - SQL Implementation
-- Version: 1.0
-- Author: Manus AI
-- Date: 2025-11-27
-- =====================================================

-- This script calculates and stores the popularity score for all leagues.
-- It is designed to be run periodically (e.g., via a cron job).

-- Step 1: Create or ensure the existence of the scoring table
-- =====================================================
CREATE TABLE IF NOT EXISTS league_popularity_scores (
    category_id VARCHAR(50) PRIMARY KEY,
    category_name VARCHAR(255),
    sport_name VARCHAR(255),
    country_code VARCHAR(10),
    
    -- Statistical Data
    total_events INT DEFAULT 0,
    avg_market_count NUMERIC(10, 2) DEFAULT 0,
    max_market_count INT DEFAULT 0,
    
    -- Scoring Data
    league_tier_score INT DEFAULT 0,
    market_depth_score INT DEFAULT 0,
    final_popularity_score NUMERIC(5, 2) DEFAULT 0,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_lps_final_score ON league_popularity_scores(final_popularity_score DESC);
CREATE INDEX IF NOT EXISTS idx_lps_sport_name ON league_popularity_scores(sport_name);

-- Step 2: Calculate and upsert league popularity scores
-- =====================================================

-- Use a Common Table Expression (CTE) to build the data step-by-step
WITH 
-- 2.1 Calculate market statistics for each league
league_market_stats AS (
    SELECT 
        c.id as category_id,
        c.name as category_name,
        c.country_code,
        s.name as sport_name,
        COUNT(DISTINCT te.event_id) as total_events,
        COALESCE(AVG(mc.market_count), 0) as avg_markets,
        COALESCE(MAX(mc.market_count), 0) as max_markets
    FROM categories c
    INNER JOIN sports s ON c.sport_id = s.id
    -- Use INNER JOIN to only include events with a valid category and sport
    INNER JOIN tracked_events te ON c.id = te.category_id AND te.sport_id = s.id
    -- Use LEFT JOIN for markets, as some events may have 0 markets
    LEFT JOIN (
        SELECT event_id, COUNT(*) as market_count
        FROM markets
        GROUP BY event_id
    ) mc ON te.event_id = mc.event_id
    GROUP BY c.id, c.name, c.country_code, s.name
    -- Filter out leagues with no events
    HAVING COUNT(DISTINCT te.event_id) > 0
),

-- 2.2 Assign scores based on the defined rules
scored_leagues AS (
    SELECT 
        *,
        -- Calculate League Tier Score (1-10)
        CASE 
            WHEN category_name ILIKE '%World Cup%' OR category_name ILIKE '%Olympics%' OR category_name ILIKE '%World Championship%' THEN 10
            WHEN category_name ILIKE '%Champions League%' OR category_name ILIKE '%Europa League%' OR category_name ILIKE '%Copa America%' THEN 9
            WHEN category_name ILIKE '%Premier League%' OR category_name ILIKE '%La Liga%' OR category_name ILIKE '%Serie A%' OR category_name ILIKE '%Bundesliga%' OR category_name ILIKE '%NBA%' THEN 8
            WHEN category_name ILIKE '%International Clubs%' OR category_name ILIKE '%International%' THEN 7
            WHEN country_code IS NOT NULL AND country_code != '' THEN 6
            ELSE 4
        END as league_tier_score,
        
        -- Calculate Market Depth Score (1-10)
        CASE 
            WHEN avg_markets > 300 THEN 10
            WHEN avg_markets > 200 THEN 9
            WHEN avg_markets > 150 THEN 8
            WHEN avg_markets > 100 THEN 7
            WHEN avg_markets > 50 THEN 6
            WHEN avg_markets > 30 THEN 5
            WHEN avg_markets > 20 THEN 4
            WHEN avg_markets > 10 THEN 3
            WHEN avg_markets > 5 THEN 2
            ELSE 1
        END as market_depth_score
    FROM league_market_stats
)

-- 2.3 Insert or update the final scores into the destination table
INSERT INTO league_popularity_scores (
    category_id,
    category_name,
    sport_name,
    country_code,
    total_events,
    avg_market_count,
    max_market_count,
    league_tier_score,
    market_depth_score,
    final_popularity_score,
    updated_at
)
SELECT 
    category_id,
    category_name,
    sport_name,
    country_code,
    total_events,
    ROUND(avg_markets, 2),
    max_markets,
    league_tier_score,
    market_depth_score,
    -- Calculate the final weighted score
    ROUND((league_tier_score * 0.5 + market_depth_score * 0.5), 2) as final_popularity_score,
    CURRENT_TIMESTAMP
FROM scored_leagues
-- Use ON CONFLICT to perform an upsert operation
ON CONFLICT (category_id) 
DO UPDATE SET
    category_name = EXCLUDED.category_name,
    sport_name = EXCLUDED.sport_name,
    country_code = EXCLUDED.country_code,
    total_events = EXCLUDED.total_events,
    avg_market_count = EXCLUDED.avg_market_count,
    max_market_count = EXCLUDED.max_market_count,
    league_tier_score = EXCLUDED.league_tier_score,
    market_depth_score = EXCLUDED.market_depth_score,
    final_popularity_score = EXCLUDED.final_popularity_score,
    updated_at = CURRENT_TIMESTAMP;

-- Step 3: (Optional) Query the results for verification
-- =====================================================

-- Uncomment the following lines to see the top 50 leagues after the script runs
/*
SELECT 
    category_name,
    sport_name,
    country_code,
    total_events,
    avg_market_count,
    final_popularity_score
FROM league_popularity_scores
ORDER BY final_popularity_score DESC, avg_market_count DESC
LIMIT 50;
*/

-- =====================================================
-- End of Script
-- =====================================================
