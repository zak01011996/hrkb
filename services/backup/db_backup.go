package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

type DBConf struct {
	Host   string // this for DB host
	User   string // this for DB user
	Pass   string // this for DB pass
	DbName string // this for DB name
}

type Aws struct {
	Auth   aws.Auth // this we need to push backup to Amazon S3 Server
	Bucket string   // this is the our amazon S3 bucket name
	Dir    string   // this is the our amazon S3 dir
}

type Backup struct {
	DB        *DBConf       // see DBConf struct
	Aws       *Aws          // see Aws struct
	BackupDir string        // this path for backup directory
	RLimit    int           // this for rotations limit (file counter in bckp dir), if RLimit = 0, old files wouldn't be removed
	AwsCon    *s3.S3        // this handles our S3
	Quit      chan struct{} // this we need to stop our service
	Started   chan struct{} // this we need to stop our service
	Error     chan error    // this error channel, all errors after start pooling will be written to this channel
}

/**
 This func creates .pgpass in user home dir.
 We need it to run pg_dump, without entering the password
**/
func (b *Backup) createPgPass() (err error) {
	fpath := os.Getenv("HOME") + "/.pgpass"
	if _, err = os.Stat(fpath); err == nil {
		if err = b.removePgPass(); err != nil {
			return
		}
	}
	f, err := os.Create(fpath)

	if err != nil {
		return
	}

	f.Chmod(0600)

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s:*:%s:%s:%s", b.DB.Host, b.DB.DbName, b.DB.User, b.DB.Pass))
	return
}

/**
 This func removes .pgpass from user home dir.
**/
func (b *Backup) removePgPass() (err error) {
	fpath := os.Getenv("HOME") + "/.pgpass"
	err = os.Remove(fpath)
	return
}

/**
 This func makes file rotation by RLimit
**/
func (b *Backup) doRotation() (err error) {
	mybucket := b.AwsCon.Bucket(b.Aws.Bucket)
	var files []os.FileInfo
	files, err = ioutil.ReadDir(b.BackupDir)

	if err != nil {
		return
	}

	count_files := len(files)
	if count_files > b.RLimit && b.RLimit > 0 {
		for _, file := range files[:count_files-b.RLimit] {
			if err = os.Remove(b.BackupDir + file.Name()); err != nil {
				return
			}
			if err = mybucket.Del(b.Aws.Dir + file.Name()); err != nil {
				return
			}

		}
	}

	return
}

/**
 This func you can use to display all files in your S3 bucket
**/
func (b *Backup) BucketList() (err error, list []s3.Key) {
	mybucket := b.AwsCon.Bucket(b.Aws.Bucket)
	res, err := mybucket.List(b.Aws.Dir, "", "", 1000)
	if err != nil {
		return
	}
	list = res.Contents

	return
}

/**
 This func saves our backups to amazon S3 bucket
**/
func (b *Backup) BackupToS3(filename string) (err error) {
	mybucket := b.AwsCon.Bucket(b.Aws.Bucket)

	data, err := ioutil.ReadFile(b.BackupDir + filename)
	if err != nil {
		return
	}

	err = mybucket.Put(b.Aws.Dir+filename, data, "text/plain", s3.BucketOwnerFull, s3.Options{})
	if err != nil {
		return
	}
	return
}

/**
 This is the main func, which creats backups of PostgreSQL database to the backup dir using pg_dump util
**/
func (b *Backup) doBackup() (err error) {

	v := time.Now()
	filename := fmt.Sprintf("%d_%d_%d_%d_%d_%d.sql", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second())

	_, err = exec.Command("pg_dump", "-h", b.DB.Host, "-U", b.DB.User, "-w", "-d", b.DB.DbName, "-f", b.BackupDir+filename).Output()
	if err != nil {
		return
	}

	err = b.BackupToS3(filename)
	if err != nil {
		return
	}

	err = b.doRotation()
	if err != nil {
		return
	}

	// Here we creating map where key is bucket file name and value is empty struct{}.
	// I'm putting there empty struct because we don't need map to have value, we just need the key to compare
	// file names from buacket and out backup dir
	bfiles := make(map[string]struct{})
	err, list := b.BucketList()
	if err != nil {
		return
	}
	for _, v := range list {
		bfiles[v.Key] = struct{}{}
	}

	files, err := ioutil.ReadDir(b.BackupDir)

	if err != nil {
		return
	}

	for _, file := range files {
		if _, ok := bfiles[b.Aws.Dir+file.Name()]; !ok {
			err = b.BackupToS3(file.Name())
			if err != nil {
				return
			}

		}
	}
	return
}

/**
 This func starts polling. Here by interval we call doBackup() func
**/
func (b *Backup) startPolling(backup bool, iv time.Duration) {
	var err error
	b.Started <- struct{}{}

	if backup {
		err = b.doBackup()
		if err != nil {
			b.Error <- err
		}
	}

	tick := time.NewTicker(iv)
	for {
		select {
		case <-b.Quit:
			tick.Stop()
			return

		case <-tick.C:
			err = b.doBackup()
			if err != nil {
				b.Error <- err
			}
		}
	}
}

/**
 This func stops out service
**/
func (b *Backup) Stop() {
	close(b.Quit)
}

/**
 This func starts service
**/
func (b *Backup) Run(backup bool, iv time.Duration) {
	go func() {
		err := b.createPgPass()

		if err != nil {
			b.Error <- err
		}

		_, err = os.Stat(b.BackupDir)
		if os.IsNotExist(err) {
			if err = os.Mkdir(b.BackupDir, os.ModePerm); err != nil {
				b.Error <- err
			}
		}

		go b.startPolling(backup, iv)
		<-b.Quit
	}()
}

/**
 This func we need to initialize our backup service in main programm
**/
func NewBackup(db *DBConf, awsc *Aws, backupdir string, rlimit, cTimeout, wTimeout, rTimeout int) *Backup {
	quit := make(chan struct{})
	started := make(chan struct{})
	err := make(chan error)

	zone := aws.EUWest
	connection := s3.New(awsc.Auth, zone)
	connection.ConnectTimeout = time.Duration(cTimeout) * time.Second
	connection.WriteTimeout = time.Duration(wTimeout) * time.Second
	connection.RequestTimeout = time.Duration(rTimeout) * time.Second

	return &Backup{db, awsc, backupdir, rlimit, connection, quit, started, err}
}
