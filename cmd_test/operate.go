package main

import (
	"Aorm"
	"Aorm/log"
	_ "github.com/go-sql-driver/mysql"
)

type AA struct {
	Id     int64 `aorm:"PRIMARY KEY"`
	Name   string
	Age    int
}

func main() {
	engine, err := Aorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/aya?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Info(err)
		return
	}
	defer engine.Close()
	e := engine.NewSession()
	a := AA{}
	_ = e.Model(a).RefTable()
	//err = e.CreateTable()
	//if err != nil {
	//	log.Info(err)
	//	return
	//}
	//e.Begin()
	err = engine.Migrate(a)
	if err != nil {
		log.Error(err)
		//e.Rollback()
		return
	}
	//e.Commit()
}
