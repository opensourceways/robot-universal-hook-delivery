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
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/robot-framework-lib/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"unsafe"
)

type delivery struct {
	wg        sync.WaitGroup
	hmac      []byte
	topic     string
	userAgent string
	log       *logrus.Entry
}

func (d *delivery) wait() {
	d.wg.Wait()
}

func (d *delivery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := &client.Authentication{Secret: d.hmac}
	err, payload, eventType, eventGUID := auth.DoAuthentication(w, r)
	if err != nil {
		d.log.WithError(err).Error("request authenticate failed")
		return
	}

	r.Header.Set("User-Agent", d.userAgent)

	m := (*map[string][]string)(unsafe.Pointer(&r.Header))
	header := make(map[string]string, len(*m))
	for k := range r.Header {
		header[k] = r.Header.Get(k)
	}

	d.wg.Add(1)
	d.publish(payload, &header, eventType, eventGUID)
}

func (d *delivery) publish(payload *bytes.Buffer, h *map[string]string, eventType, eventGUID string) {
	defer d.wg.Done()
	l := d.log.WithFields(
		logrus.Fields{
			"event-type": eventType,
			"event-id":   eventGUID,
		},
	)
	if err := kafka.Publish(d.topic, *h, payload.Bytes()); err != nil {
		l.WithError(err).Errorf("failed to send the request to the MQ")
	} else {
		l.Info("the request is successfully sent to the MQ topic[" + d.topic + "]")
	}
}
