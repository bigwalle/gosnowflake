// Copyright © 2014 Terry Mao All rights reserved.
// This file is part of gosnowflake.

// gosnowflake is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// gosnowflake is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with gosnowflake.  If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"flag"
	"fmt"
	"github.com/Terry-Mao/goconf"
	"testing"
	"time"
)

func init() {
	flag.StringVar(&confPath, "conf", "./test.conf", " set gosnowflake config file path")
}

var (
	// global config object
	goConf   = goconf.New()
	MyConf   *Config
	confPath string
)

type Config struct {
	RPCAddr   string        `goconf:"base:rpc.addr:,"`
	WorkerId  int64         `goconf:"base:worker"`
	ZKServers []string      `goconf:"zookeeper:addr:,"`
	ZKPath    string        `goconf:"zookeeper:path"`
	ZKTimeout time.Duration `goconf:"zookeeper:timeout:time"`
}

// Init init the configuration file.
func InitConfig() error {
	MyConf = &Config{
		RPCAddr:   "localhost:8080",
		WorkerId:  int64(0),
		ZKServers: []string{"localhost:2181"},
		ZKPath:    "/gosnowflake-servers",
		ZKTimeout: time.Second * 15,
	}
	if err := goConf.Parse(confPath); err != nil {
		return err
	}
	if err := goConf.Unmarshal(MyConf); err != nil {
		return err
	}
	return nil
}

func Test(t *testing.T) {
	if err := InitConfig(); err != nil {
		t.Error(err)
	}
	if err := Init(MyConf.ZKServers, MyConf.ZKPath, MyConf.ZKTimeout); err != nil {
		t.Error(err)
	}
	c := NewClient(MyConf.WorkerId)
	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		id, err := c.Id()
		if err != nil {
			t.Error(err)
		}
		ids, err := c.Ids(5)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("gosnwoflake id: %d\n", id)
		fmt.Printf("gosnwoflake ids: %d\n", ids)
	}
	c.Close()
	// check global cache map
	if _, ok := workerIdMap[MyConf.WorkerId]; ok {
		t.Error("workerId exists")
	}
}
