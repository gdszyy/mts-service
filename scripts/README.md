# Scripts

This directory contains SQL scripts and Python utilities for the popularity scoring system.

## SQL Scripts

### Event Popularity Scoring
- **`calculate_event_popularity.sql`**: Calculates and stores popularity scores for individual events based on market count.

### Tournament Popularity Scoring
- **`calculate_tournament_popularity.sql`**: Calculates and stores popularity scores for tournaments based on tier and market depth.

### Legacy
- **`calculate_league_popularity.sql`**: ⚠️ Deprecated - Use `calculate_tournament_popularity.sql` instead.

## Python Utilities

### Correlation Analysis
- **`analyze_event_tournament_correlation.py`**: Analyzes the correlation between event market counts and tournament tiers to validate the single-dimension model design.

## Usage

### Execute SQL Scripts

```bash
# Set database URL
export DATABASE_URL="postgresql://user:password@host:port/database"

# Execute event popularity scoring
psql $DATABASE_URL -f calculate_event_popularity.sql

# Execute tournament popularity scoring
psql $DATABASE_URL -f calculate_tournament_popularity.sql
```

### Run Correlation Analysis

```bash
# Install dependencies (if not already installed)
pip3 install psycopg2-binary matplotlib numpy

# Run the analysis
python3 analyze_event_tournament_correlation.py
```

The analysis will:
1. Calculate correlation coefficient between tournament scores and event market counts
2. Analyze market count distribution across tournament tiers
3. Calculate coefficient of variation within tournaments
4. Generate visualization charts
5. Provide data-driven recommendations

### Output

The analysis script generates:
- Console output with detailed statistics
- Visualization chart: `event_tournament_correlation_analysis.png`

## Scheduling

For production use, it's recommended to run the scoring scripts periodically:

```bash
# Example cron job (daily at 2 AM)
0 2 * * * psql $DATABASE_URL -f /path/to/calculate_event_popularity.sql
0 2 * * * psql $DATABASE_URL -f /path/to/calculate_tournament_popularity.sql
```

## Dependencies

- PostgreSQL client (`psql`)
- Python 3.7+
- Python packages: `psycopg2-binary`, `matplotlib`, `numpy`
