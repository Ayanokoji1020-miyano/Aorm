package clause

import "strings"

/*
 * 实现结构体 Clause 拼接各个独立的子句
 */

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set 方法根据 Type 调用对应的 generator
// 生成该子句对应的 SQL 语句。
// 并存入结构体中
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句
func (c *Clause) Build(orders ...Type) (string,[]interface{}) {
	var sqls []string
	var vars []interface{}
	for _, ord := range orders {
		if s, ok := c.sql[ord];ok {
			sqls = append(sqls,s)
			vars = append(vars,c.sqlVars[ord]...)
		}
	}
	return strings.Join(sqls," "),vars
}