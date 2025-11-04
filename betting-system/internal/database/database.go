package database

import (
	"fmt"
	"log"

	"github.com/gdsZyy/betting-system/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// Config 数据库配置
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect 连接数据库
func Connect(cfg *Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established")

	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Market{},
		&models.Outcome{},
		&models.Bet{},
		&models.BetSelection{},
		&models.BetLeg{},
		&models.BetLegSelection{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")

	return nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateIndexes 创建索引
func CreateIndexes() error {
	// 用户表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return err
	}

	// 赛事表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_external_id ON events(external_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_status ON events(status)").Error; err != nil {
		return err
	}

	// 盘口表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_markets_event_id ON markets(event_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_markets_status ON markets(status)").Error; err != nil {
		return err
	}

	// 结果选项表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_outcomes_market_id ON outcomes(market_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_outcomes_status ON outcomes(status)").Error; err != nil {
		return err
	}

	// 投注表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bets_user_id ON bets(user_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bets_status ON bets(status)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bets_placed_at ON bets(placed_at)").Error; err != nil {
		return err
	}

	// 投注选项表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bet_selections_bet_id ON bet_selections(bet_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bet_selections_outcome_id ON bet_selections(outcome_id)").Error; err != nil {
		return err
	}

	// 投注组合腿表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bet_legs_bet_id ON bet_legs(bet_id)").Error; err != nil {
		return err
	}

	// 组合腿选项关联表索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bet_leg_selections_leg_id ON bet_leg_selections(leg_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_bet_leg_selections_selection_id ON bet_leg_selections(selection_id)").Error; err != nil {
		return err
	}

	log.Println("Database indexes created")

	return nil
}

