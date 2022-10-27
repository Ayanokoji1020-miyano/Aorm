package dialect

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type mysql struct {
	InfoDB
}

// 校验 mysql 结构体是否已经实现 Dialect 接口的所有方法
var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (db *mysql) Init(dbInfo string) {
	// 现阶段暂时只存储 数据库名称
	db.dbName = dbInfo
}

func (db *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint, reflect.Uint32:
		return "integer unsigned"
	case reflect.Uint64:
		return "integer unsigned"
	case reflect.Bool:
		return "bool"
	case reflect.String:
		return "varchar(255)"
	case reflect.Float32:
		return "double precision"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime" /* 还是 timestamp */
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (db *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{ db.dbName, tableName}
	sql := "SELECT `TABLE_NAME` from `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? and `TABLE_NAME`=?"
	return sql, args
}


func parseDBName(source string) (dbName string) {
	s1 := strings.Split(source,"/")
	s2 := strings.Split(s1[len(s1)-1],"?")
	return s2[0]
}