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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testConfig struct {
	Header sectionTestConfig
}

type sectionTestConfig struct {
	ID      int
	Name    string
	YesOrNo bool
	Pi      float64
	List    string
}

var testConf = testConfig{
	Header: sectionTestConfig{
		ID:      1,
		Name:    "test",
		YesOrNo: true,
		Pi:      3.14,
		List:    "1, 2",
	},
}

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

	configPath := "test.toml"
	expectedConf, _ := ReadConfig(configPath, testLogger)

	assert.NotNil(t, expectedConf)
}

func TestReadConfigEmptyPath(t *testing.T) {
	t.Log("Testing ReadConfig")

	configPath := ""
	expectedConf, _ := ReadConfig(configPath, testLogger)

	assert.NotNil(t, expectedConf)
}

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

func TestGetEtcPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-volume-interface/etc"

	etcPath := GetEtcPath()

	assert.Equal(t, expectedEtcPath, etcPath)
}

func TestGetConfPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-volume-interface/etc/libconfig.toml"

	defaultEtcPath := GetConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetConfPathWithEnv(t *testing.T) {
	t.Log("Testing GetEtcPath")
	err := os.Setenv("SECRET_CONFIG_PATH", "src/github.com/IBM/ibmcloud-volume-interface/etc")
	assert.Nil(t, err)

	expectedEtcPath := "src/github.com/IBM/ibmcloud-volume-interface/etc/libconfig.toml"

	defaultEtcPath := GetConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetDefaultConfPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-volume-interface/etc/libconfig.toml"

	defaultEtcPath := GetDefaultConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetConfPathDir(t *testing.T) {
	t.Log("Testing GetConfPathDir")
	err := os.Setenv("SECRET_CONFIG_PATH", "src/github.com/IBM/ibmcloud-volume-interface/etc/libconfig.toml")
	assert.Nil(t, err)

	expectedEtcPath := "src/github.com/IBM/ibmcloud-volume-interface/etc/libconfig.toml"
	confPath := GetConfPathDir()
	assert.Equal(t, confPath, expectedEtcPath)

	err = os.Unsetenv("SECRET_CONFIG_PATH")
	assert.Nil(t, err)

	err = os.Unsetenv("GOPATH")
	assert.Nil(t, err)

	confPath = GetConfPathDir()
	expectedEtcPath = "src/github.com/IBM/ibmcloud-volume-interface/etc"
	assert.Equal(t, confPath, expectedEtcPath)
}
