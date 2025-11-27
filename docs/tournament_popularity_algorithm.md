# Tournament (League) Popularity Scoring Algorithm

**Version**: 1.1
**Author**: Manus AI
**Date**: 2025-11-27

## 1. Overview

This document outlines the algorithm for calculating a popularity score for sports tournaments (leagues/seasons). The goal is to quantify the social popularity and market significance of a tournament, distinguishing major international events like the FIFA World Cup from smaller, regional competitions.

This algorithm is based on the `tournaments` table and its relationship with events and markets.

## 2. Scoring Model

The popularity score is calculated using a two-dimensional weighted model that combines a tournament's **inherent prestige** (Tier Score) with its **market activity** (Market Depth Score).

### Final Score Formula

The final popularity score is a value between 1 and 10, calculated as follows:

```
Final Score = (Tournament Tier Score * 0.5) + (Market Depth Score * 0.5)
```

- **Weighting**: Both dimensions are given equal weight (50%) to balance a tournament's reputation with its actual market engagement.

## 3. Scoring Dimensions

### 3.1. Tournament Tier Score (1-10 points)

This score represents the inherent prestige and global standing of a tournament. It is determined by classifying tournaments based on their `name` and the `name` of their associated `category`.

| Tier | Score | SQL `ILIKE` Rule | Examples |
| :--- | :--- | :--- | :--- |
| T1 | 10 | `t.name` or `c.name` LIKE 
%World Cup%
, 
%Olympics%
, 
%World Championship%
 | FIFA World Cup, Olympics |
| T2 | 9 | `t.name` or `c.name` LIKE 
%Champions League%
, 
%Europa League%
, 
%Copa America%
 | UEFA Champions League |
| T3 | 8 | `t.name` or `c.name` LIKE 
%Premier League%
, 
%La Liga%
, 
%Serie A%
, 
%Bundesliga%
, 
%NBA%
 | Premier League, NBA |
| T4 | 7 | `c.name` LIKE 
%International Clubs%
, 
%International%
 | International Clubs (Soccer) |
| T5 | 6 | `c.country_code` IS NOT NULL AND `c.country_code` != 

 | Any national league |
| T6 | 4 | ELSE (All other cases) | eSoccer, regional leagues |

### 3.2. Market Depth Score (1-10 points)

This score reflects the market's interest and engagement with a tournament, measured by the **average number of betting markets available per event** within that tournament.

This is calculated by:
1. Counting the number of markets for each `event_id`.
2. Averaging these counts across all events associated with a specific `tournament_id`.

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

The algorithm is implemented through a SQL script that performs the following steps:

1.  **Creates a `tournament_popularity_scores` table** to store the results.
2.  **Calculates statistics**: Computes `avg_market_count` for each tournament.
3.  **Assigns scores**: Determines `tournament_tier_score` and `market_depth_score`.
4.  **Calculates final score**: Applies the weighted formula.
5.  **Upserts data**: Inserts or updates the `tournament_popularity_scores` table.

For the complete SQL script, see `scripts/calculate_tournament_popularity.sql`.

## 5. Data Dependencies

The algorithm relies on the following tables and fields:

- **`tournaments`**: `id`, `name`, `category_id`
- **`categories`**: `id`, `name`, `country_code`
- **`tracked_events`**: `event_id`, `tournament_id`
- **`markets`**: `id`, `event_id`
