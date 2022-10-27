package clause

import (
	"fmt"
	"strings"
)

/*
 * 实现各个 SQL 子句的生成规则
 */

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 根据参数数量生成 ? 占位符字符串 （?,?,...）
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}

	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder // values 语句拼接
	var vars []interface{}  // values 语句参数
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %v", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// ORDER BY $field
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _update 设计入参是2个，第一个参数是表名(table)，第二个参数是 map 类型，表示待更新的键值对
func _update(values ...interface{}) (string,[]interface{}) {
	// UPDATE $tableName SET $field = $value
	tableName := values[0]
	fields := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k,v := range fields {
		keys = append(keys,k+" = ?")
		vars = append(vars,v)
	}
	return fmt.Sprintf("UPDATE %s SET %s",tableName,strings.Join(keys,", ")),vars
}


func _delete(values ...interface{}) (string,[]interface{}) {
	// DELETE FROM $tableName
	return fmt.Sprintf("DELETE FROM %s",values[0]),[]interface{}{}
}


func _count(values ...interface{}) (string,[]interface{}) {
	return _select(values[0],[]string{"COUNT(*)"})
}