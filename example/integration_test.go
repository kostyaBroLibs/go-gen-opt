// Copyright 2021 kostyaBro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package example stores the integrations tests for showing how programm work.
package example_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Mock for future code
// #TODO update this later
type Worker interface {
	Work(orig []byte) ([]byte, error)
}

type defaultWorker struct{}

func (dw *defaultWorker) Work(orig []byte) ([]byte, error) {
	return orig, nil
}

var tempWorker = &defaultWorker{}

// Mock end

const testCasesFolder = "test_cases"

type integrationTest struct {
	suite.Suite
}

func (t *integrationTest) TestFull() {
	origFileName := "it_original_file.go"
	expectFileName := "it_expected_file.go"
	testcases := t.mustProvideTestCases()

	actual, err := tempWorker.Work(testcases[origFileName])
	t.Require().Nil(err)
	t.Assert().Equal(testcases[expectFileName], actual)
}

func TestFull(t *testing.T) {
	// #TODO: enable tests
	// suite.Run(t, new(integrationTest))
}

func (t *integrationTest) mustProvideTestCases() map[string][]byte {
	pwd, err := os.Getwd()
	t.Require().Nil(err, "can't load pwd")

	testCasesFullPath := fmt.Sprintf("%s/%s", pwd, testCasesFolder)
	dir, err := os.Open(testCasesFullPath)
	t.Require().Nil(err, "can't open testcases dir")

	files, err := dir.ReadDir(-1)
	t.Require().Nil(err, "can't read files in tesctace folder")

	output := make(map[string][]byte, len(files))
	for _, file := range files {
		fileFullPath := fmt.Sprintf("%s/%s", testCasesFullPath, file.Name())
		fileBytes, err := ioutil.ReadFile(fileFullPath)
		t.Require().Nil(err)

		output[file.Name()] = fileBytes
	}

	return output
}
