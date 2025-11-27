# Event (Match) Popularity Scoring Algorithm

**Version**: 1.0
**Author**: Manus AI
**Date**: 2025-11-27

## 1. Overview

This document outlines the algorithm for calculating a popularity score for individual sports matches (events). The goal is to quantify the market interest and significance of a single event, distinguishing a World Cup final from a regular-season regional match.

This algorithm is designed to be independent of the league-level scoring system and focuses solely on event-specific metrics available in the `betradar-uof-service` database.

## 2. Scoring Model

The popularity of a single match is most directly reflected by the number of betting markets offered for it. A higher number of markets indicates that the operator (and the market in general) expects higher engagement and betting volume.

Therefore, the model is a **single-dimension model** based entirely on market count.

### Final Score Formula

The final popularity score is a value between 1 and 10, calculated by normalizing the market count for a given event.

```
Final Score = Market Count Score
```

## 3. Scoring Dimension

### 3.1. Market Count Score (1-10 points)

This score is derived by counting the number of records in the `markets` table associated with a specific `event_id`.

The score is assigned based on a threshold system designed to handle the wide distribution of market counts, where a few major events have a very high number of markets.

| Score | Market Count (`market_count`) | Market Level |
| :--- | :--- | :--- |
| 10 | `market_count` > 400 | Global Event (e.g., World Cup Final) |
| 9 | 300 < `market_count` <= 400 | Major International Event |
| 8 | 200 < `market_count` <= 300 | Top-Tier National Event |
| 7 | 150 < `market_count` <= 200 | High-Profile Match |
| 6 | 100 < `market_count` <= 150 | Standard High-Level Match |
| 5 | 50 < `market_count` <= 100 | Average Match |
| 4 | 30 < `market_count` <= 50 | Below Average Match |
| 3 | 10 < `market_count` <= 30 | Low-Interest Match |
| 2 | 5 < `market_count` <= 10 | Very Low-Interest Match |
| 1 | `market_count` <= 5 | Minimal Interest |

## 4. SQL Implementation

The algorithm is implemented through a SQL script that performs the following steps:

1.  **Creates a `event_popularity_scores` table** to store the results.
2.  **Calculates market count**: For each `event_id` in `tracked_events`, it counts the associated markets.
3.  **Assigns score**: Determines the `market_count_score` based on the rules above.
4.  **Upserts data**: Inserts or updates the `event_popularity_scores` table with the latest scores.

This script can be run periodically to score new and ongoing events.

For the complete SQL script, see `scripts/calculate_event_popularity.sql`.

## 5. Data Dependencies

The algorithm relies on the following tables and fields:

- **`tracked_events`**: `event_id`
- **`markets`**: `id`, `event_id`

This model is robust as it has minimal dependencies and uses the most reliable indicator of event-specific market interest.
