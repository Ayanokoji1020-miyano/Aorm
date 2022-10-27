package dialect

import (
	"github.com/Ayanokoji1020-miyano/Aorm/log"
	"reflect"
)

// dialect(方言)
var dialectMap = make(map[string]Dialect)

// Dialect
// DataTypeOf 用于将 Go 语言的类型转换为该数据库的数据类型。
// TableExistSQL 返回某个表是否存在的 SQL 语句
type Dialect interface {
	Init(dbInfo string)
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string,[]interface{})
}

// RegisterDialect 注册 dialect 实例
func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

// GetDialect 获取 dialect 实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}

type InfoDB struct {
	dbName string
}

func OpenDialect(driver, source string) (dialect Dialect) {
	// 确保特定方言的存在
	dial, ok := GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	// TODO 一下方式暂定
	dbName := parseDBName(source)
	dial.Init(dbName)
	return dial
}
