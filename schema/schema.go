package schema

import (
	"go/ast"
	"reflect"

	"github.com/Ayanokoji1020-miyano/Aorm/dialect"
)

type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

type Schema struct {
	// 被映射的对象
	Model interface{}

	// 表名
	Name string

	// 字段
	Fields []*Field

	// FieldNames 包含所有的字段名(列名)
	FieldNames []string

	// fieldMap 记录字段名和 Field 的映射关系
	fieldMap map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// reflect.Indirect 获取 dest 值指针指向的实例
	// reflect.ValueOf 获取 dest 的值
	// 为什么不用 reflect 的 Elem():
	// 传入的类型为 Interface，此时可以传入地址也可以传入结构体
	// 如果采用 Elem 函数在对于传入结构体的情况下会发生 Panic，而 Indirect 则可以正常执行
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,

		// modelType.Name 获取结构体名称作为表名
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// 通过下表获取结构体某字段
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			field := &Field{
				Name: f.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}

			if v, ok := f.Tag.Lookup("aorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, f.Name)
			schema.fieldMap[f.Name] = field
		}
	}
	return schema
}

// RecordValue
// 反射，通过字段名返回对应字段值
// u1 := &User{Name: "Tom", Age: 18}
// u2 := &User{Name: "Sam", Age: 25}
// s.Insert(u1, u2, ...)
// u1、u2 转换为 ("Tom", 18), ("Same", 25)
func (schema *Schema) RecordValue(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues,destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}