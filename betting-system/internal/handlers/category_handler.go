package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"betting-system/internal/models"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// CategoryResponse 分类响应结构
type CategoryResponse struct {
	ID         int64  `json:"id"`
	ExternalID string `json:"external_id"`
	SportID    int64  `json:"sport_id"`
	SportName  string `json:"sport_name"`
	Name       string `json:"name"`
	MatchCount int64  `json:"match_count"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Description 根据体育类型获取分类列表,支持分页和排序
// @Tags categories
// @Accept json
// @Produce json
// @Param sport_ids query string false "体育类型ID列表,逗号分隔,如: sr:sport:1,sr:sport:2"
// @Param page query int false "页码,默认1"
// @Param page_size query int false "每页大小,默认20,最大100"
// @Param sort query string false "排序方式: name_asc(默认), name_desc, match_count_desc"
// @Success 200 {object} map[string]interface{}
// @Router /api/categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// 获取查询参数
	sportIDsParam := c.Query("sport_ids")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	sortBy := c.DefaultQuery("sort", "name_asc")

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

	// 构建基础查询
	query := h.db.Model(&models.Category{}).
		Joins("LEFT JOIN sports ON categories.sport_id = sports.id")

	// 处理 sport_ids 过滤
	if sportIDsParam != "" {
		sportIDs := strings.Split(sportIDsParam, ",")
		// 清理空格
		for i, id := range sportIDs {
			sportIDs[i] = strings.TrimSpace(id)
		}
		query = query.Where("sports.external_id IN ?", sportIDs)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count categories",
		})
		return
	}

	// 构建查询,包含比赛数量统计
	var categories []struct {
		models.Category
		SportName  string `json:"sport_name"`
		MatchCount int64  `json:"match_count"`
	}

	// 子查询统计每个分类下的比赛数量
	subQuery := h.db.Model(&models.Event{}).
		Select("tournaments.category_id, COUNT(*) as match_count").
		Joins("LEFT JOIN tournaments ON events.tournament_id = tournaments.id").
		Where("tournaments.category_id IS NOT NULL").
		Group("tournaments.category_id")

	// 主查询
	mainQuery := h.db.Table("categories").
		Select("categories.*, sports.name as sport_name, COALESCE(match_counts.match_count, 0) as match_count").
		Joins("LEFT JOIN sports ON categories.sport_id = sports.id").
		Joins("LEFT JOIN (?) as match_counts ON categories.id = match_counts.category_id", subQuery)

	// 应用 sport_ids 过滤
	if sportIDsParam != "" {
		sportIDs := strings.Split(sportIDsParam, ",")
		for i, id := range sportIDs {
			sportIDs[i] = strings.TrimSpace(id)
		}
		mainQuery = mainQuery.Where("sports.external_id IN ?", sportIDs)
	}

	// 应用排序
	switch sortBy {
	case "name_desc":
		mainQuery = mainQuery.Order("categories.name DESC")
	case "match_count_desc":
		mainQuery = mainQuery.Order("match_count DESC, categories.name ASC")
	default: // name_asc
		mainQuery = mainQuery.Order("categories.name ASC")
	}

	// 应用分页
	if err := mainQuery.Offset(offset).Limit(pageSize).Scan(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	// 构建响应
	var responseData []CategoryResponse
	for _, cat := range categories {
		responseData = append(responseData, CategoryResponse{
			ID:         cat.ID,
			ExternalID: cat.ExternalID,
			SportID:    cat.SportID,
			SportName:  cat.SportName,
			Name:       cat.Name,
			MatchCount: cat.MatchCount,
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
