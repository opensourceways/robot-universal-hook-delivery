// Copyright 2024 Chao Feng
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
	"bytes"
	"encoding/gob"
	kafka "github.com/opensourceways/kafka-lib/agent"
	"net/http"
	"path/filepath"
	"reflect"
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

	absPath, err := filepath.Abs("./testdata/config.yaml")
	if err != nil {
		t.Errorf("mismatch data type, 2")
	}

	got, err := loadConfig(absPath)
	if err != nil {
		t.Errorf("mismatch data type, %+v", err)
	}

	assertData(t, got, want)

	absPath, err = filepath.Abs("./testdata/config1.yaml")
	if err != nil {
		t.Errorf("mismatch data type, 2")
	}

	_, err = loadConfig(absPath)
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assertData(t, err.Error(), "missing topic")
	}

	absPath, err = filepath.Abs("./testdata/config2.yaml")
	if err != nil {
		t.Errorf("mismatch data type, 2")
	}

	_, err = loadConfig(absPath)
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assertData(t, err.Error(), "missing user_agent")
	}

	absPath, err = filepath.Abs("./testdata/config3.yaml")
	if err != nil {
		t.Errorf("mismatch data type, 2")
	}

	_, err = loadConfig(absPath)
	if err == nil {
		t.Errorf("mismatch check")
	} else {
		assertData(t, err.Error(), "invalid mq address")
	}
}

func assertMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func assertData(t *testing.T, got any, want any) {
	t.Helper()
	t1, t2 := reflect.TypeOf(got), reflect.TypeOf(want)
	if t1.String() != t2.String() {
		t.Errorf("mismatch data type, got: %v, want %v", t1.String(), t2.String())
	}

	if t1.String() == "string" {
		if got != want {
			t.Errorf("string data different, got: %s, want %s", got, want)
		}
		return
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(got); err != nil {
		t.Errorf("got: %v", err)
	}
	gotBytes := buf.Bytes()

	buf1 := new(bytes.Buffer)
	enc1 := gob.NewEncoder(buf1)
	if err := enc1.Encode(want); err != nil {
		t.Errorf("want: %v", err)
	}
	wantBytes := buf1.Bytes()

	if len(gotBytes) != len(wantBytes) {
		t.Errorf("mismatch data length, got: %v, want %v", len(gotBytes), len(wantBytes))
	}

	for i := 0; i < len(wantBytes); i++ {
		if gotBytes[i] != wantBytes[i] {
			t.Errorf("data different, got: %v, want %v", gotBytes, wantBytes)
			break
		}
	}
}

//
//func FuzzName(f *testing.F) {
//	f.Fuzz(func(t *testing.T) {
//
//	})
//}
