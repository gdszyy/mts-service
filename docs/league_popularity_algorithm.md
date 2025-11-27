# League Popularity Scoring Algorithm

**Version**: 1.0
**Author**: Manus AI
**Date**: 2025-11-27

## 1. Overview

This document outlines the algorithm for calculating a popularity score for sports leagues and tournaments. The goal is to quantify the social popularity and market significance of a league, distinguishing major international events like the FIFA World Cup from smaller, regional competitions.

The algorithm is designed to work with the data structure of the `betradar-uof-service` project, leveraging existing data from the Sportradar Unified Odds Feed (UOF).

## 2. Scoring Model

The popularity score is calculated using a two-dimensional weighted model that combines a league's **inherent prestige** (Tier Score) with its **market activity** (Market Depth Score).

### Final Score Formula

The final popularity score is a value between 1 and 10, calculated as follows:

```
Final Score = (League Tier Score * 0.5) + (Market Depth Score * 0.5)
```

- **Weighting**: Both dimensions are given equal weight (50%) to balance a league's reputation with its actual market engagement.

## 3. Scoring Dimensions

### 3.1. League Tier Score (1-10 points)

This score represents the inherent prestige and global standing of a league. It is determined by classifying leagues into tiers based on their name and scope, using keyword matching on the `categories.name` and `categories.country_code` fields.

| Tier | Score | SQL `ILIKE` Rule | Examples |
| :--- | :--- | :--- | :--- |
| T1 | 10 | `name` LIKE '%World Cup%' OR `name` LIKE '%Olympics%' OR `name` LIKE '%World Championship%' | FIFA World Cup, Olympics |
| T2 | 9 | `name` LIKE '%Champions League%' OR `name` LIKE '%Europa League%' OR `name` LIKE '%Copa America%' | UEFA Champions League |
| T3 | 8 | `name` LIKE '%Premier League%' OR `name` LIKE '%La Liga%' OR `name` LIKE '%Serie A%' OR `name` LIKE '%Bundesliga%' OR `name` LIKE '%NBA%' | Premier League, NBA |
| T4 | 7 | `name` LIKE '%International Clubs%' OR `name` LIKE '%International%' | International Clubs (Soccer) |
| T5 | 6 | `country_code` IS NOT NULL AND `country_code` != '' | Any national league |
| T6 | 4 | ELSE (All other cases) | eSoccer, regional leagues |

### 3.2. Market Depth Score (1-10 points)

This score reflects the market's interest and engagement with a league, measured by the average number of betting markets available per event.

This is calculated by:
1. Counting the number of markets for each `event_id` in the `markets` table.
2. Averaging these counts across all events within a specific `category_id`.

The score is then assigned based on the following thresholds:

| Score | Average Markets (`avg_markets`) | Market Level |
| :--- | :--- | :--- |
| 10 | `avg_markets` > 300 | Extremely High |
| 9 | 200 < `avg_markets` <= 300 | Very High |
| 8 | 150 < `avg_markets` <= 200 | High |
| 7 | 100 < `avg_markets` <= 150 | Moderately High |
| 6 | 50 < `avg_markets` <= 100 | Above Average |
| 5 | 30 < `avg_markets` <= 50 | Average |
| 4 | 20 < `avg_markets` <= 30 | Below Average |
| 3 | 10 < `avg_markets` <= 20 | Low |
| 2 | 5 < `avg_markets` <= 10 | Very Low |
| 1 | `avg_markets` <= 5 | Minimal |

## 4. SQL Implementation

The algorithm is implemented through a single, comprehensive SQL script that performs the following steps:

1.  **Creates a `league_popularity_scores` table** to store the results.
2.  **Calculates statistics**: Computes `avg_market_count` for each league.
3.  **Assigns scores**: Determines `league_tier_score` and `market_depth_score` based on the rules above.
4.  **Calculates final score**: Applies the weighted formula.
5.  **Upserts data**: Inserts or updates the `league_popularity_scores` table with the latest scores.

This script is designed to be run periodically (e.g., daily) to keep the popularity scores up-to-date.

For the complete SQL script, see `scripts/calculate_league_popularity.sql`.

## 5. Data Dependencies

The algorithm relies on the following tables and fields:

- **`categories`**: `id`, `name`, `sport_id`, `country_code`
- **`sports`**: `id`, `name`
- **`tracked_events`**: `event_id`, `category_id`, `sport_id`
- **`markets`**: `id`, `event_id`

**Data Integrity Note**: The accuracy of the algorithm is highly dependent on the correctness of the `sport_id` and `category_id` associations in the `tracked_events` table. Recent data quality improvements have been crucial for the algorithm's effectiveness.
