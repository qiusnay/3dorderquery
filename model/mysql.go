package model
import (
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/util"
)

const (
	dbPingInterval = 90 * time.Second
	dbMaxLiftTime  = 2 * time.Hour
)

type Database struct {
	Addr     string `toml:"addr"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DbName   string `toml:"dbname"`
	MaxIdle  int    `toml:"max_idle"`
	MaxOpen  int    `toml:"max_open"`
}
type	config struct {
	Master Database   `toml:"master"`
	Slaves []Database `toml:"slave"`
}
var conf config

var DB *gorm.DB

func createConnectionURL(username, password, addr, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, addr, dbName)
}

func DbStart()(*gorm.DB, error) {
	util.Config().Bind("conf", "database", &conf)
	
	logger.Infof("database connect erro : %s", conf.Master.Addr)
	db, err := gorm.Open("mysql", createConnectionURL(conf.Master.Username, conf.Master.Password, conf.Master.Addr, conf.Master.DbName))
	if err != nil {
		logger.Infof("database connect erro : %s", err)
		return db, err
		//panic("连接数据库失败")
	}
	DB = db
	// db.LogMode(true)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(conf.Master.MaxIdle)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(conf.Master.MaxOpen)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Hour)

	go keepDbAlived(db)
	go Automigrate()

	// defer db.Close()

	return db, err
}

func Automigrate() {
	if !DB.HasTable("tb_jd_original_order") {
		DB.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 comment '京东原始订单数据表{裘年宝}'").CreateTable(&JdOriginalOrder{})
	} else {
		// fmt.Println("检查更新.......")
		DB.AutoMigrate(&JdOriginalOrder{})
		// fmt.Println("数据已更新!")
	}
}

func keepDbAlived(db *gorm.DB) {
	t := time.Tick(dbPingInterval)
	var err error
	for {
		<-t
		err = db.DB().Ping()
		if err != nil {
			logger.Infof("database ping: %s", err)
		} else {
			logger.Infof("database ping sucess")
		}
	}
}