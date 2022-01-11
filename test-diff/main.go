package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

var testCaseRE = regexp.MustCompile(`func Test(\w+)\(t \*testing\.T\) {`)
var testFile = regexp.MustCompile(`tencentcloud/(\w+)_test\.go$`)
var moduleRE = regexp.MustCompile(`tencentcloud/(data_source|resource)_tc_(\w+)\.go$`)

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func getDiffFiles() []string {
	tag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	assertError(err)

	cmd := exec.Command("git", "diff", "--name-only", "HEAD", strings.TrimSpace(string(tag)))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return []string{}
	}

	return strings.Split(strings.TrimSpace(out.String()), "\n")
}

func main() {
	var testCasePattern string

	files := getDiffFiles()

	func() {
		for i := range files {
			var (
				fileName     = files[i]
				testFilePath string
			)
			// isTestFile
			if testFile.MatchString(fileName) {
				testFilePath = fmt.Sprintf("../%s", fileName)
				// is Datasource/Resource *.go
			} else if moduleRE.MatchString(fileName) {
				testFilePath = fmt.Sprintf("../%s_test.go", strings.TrimSuffix(fileName, ".go"))
			}

			if testFilePath == "" {
				continue
			}

			b, err := ioutil.ReadFile(testFilePath)
			assertError(err)
			result := testCaseRE.FindAllSubmatch(b, -1)
			for _, i := range result {
				if len(i) < 1 {
					continue
				}
				testCasePattern += strings.TrimPrefix(string(i[1]), "AccTencentCloud")
				testCasePattern += "|"
			}

		}
	}()

	fmt.Println(strings.TrimSuffix(testCasePattern, "|"))
}
