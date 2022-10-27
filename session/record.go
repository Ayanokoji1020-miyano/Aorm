package session

import (
	"errors"
	"github.com/Ayanokoji1020-miyano/Aorm/clause"
	"reflect"
)

/*
	后续所有构造 SQL 语句的方式都将与 Insert 中构造 SQL 语句的方式一致。分两步：
		1）多次调用 clause.Set() 构造好每一个子句。
		2）调用一次 clause.Build() 按照传入的顺序构造出最终的 SQL 语句。
	构造完成后，调用 Raw().Exec() 方法执行
 */

// Insert
//　s := aorm.NewEngine("mysql", "gee.db").NewSession()
//　u1 := &User{Name: "Tom", Age: 18}
//　u2 := &User{Name: "Sam", Age: 25}
//　s.Insert(u1, u2, ...)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	var recordValues []interface{}
	for _, v := range values {
		// 获取模型信息
		table := s.Model(v).RefTable()
		// 构建插入 INSERT INTO $tableName ($fields)
		s.clause.Set(clause.INSERT,table.Name,table.FieldNames)
		recordValues = append(recordValues,table.RecordValue(v))
	}

	// 构建 VALUES ($v1), ($v2), ...
	s.clause.Set(clause.VALUES,recordValues...)
	// 合并 SQL　语句
	sql, vars := s.clause.Build(clause.INSERT,clause.VALUES)
	result, err := s.Raw(sql,vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find
// s.Find(&users)
func (s *Session) Find(value interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(value))
	// destSlice.Type().Elem() 获取切片的单个元素的类型 destType
	// 使用 reflect.New() 方法创建一个 destType 的实例
	// 作为 Model() 的入参，映射出表结构 RefTable()。
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	// 根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
	s.clause.Set(clause.SELECT,table.Name,table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT,clause.WHERE,clause.ORDERBY,clause.LIMIT)
	rows, err := s.Raw(sql,vars).QueryRows()
	if err != nil {
		return err
	}
	// 遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var value []interface{}
		for _, name := range table.FieldNames {
			value = append(value,dest.FieldByName(name).Addr().Interface())
		}
		// 调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段
		if err := rows.Scan(value...); err != nil {
			return err
		}
		// 将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中
		destSlice.Set(reflect.Append(destSlice,dest))
	}
	return rows.Close()
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface());err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}


func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE,s.RefTable().Name,m)
	sql, vars := s.clause.Build(clause.UPDATE,clause.WHERE)
	result, err := s.Raw(sql,vars).Exec()
	if err != nil {
		return 0,err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE,s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE,clause.WHERE)
	result, err := s.Raw(sql,vars).Exec()
	if err != nil {
		return 0,err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT)
	sql, vars := s.clause.Build(clause.COUNT,clause.WHERE)
	row := s.Raw(sql,vars).QueryRow()

	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp,nil
}

/*
 * 链式调用(chain)
 */

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT,num)
	return s
}


func (s *Session) Where(desc string,args ...interface{}) *Session {
	var vars []interface{}
	vars = append(vars,desc)
	s.clause.Set(clause.WHERE,append(vars,args...)...)
	return s
}

func (s *Session) OrderBy(dest string) *Session {
	s.clause.Set(clause.ORDERBY,dest)
	return s
}
