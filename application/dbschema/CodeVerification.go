//Do not edit this file, which is automatically generated by the generator.
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
	"time"
)

type CodeVerification struct {
	param   *factory.Param
	trans	*factory.Transaction
	objects []*CodeVerification
	
	Id      	uint    	`db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Code    	string  	`db:"code" bson:"code" comment:"验证码" json:"code" xml:"code"`
	Created 	uint    	`db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Uid     	uint    	`db:"uid" bson:"uid" comment:"创建者" json:"uid" xml:"uid"`
	Used    	uint    	`db:"used" bson:"used" comment:"使用时间" json:"used" xml:"used"`
	Purpose 	string  	`db:"purpose" bson:"purpose" comment:"目的" json:"purpose" xml:"purpose"`
	Start   	uint    	`db:"start" bson:"start" comment:"有效时间" json:"start" xml:"start"`
	End     	uint    	`db:"end" bson:"end" comment:"失效时间" json:"end" xml:"end"`
	Disabled	string  	`db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
}

func (this *CodeVerification) Trans() *factory.Transaction {
	return this.trans
}

func (this *CodeVerification) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *CodeVerification) Objects() []*CodeVerification {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *CodeVerification) NewObjects() *[]*CodeVerification {
	this.objects = []*CodeVerification{}
	return &this.objects
}

func (this *CodeVerification) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("code_verification").SetModel(this)
}

func (this *CodeVerification) SetParam(param *factory.Param) factory.Model {
	this.param = param
	return this
}

func (this *CodeVerification) Param() *factory.Param {
	if this.param == nil {
		return this.NewParam()
	}
	return this.param
}

func (this *CodeVerification) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *CodeVerification) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *CodeVerification) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *CodeVerification) Add() (pk interface{}, err error) {
	this.Created = uint(time.Now().Unix())
	this.Id = 0
	pk, err = this.Param().SetSend(this).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			this.Id = v
		} else if v, y := pk.(int64); y {
			this.Id = uint(v)
		}
	}
	return
}

func (this *CodeVerification) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Update()
}

func (this *CodeVerification) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		
	},func(){
		this.Created = uint(time.Now().Unix())
	this.Id = 0
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			this.Id = v
		} else if v, y := pk.(int64); y {
			this.Id = uint(v)
		}
	}
	return 
}

func (this *CodeVerification) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetArgs(args...).SetMiddleware(mw).Delete()
}

func (this *CodeVerification) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return this.Param().SetArgs(args...).SetMiddleware(mw).Count()
}