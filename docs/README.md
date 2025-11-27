# League Popularity Scoring System

This directory contains the documentation and implementation for the **League Popularity Scoring Algorithm**, which quantifies the social popularity and market significance of sports leagues.

## Overview

The algorithm assigns a popularity score (1-10) to each league based on two key dimensions:
1. **League Tier Score**: The inherent prestige and global standing of the league
2. **Market Depth Score**: The market's engagement, measured by the average number of betting markets per event

## Documentation

- **[league_popularity_algorithm.md](./league_popularity_algorithm.md)**: Complete algorithm specification, including scoring rules, formulas, and data dependencies.

## Implementation

The algorithm is implemented as a SQL script that can be executed periodically to keep scores up-to-date:

- **SQL Script**: `../scripts/calculate_league_popularity.sql`
- **Database Table**: `league_popularity_scores`

## Quick Start

### 1. Execute the Algorithm

Run the SQL script to calculate and store popularity scores for all leagues:

```bash
psql $DATABASE_URL -f scripts/calculate_league_popularity.sql
```

Or use the Python helper script:

```bash
python3 scripts/execute_league_popularity.py
```

### 2. Query the Results

View the top leagues by popularity:

```sql
SELECT 
    category_name,
    sport_name,
    total_events,
    avg_market_count,
    final_popularity_score
FROM league_popularity_scores
ORDER BY final_popularity_score DESC
LIMIT 20;
```

### 3. Filter by Sport

Get the most popular soccer leagues:

```sql
SELECT 
    category_name,
    country_code,
    final_popularity_score
FROM league_popularity_scores
WHERE sport_name = 'Soccer'
ORDER BY final_popularity_score DESC
LIMIT 10;
```

## Score Tiers

| Tier | Score Range | Description |
| :--- | :--- | :--- |
| S-Tier | 9.0 - 10.0 | Super Popular (e.g., World Cup) |
| A-Tier | 8.0 - 8.9 | Very Popular (e.g., Premier League, Champions League) |
| B-Tier | 7.0 - 7.9 | Popular (e.g., Major national leagues) |
| C-Tier | 6.0 - 6.9 | Average (e.g., Standard national leagues) |
| D-Tier | < 6.0 | Below Average (e.g., Regional leagues) |

## Maintenance

It is recommended to run the scoring script:
- **Daily**: For production systems with high data turnover
- **Weekly**: For systems with stable data

The script uses `INSERT ... ON CONFLICT ... DO UPDATE` to ensure idempotent execution.

## Data Quality

The accuracy of the algorithm depends on the integrity of the following data:
- Correct `sport_id` associations in `tracked_events`
- Valid `category_id` references
- Up-to-date market data in the `markets` table

## Version History

- **v1.0** (2025-11-27): Initial release
  - Two-dimensional scoring model (Tier + Market Depth)
  - SQL implementation with automatic upsert
  - Comprehensive documentation

## Author

**Manus AI**

For questions or suggestions, please open an issue in the repository.
