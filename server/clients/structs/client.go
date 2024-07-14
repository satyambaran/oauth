package structs

import (
    "time"

    "gorm.io/gorm"
)

type Client struct {
    ID        int            `gorm:"primaryKey;autoIncrement"`
    Name      string         `gorm:"not null"`
    Email     string         `gorm:"not null"`
    ClientID  string         `json:"client_id" gorm:"unique;not null"`
    Password  string         `gorm:"not null"`
    Salt      string         `gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // IMP Records with a non-null DeletedAt field are considered deleted, but they remain in the database.
}
