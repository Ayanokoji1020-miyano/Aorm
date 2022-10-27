package session

import (
	"Aorm/log"
	"Aorm/schema"
	"fmt"
	"reflect"
	"strings"
)

// Model 用于给 refTable 赋值
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value,s.dialect)
	}
	return s
}


// RefTable 返回 refTable 的值
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}


func (s *Session) CreateTable() error {
	table := s.RefTable()
	var column []string
	for _, field := range table.Fields {
		column = append(column,fmt.Sprintf("`%s` `%s` `%s`",field.Name,field.Type,field.Tag))
	}
	desc := strings.Join(column,"`, `")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE `%s` (`%s`);",table.Name,desc)).Exec()
	return err
}


func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;",s.refTable.Name)).Exec()
	return err
}


func (s *Session) HasTable() bool {
	sql,values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql,values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)

	return strings.ToLower(tmp) == strings.ToLower(s.RefTable().Name)
}