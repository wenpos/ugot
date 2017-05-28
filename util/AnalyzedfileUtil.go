package util

import (
	"ugot/logger"

	"io/ioutil"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
	"unicode"
	"github.com/wenpos/ugot/util"
)

//注意关闭资源，使用defer
func ReadFileWithOsOpen(file_path string) {
	f, err := os.Open(file_path)
	defer f.Close()
	check(err, file_path)
	buffer1 := make([]byte, 1000)
	n1, err := f.Read(buffer1)
	check(err, file_path)
	fmt.Printf("%d bytes: %s\n", n1, string(buffer1))
}

//一次性读取全部内容，当文件较大时，会占用较多内存
func ReadFileWithIoUtil(file_path string) {
	data, err := ioutil.ReadFile(file_path)
	check(err, file_path)
	fmt.Println(string(data))
}

func PrintFileByLine(file_path string){
	f, err := os.Open(file_path)
	defer f.Close()
	check(err, file_path)
	b := bufio.NewReader(f)
	line, err := b.ReadString('\n')
	for ; err == nil; line, err = b.ReadString('\n') {
		fmt.Print(line)
	}
}

func GetGoFilesLineCount(file_path string, package_str string) (map[string]int, map[string]int) {
	var file2Count = make(map[string]int)
	var totalCount = make(map[string]int)
	f, err := os.Open(file_path)
	defer f.Close()
	check(err, file_path)
	b := bufio.NewReader(f)
	line, err := b.ReadString('\n')
	for ; err == nil; line, err = b.ReadString('\n') {
		if strings.Contains(line, package_str) {
			fileLineCounts := strings.Split(line, package_str)
			fileLineCount := strings.Split(fileLineCounts[1], ":")
			fileName := fileLineCount[0]
			counts := strings.Split(fileLineCount[1], ",")
			endStr := strings.Split(counts[1], ".")[0]
			startStr := strings.Split(counts[0], ".")[0]
			end := convertStr2Int(endStr)
			start := convertStr2Int(startStr)
			if file2Count[fileName] == 0 {
				file2Count[fileName] = end - start + 1
			} else {
				file2Count[fileName] = end - start + 1 + file2Count[fileName]
			}
		}
	}
	if err == io.EOF {
		fmt.Print(line)
	} else {
		panic("read occur error! " + file_path + " is not a file path")
	}
	total := 0;
	for _, count := range file2Count {
		total = count + total
		totalCount["total"] = total
	}
	return file2Count, totalCount
}

func GetGoCovResultTotalCoverage(file_path string) float64 {
	f, err := os.Open(file_path)
	defer f.Close()
	check(err, file_path)
	b := bufio.NewReader(f)
	line, err := b.ReadString('\n')
	total := 0.0
	for ; err == nil; line, err = b.ReadString('\n') {
		if strings.HasPrefix(line, "total:") {
			resultStr := strings.Split(line, `(statements)`)[1]
			result := strings.FieldsFunc(resultStr, unicode.IsSpace)[0]
			total, _ = strconv.ParseFloat(strings.Split(result, "%")[0], 3)
		}
	}
	return total

}

func ParseAnalysisFile(file_path string)  {
	f, err := os.Open(file_path)
	defer f.Close()
	check(err, file_path)
	b := bufio.NewReader(f)
	line, err := b.ReadString('\n')
	totalLines := 0
	totalCoveredLines := 0
	for ; err == nil; line, err = b.ReadString('\n') {
		if len(line) != 0 && strings.Contains(line, util.PathAdapterSystem("/")) {
			parse := strings.Split(line, ":")
			tempTotal, _ := strconv.Atoi(parse[1])
			totalLines = totalLines + tempTotal
			tempCoveredTotal, _ := strconv.Atoi(parse[2])
			totalCoveredLines = totalCoveredLines + tempCoveredTotal
		}
	}
	coverage := float64(totalCoveredLines) * float64(100) / float64(totalLines)
	WriteStringFile(file_path,"total_lines:"+strconv.Itoa(totalLines))
	WriteStringFile(file_path,"total_covered_lines:"+strconv.Itoa(totalCoveredLines))
	WriteStringFile(file_path,"total_coverage:"+strconv.FormatFloat(coverage, 'f', 1, 64) + "%")
	fmt.Println("[package]:[total.lines]:[total.covered.lines]:[package.coverage]:")
	PrintFileByLine(file_path)
}

func convertStr2Int(str string) int {
	intValue, err := strconv.Atoi(str)
	logger.CheckError(err, "Failed to convert str ["+str+"] to int type")
	return intValue
}

func check(err error, filePath string) {
	if err != nil {
		logger.GetLogger().Error("Read file(s) in path [" + filePath + "] failed")
	}
}

func WriteBytesFileAfterClean(file_path string, outs []byte) bool {
	os.Remove(file_path)
	return WriteBytesFile(file_path, outs)
}

func WriteBytesFile(file_path string, outs []byte) bool {
	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		logger.GetLogger().Error("Read file(s) in path [" + file_path + "] failed")
		return false
	}
	_, errByte := file.Write(outs)
	checkWrite(errByte, file_path)
	_, errStr := file.WriteString("\n")
	if checkWrite(errByte, file_path) && checkWrite(errStr, file_path) {
		return true
	}
	return false
}

func WriteStringFile(file_path string, outs string) bool {
	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		logger.GetLogger().Error("Read file(s) in path [" + file_path + "] failed")
		return false
	}
	_, errByte := file.WriteString(outs)
	checkWrite(errByte, file_path)
	_, errStr := file.WriteString("\n")
	if checkWrite(errByte, file_path) && checkWrite(errStr, file_path) {
		return true
	}
	return false
}

func checkWrite(err error, filePath string) bool {
	if err != nil {
		logger.GetLogger().Error("Write file in path [" + filePath + "] failed")
		return false
	}
	return true
}
