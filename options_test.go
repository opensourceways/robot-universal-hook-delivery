// Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGatherOptions(t *testing.T) {

	args := []string{
		"***",
		"--port=8511",
		"--config-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"config4.yaml"),
		"--hmac-secret-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"hmac"),
		"--enable-debug=true",
		"--del-secret=false",
		"--handle-path=gitcode-hook",
	}

	o := new(robotOptions)
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	_, _ = o.gatherOptions(fs, args[1:]...)
	assert.Equal(t, true, o.interrupt)
}

func TestGatherOptions1(t *testing.T) {

	args := []string{
		"***",
		"--port=8511",
		"--config-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"config.yaml"),
		"--hmac-secret-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"hmac1"),
		"--enable-debug=true",
		"--del-secret=false",
		"--handle-path=gitcode-hook",
	}

	o := new(robotOptions)
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	_, _ = o.gatherOptions(fs, args[1:]...)
	assert.Equal(t, true, o.interrupt)

}

func TestGatherOptions2(t *testing.T) {

	args := []string{
		"***",
		"--port=8511",
		"--config-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"config.yaml"),
		"--hmac-secret-file=" + findTestdata(t, "testdata"+string(os.PathSeparator)+"hmac"),
		"--enable-debug=true",
		"--del-secret=false",
		"--handle-path=gitcode-hook",
	}

	o := new(robotOptions)
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	cfg, hmac := o.gatherOptions(fs, args[1:]...)
	assert.Equal(t, false, o.interrupt)

	assert.Equal(t, "127.0.0.1:9092", cfg.Kafka.Address)
	assert.Equal(t, "2.12.0", cfg.Kafka.Version)
	assert.Equal(t, "metadata_webhook_gitcode", cfg.Topic)
	assert.Equal(t, "robot-1", cfg.UserAgent)
	assert.Equal(t, "32123uiyfgdiasd", string(hmac))
}

func findTestdata(t *testing.T, path string) string {

	i := 0
retry:
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Error(path + " not found")
		return ""
	}
	if _, err = os.Stat(absPath); !os.IsNotExist(err) {
		return absPath
	} else {
		i++
		path = ".." + string(os.PathSeparator) + path
		if i <= 3 {
			goto retry
		}
	}

	t.Log(path + " not found")
	return ""
}
