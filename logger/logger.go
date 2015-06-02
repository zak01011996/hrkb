package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	"hrkb/conf"
)

func Initialize(c conf.Logger) (err error) {

	switch c.Dist {
	case "file":

		if c.File.Minutes <= 0 {
			return
		}

		loggerErrChan := make(chan error)

		now := time.Now()

		m := now.Minute()
		i := (c.File.Minutes-(m-(m/c.File.Minutes*c.File.Minutes)))*60 - now.Second()

		go func(i int) {
			//Pre declaration need this for excluding redeclare any time variables in unlimited loop
			var dir, s string
			var err error
			var now time.Time

			for {

				now = time.Now()

				dir = c.File.Folder + fmt.Sprintf("%02d_%02d_%02d", now.Day(), now.Month(), now.Year())

				err = os.MkdirAll(dir, os.ModePerm)

				if err == nil {
					s = dir + "/" + fmt.Sprintf("%02d_%02d", now.Hour(), now.Minute()/c.File.Minutes*c.File.Minutes)
					if c.File.Ext != "" {
						s += "." + c.File.Ext
					}
					err = beego.SetLogger("file", `{"filename":"`+s+`"}`)
				}

				if err != nil {
					loggerErrChan <- err
				}

				time.Sleep(time.Millisecond * 1000 * time.Duration(i))
				i = c.File.Minutes * 60
			}
		}(i)

		go func() {
			for {
				beego.Error(<-loggerErrChan)
			}
		}()

	}

	return
}
