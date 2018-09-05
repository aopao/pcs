package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var tablePrefix = ""

var o = orm.Ormer(nil)

/*
* Database Init Function
 */
func Init() {
	dbHost := beego.AppConfig.String("db_host")
	dbPort := beego.AppConfig.String("db_port")
	dbName := beego.AppConfig.String("db_name")
	dbUsername := beego.AppConfig.String("db_username")
	dbPassword := beego.AppConfig.String("db_password")
	tablePrefix = beego.AppConfig.String("table_prefix")
	logs.Info("Init db connecting...dbHost:", dbHost)
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	maxIdle := 30
	maxConn := 100
	orm.RegisterDataBase("default", "mysql", dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&allowOldPasswords=1", maxIdle, maxConn)

	RegisterModel()

	o = orm.NewOrm()
}

//返回带前缀的表名
func TableName(str string) string {
	return tablePrefix + str
}

/*
* 注册Model
 */
func RegisterModel() {
	orm.RegisterModel(new(BaseCategory), new(MixCategory), new(HollandQuestion), new(HollandAnswerResult), new(User), new(Token))
	//orm.RegisterModelWithPrefix("pcs_", new(BaseCategory), new(MixCategory))
}
