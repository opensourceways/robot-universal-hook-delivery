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
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	want := configuration{
		Kafka: kafka.Config{
			Address: "127.0.0.1:9092",
			Version: "2.12.0",
		},
		Topic:     "metadata_webhook_gitcode",
		UserAgent: "robot-1",
	}

	got, err := loadConfig(findTestdata(t, "testdata"+string(os.PathSeparator)+"config.yaml"))
	if err != nil {
		t.Errorf("mismatch data type, %+v", err)
	}
	assert.Equal(t, want, got)

	_, err = loadConfig(findTestdata(t, "testdata"+string(os.PathSeparator)+"config1.yaml"))
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assert.Equal(t, "missing topic", err.Error())
	}

	_, err = loadConfig(findTestdata(t, "testdata"+string(os.PathSeparator)+"config2.yaml"))
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assert.Equal(t, "missing user_agent", err.Error())
	}

	_, err = loadConfig(findTestdata(t, "testdata"+string(os.PathSeparator)+"config3.yaml"))
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assert.Equal(t, "invalid mq address", err.Error())
	}
}
