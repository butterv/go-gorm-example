package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=dev dbname=postgres password=pass sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	now := time.Now()
	dropTable(db)
	timeAfterDropTable := time.Now()
	fmt.Printf("dropTable: %s\n", timeAfterDropTable.Sub(now).String())

	migrate(db)
	timeAfterMigrate := time.Now()
	fmt.Printf("migrate: %s\n", timeAfterMigrate.Sub(timeAfterDropTable).String())

	insert(db)
	timeAfterInsert := time.Now()
	fmt.Printf("insert: %s\n", timeAfterInsert.Sub(timeAfterMigrate).String())

	selectAndUpdate(db)
	timeAfterSelectAndUpdate := time.Now()
	fmt.Printf("selectAndUpdate: %s\n", timeAfterSelectAndUpdate.Sub(timeAfterInsert).String())

	selectAndDelete(db)
	timeAfterSelectAndDelete := time.Now()
	fmt.Printf("selectAndDelete: %s\n", timeAfterSelectAndDelete.Sub(timeAfterSelectAndUpdate).String())
}

func dropTable(db *gorm.DB) {
	// テーブルを削除
	db.DropTable(&User{})
}

func migrate(db *gorm.DB) {
	// テーブルを作成
	db.AutoMigrate(&User{})
}

func insert(db *gorm.DB) {
	for i := 1; i <= 100; i++ {
		u := User{
			Name: fmt.Sprintf("gorm_test_user_%03d", i),
		}
		// データを登録
		db.Create(&u)
	}
}

func selectAndUpdate(db *gorm.DB) {
	for i := 1; i <= 100; i++ {
		var u User
		// データを取得
		if result := db.First(&u, "id = ?", i); result.Error != nil {
			if result.RecordNotFound() {
				fmt.Printf("record not found(id = %d)\n", i)
			} else {
				fmt.Printf("err(id = %d): %s\n", i, result.Error.Error())
			}
			continue
		}
		// データを更新
		db.Model(&u).Update("name", fmt.Sprintf("gorm_test_user_%03d_updated", i))
	}
}

func selectAndDelete(db *gorm.DB) {
	for i := 1; i <= 100; i++ {
		var u User
		// データを取得
		if result := db.First(&u, "id = ?", i); result.Error != nil {
			if result.RecordNotFound() {
				fmt.Printf("record not found(id = %d)\n", i)
			} else {
				fmt.Printf("err(id = %d): %s\n", i, result.Error.Error())
			}
			continue
		}
		// データを物理削除
		db.Unscoped().Delete(&u)
	}
}
