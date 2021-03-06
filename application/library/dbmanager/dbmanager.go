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
package dbmanager

import (
	"errors"

	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/echo"
)

func New(ctx echo.Context, auth *driver.DbAuth) *dbManager {
	return &dbManager{
		Context: ctx,
		DbAuth:  auth,
	}
}

type dbManager struct {
	echo.Context
	*driver.DbAuth
	GenURL func(string, ...string) string
}

func (d *dbManager) Driver(typeName string) (driver.Driver, error) {
	dv, ok := driver.Get(typeName)
	if ok {
		dv.Init(d.Context, d.DbAuth)
		return dv, nil
	}
	return nil, errors.New(d.T(`很抱歉，暂时不支持%v`, typeName))
}

func (d *dbManager) Run(typeName string, operation string) error {
	drv, err := d.Driver(typeName)
	if err != nil {
		return err
	}
	if !drv.IsSupported(operation) {
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
	defer drv.SaveResults()
	drv.SetURLGenerator(d.GenURL)
	d.SetFunc(`Results`, drv.SavedResults)
	switch operation {
	case `login`:
		return drv.Login()
	case `logout`:
		return drv.Logout()
	case `processList`:
		return drv.ProcessList()
	case `privileges`:
		return drv.Privileges()
	case `info`:
		return drv.Info()
	case `createDb`:
		return drv.CreateDb()
	case `modifyDb`:
		return drv.ModifyDb()
	case `listDb`:
		return drv.ListDb()
	case `createTable`:
		return drv.CreateTable()
	case `modifyTable`:
		return drv.ModifyTable()
	case `listTable`:
		return drv.ListTable()
	case `viewTable`:
		return drv.ViewTable()
	case `listData`:
		return drv.ListData()
	case `createData`:
		return drv.CreateData()
	case `indexes`:
		return drv.Indexes()
	case `foreign`:
		return drv.Foreign()
	case `trigger`:
		return drv.Trigger()
	case `runCommand`:
		return drv.RunCommand()
	case `import`:
		return drv.Import()
	case `export`:
		return drv.Export()
	default:
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
}
