package main

import (
	"fmt"
	"github.com/Ayanokoji1020-miyano/Aorm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connectT()
}


func connectT() {
	engine, err := Aorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/aya?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return
	}
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS Userrr;").Exec()
	_, _ = s.Raw("CREATE TABLE Userrr(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE Userrr(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO Userrr(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}