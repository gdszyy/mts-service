package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"betting-system/internal/models"
)

// TournamentHandler 联赛处理器
type TournamentHandler struct {
	db *gorm.DB
}

// NewTournamentHandler 创建联赛处理器
func NewTournamentHandler(db *gorm.DB) *TournamentHandler {
	return &TournamentHandler{db: db}
}

// TournamentResponse 联赛响应结构
type TournamentResponse struct {
	ID           int64      `json:"id"`
	ExternalID   string     `json:"external_id"`
	CategoryID   int64      `json:"category_id"`
	CategoryName string     `json:"category_name"`
	Name         string     `json:"name"`
	Scheduled    *time.Time `json:"scheduled,omitempty"`
	ScheduledEnd *time.Time `json:"scheduled_end,omitempty"`
	MatchCount   int64      `json:"match_count"`
}

// GetTournaments 获取联赛列表
// @Summary 获取联赛列表
// @Description 根据分类获取联赛列表,支持分页和排序
// @Tags tournaments
// @Accept json
// @Produce json
// @Param category_id query int true "分类ID"
// @Param page query int false "页码,默认1"
// @Param page_size query int false "每页大小,默认20,最大100"
// @Param sort query string false "排序方式: name_asc(默认), name_desc, match_count_desc"
// @Success 200 {object} map[string]interface{}
// @Router /api/tournaments [get]
func (h *TournamentHandler) GetTournaments(c *gin.Context) {
	// 获取查询参数
	categoryIDStr := c.Query("category_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	sortBy := c.DefaultQuery("sort", "name_asc")

	// 验证必填参数
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "category_id is required",
		})
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category_id",
		})
		return
	}

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 验证分类是否存在
	var category models.Category
	if err := h.db.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch category",
		})
		return
	}

	// 构建基础查询
	query := h.db.Model(&models.Tournament{}).Where("category_id = ?", categoryID)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count tournaments",
		})
		return
	}

	// 构建查询,包含比赛数量统计
	var tournaments []struct {
		models.Tournament
		CategoryName string `json:"category_name"`
		MatchCount   int64  `json:"match_count"`
	}

	// 子查询统计每个联赛下的比赛数量
	subQuery := h.db.Model(&models.Event{}).
		Select("tournament_id, COUNT(*) as match_count").
		Where("tournament_id IS NOT NULL").
		Group("tournament_id")

	// 主查询
	mainQuery := h.db.Table("tournaments").
		Select("tournaments.*, categories.name as category_name, COALESCE(match_counts.match_count, 0) as match_count").
		Joins("LEFT JOIN categories ON tournaments.category_id = categories.id").
		Joins("LEFT JOIN (?) as match_counts ON tournaments.id = match_counts.tournament_id", subQuery).
		Where("tournaments.category_id = ?", categoryID)

	// 应用排序
	switch sortBy {
	case "name_desc":
		mainQuery = mainQuery.Order("tournaments.name DESC")
	case "match_count_desc":
		mainQuery = mainQuery.Order("match_count DESC, tournaments.name ASC")
	default: // name_asc
		mainQuery = mainQuery.Order("tournaments.name ASC")
	}

	// 应用分页
	if err := mainQuery.Offset(offset).Limit(pageSize).Scan(&tournaments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tournaments",
		})
		return
	}

	// 构建响应
	var responseData []TournamentResponse
	for _, tour := range tournaments {
		responseData = append(responseData, TournamentResponse{
			ID:           tour.ID,
			ExternalID:   tour.ExternalID,
			CategoryID:   tour.CategoryID,
			CategoryName: tour.CategoryName,
			Name:         tour.Name,
			Scheduled:    tour.Scheduled,
			ScheduledEnd: tour.ScheduledEnd,
			MatchCount:   tour.MatchCount,
		})
	}

	// 计算总页数
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseData,
		"pagination": PaginationResponse{
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	})
}
