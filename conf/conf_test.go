package conf

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
	"testing"
)

type Conf struct {
	Boot   bool    `conf:"block::key1"`
	Int    int     `conf:"key2"`
	Int64  int64   `conf:"key3"`
	Float  float64 `conf:"key4"`
	String string  `conf:"key5"`
	Except string
}

type ConfA struct {
	Diff Conf `conf:"key1"`
}

func TestParseConfig(t *testing.T) {

	beeConf, err := config.NewConfig("ini", "test.conf")

	if err != nil {
		log.Fatal(err)
	}

	s := Conf{}

	if err := ParseConfig(beeConf, &s); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestNotNilErr(t *testing.T) {
	beeConf, err := config.NewConfig("ini", "test.conf")

	if err != nil {
		log.Fatal(err)
	}

	s := Conf{}

	err = ParseConfig(beeConf, s)

	if fmt.Sprintf("%s", err) != NOT_PTR {
		t.Errorf("Expected %s, got %s", NOT_PTR, err)
	}
}

func TestUnsupported(t *testing.T) {
	beeConf, err := config.NewConfig("ini", "test.conf")

	if err != nil {
		log.Fatal(err)
	}

	s := ConfA{}
	err = ParseConfig(beeConf, &s)

	if fmt.Sprintf("%s", err) != UNSUPPORTED_TYPE {
		t.Errorf("Expected %s, got %s", UNSUPPORTED_TYPE, err)
	}
}
