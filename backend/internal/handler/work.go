package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"aiweekly/backend/internal/llm"
	"aiweekly/backend/internal/middleware"
	"aiweekly/backend/internal/model"
	"aiweekly/backend/internal/response"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type WorkHandler struct {
	DB  *gorm.DB
	LLM *llm.Service
}

type createWorkReq struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	RecordDate string `json:"record_date" binding:"required"`
	StartHour  int    `json:"start_hour" binding:"required,min=0,max=23"`
	EndHour    int    `json:"end_hour" binding:"required,min=0,max=23"`
}

type updateWorkReq struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	StartHour int    `json:"start_hour" binding:"required,min=0,max=23"`
	EndHour   int    `json:"end_hour" binding:"required,min=0,max=23"`
}

func (h *WorkHandler) GetRecords(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var records []model.WorkRecord
	if err := h.DB.Where("user_id = ? AND record_date = ?", user.ID, date).
		Order("start_hour ASC").Find(&records).Error; err != nil {
		response.Fail(c, 500, "获取工作记录失败")
		return
	}

	var dailyReport *model.DailyReport
	h.DB.Where("user_id = ? AND report_date = ?", user.ID, date).First(&dailyReport)

	// 计算当日统计
	totalTasks := len(records)
	totalHours := 0
	earliestStart := 0
	latestEnd := 0
	morningHours := 0   // 6:00-12:00
	afternoonHours := 0 // 12:00-18:00
	eveningHours := 0   // 18:00-24:00

	if totalTasks > 0 {
		earliestStart = records[0].StartHour
		latestEnd = records[0].EndHour
		for _, r := range records {
			totalHours += r.EndHour - r.StartHour
			if r.StartHour < earliestStart {
				earliestStart = r.StartHour
			}
			if r.EndHour > latestEnd {
				latestEnd = r.EndHour
			}
			// 计算时段 - 正确处理跨时段的情况
			start := r.StartHour
			end := r.EndHour

			// 计算早晨时段 (6-12)
			if start < 12 && end > 6 {
				mStart := max(start, 6)
				mEnd := min(end, 12)
				if mEnd > mStart {
					morningHours += mEnd - mStart
				}
			}

			// 计算下午时段 (12-18)
			if start < 18 && end > 12 {
				aStart := max(start, 12)
				aEnd := min(end, 18)
				if aEnd > aStart {
					afternoonHours += aEnd - aStart
				}
			}

			// 计算晚上时段 (18-24)
			if start < 24 && end > 18 {
				eStart := max(start, 18)
				eEnd := min(end, 24)
				if eEnd > eStart {
					eveningHours += eEnd - eStart
				}
			}
		}
	}

	response.OK(c, gin.H{
		"date":         date,
		"records":      records,
		"daily_report": dailyReport,
		"stats": gin.H{
			"total_tasks":     totalTasks,
			"total_hours":     totalHours,
			"earliest_start":  earliestStart,
			"latest_end":      latestEnd,
			"morning_hours":   morningHours,
			"afternoon_hours": afternoonHours,
			"evening_hours":   eveningHours,
		},
	})
}

func (h *WorkHandler) CreateRecord(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	var req createWorkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	if req.StartHour >= req.EndHour {
		response.Fail(c, 400, "结束时间必须大于开始时间")
		return
	}

	record := model.WorkRecord{
		UserID:     user.ID,
		Title:      req.Title,
		Content:    req.Content,
		RecordDate: req.RecordDate,
		StartHour:  req.StartHour,
		EndHour:    req.EndHour,
	}

	if err := h.DB.Create(&record).Error; err != nil {
		response.Fail(c, 500, "创建工作记录失败")
		return
	}

	response.OK(c, record)
}

