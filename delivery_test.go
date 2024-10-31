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
	"bytes"
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	headerEventType            = "X-GitCode-Event"
	headerEventTypeValue       = "Note Hook"
	headerEventGUID            = "X-GitCode-Delivery"
	headerEventGUIDValue       = "gsadiuoady"
	headerEventToken           = "X-GitCode-Token"
	headerEventTokenValue      = "1234"
	headerUserAgent            = "User-Agent"
	headerUserAgentValue       = "git-gitcode-hook"
	headerContentTypeName      = "Content-Type"
	headerContentTypeJsonValue = "application/json"
)

func TestDelivery(t *testing.T) {
	d := delivery{
		topic:     "",
		userAgent: "gitcode-hook",
		hmac:      []byte("fgiuagyds"),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/case1", bytes.NewBufferString("fihoagdshajbolkhasdb"))
	req.Header.Set(headerUserAgent, headerUserAgentValue)
	req.Header.Set(headerContentTypeName, headerContentTypeJsonValue)
	req.Header.Set(headerEventType, headerEventTypeValue)
	req.Header.Set(headerEventToken, headerEventTokenValue)
	req.Header.Set(headerEventGUID, headerEventGUIDValue)
	d.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	patch := gomonkey.ApplyFunc(kafka.Publish, func(topic string, header map[string]string, msg []byte) error {
		return nil
	})
	defer patch.Reset()
	d.hmac = []byte(headerEventTokenValue)
	logrus.SetLevel(logrus.DebugLevel)
	time.Sleep(2 * time.Second)
	d.ServeHTTP(w, req)
	d.wait()
}

func TestDeliveryError(t *testing.T) {
	d := delivery{
		topic:     "",
		userAgent: "gitcode-hook",
		hmac:      []byte(headerEventTokenValue),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/case2", bytes.NewBufferString("981237096471631"))
	req.Header.Set(headerUserAgent, headerUserAgentValue)
	req.Header.Set(headerContentTypeName, headerContentTypeJsonValue)
	req.Header.Set(headerEventType, headerEventTypeValue)
	req.Header.Set(headerEventToken, headerEventTokenValue)
	req.Header.Set(headerEventGUID, headerEventGUIDValue)

	patch := gomonkey.ApplyFunc(kafka.Publish, func(topic string, header map[string]string, msg []byte) error {
		return errors.New("jhgvkdashgvkhfasda")
	})
	defer patch.Reset()

	logrus.SetLevel(logrus.DebugLevel)
	time.Sleep(2 * time.Second)
	d.ServeHTTP(w, req)
	d.wait()
}
