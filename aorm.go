package Aorm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Ayanokoji1020-miyano/Aorm/dialect"
	"github.com/Ayanokoji1020-miyano/Aorm/log"
	"github.com/Ayanokoji1020-miyano/Aorm/session"
)

// TODO 都与 scheme 相关
// TODO 使用 Aorm 时 需要主动设置 主键自增 tag
// TODO 模板字段解析比较生硬

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine,err error) {
	db, err := sql.Open(driver,source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
	}
	dial := dialect.OpenDialect(driver,source)

	e = &Engine{
			db:db,
			dialect: dial,
		}
	log.Info("Connect database success")
	return
}



func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}

// NewSession
// 提供 NewSession() 方法，可以通过 Engine 实例创建会话，进而与数据库进行交互了
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db,engine.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err = s.Begin(); err != nil {
		return nil,err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}

// Migrate 数据库迁移
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{},err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist",s.RefTable().Name)
			return nil,s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM `%s` LIMIT 1",table.Name)).QueryRows()
		columns, _ := rows.Columns()
		err = rows.Close()
		if err != nil {
			log.Error("database has not disconnect")
			return
		}
		addCols := difference(table.FieldNames,columns)
		delCols := difference(columns,table.FieldNames)
		log.Infof("added cols %v, deleted cols %v",addCols,delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` `%s`",table.Name,f.Name,f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames,"`, `")
		_, err = s.Raw(fmt.Sprintf("CREATE TABLE `%s` AS SELECT `%s` FROM `%s`;",tmp,fieldStr,table.Name)).Exec()
		if err != nil {
			return
		}
		_, err = s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;",table.Name)).Exec()
		if err != nil {
			return
		}
		_, err = s.Raw(fmt.Sprintf("ALTER TABLE `%s` RENAME TO `%s`;",tmp,table.Name)).Exec()
		return
	})
	return err
}

// 求 a,b 差集 a - b
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]struct{})
	for _,v := range b {
		mapB[v] = struct{}{}
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff,v)
		}
	}
	return
}