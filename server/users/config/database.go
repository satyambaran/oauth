package config

import (
    "os"
    "time"

    "github.com/joho/godotenv"
    "github.com/satyambaran/oauth/server/users/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    gorm_logger "gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        panic("unable to load env file")
    }
    dbUrl := os.Getenv("SERVER_USER_DB")
    if dbUrl == "" {
        panic("SERVER_USER_DB not set")
    }
    db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
        Logger: gorm_logger.Default.LogMode(gorm_logger.Info),
    })
    if err != nil {
        panic("failed to connect to db")
    }

    migrate(db)
    // db.Clauses(clause.OnConflict{
    //     Columns: []clause.Column{{Name: "id"}},
    //     DoUpdates: clause.AssignmentColumns([]string{
    //         "name", "email", "role", "password", "salt",
    //     }),
    // })
    sqlDB, err := db.DB()
    if err != nil {
        panic("Failed to get database connection pool:" + err.Error())
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)           // Set the maximum number of connections in the idle connection pool
    sqlDB.SetMaxOpenConns(100)          // Set the maximum number of open connections to the database
    sqlDB.SetConnMaxLifetime(time.Hour) // Set the maximum amount of time a connection may be reused

    return db
}
func migrate(db *gorm.DB) {
    for _, model := range models.Models {
        err := db.AutoMigrate(model)
        if err != nil {
            panic("failed to migrate model")
        }
    }
}
