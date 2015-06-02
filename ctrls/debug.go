package ctrls

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type Debug struct {
	BaseController
}

// Action to show latest log
// in html page
func (c *Debug) LastLog() {
	// get log dirs
	dir := beego.AppConfig.String("log::folder")

	// today date in logger directory format 02_20_2015
	now := time.Now()
	today := fmt.Sprintf("%02d_%02d_%02d", now.Day(), now.Month(), now.Year())

	// read files from todays log folder
	lastDirPath := dir + today
	files, err := ioutil.ReadDir(lastDirPath)
	if c.CheckErr(err, internalErr, "could not read dir:"+lastDirPath) {
		return
	}
	if len(files) == 0 {
		c.CheckErr(notNilErr, internalErr, "no log files in dir")
		return
	}
	// read last log file
	lastFile := files[len(files)-1].Name()
	if c.CheckErr(err, internalErr) {
		return
	}
	// read file
	filePath := lastDirPath + "/" + lastFile
	buf, err := ioutil.ReadFile(filePath)
	if c.CheckErr(err, internalErr, "File Path:"+filePath) {
		return
	}
	// convert to string from bytes
	str := string(buf)
	// split lines
	logs := strings.Split(str, "\n")
	// number of lines
	n := len(logs)
	reverseLogs := make([]string, n)
	// iterate and populate in reverse
	for i := 0; i < n; i++ {
		reverseLogs[n-1-i] = logs[i]
	}
	c.Data["logs"] = reverseLogs
}
