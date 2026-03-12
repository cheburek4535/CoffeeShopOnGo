package service

import (
	"time"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
    ID int `gorm:"primaryKey"`
    Name string
    Price float64
    Category string
    VAT int
    IsActive bool
    
    Sales []Sale
}

type Staff struct {
    ID int `gorm:"primaryKey"`
    FullName string
    Salary float64
    Age int
    Position string

    Sales []Sale
}

type Sale struct {
   ID int `gorm:"primaryKey"`
   Amount float64
   CreatedAt time.Time  
   ProductId int
   StaffId int
   Quantity int

   Product *Product
   Staff *Staff
}

type Warehouse struct {
    ID int `gorm:"primaryKey"`
    ProductId int
    Stock int
    Status string

    Product *Product
}

// var DB = *&gorm.DB{}
// func CreateDB() {
//     db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Ошибка подключения к БД: ", err)
// 	}

// 	err = db.AutoMigrate(&Product{}, &Staff{}, &Sale{}, &Warehouse{})
// 	if err != nil {
// 		log.Fatal("Ошибка миграции: ", err)
// 	}

//     DB = db
//     return nil
// }
var DB *gorm.DB

func Init() error {
    db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    if err != nil {
        return err
    }
    err = db.AutoMigrate(&Product{}, &Staff{}, &Sale{}, &Warehouse{})
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}
    DB = db
    return nil
}