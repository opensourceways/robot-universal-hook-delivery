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
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

const component = "robot-gitcode-hook-delivery"

type robotOptions struct {
	service           options.ServiceOptions
	enableDebug       bool
	hmacSecretFile    string
	delHmacSecretFile bool
	handlePath        string
	platform          string
}

func (o *robotOptions) Validate() error {

	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (*robotOptions, error) {
	o := new(robotOptions)

	o.service.AddFlags(fs)
	if err := o.service.Validate(); err != nil {
		return nil, err
	}

	fs.BoolVar(
		&o.enableDebug, "enable-debug", false,
		"whether to enable debug model.",
	)

	fs.BoolVar(
		&o.delHmacSecretFile, "del-secret", true,
		"whether to delete HMAC secret file.",
	)

	fs.StringVar(
		&o.hmacSecretFile, "hmac-secret-file", "/etc/webhook/hmac",
		"Path to the file containing the HMAC secret.",
	)

	fs.StringVar(
		&o.handlePath, "handle-path", "webhook",
		"http server handle interface path",
	)

	fs.StringVar(
		&o.platform, "platform", "gitcode",
		"codehub platform.",
	)

	err := fs.Parse(args)
	return o, err
}

func logger() *logrus.Entry {
	logrusutil.ComponentInit(component)
	return logrus.NewEntry(logrus.StandardLogger())
}

func main() {
	log := logger()
	os.Args = append(os.Args,
		"--port=8511",
		"--config-file=D:\\B\\local\\config-gitcode-hook-delivery.yaml",
		"--hmac-secret-file=D:\\B\\local\\gitcode-secret",
		"--enable-debug=true",
		"--del-secret=false",
		"--handle-path=gitcode-hook",
	) // TODO
	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Errorf("invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config, err:%s", err.Error())
		return
	}

	// init kafka
	if err := kafka.Init(&cfg.Kafka, log, nil, "", false); err != nil {
		logrus.Errorf("init kafka, err:%s", err.Error())
		return
	}

	defer kafka.Exit()

	hmac, err := secret.LoadSingleSecret(o.hmacSecretFile)
	if err != nil {
		logrus.Errorf("load hmac, err:%s", err.Error())

		return
	}

	if o.delHmacSecretFile {
		if err = os.Remove(o.hmacSecretFile); err != nil {
			logrus.Errorf("remove hmac, err:%s", err.Error())
			return
		}
	}

	// server
	d := delivery{
		topic:     cfg.Topic,
		userAgent: cfg.UserAgent,
		hmac:      hmac,
		platform:  o.platform,
	}

	defer d.wait()

	run(&d, o)
}

func run(d *delivery, o robotOptions) {
	defer interrupts.WaitForGracefulShutdown()

	// Return 200 on / for health checks.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// For /**-hook, handle a webhook normally.
	http.Handle("/"+o.handlePath, d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.service.Port)}

	interrupts.ListenAndServe(httpServer, o.service.GracePeriod)
}
