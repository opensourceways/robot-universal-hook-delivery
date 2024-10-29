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
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/robot-framework-lib/config"
	"github.com/opensourceways/robot-framework-lib/framework"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

const component = "robot-universal-hook-delivery"

func main() {
	logrusutil.ComponentInit(component)
	opt := new(robotOptions)
	cfg, hmac := opt.gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if opt.shutdown {
		return
	}

	lgr := logrus.NewEntry(logrus.StandardLogger())
	// init kafka
	if err := kafka.Init(&cfg.Kafka, lgr, nil, "", false); err != nil {
		logrus.Errorf("init kafka, err:%s", err.Error())
		return
	}

	// server
	handler := &delivery{
		topic:     cfg.Topic,
		userAgent: cfg.UserAgent,
		hmac:      hmac,
	}
	interrupts.OnInterrupt(func() {
		kafka.Exit()
		handler.wait()
	})

	// Return 200 on / for health checks.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	// For /**-hook, handle a webhook normally.
	http.Handle("/"+opt.handlePath, handler)
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(opt.service.Port)}

	framework.StartupServer(httpServer, opt.service, config.ServerAdditionOptions{})
}
