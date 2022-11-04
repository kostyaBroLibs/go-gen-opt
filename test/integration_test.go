package test

import (
	"fmt"
	"os"
	"strings"
)

const testDataFolder = "testdata"

type testCase struct {
	sourceFilePath   string
	expectedFilePath string
}

// testCases is the map of test cases where key is the testCase name.
type testCases map[string]*testCase

func (ts testCases) Validate() error {
	for name, testCase := range ts {
		if testCase.sourceFilePath == "" {
			return fmt.Errorf("test case %s has no source file", name)
		}

		if testCase.expectedFilePath == "" {
			return fmt.Errorf("test case %s has no expected file", name)
		}
	}

	return nil
}

func loadTestCases() (testCases, error) {
	files, err := os.ReadDir(testDataFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read testdata directory: %w", err)
	}

	output := make(testCases)

	for i, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".go") {
			continue
		}

		fileName = fileName[:len(fileName)-3]

		separatedName := strings.Split(fileName, "_")
		testCaseName := strings.Join(separatedName[:len(separatedName)-1], "_")

		if output[testCaseName] == nil {
			output[testCaseName] = &testCase{}
		}

		switch {
		case strings.HasSuffix(fileName, "_source"):
			output[testCaseName].sourceFilePath =
				testDataFolder + string(os.PathSeparator) + files[i].Name()
		case strings.HasSuffix(fileName, "_result"):
			output[testCaseName].expectedFilePath =
				testDataFolder + string(os.PathSeparator) + files[i].Name()
		}
	}

	if err = output.Validate(); err != nil {
		return nil, fmt.Errorf("folder with testcases not valid: %w", err)
	}

	return output, nil
}
