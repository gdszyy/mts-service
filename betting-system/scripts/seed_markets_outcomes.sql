-- 为测试数据添加 markets 和 outcomes
-- 确保先运行 seed_test_data.sql

-- 为 Premier League 比赛添加 markets 和 outcomes

-- 获取 event_id
DO $$
DECLARE
    event_id_1 BIGINT;
    event_id_2 BIGINT;
    market_id BIGINT;
BEGIN
    -- 获取第一场比赛的 ID
    SELECT id INTO event_id_1 FROM events WHERE external_id = 'sr:match:1001';
    
    IF event_id_1 IS NOT NULL THEN
        -- 为第一场比赛添加 1x2 market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_1, '1x2', NULL, 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        -- 添加 outcomes
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, '1', 2.50, 'active', NOW(), NOW()),
            (market_id, 'x', 3.20, 'active', NOW(), NOW()),
            (market_id, '2', 2.80, 'active', NOW(), NOW());
        
        -- 添加 handicap market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_1, 'handicap', 'handicap=-1.5', 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, '1', 1.85, 'active', NOW(), NOW()),
            (market_id, '2', 1.95, 'active', NOW(), NOW());
        
        -- 添加 totals market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_1, 'totals', 'total=2.5', 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, 'over', 1.90, 'active', NOW(), NOW()),
            (market_id, 'under', 1.90, 'active', NOW(), NOW());
    END IF;
    
    -- 获取第二场比赛的 ID
    SELECT id INTO event_id_2 FROM events WHERE external_id = 'sr:match:1002';
    
    IF event_id_2 IS NOT NULL THEN
        -- 为第二场比赛添加 1x2 market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_2, '1x2', NULL, 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, '1', 2.10, 'active', NOW(), NOW()),
            (market_id, 'x', 3.40, 'active', NOW(), NOW()),
            (market_id, '2', 3.50, 'active', NOW(), NOW());
        
        -- 添加 handicap market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_2, 'handicap', 'handicap=-1.0', 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, '1', 1.75, 'active', NOW(), NOW()),
            (market_id, '2', 2.05, 'active', NOW(), NOW());
        
        -- 添加 totals market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_2, 'totals', 'total=3.0', 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, 'over', 2.00, 'active', NOW(), NOW()),
            (market_id, 'under', 1.80, 'active', NOW(), NOW());
        
        -- 添加 both_teams_to_score market
        INSERT INTO markets (event_id, market_type, specifier, status, created_at, updated_at)
        VALUES (event_id_2, 'both_teams_to_score', NULL, 'active', NOW(), NOW())
        RETURNING id INTO market_id;
        
        INSERT INTO outcomes (market_id, outcome_id, odds, status, created_at, updated_at)
        VALUES 
            (market_id, 'yes', 1.65, 'active', NOW(), NOW()),
            (market_id, 'no', 2.20, 'active', NOW(), NOW());
    END IF;
END $$;

-- 验证数据
SELECT 'Markets Count:' as info, COUNT(*) as count FROM markets
UNION ALL
SELECT 'Outcomes Count:', COUNT(*) FROM outcomes
UNION ALL
SELECT 'Events with Markets:', COUNT(DISTINCT event_id) FROM markets;

-- 显示每个 event 的 market 统计
SELECT 
    e.external_id,
    e.home_team,
    e.away_team,
    COUNT(m.id) as market_count,
    STRING_AGG(DISTINCT m.market_type, ', ' ORDER BY m.market_type) as market_types
FROM events e
LEFT JOIN markets m ON e.id = m.event_id
GROUP BY e.id, e.external_id, e.home_team, e.away_team
ORDER BY e.external_id;
