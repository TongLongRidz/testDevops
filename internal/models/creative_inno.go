package models

import (
	"time"
)

type CreativityInnovation struct {
    CreativityInID   uint      `gorm:"primaryKey;column:creativity_in_id" json:"creativity_in_id"`
    FormID           uint      `gorm:"column:form_id" json:"form_id"`
    DateReceived     time.Time `gorm:"column:date_received" json:"date_received"`
    TeamName         string    `gorm:"column:team_name" json:"team_name"`
    ProjectTitle     string    `gorm:"column:project_title" json:"project_title"`
    Prize            string    `gorm:"column:prize" json:"prize"`
    OrganizedBy      string    `gorm:"column:organized_by" json:"organized_by"`
    CompetitionLevel string    `gorm:"column:competition_level" json:"competition_level"`
    ActivityCategory string    `gorm:"column:activity_category" json:"activity_category"`
}

func (CreativityInnovation) TableName() string {
    return "Creativity_Innovation"
}