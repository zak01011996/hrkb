package main

import (
	"log"
	"time"

	"hrkb/conf"
	"hrkb/ctrls"
	"hrkb/logger"
	M "hrkb/models"
	_ "hrkb/routers"
	"hrkb/services/backup"

	"github.com/astaxie/beego"
	"github.com/goamz/goamz/aws"
)

var (
	AppConfig conf.App
)

func init() {
	var err error
	// read main config file
	AppConfig, err = conf.Initialize(beego.AppConfig.String("runmode"), beego.AppConfig)
	if err != nil {
		log.Fatal(err)
	}
	// Open connection to DB and pass connection to DataMapper
	if err := M.DbOpen(AppConfig.Db); err != nil {
		log.Fatal(err)
	}

	// Register and Prepare app models in our DataMapper
	M.PrepareTables(&M.User{}, &M.Cand{}, &M.Rat{}, &M.Crit{}, &M.Dep{}, &M.Comment{}, &M.Role{}, &M.File{}, &M.Contact{}, &M.Mail{}, &M.Lang{})

	if err := logger.Initialize(AppConfig.Log); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initBackUpService()
	// Init and Load existing languages to i18n library
	ctrls.LoadLangs()
	// Init Controllers package
	ctrls.CtrlInit()
	// start Application

	beego.Run()
	M.DbClose()
}

// Init BackUp service if everything ok, otherwise throw fatal err
func initBackUpService() {
	// Getting data from config
	var accessKey, secretKey, bucket, bucketDir, bckpDir string
	var interval, rotLimit, cTimeout, rTimeout, wTimeout int
	var err error
	accessKey = beego.AppConfig.String("aws::accesskey")
	secretKey = beego.AppConfig.String("aws::secretkey")
	bucket = beego.AppConfig.String("aws::bucket")
	bucketDir = beego.AppConfig.String("aws::bucketdir")
	cTimeout, err = beego.AppConfig.Int("aws::contimeout")
	rTimeout, err = beego.AppConfig.Int("aws::reqtimeout")
	wTimeout, err = beego.AppConfig.Int("aws::writetimeout")
	bckpDir = beego.AppConfig.String("bckp::dir")
	interval, err = beego.AppConfig.Int("bckp::interval")
	rotLimit, err = beego.AppConfig.Int("bckp::rotlimit")

	// fatal if there is configuration error
	if err != nil {
		log.Fatal("BackUpService read conf error", err)
	}
	// Auth data for Amazon
	auth := aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	// db config from main configuration
	dbConf := AppConfig.Db

	db := service.DBConf{dbConf.Postgres.Host, dbConf.Postgres.User, dbConf.Postgres.Pass, dbConf.Postgres.Database}
	aws := service.Aws{auth, bucket, bucketDir}

	bp := service.NewBackup(&db, &aws, bckpDir, rotLimit, cTimeout, wTimeout, rTimeout)
	bp.Run(true, time.Duration(interval)*time.Hour)
	go func() {
		for {
			select {
			case err := <-bp.Error:
				beego.Error(err)
			case <-bp.Started:
				beego.Info("Backup service started")
			}
		}

	}()
}