func (h *WorkHandler) UpdateRecord(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的记录ID")
		return
	}

	var req updateWorkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	if req.StartHour >= req.EndHour {
		response.Fail(c, 400, "结束时间必须大于开始时间")
		return
	}

	var record model.WorkRecord
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&record).Error; err != nil {
		response.Fail(c, 404, "记录不存在")
		return
	}

	record.Title = req.Title
	record.Content = req.Content
	record.StartHour = req.StartHour
	record.EndHour = req.EndHour

	if err := h.DB.Save(&record).Error; err != nil {
		response.Fail(c, 500, "更新工作记录失败")
		return
	}

	response.OK(c, record)
}

func (h *WorkHandler) DeleteRecord(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的记录ID")
		return
	}

	result := h.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&model.WorkRecord{})
	if result.Error != nil {
		response.Fail(c, 500, "删除工作记录失败")
		return
	}

	if result.RowsAffected == 0 {
		response.Fail(c, 404, "记录不存在")
		return
	}

	response.OK(c, gin.H{})
}

func (h *WorkHandler) GetMonthlyStats(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	println("GetMonthlyStats called, year:", yearStr, "month:", monthStr, "user ID:", user.ID)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.Fail(c, 400, "无效的年份")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		response.Fail(c, 400, "无效的月份")
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	println("Date range:", startDate.Format("2006-01-02"), "to", endDate.Format("2006-01-02"))

	type recordStat struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var stats []recordStat
	result := h.DB.Table("work_records").
		Select("record_date as date, COUNT(*) as count").
		Where("user_id = ? AND record_date >= ? AND record_date <= ?", user.ID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("record_date").
		Scan(&stats)

	if result.Error != nil {
		println("Error:", result.Error.Error())
		response.Fail(c, 500, "获取月度统计失败")
		return
	}

	println("Found records:", len(stats))
	for i, stat := range stats {
		println(i, ":", stat.Date, "count:", stat.Count)
	}

	response.OK(c, stats)
}

type getWeeklyReportReq struct {
	WeekStart string `json:"week_start" binding:"required"`
	WeekEnd   string `json:"week_end" binding:"required"`
}

func (h *WorkHandler) GetWeeklyReport(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	weekStart := c.Query("week_start")
	weekEnd := c.Query("week_end")

	if weekStart == "" || weekEnd == "" {
		response.Fail(c, 400, "参数错误")
		return
	}

	var report model.WeeklyReport
	if err := h.DB.Where("user_id = ? AND week_start = ? AND week_end = ?", user.ID, weekStart, weekEnd).First(&report).Error; err != nil {
		// 没找到周报，返回空而不是错误
		response.OK(c, nil)
		return
	}

	response.OK(c, report)
}

type generateWeeklyReportReq struct {
	WeekStart string `json:"week_start" binding:"required"`
	WeekEnd   string `json:"week_end" binding:"required"`
}

func (h *WorkHandler) GetMonthlySummary(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.Fail(c, 400, "无效的年份")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		response.Fail(c, 400, "无效的月份")
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	// 计算总工时和任务数
	type workStats struct {
		TotalHours int
		TotalTasks int
	}

	var stats workStats
	h.DB.Table("work_records").
		Select("IFNULL(SUM(end_hour - start_hour), 0) as total_hours, IFNULL(COUNT(*), 0) as total_tasks").
		Where("user_id = ? AND record_date >= ? AND record_date <= ?", user.ID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Scan(&stats)

	// 计算周报数 - 跨月份周的周报也需要统计
	var reportCount int64
	h.DB.Model(&model.WeeklyReport{}).
		Where("user_id = ? AND week_start <= ? AND week_end >= ?", user.ID, endDate.Format("2006-01-02"), startDate.Format("2006-01-02")).
		Count(&reportCount)

	// 计算周数（该月涉及多少周）
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)

	// 计算该月第一天是星期几 (Monday=1, Sunday=0)
	firstWeekday := firstDayOfMonth.Weekday()
	if firstWeekday == time.Sunday {
		firstWeekday = 7
	}

	// 计算该月最后一天是星期几
	lastWeekday := lastDayOfMonth.Weekday()
	if lastWeekday == time.Sunday {
		lastWeekday = 7
	}

	// 计算总天数
	totalDays := lastDayOfMonth.Day()

	// 计算涉及的周数
	weekCount := (totalDays + int(firstWeekday) - 1) / 7
	if (totalDays+int(firstWeekday)-1)%7 != 0 {
		weekCount++
	}

	// 获取繁忙日期
	type busyDay struct {
		Date  string
		Count int
	}
	var busyDays []busyDay
	h.DB.Table("work_records").
		Select("record_date as date, COUNT(*) as count").
		Where("user_id = ? AND record_date >= ? AND record_date <= ?", user.ID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("record_date").
		Having("COUNT(*) >= 7").
		Scan(&busyDays)

	busyDates := make([]string, len(busyDays))
	for i, d := range busyDays {
		busyDates[i] = d.Date
	}

	response.OK(c, gin.H{
		"total_hours":  stats.TotalHours,
		"total_tasks":  stats.TotalTasks,
		"week_count":   weekCount,
		"report_count": reportCount,
		"busy_days":    busyDates,
	})
}

// ===== 日报相关 =====

type generateDailyReportReq struct {
	ReportDate string `json:"report_date" binding:"required"`
}

func (h *WorkHandler) GenerateDailyReport(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	var req generateDailyReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 获取当天的工作记录
	var records []model.WorkRecord
	if err := h.DB.Where("user_id = ? AND record_date = ?", user.ID, req.ReportDate).
		Order("start_hour ASC").Find(&records).Error; err != nil {
		response.Fail(c, 500, "获取工作记录失败")
		return
	}

	if len(records) == 0 {
		response.Fail(c, 400, "当天没有工作记录")
		return
	}

	// 调用 LLM 生成日报
	var reportOutput *llm.DailyReportOutput
	var err error
	if h.LLM != nil && h.LLM.APIKey != "" {
		inputs := make([]llm.WorkRecordInput, len(records))
		for i, r := range records {
			inputs[i] = llm.WorkRecordInput{
				Title:     r.Title,
				Content:   r.Content,
				StartHour: r.StartHour,
				EndHour:   r.EndHour,
			}
		}
		reportOutput, err = h.LLM.GenerateDailyReport(inputs)
		if err != nil {
			response.Fail(c, 500, "生成日报失败: "+err.Error())
			return
		}
	} else {
		// 如果没有配置 LLM，使用简单版本
		totalHours := 0
		content := "今日工作记录：\n\n"
		for i, r := range records {
			content += strconv.Itoa(i+1) + ". " + r.Title + "\n   " + r.Content + "\n   "
			content += strconv.Itoa(r.StartHour) + ":00 - " + strconv.Itoa(r.EndHour) + ":00\n\n"
			totalHours += r.EndHour - r.StartHour
		}
		reportOutput = &llm.DailyReportOutput{
			Title:   req.ReportDate + " 工作总结",
			Content: content,
		}
	}

	// 保存到数据库
	var existingReport model.DailyReport
	result := h.DB.Where("user_id = ? AND report_date = ?", user.ID, req.ReportDate).First(&existingReport)
	if result.Error == nil {
		// 更新已有的
		existingReport.Title = reportOutput.Title
		existingReport.Content = reportOutput.Content
		existingReport.CreatedAt = time.Now()
		if err := h.DB.Save(&existingReport).Error; err != nil {
			response.Fail(c, 500, "更新日报失败")
			return
		}
		response.OK(c, existingReport)
		return
	}

	// 创建新的
	newReport := model.DailyReport{
		UserID:     user.ID,
		Title:      reportOutput.Title,
		Content:    reportOutput.Content,
		Type:       "daily",
		ReportDate: req.ReportDate,
		CreatedAt:  time.Now(),
	}
	if err := h.DB.Create(&newReport).Error; err != nil {
		response.Fail(c, 500, "创建日报失败")
		return
	}
	response.OK(c, newReport)
}

// ===== 周报相关更新 =====

func (h *WorkHandler) GenerateWeeklyReport(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	var req generateWeeklyReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 获取该周所有日报
	var dailyReports []model.DailyReport
	if err := h.DB.Where("user_id = ? AND report_date >= ? AND report_date <= ?",
		user.ID, req.WeekStart, req.WeekEnd).Find(&dailyReports).Error; err != nil {
		response.Fail(c, 500, "获取日报失败")
		return
	}

	// 获取用户设置的风格
	style := "professional"
	var setting model.UserSetting
	if err := h.DB.Where("user_id = ?", user.ID).First(&setting).Error; err == nil {
		style = setting.WeeklyReportStyle
	}

	// 调用 LLM 生成周报
	var reportOutput *llm.WeeklyReportOutput
	var err error
	if h.LLM != nil && h.LLM.APIKey != "" {
		inputs := make([]llm.DailyReportInput, len(dailyReports))
		for i, r := range dailyReports {
			inputs[i] = llm.DailyReportInput{
				Date:    r.ReportDate,
				Title:   r.Title,
				Content: r.Content,
			}
		}
		reportOutput, err = h.LLM.GenerateWeeklyReport(inputs, style)
		if err != nil {
			response.Fail(c, 500, "生成周报失败: "+err.Error())
			return
		}
	} else {
		// 如果没有配置 LLM，使用简单版本
		content := "本周工作总结：\n\n"
		if len(dailyReports) > 0 {
			for _, r := range dailyReports {
				content += r.ReportDate + "：" + r.Title + "\n"
			}
		} else {
			content = "本周暂无日报记录"
		}
		reportOutput = &llm.WeeklyReportOutput{
			Title:   req.WeekStart + " ~ " + req.WeekEnd + " 周报",
			Summary: "本周工作已总结",
			Content: content,
		}
	}

	// 保存到数据库
	var existingReport model.WeeklyReport
	result := h.DB.Where("user_id = ? AND week_start = ? AND week_end = ?", user.ID, req.WeekStart, req.WeekEnd).First(&existingReport)
	if result.Error == nil {
		existingReport.Title = reportOutput.Title
		existingReport.Summary = reportOutput.Summary
		existingReport.Content = reportOutput.Content
		existingReport.CreatedAt = time.Now()
		if err := h.DB.Save(&existingReport).Error; err != nil {
			response.Fail(c, 500, "更新周报失败")
			return
		}
		response.OK(c, existingReport)
		return
	}

	newReport := model.WeeklyReport{
		UserID:    user.ID,
		Title:     reportOutput.Title,
		Summary:   reportOutput.Summary,
		Content:   reportOutput.Content,
		WeekStart: req.WeekStart,
		WeekEnd:   req.WeekEnd,
		CreatedAt: time.Now(),
	}
	if err := h.DB.Create(&newReport).Error; err != nil {
		response.Fail(c, 500, "创建周报失败")
		return
	}
	response.OK(c, newReport)
}

// ===== 用户设置相关 =====

type updateUserSettingReq struct {
	WeeklyReportStyle string `json:"weekly_report_style" binding:"required,oneof=professional relaxed concise"`
}

func (h *WorkHandler) GetUserSetting(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	var setting model.UserSetting
	if err := h.DB.Where("user_id = ?", user.ID).First(&setting).Error; err != nil {
		// 如果不存在，返回默认值
		response.OK(c, model.UserSetting{
			UserID:            user.ID,
			WeeklyReportStyle: "professional",
		})
		return
	}
	response.OK(c, setting)
}

func (h *WorkHandler) UpdateUserSetting(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	var req updateUserSettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	var setting model.UserSetting
	result := h.DB.Where("user_id = ?", user.ID).First(&setting)
	if result.Error == nil {
		// 更新
		setting.WeeklyReportStyle = req.WeeklyReportStyle
		if err := h.DB.Save(&setting).Error; err != nil {
			response.Fail(c, 500, "更新设置失败")
			return
		}
	} else {
		// 创建
		setting = model.UserSetting{
			UserID:            user.ID,
			WeeklyReportStyle: req.WeeklyReportStyle,
		}
		if err := h.DB.Create(&setting).Error; err != nil {
			response.Fail(c, 500, "创建设置失败")
			return
		}
	}
	response.OK(c, setting)
}
