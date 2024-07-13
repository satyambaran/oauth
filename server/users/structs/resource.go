package structs

import (
    "time"

    "gorm.io/gorm"
)

type Resource struct {
    ID        int            `gorm:"primaryKey;autoIncrement"`
    Type      int            `gorm:"not null"`
    Name      string         `gorm:"not null"`
    UserId    int            `gorm:"not null" json:"user_id"`
    URI       string         `gorm:"not null" json:"uri"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // IMP Records with a non-null DeletedAt field are considered deleted, but they remain in the database.
}
