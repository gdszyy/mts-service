# Popularity Scoring System

This directory contains the documentation and implementation for **two independent popularity scoring algorithms**:

1. **Event (Match) Popularity Scoring**: Quantifies the market interest in individual matches
2. **Tournament (League) Popularity Scoring**: Quantifies the social popularity and market significance of sports leagues

## Overview

### Event Popularity

The event algorithm assigns a popularity score (1-10) to each match based on:
- **Market Count**: The number of betting markets offered for the event

### Tournament Popularity

The tournament algorithm assigns a popularity score (1-10) to each league based on two key dimensions:
1. **Tournament Tier Score**: The inherent prestige and global standing of the league
2. **Market Depth Score**: The market's engagement, measured by the average number of betting markets per event

## Documentation

- **[event_popularity_algorithm.md](./event_popularity_algorithm.md)**: Event (match) popularity scoring specification
- **[tournament_popularity_algorithm.md](./tournament_popularity_algorithm.md)**: Tournament (league) popularity scoring specification
- **[league_popularity_algorithm.md](./league_popularity_algorithm.md)**: ⚠️ Deprecated - Use tournament_popularity_algorithm.md instead
- **[popularity_analysis_report.md](./popularity_analysis_report.md)**: Comprehensive analysis report with execution results

## Implementation

Both algorithms are implemented as SQL scripts that can be executed periodically:

### Event Popularity
- **SQL Script**: `../scripts/calculate_event_popularity.sql`
- **Database Table**: `event_popularity_scores`
- **Records**: 16,707 events scored

### Tournament Popularity
- **SQL Script**: `../scripts/calculate_tournament_popularity.sql`
- **Database Table**: `tournament_popularity_scores`
- **Records**: 548 tournaments scored

## Quick Start

### 1. Execute Both Algorithms

Run both SQL scripts to calculate and store popularity scores:

```bash
# Execute event popularity scoring
psql $DATABASE_URL -f scripts/calculate_event_popularity.sql

# Execute tournament popularity scoring
psql $DATABASE_URL -f scripts/calculate_tournament_popularity.sql
```

Or use the Python helper script to run both:

```bash
python3 scripts/execute_both_algorithms.py
```

### 2. Query Event Results

View the top events by popularity:

```sql
SELECT 
    event_id,
    market_count,
    popularity_score
FROM event_popularity_scores
ORDER BY popularity_score DESC, market_count DESC
LIMIT 20;
```

### 3. Query Tournament Results

View the top tournaments by popularity:

```sql
SELECT 
    tournament_name,
    category_name,
    sport_name,
    total_events,
    avg_market_count,
    final_popularity_score
FROM tournament_popularity_scores
ORDER BY final_popularity_score DESC
LIMIT 20;
```

## Score Tiers

### Event Score Tiers

| Score | Market Count | Description |
| :--- | :--- | :--- |
| 10 | >400 | Global Event (e.g., World Cup Final) |
| 9 | 300-400 | Major International Event |
| 8 | 200-300 | Top-Tier National Event |
| 7 | 150-200 | High-Profile Match |
| 6 | 100-150 | Standard High-Level Match |
| 1-5 | <100 | Below Average to Minimal Interest |

### Tournament Score Tiers

| Tier | Score Range | Description |
| :--- | :--- | :--- |
| S-Tier | 9.0 - 10.0 | Super Popular (e.g., Champions League, NBA) |
| A-Tier | 8.0 - 8.9 | Very Popular (e.g., Top national leagues) |
| B-Tier | 7.0 - 7.9 | Popular (e.g., Major national leagues) |
| C-Tier | 6.0 - 6.9 | Average (e.g., Standard national leagues) |
| D-Tier | < 6.0 | Below Average (e.g., Regional leagues) |

## Maintenance

It is recommended to run the scoring script:
- **Daily**: For production systems with high data turnover
- **Weekly**: For systems with stable data

The script uses `INSERT ... ON CONFLICT ... DO UPDATE` to ensure idempotent execution.

## Data Quality

The accuracy of both algorithms depends on the integrity of the following data:
- Correct `sport_id` associations in `tracked_events`
- Valid `category_id` and `tournament_id` references
- Up-to-date market data in the `markets` table
- Accurate tournament and category metadata

## Version History

- **v1.0** (2025-11-27): Initial release
  - Event popularity scoring (single-dimension: market count)
  - Tournament popularity scoring (two-dimension: tier + market depth)
  - SQL implementation with automatic upsert
  - Comprehensive documentation and analysis report
  - Successfully scored 16,707 events and 548 tournaments

## Author

**Manus AI**

For questions or suggestions, please open an issue in the repository.
