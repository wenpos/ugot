package main

import (
	"os"
	"fmt"
	"os/exec"
	"ugot/logger"
	"path/filepath"
	"strings"
	"io/ioutil"
	"runtime"
	"strconv"
	"bytes"
	"flag"
	"log"
	"ugot/util"
)

func TestAndAnalyzePackageCoverage(path string) {
	resultFile := getCoverageResultFilePath(PathAdapterSystem(path))
	os.Remove(resultFile)
	filepath.Walk(PathAdapterSystem(path), analyzePackageCoverage)
	util.ParseAnalysisFile(PathAdapterSystem(resultFile))
	CleanAnalyzedPackageFile(PathAdapterSystem(path))
}

func CleanAnalyzedPackageFile(path string) {
	filepath.Walk(path, cleanAnalyzedOutFiles)
}

func cleanAnalyzedOutFiles(path string, info os.FileInfo, err error) error {
	err = matchAndDeleteFiles(err, info, path, `*.out`)
	err = matchAndDeleteFiles(err, info, path, `*.result`)
	return err
}
func matchAndDeleteFiles(err error, info os.FileInfo, path string, pattern string) error {
	ok, err := filepath.Match(pattern, info.Name())
	if ok {
		os.Remove(path)
	}
	return err
}

func analyzePackageCoverage(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		if hasSpecificFiles(path, "_test.go") {
			resultFile := getCoverageResultFilePath(path)
			//go test to generate .out file for analysis
			execGoTestCoverProfile(path, info.Name())
			if hasSpecificFiles(path, ".out") {
				//return code line count of current package
				_, packageLineCountStr := util.GetGoFilesLineCount(PathAdapterSystem(path+"/"+info.Name()+".out"), info.Name())
				//go tool generate every file coverage and write the result for analysis
				execGoToolCover(path, info.Name())
				packageLineCount := float64(packageLineCountStr["total"])
				if hasSpecificFiles(path, ".result") {
					//every package coverage
					packageLineCoverage := util.GetGoCovResultTotalCoverage(PathAdapterSystem(path + "/" + info.Name() + ".result"))
					coveredLineCount := packageLineCount * (packageLineCoverage / 100)
					util.WriteStringFile(resultFile, SplitPath(path, "src/")[1] + ":"+
						strconv.Itoa(packageLineCountStr["total"])+ ":"+
						strconv.FormatFloat(coveredLineCount, 'f', 0, 64)+ ":"+
						strconv.FormatFloat(packageLineCoverage, 'f', 1, 64)+ "%")
				}

			} else {
				logger.GetLogger().Error("Analysis Failed in package [" + path + "]")
			}
		}
	}
	return err
}

func execGoTestCoverProfile(path string, covName string) {
	if runtime.GOOS == "windows" {
		executeGoTestProfileCmd("cmd", PathAdapterSystem(path), "/C", `go`, "test", "-v", "-coverprofile="+covName+".out")
	} else if runtime.GOOS == "linux" {
		executeGoTestProfileCmd("go", PathAdapterSystem(path), "test", "-v", "-coverprofile="+covName+".out")
	}
}

func execGoToolCover(path string, covName string) {
	if runtime.GOOS == "windows" {
		execGoToolCmdAndWriteResultAfterClean("cmd", PathAdapterSystem(path), PathAdapterSystem(path+"/"+covName+".result"), "/C", `go`, "tool", "cover", "-func="+covName+".out")
	} else if runtime.GOOS == "linux" {
		execGoToolCmdAndWriteResultAfterClean("go", PathAdapterSystem(path), PathAdapterSystem(path+"/"+covName+".result"), "tool", "cover", "-func="+covName+".out")
	}
}

func executeGoTestProfileCmd(cmdName string, cmdExePath string, args ... string) bool {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = cmdExePath
	output, err := cmd.CombinedOutput()
	ignore := flag.Arg(1)
	if len(ignore) != 0 && ignore == "--ignore" {
		printIgnoreIfTestFails(output)
	} else {
		printPanicIfTestFails(output, cmdExePath)
	}
	if err != nil {
		//logger.CheckError(err, "Failed to execute command ["+cmdName+"] ")
		return false
	}
	return true
}

func execGoToolCmdAndWriteResultAfterClean(cmdName string, cmdExcPath string, file_path string, args ... string) bool {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = cmdExcPath
	output, err := cmd.CombinedOutput()
	util.WriteBytesFileAfterClean(file_path, output)
	printOutput(output)
	if err != nil {
		logger.CheckError(err, "Failed to execute command ["+cmdName+"] ")
		return false
	}
	return true
}

func getCoverageResultFilePath(path string) string {
	packagePathes := SplitPath(path, "src/")
	moduleName := SplitPath(packagePathes[1], "/")[0]
	resultFile := PathAdapterSystem(packagePathes[0] + "src/" + moduleName + "_analysis")
	return resultFile
}

func hasSpecificFiles(path string, suffix string) bool {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), suffix) {
			return true
		}
	}
	return false
}

func PathAdapterSystem(path string) string {
	if runtime.GOOS == "windows" {
		return filepath.FromSlash(path)
	} else if runtime.GOOS == "linux" {
		return filepath.ToSlash(path)
	} else {
		return path
	}
}

func SplitPath(path string, sep string) []string {
	if runtime.GOOS == "windows" {
		return strings.Split(path, filepath.FromSlash(sep))
	} else if runtime.GOOS == "linux" {
		return strings.Split(path, filepath.ToSlash(sep))
	} else {
		return strings.Split(path, sep)
	}
}

func PathAppend(path ... string) string {
	var buffer bytes.Buffer
	for _, v := range path {
		buffer.WriteString(v)
	}
	return buffer.String()
}

func panicIfTestFails(outs []byte, cmdExePath string) {
	b := bytes.NewBuffer(outs)
	line, err := b.ReadString('\n')
	fmt.Println(bytes.NewReader(outs).ReadAt(outs, int64(1)))
	for ; err == nil; line, err = b.ReadString('\n') {
		if !strings.Contains(line, "ok") && !strings.Contains(line, cmdExePath) {
			failUts(line)
		}
	}
}

func failUts(line string) {
	if strings.Contains(line, "--- FAIL") {
		log.Fatalf("UT failed: %s", line)
	}
}

func printIgnoreIfTestFails(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("%s\n", string(outs))
	}
}

func printPanicIfTestFails(outs []byte, cmdExePath string) {
	if len(outs) > 0 {
		fmt.Printf("%s\n", string(outs))
		panicIfTestFails(outs, cmdExePath)
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("%s\n", string(outs))
	}
}
