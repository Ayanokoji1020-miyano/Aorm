package session

import (
	"Aorm/dialect"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/aya?charset=utf8mb4&parseTime=True&loc=Local")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

type Use struct {
	Id   int `aorm:"PRIMARY KEY"`
	Name string
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&Use{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
