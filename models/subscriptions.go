package models

import "time"

type Subscription struct {
	ID        string    `gorm:"primary_key" json:"id,omitempty"`
	UserId    string    `json:"userId,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
}
