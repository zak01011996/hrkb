package service

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/goamz/goamz/testutil"
	"hrkb/conf"
)

var bp *Backup
var testServer = testutil.NewHTTPServer()

// This we need to emulate response of bucket list
var GetListResultDump1 = ` 
<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01">
  <Name>quotes</Name>
  <Prefix>N</Prefix>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>Nelson</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>5</Size>
    <StorageClass>STANDARD</StorageClass>
    <Owner>
      <ID>bcaf161ca5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
     </Owner>
  </Contents>
  <Contents>
    <Key>Test_file</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>5</Size>
    <StorageClass>STANDARD</StorageClass>
    <Owner>
      <ID>bcaf161ca5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
     </Owner>
  </Contents>
  <Contents>
    <Key>Test_file_2</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>5</Size>
    <StorageClass>STANDARD</StorageClass>
    <Owner>
      <ID>bcaf161ca5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
     </Owner>
  </Contents>

  <Contents>
    <Key>Neo</Key>
    <LastModified>2006-01-01T12:00:00.000Z</LastModified>
    <ETag>&quot;828ef3fdfa96f00ad9f27c383fc9ac7f&quot;</ETag>
    <Size>4</Size>
    <StorageClass>STANDARD</StorageClass>
     <Owner>
      <ID>bcaf1ffd86a5fb16fd081034f</ID>
      <DisplayName>webfile</DisplayName>
    </Owner>
 </Contents>
</ListBucketResult>
`

func init() {
	testServer.Start()

	auth := aws.Auth{
		AccessKey: "abc",
		SecretKey: "123",
	}

	c, err := config.NewConfig("ini", "../../conf/app.conf")

	if err != nil {
		log.Fatal(err)
	}

	TestConfig, err := conf.Initialize("test", c)

	if err != nil {
		log.Fatal(err)
	}

	dbConf := TestConfig.Db

	db := DBConf{dbConf.Postgres.Host, dbConf.Postgres.User, dbConf.Postgres.Pass, dbConf.Postgres.Database}
	awsl := Aws{auth, "git-lab", "hrkb_backups/"}

	bp = NewBackup(&db, &awsl, "test_backup/", 3, 10, 15, 10)

	err = removeTestData()
	if err != nil {
		log.Fatal(err)
	}
	bp.AwsCon = s3.New(awsl.Auth, aws.Region{Name: "faux-region-1", S3Endpoint: testServer.URL})
	err = createTestData()
	if err != nil {
		log.Fatal(err)
	}
}

func createTestData() (err error) {
	_, err = os.Stat(bp.BackupDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(bp.BackupDir, os.ModePerm); err != nil {
			return
		}
	}

	files := []string{"Nelson", "Test_file", "Test_file_1", "Neo"}
	for _, file := range files {
		_, err = os.Create(bp.BackupDir + file)
		if err != nil {
			return
		}
	}
	return
}

func removeTestData() error {
	return os.RemoveAll(bp.BackupDir)
}

func TestCreatePgPass(t *testing.T) {
	err := bp.createPgPass()
	if err != nil {
		t.Error(err)
	}

}

func TestBackupToS3(t *testing.T) {
	testServer.Response(200, nil, "")
	err := bp.BackupToS3("Test_file")
	if err != nil {
		t.Error(err)
	}
	testServer.WaitRequest()
}

func TestDoRotation(t *testing.T) {
	testServer.Response(200, nil, "")
	err := bp.doRotation()
	if err != nil {
		t.Error(err)
	}
	testServer.WaitRequest()
}

func TestDoBackup(t *testing.T) {
	// testServer.Response we need to emulate response for all our requests from doBackup function
	// This for put
	testServer.Response(200, nil, "")
	// This for Rotation
	testServer.Response(200, nil, "")
	// This for BucketList
	testServer.Response(200, nil, GetListResultDump1)
	// This for upload files
	testServer.Response(200, nil, "")
	testServer.Response(200, nil, "")
	testServer.Response(200, nil, "")
	err := bp.doBackup()
	if err != nil {
		t.Error(err)
	}

	testServer.WaitRequest()
}

func TestBucketList(t *testing.T) {
	testServer.Response(200, nil, GetListResultDump1)

	err, _ := bp.BucketList()

	if err != nil {
		t.Error(err)
	}
	testServer.WaitRequest()
}

func TestRun(t *testing.T) {
	bp.Run(false, 1*time.Second)
	for {
		select {
		case <-bp.Started:
			bp.Stop()
			return
		case err := <-bp.Error:
			t.Error(err)
			bp.Stop()
			return
		}
	}
}

func TestRemovePgPass(t *testing.T) {
	err := bp.removePgPass()
	if err != nil {
		t.Error(err)
	}
}
