/**
 * Copyright 2020 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package config ...
package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM/secret-utils-lib/pkg/k8s_utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getContextLogger() (*zap.Logger, zap.AtomicLevel) {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	traceLevel := zap.NewAtomicLevel()
	traceLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleDebugging, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (lvl >= traceLevel.Level()) && (lvl < zapcore.ErrorLevel)
		})),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, traceLevel
}

var testLogger, _ = getContextLogger()

func TestReadConfig(t *testing.T) {
	t.Log("Testing ReadConfig")

	testcases := []struct {
		testcasename string
		configPath   string
		expectedErr  error
	}{
		{
			testcasename: "Valid secret config",
			configPath:   "etc/libconfig.toml",
			expectedErr:  nil,
		},
		{
			testcasename: "Non existing secret",
			configPath:   "etc/non-exist.toml",
			expectedErr:  errors.New("not nil"),
		},
		{
			testcasename: "Invalid secret config",
			configPath:   "etc/invalid-config.toml",
			expectedErr:  errors.New("not nil"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.testcasename, func(t *testing.T) {
			kc, err := k8s_utils.FakeGetk8sClientSet()
			if err != nil {
				t.Errorf("Error getting clientset. Error: %v", err)
			}
			pwd, err := os.Getwd()
			if err != nil {
				t.Errorf("Failed to get current working directory, test will fail, error: %v", err)
			}
			secretDataPath := filepath.Join(pwd, "..", testcase.configPath)
			_ = k8s_utils.FakeCreateSecret(kc, "DEFAULT", secretDataPath)
			_, err = ReadConfig(kc, testLogger)
			if testcase.expectedErr != nil {
				assert.NotNil(t, err, testcase.expectedErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

/*
func TestParseConfig(t *testing.T) {
	t.Log("Testing config parsing")
	var testParseConf testConfig

	configPath := "test.toml"
	err := ParseConfig(configPath, &testParseConf, testLogger)
	assert.Nil(t, err)

	expected := testConf
	assert.Exactly(t, expected, testParseConf)
}

func TestParseConfigNoMatch(t *testing.T) {
	t.Log("Testing config parsing false positive")
	var testParseConf testConfig

	configPath := "test.toml"
	err := ParseConfig(configPath, &testParseConf, testLogger)
	assert.Nil(t, err)

	expected := testConfig{
		Header: sectionTestConfig{
			ID:      1,
			Name:    "testnomatch",
			YesOrNo: true,
			Pi:      3.14,
			List:    "1, 2",
		}}

	assert.NotEqual(t, expected, testParseConf)
}

func TestParseConfigNoMatchTwo(t *testing.T) {
	t.Log("Testing config parsing false positive")
	var testParseConf testConfig

	configPath := "test1.toml"
	err := ParseConfig(configPath, &testParseConf, testLogger)
	assert.Nil(t, err)

	expected := testConfig{
		Header: sectionTestConfig{
			ID:      1,
			Name:    "testnomatch",
			YesOrNo: true,
			Pi:      3.14,
			List:    "1, 2",
		}}

	assert.NotEqual(t, expected, testParseConf)
}
*/
func TestGetGoPath(t *testing.T) {
	t.Log("Testing getting GOPATH")
	goPath := "/tmp"
	err := os.Setenv("GOPATH", goPath)
	assert.Nil(t, err)

	path := GetGoPath()

	assert.Equal(t, goPath, path)
}

func TestGetEnv(t *testing.T) {
	t.Log("Testing getting ENV")
	goPath := "/tmp"
	err := os.Setenv("ENVTEST", goPath)
	assert.Nil(t, err)

	path := getEnv("ENVTEST")

	assert.Equal(t, goPath, path)
}

func TestGetGoPathNullPath(t *testing.T) {
	t.Log("Testing getting GOPATH NULL Path")
	goPath := ""
	err := os.Setenv("GOPATH", goPath)
	assert.Nil(t, err)

	path := GetGoPath()

	assert.Equal(t, goPath, path)
}
