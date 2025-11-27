-- =====================================================
-- Event (Match) Popularity Scoring Algorithm - SQL Implementation
-- Version: 1.0
-- Author: Manus AI
-- Date: 2025-11-27
-- =====================================================

-- This script calculates and stores the popularity score for individual events.

-- Step 1: Create or ensure the existence of the scoring table
-- =====================================================
CREATE TABLE IF NOT EXISTS event_popularity_scores (
    event_id VARCHAR(50) PRIMARY KEY,
    tournament_id VARCHAR(50),
    category_id VARCHAR(50),
    sport_id VARCHAR(50),
    market_count INT DEFAULT 0,
    popularity_score INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_eps_popularity_score ON event_popularity_scores(popularity_score DESC);
CREATE INDEX IF NOT EXISTS idx_eps_tournament_id ON event_popularity_scores(tournament_id);

-- Step 2: Calculate and upsert event popularity scores
-- =====================================================

WITH 
-- 2.1 Calculate market count for each event
event_market_counts AS (
    SELECT 
        te.event_id,
        te.tournament_id,
        te.category_id,
        te.sport_id,
        COUNT(m.id) as market_count
    FROM tracked_events te
    -- Use LEFT JOIN to include events with 0 markets
    LEFT JOIN markets m ON te.event_id = m.event_id
    GROUP BY te.event_id, te.tournament_id, te.category_id, te.sport_id
),

-- 2.2 Assign scores based on market count
scored_events AS (
    SELECT 
        *,
        CASE 
            WHEN market_count > 400 THEN 10
            WHEN market_count > 300 THEN 9
            WHEN market_count > 200 THEN 8
            WHEN market_count > 150 THEN 7
            WHEN market_count > 100 THEN 6
            WHEN market_count > 50 THEN 5
            WHEN market_count > 30 THEN 4
            WHEN market_count > 10 THEN 3
            WHEN market_count > 5 THEN 2
            ELSE 1
        END as popularity_score
    FROM event_market_counts
)

-- 2.3 Insert or update the scores into the destination table
INSERT INTO event_popularity_scores (
    event_id,
    tournament_id,
    category_id,
    sport_id,
    market_count,
    popularity_score,
    updated_at
)
SELECT 
    event_id,
    tournament_id,
    category_id,
    sport_id,
    market_count,
    popularity_score,
    CURRENT_TIMESTAMP
FROM scored_events
-- Use ON CONFLICT to perform an upsert operation
ON CONFLICT (event_id) 
DO UPDATE SET
    tournament_id = EXCLUDED.tournament_id,
    category_id = EXCLUDED.category_id,
    sport_id = EXCLUDED.sport_id,
    market_count = EXCLUDED.market_count,
    popularity_score = EXCLUDED.popularity_score,
    updated_at = CURRENT_TIMESTAMP;

-- =====================================================
-- End of Script
-- =====================================================
