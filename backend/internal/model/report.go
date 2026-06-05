package model

import "time"

type WorkRecord struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	Title      string    `gorm:"type:varchar(100);not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	RecordDate string    `gorm:"type:varchar(10);index;not null" json:"record_date"`
	StartHour  int       `gorm:"not null" json:"start_hour"`
	EndHour    int       `gorm:"not null" json:"end_hour"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DailyReport struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"index;not null" json:"user_id"`
	Title      string    `gorm:"type:varchar(200)" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Type       string    `gorm:"type:varchar(20);not null" json:"type"`
	ReportDate string    `gorm:"type:varchar(10);index;not null" json:"report_date"`
	CreatedAt  time.Time `json:"created_at"`
}

type WeeklyReport struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"type:varchar(200)" json:"title"`
	Summary   string    `gorm:"type:text" json:"summary"`
	Content   string    `gorm:"type:text" json:"content"`
	WeekStart string    `gorm:"type:varchar(10);not null" json:"week_start"`
	WeekEnd   string    `gorm:"type:varchar(10);not null" json:"week_end"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSetting struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	WeeklyReportStyle string    `gorm:"type:varchar(20);default:'professional'" json:"weekly_report_style"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
