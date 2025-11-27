#!/usr/bin/env python3
"""
Event-Tournament Correlation Analysis

This script analyzes the correlation between event market counts and tournament tiers
to validate the design decision of using a single-dimension model for event popularity scoring.

Author: Manus AI
Date: 2025-11-27
"""

import psycopg2
import matplotlib.pyplot as plt
import matplotlib
import numpy as np
import os

# Configure Chinese font for matplotlib
matplotlib.rcParams['font.sans-serif'] = ['Noto Sans CJK SC', 'DejaVu Sans']
matplotlib.rcParams['axes.unicode_minus'] = False

# Database connection URL - should be set via environment variable in production
DB_URL = os.getenv('DATABASE_URL', 'postgresql://postgres:qcriEvdpsnxvfPLaGuCuTqtivHpKoodg@turntable.proxy.rlwy.net:48608/railway')

def analyze_correlation():
    """
    Main analysis function that:
    1. Calculates correlation between tournament scores and event market counts
    2. Analyzes market count distribution across tournament tiers
    3. Calculates coefficient of variation within tournaments
    4. Generates visualization charts
    5. Provides recommendations
    """
    conn = psycopg2.connect(DB_URL)
    cur = conn.cursor()
    
    print("\n" + "="*100)
    print("Event Market Count vs Tournament Tier Correlation Analysis")
    print("="*100)
    
    # Query 1: Get event market counts with tournament information
    query = """
    SELECT 
        eps.event_id,
        eps.market_count,
        eps.popularity_score as event_score,
        tps.tournament_name,
        tps.final_popularity_score as tournament_score,
        tps.tournament_tier_score
    FROM event_popularity_scores eps
    INNER JOIN tournament_popularity_scores tps 
        ON eps.tournament_id = tps.tournament_id
    WHERE eps.market_count > 0
    ORDER BY eps.market_count DESC;
    """
    
    cur.execute(query)
    rows = cur.fetchall()
    
    print(f"\n📊 Successfully matched events: {len(rows)}")
    
    if not rows:
        print("❌ No data available for analysis")
        cur.close()
        conn.close()
        return
    
    # Extract data for analysis
    market_counts = [row[1] for row in rows]
    event_scores = [row[2] for row in rows]
    tournament_scores = [float(row[4]) for row in rows]
    tournament_tiers = [row[5] for row in rows]
    
    # Query 2: Market count statistics by tournament tier
    print("\n" + "="*100)
    print("Market Count Statistics by Tournament Tier")
    print("="*100)
    
    tier_query = """
    SELECT 
        tps.tournament_tier_score,
        COUNT(eps.event_id) as event_count,
        AVG(eps.market_count) as avg_market_count,
        MIN(eps.market_count) as min_market_count,
        MAX(eps.market_count) as max_market_count,
        STDDEV(eps.market_count) as stddev_market_count
    FROM event_popularity_scores eps
    INNER JOIN tournament_popularity_scores tps 
        ON eps.tournament_id = tps.tournament_id
    WHERE eps.market_count > 0
    GROUP BY tps.tournament_tier_score
    ORDER BY tps.tournament_tier_score DESC;
    """
    
    cur.execute(tier_query)
    tier_stats = cur.fetchall()
    
    print(f"\n{'Tier Score':<12} {'Events':<10} {'Avg Markets':<15} {'Min':<10} {'Max':<10} {'Std Dev':<15}")
    print("-" * 80)
    
    for row in tier_stats:
        tier, count, avg, min_val, max_val, stddev = row
        print(f"{tier:<12} {count:<10} {float(avg):<15.2f} {min_val:<10} {max_val:<10} {float(stddev) if stddev else 0:<15.2f}")
    
    # Query 3: Coefficient of variation within tournaments
    print("\n" + "="*100)
    print("Coefficient of Variation (CV) Within Tournaments")
    print("="*100)
    
    cv_query = """
    SELECT 
        tps.tournament_name,
        tps.final_popularity_score,
        COUNT(eps.event_id) as event_count,
        AVG(eps.market_count) as avg_market_count,
        STDDEV(eps.market_count) as stddev_market_count,
        CASE 
            WHEN AVG(eps.market_count) > 0 THEN STDDEV(eps.market_count) / AVG(eps.market_count)
            ELSE 0
        END as coefficient_of_variation
    FROM event_popularity_scores eps
    INNER JOIN tournament_popularity_scores tps 
        ON eps.tournament_id = tps.tournament_id
    WHERE eps.market_count > 0
    GROUP BY tps.tournament_id, tps.tournament_name, tps.final_popularity_score
    HAVING COUNT(eps.event_id) >= 3
    ORDER BY coefficient_of_variation DESC
    LIMIT 20;
    """
    
    cur.execute(cv_query)
    cv_stats = cur.fetchall()
    
    print(f"\n{'Tournament Name':<40} {'Score':<10} {'Events':<10} {'Avg Markets':<12} {'CV':<10}")
    print("-" * 90)
    
    for row in cv_stats:
        name, score, count, avg, stddev, cv = row
        print(f"{str(name)[:40]:<40} {float(score):<10.2f} {count:<10} {float(avg):<12.2f} {float(cv) if cv else 0:<10.2f}")
    
    # Generate visualizations
    print("\n📈 Generating visualization charts...")
    
    fig, axes = plt.subplots(2, 2, figsize=(16, 12))
    
    # Chart 1: Tournament tier vs average market count (bar chart)
    tier_scores = [row[0] for row in tier_stats]
    tier_avg_markets = [float(row[2]) for row in tier_stats]
    
    axes[0, 0].bar(tier_scores, tier_avg_markets, color='steelblue', alpha=0.7)
    axes[0, 0].set_xlabel('Tournament Tier Score', fontsize=12)
    axes[0, 0].set_ylabel('Average Market Count', fontsize=12)
    axes[0, 0].set_title('Tournament Tier vs Event Average Market Count', fontsize=14, fontweight='bold')
    axes[0, 0].grid(axis='y', alpha=0.3)
    
    # Chart 2: Scatter plot - tournament score vs event market count
    axes[0, 1].scatter(tournament_scores, market_counts, alpha=0.3, s=20)
    axes[0, 1].set_xlabel('Tournament Popularity Score', fontsize=12)
    axes[0, 1].set_ylabel('Event Market Count', fontsize=12)
    axes[0, 1].set_title('Tournament Score vs Event Market Count (Scatter)', fontsize=14, fontweight='bold')
    axes[0, 1].grid(alpha=0.3)
    
    # Calculate and display correlation coefficient
    correlation = np.corrcoef(tournament_scores, market_counts)[0, 1]
    axes[0, 1].text(0.05, 0.95, f'Correlation: {correlation:.3f}', 
                    transform=axes[0, 1].transAxes, 
                    bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5),
                    verticalalignment='top')
    
    # Chart 3: Box plot - market count distribution by tournament tier
    tier_data = {}
    for mc, tier in zip(market_counts, tournament_tiers):
        if tier not in tier_data:
            tier_data[tier] = []
        tier_data[tier].append(mc)
    
    sorted_tiers = sorted(tier_data.keys(), reverse=True)
    box_data = [tier_data[tier] for tier in sorted_tiers]
    
    axes[1, 0].boxplot(box_data, tick_labels=sorted_tiers)
    axes[1, 0].set_xlabel('Tournament Tier Score', fontsize=12)
    axes[1, 0].set_ylabel('Event Market Count', fontsize=12)
    axes[1, 0].set_title('Market Count Distribution by Tournament Tier (Box Plot)', fontsize=14, fontweight='bold')
    axes[1, 0].grid(axis='y', alpha=0.3)
    
    # Chart 4: Score distribution comparison (current vs simulated with tier)
    simulated_scores = []
    for mc, tier in zip(market_counts, tournament_tiers):
        # Current score (market count only)
        if mc > 400: current = 10
        elif mc > 300: current = 9
        elif mc > 200: current = 8
        elif mc > 150: current = 7
        elif mc > 100: current = 6
        elif mc > 50: current = 5
        elif mc > 30: current = 4
        elif mc > 10: current = 3
        elif mc > 5: current = 2
        else: current = 1
        
        # Simulated score (50% market + 50% tier)
        market_score = current
        simulated = (market_score + tier) / 2
        simulated_scores.append(simulated)
    
    axes[1, 1].hist([event_scores, simulated_scores], bins=20, 
                    label=['Current (Market Only)', 'Simulated (Market + Tier)'], 
                    alpha=0.6)
    axes[1, 1].set_xlabel('Event Score', fontsize=12)
    axes[1, 1].set_ylabel('Event Count', fontsize=12)
    axes[1, 1].set_title('Event Score Distribution Comparison', fontsize=14, fontweight='bold')
    axes[1, 1].legend()
    axes[1, 1].grid(axis='y', alpha=0.3)
    
    plt.tight_layout()
    output_path = '/home/ubuntu/event_tournament_correlation_analysis.png'
    plt.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"✅ Chart saved: {output_path}")
    
    # Calculate key metrics
    print("\n" + "="*100)
    print("Key Metrics Summary")
    print("="*100)
    
    print(f"\n1. Correlation coefficient (tournament score vs event market count): {correlation:.3f}")
    
    if correlation > 0.7:
        print("   ➜ Strong positive correlation: Tournament tier significantly affects event market count")
    elif correlation > 0.4:
        print("   ➜ Moderate positive correlation: Tournament tier has some effect on event market count")
    else:
        print("   ➜ Weak correlation: Tournament tier has minimal effect on event market count")
    
    # Calculate average coefficient of variation
    avg_cv = np.mean([float(row[5]) if row[5] else 0 for row in cv_stats])
    print(f"\n2. Average coefficient of variation within tournaments: {avg_cv:.3f}")
    
    if avg_cv > 0.5:
        print("   ➜ High variation: Market counts vary significantly within the same tournament")
    elif avg_cv > 0.3:
        print("   ➜ Moderate variation: Market counts have some variation within the same tournament")
    else:
        print("   ➜ Low variation: Market counts are relatively stable within the same tournament")
    
    # Provide recommendation
    print("\n" + "="*100)
    print("Conclusion and Recommendation")
    print("="*100)
    
    if correlation > 0.6 and avg_cv < 0.4:
        print("\n✅ RECOMMENDATION: No need to add tournament tier dimension to event scoring")
        print("   Reasons:")
        print("   1. Strong correlation means market count already captures tournament tier information")
        print("   2. Relatively stable market counts within tournaments mean tier won't add discrimination")
        print("   3. Current single-dimension model is simple, efficient, and follows Occam's Razor principle")
    elif correlation < 0.4 or avg_cv > 0.5:
        print("\n⚠️ RECOMMENDATION: Consider adding tournament tier dimension to event scoring")
        print("   Reasons:")
        print("   1. Weak correlation means market count doesn't fully reflect tournament tier")
        print("   2. High variation within tournaments means tier dimension could help balance scores")
        print("   3. Two-dimension model could provide more accurate event popularity assessment")
    else:
        print("\n🤔 RECOMMENDATION: Optional to add tournament tier dimension")
        print("   Reasons:")
        print("   1. Correlation and variation are at moderate levels")
        print("   2. Decision depends on specific business requirements")
        print("   3. Consider A/B testing to validate effectiveness")
    
    cur.close()
    conn.close()
    
    print("\n" + "="*100)
    print("Analysis Complete")
    print("="*100)

if __name__ == "__main__":
    try:
        analyze_correlation()
    except Exception as e:
        print(f"\n❌ Error during analysis: {e}")
        import traceback
        traceback.print_exc()
