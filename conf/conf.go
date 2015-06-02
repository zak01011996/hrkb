package conf

import (
	"errors"
	"github.com/astaxie/beego/config"
	"reflect"
)

const (
	CONF_TAG         string = "conf"
	NOT_PTR          string = "Interface is not a pointer"
	UNSUPPORTED_TYPE string = "Unsupported type"
)

type Database struct {
	Driver   string
	Postgres struct {
		Host, Database, User, Pass string
	}
}

type Logger struct {
	Dist string
	File struct {
		Folder, Ext string
		Minutes     int
	}
}

type App struct {
	Db  Database
	Log Logger
}

//Parses Config to struct with beego AppConfig
func ParseConfig(beeConf config.ConfigContainer, obj interface{}) error {

	val := reflect.ValueOf(obj)

	if !val.CanAddr() && val.Kind() != reflect.Ptr {
		return errors.New(NOT_PTR)
	}

	val = val.Elem()
	tp := val.Type()

	for i := 0; i < val.NumField(); i++ {
		//Get Config tag
		tag := tp.Field(i).Tag.Get(CONF_TAG)

		if tag == "" {
			continue
		}

		v := val.Field(i)
		//Type casting
		switch v.Interface().(type) {
		case bool:
			b, err := beeConf.Bool(tag)

			if err != nil {
				return err
			}

			v.SetBool(b)
		case int, int64:
			b, err := beeConf.Int(tag)

			if err != nil {
				return err
			}
			v.SetInt(int64(b))
		case float64:
			b, err := beeConf.Float(tag)

			if err != nil {
				return err
			}
			v.SetFloat(b)
		case string:
			v.SetString(beeConf.String(tag))
		default:
			return errors.New(UNSUPPORTED_TYPE)
		}
	}
	return nil
}

func Initialize(env string, cc config.ConfigContainer) (r App, err error) {

	r.Db.Driver = cc.String(env + "::driver")
	r.Db.Postgres.Host = cc.String(env + "::pghost")
	r.Db.Postgres.Database = cc.String(env + "::pgdb")
	r.Db.Postgres.User = cc.String(env + "::pguser")
	r.Db.Postgres.Pass = cc.String(env + "::pgpass")

	r.Log.Dist = cc.String("log::dist")

	if r.Log.Dist == "file" {
		r.Log.File.Minutes, err = cc.Int("log::minutes")
		if err != nil {
			return
		}
		r.Log.File.Ext = cc.String("log::ext")
		r.Log.File.Folder = cc.String("log::folder")
	}

	return
}
