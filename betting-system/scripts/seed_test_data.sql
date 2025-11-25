-- 测试数据种子脚本
-- 用于测试 Category 和 Tournament 接口

-- 插入测试 Sports 数据
INSERT INTO sports (external_id, name, created_at, updated_at) VALUES
('sr:sport:1', 'Football', NOW(), NOW()),
('sr:sport:2', 'Basketball', NOW(), NOW()),
('sr:sport:5', 'Tennis', NOW(), NOW())
ON CONFLICT (external_id) DO NOTHING;

-- 插入测试 Categories 数据
INSERT INTO categories (external_id, sport_id, name, created_at, updated_at)
SELECT 'sr:category:1', s.id, 'England', NOW(), NOW()
FROM sports s WHERE s.external_id = 'sr:sport:1'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO categories (external_id, sport_id, name, created_at, updated_at)
SELECT 'sr:category:7', s.id, 'Spain', NOW(), NOW()
FROM sports s WHERE s.external_id = 'sr:sport:1'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO categories (external_id, sport_id, name, created_at, updated_at)
SELECT 'sr:category:31', s.id, 'Germany', NOW(), NOW()
FROM sports s WHERE s.external_id = 'sr:sport:1'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO categories (external_id, sport_id, name, created_at, updated_at)
SELECT 'sr:category:34', s.id, 'NBA', NOW(), NOW()
FROM sports s WHERE s.external_id = 'sr:sport:2'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO categories (external_id, sport_id, name, created_at, updated_at)
SELECT 'sr:category:82', s.id, 'ATP', NOW(), NOW()
FROM sports s WHERE s.external_id = 'sr:sport:5'
ON CONFLICT (external_id) DO NOTHING;

-- 插入测试 Tournaments 数据
INSERT INTO tournaments (external_id, category_id, name, scheduled, scheduled_end, created_at, updated_at)
SELECT 'sr:tournament:17', c.id, 'Premier League', '2024-08-01 00:00:00', '2025-05-31 23:59:59', NOW(), NOW()
FROM categories c WHERE c.external_id = 'sr:category:1'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO tournaments (external_id, category_id, name, scheduled, scheduled_end, created_at, updated_at)
SELECT 'sr:tournament:34', c.id, 'La Liga', '2024-08-01 00:00:00', '2025-05-31 23:59:59', NOW(), NOW()
FROM categories c WHERE c.external_id = 'sr:category:7'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO tournaments (external_id, category_id, name, scheduled, scheduled_end, created_at, updated_at)
SELECT 'sr:tournament:35', c.id, 'Bundesliga', '2024-08-01 00:00:00', '2025-05-31 23:59:59', NOW(), NOW()
FROM categories c WHERE c.external_id = 'sr:category:31'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO tournaments (external_id, category_id, name, scheduled, scheduled_end, created_at, updated_at)
SELECT 'sr:tournament:132', c.id, 'NBA Regular Season', '2024-10-01 00:00:00', '2025-04-30 23:59:59', NOW(), NOW()
FROM categories c WHERE c.external_id = 'sr:category:34'
ON CONFLICT (external_id) DO NOTHING;

-- 插入测试 Events 数据
INSERT INTO events (external_id, sport_id, tournament_id, home_team, away_team, start_time, status, created_at, updated_at)
SELECT 'sr:match:1001', 'sr:sport:1', t.id, 'Manchester United', 'Liverpool', '2024-12-01 15:00:00', 'scheduled', NOW(), NOW()
FROM tournaments t WHERE t.external_id = 'sr:tournament:17'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO events (external_id, sport_id, tournament_id, home_team, away_team, start_time, status, created_at, updated_at)
SELECT 'sr:match:1002', 'sr:sport:1', t.id, 'Chelsea', 'Arsenal', '2024-12-02 17:30:00', 'scheduled', NOW(), NOW()
FROM tournaments t WHERE t.external_id = 'sr:tournament:17'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO events (external_id, sport_id, tournament_id, home_team, away_team, start_time, status, created_at, updated_at)
SELECT 'sr:match:1003', 'sr:sport:1', t.id, 'Real Madrid', 'Barcelona', '2024-12-03 20:00:00', 'scheduled', NOW(), NOW()
FROM tournaments t WHERE t.external_id = 'sr:tournament:34'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO events (external_id, sport_id, tournament_id, home_team, away_team, start_time, status, created_at, updated_at)
SELECT 'sr:match:1004', 'sr:sport:1', t.id, 'Bayern Munich', 'Borussia Dortmund', '2024-12-04 18:30:00', 'scheduled', NOW(), NOW()
FROM tournaments t WHERE t.external_id = 'sr:tournament:35'
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO events (external_id, sport_id, tournament_id, home_team, away_team, start_time, status, created_at, updated_at)
SELECT 'sr:match:2001', 'sr:sport:2', t.id, 'Lakers', 'Warriors', '2024-12-05 19:00:00', 'scheduled', NOW(), NOW()
FROM tournaments t WHERE t.external_id = 'sr:tournament:132'
ON CONFLICT (external_id) DO NOTHING;

-- 验证数据
SELECT 'Sports Count:' as info, COUNT(*) as count FROM sports
UNION ALL
SELECT 'Categories Count:', COUNT(*) FROM categories
UNION ALL
SELECT 'Tournaments Count:', COUNT(*) FROM tournaments
UNION ALL
SELECT 'Events Count:', COUNT(*) FROM events;
