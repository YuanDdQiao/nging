/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package mysql

import (
	"database/sql"
	"fmt"

	"github.com/webx-top/com"
)

func (m *mySQL) support(feature string) bool {
	switch feature {
	case "scheme", "sequence", "type", "view_trigger":
		return true
	default:
		if com.VersionCompare(m.getVersion(), "5.1") == 1 {
			switch feature {
			case "event", "partitioning":
				return true
			}
		}
		if com.VersionCompare(m.getVersion(), "5") == 1 {
			switch feature {
			case "routine", "trigger", "view":
				return true
			}
		}
		return false
	}
}

func (m *mySQL) showVariables() ([]map[string]string, error) {
	sqlStr := "SHOW VARIABLES"
	return m.kvVal(sqlStr)
}

func (m *mySQL) killProcess(processId int64) error {
	sqlStr := fmt.Sprintf("KILL %d", processId)
	_, err := m.newParam().SetCollection(sqlStr).Exec()
	return err
}

func (m *mySQL) processList() ([]*ProcessList, error) {
	r := []*ProcessList{}
	sqlStr := "SHOW FULL PROCESSLIST"
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return r, err
	}
	n := len(cols)
	for rows.Next() {
		v := &ProcessList{}
		values := []interface{}{
			&v.Id, &v.User, &v.Host, &v.Db, &v.Command, &v.Time, &v.State, &v.Info, &v.Progress,
		}
		if n == 9 {
			err = rows.Scan(values...)
		} else {
			err = rows.Scan(values[0:n]...)
		}
		if err != nil {
			break
		}
		r = append(r, v)
	}
	return r, err
}

func (m *mySQL) showStatus() ([]map[string]string, error) {
	sqlStr := "SHOW STATUS"
	return m.kvVal(sqlStr)
}

// 获取支持的字符集
func (m *mySQL) getCollations() (*Collations, error) {
	sqlStr := `SHOW COLLATION`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := NewCollations()
	for rows.Next() {
		var v Collation
		err := rows.Scan(&v.Collation, &v.Charset, &v.Id, &v.Default, &v.Compiled, &v.Sortlen)
		if err != nil {
			return nil, err
		}
		coll, ok := ret.Collations[v.Charset.String]
		if !ok {
			coll = []Collation{}
		}
		if v.Default.Valid && len(v.Default.String) > 0 {
			ret.Defaults[v.Charset.String] = len(coll)
		}
		coll = append(coll, v)
		ret.Collations[v.Charset.String] = coll
	}
	return ret, nil
}

func (m *mySQL) getCollation(dbName string, collations *Collations) (string, error) {
	var err error
	if collations == nil {
		collations, err = m.getCollations()
		if err != nil {
			return ``, err
		}
	}
	sqlStr := "SHOW CREATE DATABASE " + quoteCol(dbName)
	row := m.newParam().SetCollection(sqlStr).QueryRow()
	var database sql.NullString
	var createDb sql.NullString
	err = row.Scan(&database, &createDb)
	if err != nil {
		return ``, err
	}
	matches := reCollate.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		return matches[1], nil
	}
	matches = reCharacter.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		if idx, ok := collations.Defaults[matches[1]]; ok {
			return collations.Collations[matches[1]][idx].Collation.String, nil
		}
	}

	return ``, nil
}

func (m *mySQL) getTableStatus(dbName string, tableName string, fast bool) (map[string]*TableStatus, []string, error) {
	sqlStr := `SHOW TABLE STATUS`
	if len(dbName) > 0 {
		sqlStr += " FROM " + quoteCol(dbName)
	}
	if len(tableName) > 0 {
		tableName = quoteVal(tableName, '_', '%')
		sqlStr += ` LIKE ` + tableName
	}
	ret := map[string]*TableStatus{}
	sorts := []string{}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return ret, sorts, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &TableStatus{}
		err := rows.Scan(&v.Name, &v.Engine, &v.Version, &v.Row_format, &v.Rows, &v.Avg_row_length, &v.Data_length, &v.Max_data_length, &v.Index_length, &v.Data_free, &v.Auto_increment, &v.Create_time, &v.Update_time, &v.Check_time, &v.Collation, &v.Checksum, &v.Create_options, &v.Comment)
		if err != nil {
			return ret, sorts, err
		}
		if v.Engine.String == `InnoDB` {
			v.Comment.String = reInnoDBComment.ReplaceAllString(v.Comment.String, `$1`)
		}
		ret[v.Name.String] = v
		sorts = append(sorts, v.Name.String)
		if len(tableName) > 0 {
			return ret, sorts, nil
		}
	}
	return ret, sorts, nil
}

func (m *mySQL) getEngines() ([]*SupportedEngine, error) {
	sqlStr := `SHOW ENGINES`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []*SupportedEngine{}
	for rows.Next() {
		v := &SupportedEngine{}
		err := rows.Scan(&v.Engine, &v.Support, &v.Comment, &v.Transactions, &v.XA, &v.Savepoints)
		if err != nil {
			return nil, err
		}
		if v.Support.String == `YES` || v.Support.String == `DEFAULT` {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (m *mySQL) getVersion() string {
	if len(m.version) > 0 {
		return m.version
	}
	row := m.newParam().SetCollection(`SELECT version()`).QueryRow()
	var v sql.NullString
	err := row.Scan(&v)
	if err != nil {
		return err.Error()
	}
	m.version = v.String
	return v.String
}

func (m *mySQL) baseInfo() error {
	if m.Get(`dbList`) == nil {
		dbList, err := m.getDatabases()
		if err != nil {
			m.fail(err.Error())
			return m.returnTo(`/db`)
		}
		m.Set(`dbList`, dbList)
	}
	if len(m.dbName) > 0 {
		tableList, err := m.getTables()
		if err != nil {
			m.fail(err.Error())
			return m.returnTo(m.GenURL(`listDb`))
		}
		m.Set(`tableList`, tableList)
	}

	m.Set(`dbVersion`, m.getVersion())
	return nil
}
