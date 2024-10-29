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
	"flag"
	"github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"
	"os"
)

type robotOptions struct {
	service           options.ServiceOptions
	enableDebug       bool
	delHmacSecretFile bool
	shutdown          bool
	hmacSecretFile    string
	handlePath        string
}

func (o *robotOptions) openDebug(fs *flag.FlagSet) func() {
	fs.BoolVar(
		&o.enableDebug, "enable-debug", false,
		"whether to enable debug model.",
	)

	return func() {
		if o.enableDebug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("debug enabled.")
		}
	}

}

func (o *robotOptions) loadSecret(fs *flag.FlagSet) func() []byte {
	fs.StringVar(
		&o.hmacSecretFile, "hmac-secret-file", "/etc/webhook/hmac",
		"Path to the file containing the HMAC secret.",
	)
	fs.BoolVar(
		&o.delHmacSecretFile, "del-secret", true,
		"whether to delete HMAC secret file.",
	)

	return func() []byte {
		hmac, err := secret.LoadSingleSecret(o.hmacSecretFile)
		if err != nil {
			logrus.Errorf("load hmac, err:%s", err.Error())
			o.shutdown = true
		}
		if o.delHmacSecretFile {
			if err = os.Remove(o.hmacSecretFile); err != nil {
				logrus.Errorf("remove hmac, err:%s", err.Error())
				o.shutdown = true
			}
		}
		return hmac
	}
}

func (o *robotOptions) Validate() error {
	return o.service.Validate()
}

func (o *robotOptions) gatherOptions(fs *flag.FlagSet, args ...string) (*configuration, []byte) {

	o.service.AddFlags(fs)
	debug := o.openDebug(fs)
	hmacFunc := o.loadSecret(fs)
	fs.StringVar(
		&o.handlePath, "handle-path", "webhook",
		"http server handle interface path",
	)

	_ = fs.Parse(args)

	if err := o.service.Validate(); err != nil {
		logrus.Errorf("invalid service options, err:%s", err.Error())
		o.shutdown = true
		return nil, nil
	}
	cfg, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config, err:%s", err.Error())
		o.shutdown = true
		return nil, nil
	}

	debug()
	hmac := hmacFunc()

	return &cfg, hmac
}
