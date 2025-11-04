package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Balance   float64        `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`
	Currency  string         `gorm:"size:3;not null;default:'CNY'" json:"currency"`
	Status    string         `gorm:"size:20;not null;default:'active'" json:"status"` // active/suspended/closed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Event 赛事模型
type Event struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID string         `gorm:"uniqueIndex;size:100" json:"external_id"`
	SportID    string         `gorm:"size:50;not null" json:"sport_id"`
	HomeTeam   string         `gorm:"size:255;not null" json:"home_team"`
	AwayTeam   string         `gorm:"size:255;not null" json:"away_team"`
	StartTime  time.Time      `gorm:"not null" json:"start_time"`
	Status     string         `gorm:"size:20;not null;default:'scheduled'" json:"status"` // scheduled/live/finished/cancelled
	HomeScore  *int           `json:"home_score,omitempty"`
	AwayScore  *int           `json:"away_score,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Market 盘口模型
type Market struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	EventID    int64          `gorm:"index;not null" json:"event_id"`
	Event      Event          `gorm:"foreignKey:EventID" json:"event,omitempty"`
	MarketType string         `gorm:"size:50;not null" json:"market_type"` // 1x2/handicap/totals等
	Specifier  string         `gorm:"size:100" json:"specifier,omitempty"` // 盘口参数
	Status     string         `gorm:"size:20;not null;default:'active'" json:"status"` // active/suspended/settled/cancelled
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Outcome 结果选项模型
type Outcome struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	MarketID  int64          `gorm:"index;not null" json:"market_id"`
	Market    Market         `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	OutcomeID string         `gorm:"size:50;not null" json:"outcome_id"` // 1/x/2, over/under等
	Odds      float64        `gorm:"type:decimal(10,4);not null" json:"odds"`
	Status    string         `gorm:"size:20;not null;default:'active'" json:"status"` // active/suspended/won/lost/void
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Bet 投注模型
type Bet struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64          `gorm:"index;not null" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BetType         string         `gorm:"size:50;not null" json:"bet_type"` // single/accumulator/trixie等
	TotalStake      float64        `gorm:"type:decimal(15,2);not null" json:"total_stake"`
	PotentialReturn float64        `gorm:"type:decimal(15,2);not null" json:"potential_return"`
	ActualReturn    *float64       `gorm:"type:decimal(15,2)" json:"actual_return,omitempty"`
	Status          string         `gorm:"size:20;not null;default:'pending'" json:"status"` // pending/won/lost/void/partially_won
	PlacedAt        time.Time      `gorm:"not null" json:"placed_at"`
	SettledAt       *time.Time     `json:"settled_at,omitempty"`
	Selections      []BetSelection `gorm:"foreignKey:BetID" json:"selections,omitempty"`
	Legs            []BetLeg       `gorm:"foreignKey:BetID" json:"legs,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BetSelection 投注选项模型
type BetSelection struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BetID     int64          `gorm:"index;not null" json:"bet_id"`
	OutcomeID int64          `gorm:"index;not null" json:"outcome_id"`
	Outcome   Outcome        `gorm:"foreignKey:OutcomeID" json:"outcome,omitempty"`
	Odds      float64        `gorm:"type:decimal(10,4);not null" json:"odds"` // 下注时的赔率
	IsBanker  bool           `gorm:"not null;default:false" json:"is_banker"`
	Status    string         `gorm:"size:20;not null;default:'pending'" json:"status"` // pending/won/lost/void
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BetLeg 投注组合腿模型
type BetLeg struct {
	ID              int64              `gorm:"primaryKey;autoIncrement" json:"id"`
	BetID           int64              `gorm:"index;not null" json:"bet_id"`
	LegType         string             `gorm:"size:50;not null" json:"leg_type"` // single/double/treble/4-fold等
	Stake           float64            `gorm:"type:decimal(15,2);not null" json:"stake"`
	Odds            float64            `gorm:"type:decimal(10,4);not null" json:"odds"`
	PotentialReturn float64            `gorm:"type:decimal(15,2);not null" json:"potential_return"`
	ActualReturn    *float64           `gorm:"type:decimal(15,2)" json:"actual_return,omitempty"`
	Status          string             `gorm:"size:20;not null;default:'pending'" json:"status"` // pending/won/lost/void
	Selections      []BetLegSelection  `gorm:"foreignKey:LegID" json:"selections,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	DeletedAt       gorm.DeletedAt     `gorm:"index" json:"-"`
}

// BetLegSelection 组合腿选项关联模型
type BetLegSelection struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	LegID       int64          `gorm:"index;not null" json:"leg_id"`
	SelectionID int64          `gorm:"index;not null" json:"selection_id"`
	Selection   BetSelection   `gorm:"foreignKey:SelectionID" json:"selection,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Event) TableName() string {
	return "events"
}

func (Market) TableName() string {
	return "markets"
}

func (Outcome) TableName() string {
	return "outcomes"
}

func (Bet) TableName() string {
	return "bets"
}

func (BetSelection) TableName() string {
	return "bet_selections"
}

func (BetLeg) TableName() string {
	return "bet_legs"
}

func (BetLegSelection) TableName() string {
	return "bet_leg_selections"
}

