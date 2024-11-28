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
	"github.com/opensourceways/robot-framework-lib/config"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"
	"os"
)

type robotOptions struct {
	service           config.FrameworkOptions
	enableDebug       bool
	delHmacSecretFile bool
	interrupt         bool
	hmacSecretFile    string
}

func (o *robotOptions) Validate() error {
	return o.service.Validate()
}

func (o *robotOptions) addFlags(fs *flag.FlagSet) {
	o.service.AddFlagsComposite(fs)
	fs.BoolVar(
		&o.enableDebug, "enable-debug", false,
		"whether to enable debug model.",
	)
	fs.StringVar(
		&o.hmacSecretFile, "hmac-secret-file", "",
		"Path to the file containing the HMAC secret.",
	)
	fs.BoolVar(
		&o.delHmacSecretFile, "del-secret", true,
		"whether to delete HMAC secret file.",
	)
}

// gatherOptions gather the necessary arguments from command line for project startup.
// It returns the configuration and the token to using for subsequent processes.
func (o *robotOptions) gatherOptions(fs *flag.FlagSet, args ...string) (*configuration, []byte) {
	o.addFlags(fs)
	_ = fs.Parse(args)
	cnf, hmacSecret := o.validateFlags()

	return cnf, hmacSecret
}

func (o *robotOptions) validateFlags() (*configuration, []byte) {
	if err := o.service.ValidateComposite(); err != nil {
		logrus.Errorf("invalid service options, err:%s", err.Error())
		o.interrupt = true
		return nil, nil
	}

	configmap, err := loadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.WithError(err).Error("fatal error occurred while loading config")
		o.interrupt = true
		return nil, nil
	}

	hmacSecret, err := secret.LoadSingleSecret(o.hmacSecretFile)
	if err != nil {
		logrus.WithError(err).Error("fatal error occurred while loading secret")
		o.interrupt = true
		return nil, nil
	}
	if o.delHmacSecretFile {
		if err = os.Remove(o.hmacSecretFile); err != nil {
			logrus.WithError(err).Error("fatal error occurred while deleting token")
			o.interrupt = true
			return nil, nil
		}
	}
	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	return &configmap, hmacSecret
}
