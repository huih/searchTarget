package main

import (
	"fmt"
	"strconv"
	"io"
	"strings"
	"bufio"
	"os"
	"github.com/gotools/logs"
	"github.com/gotools/files"
	"github.com/gotools/lists/fixedlist"
	"errors"
)

const (
	SEARCH_NO_MATCH = 1
	SEARCH_MATCH = 2
	SEARCH_NO_MATCH_UP = 3
	SEARCH_MATCH_GREP_LINE = 4
)

func SearchMatch(f *os.File, targetArray []string, checkfile string) {
	sourceBuf := bufio.NewReader(f)
	
	index := 0
	lineNo := 0
	
	resultArray := make([]string, len(targetArray))
	for {
		line, err := sourceBuf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		lineNo = lineNo + 1
		
		if strings.Contains(line, targetArray[index]) {
			line = fmt.Sprintf("%s:%d %s", checkfile, lineNo, line)
			resultArray[index] = line
			index = index + 1
			
		} else if index < len(targetArray) {
			index = 0
			continue
		} 
		
		if index == len(targetArray) {
			for _, data := range resultArray {
				logs.Debug(data)
			}
			index = 0
		}
	}
}

func SearchNoMatch(f *os.File, targetArray []string, checkfile string) {
	sourceBuf := bufio.NewReader(f)
	
	index := 0
	lineNo := 0
	
	resultArray := make([]string, len(targetArray))
	for {
		line, err := sourceBuf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		lineNo = lineNo + 1
		
		line = fmt.Sprintf("%s:%d %s", checkfile, lineNo, line)
		if strings.Contains(line, targetArray[index]) {	
			resultArray[index] = line
			index = index + 1
			
		} else if index < len(targetArray) && index > 0 {
			for i:= 0; i < index; i++ {
				logs.Debug(resultArray[i])
			}
			logs.Debug(line)
			index = 0
			continue
		} 
		
		if index >= len(targetArray) {
			index = 0
		}
		
	}
}

func SearchNoMatchUp(f *os.File, targetArray []string, checkfile string) {
	sourceBuf := bufio.NewReader(f)
	
	index := 0
	lineNo := 0
	
	if len(targetArray) < 2 {
		logs.Debug("the target string must exceed two lines or two lines.")
		return
	}
	
	resultArray := make([]string, len(targetArray))
	for {
		line, err := sourceBuf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		lineNo = lineNo + 1
		
		line = fmt.Sprintf("%s:%d %s", checkfile, lineNo, line)
		if index == 0 && strings.Contains(line, targetArray[index]) == false  {	
			resultArray[index] = line
			index = index + 1
		} else if index == 1 && strings.Contains(line, targetArray[index]) == false && strings.Contains(line, targetArray[0]) == false {
			index = 0;
			resultArray[index] = line
			index = index + 1
		} else if index == 1 && strings.Contains(line, targetArray[index]) {
			for i:= 0; i < index; i++ {
				logs.Debug(resultArray[i])
			}
			logs.Debug(line)
			index = 0
		} else {
			index = 0
		}
	}
}

func SearchMatchGrepLine(f *os.File, targetArray []string, checkfile string) {
	
	//read the number of line from args
	lineNum := 1
	if len(os.Args) >= 5 {
		lineNum,_ = strconv.Atoi(os.Args[4])
	}
	if lineNum <= 0 {
		lineNum = 1
	}
	
	sourceBuf := bufio.NewReader(f)
	
	index := 0
	lineNo := 0
	
	fixedList := fixedlist.New(len(targetArray) + lineNum)
	bakFixedList := fixedlist.New(lineNum)
	
	for {
		line, err := sourceBuf.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		lineNo = lineNo + 1
		
		line = fmt.Sprintf("%s:%d %s", checkfile, lineNo, line)
		if strings.Contains(line, targetArray[index]) {
			index = index + 1
			fixedList.Add(line)
		} else {
			fixedList.Add(line)
			index = 0
		}
		
		if index < len (targetArray) {
			continue
		}
		
		if index >= len(targetArray) {
			tmpLineNum := 0
			for {
				//read lineNum line to bakFixedList
				line, err := sourceBuf.ReadString('\n')
				if err != nil || io.EOF == err {//read finished
					break	
				}
				lineNo = lineNo + 1
				line = fmt.Sprintf("%s:%d %s", checkfile, lineNo, line)
				bakFixedList.Add(line)
				tmpLineNum = tmpLineNum + 1
				if tmpLineNum >= lineNum {
					break
				}
			}
		}
		
		if (index >= len(targetArray)) {
			//print fixedList comment
			for {
				v := fixedList.PopFront()
				if v == nil {
					break
				}
				logs.Debug("%s", v)
			}
			
			//print bakFixedList comment
			for {
				v := bakFixedList.PopFront()
				if v == nil {
					break
				}
				logs.Debug("%s", v)
				fixedList.Add(v)
			}
		}
	}
}

func HandleArg() (string, string, int, error){
	targetFile := "target.txt"
	checkFile := "source.txt"
	searchType := SEARCH_NO_MATCH_UP
	
	var err error
	
	if len(os.Args) > 1 && len(os.Args) < 4 {
		logs.Debug("usage searchTarget error, please use help")
		return "", "", 0, errors.New("usage error")
	}
	if len(os.Args) >= 4 {
		searchType, err = strconv.Atoi(os.Args[3])
		if err == nil {
			targetFile = os.Args[1]
			checkFile = os.Args[2]
		} else {
			searchType = SEARCH_NO_MATCH_UP
		}
	}
	
	typeStr := "SEARCH_NO_MATCH"
	if searchType == SEARCH_MATCH {
		typeStr = "SEARCH_MATCH"
	} else if searchType == SEARCH_NO_MATCH_UP {
		typeStr = "SEARCH_NO_MATCH_UP"
	}
	
	logs.Debug("targetFileName:%s, checkFileName:%s, searchType:%s\n", targetFile, checkFile, typeStr)	
	return targetFile, checkFile, searchType, nil
}

func Help() bool {
	if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "help") {
		logs.Debug("Usage: searchTarget [targetFileName checkFileName searchType]")
		logs.Debug("searchType have two values 1: SEARCH_NO_MATCH 2:SEARCH_MATCH 3:SEARCH_NO_MATCH_UP 4:SEARCH_MATCH_GREP_LINE")
		return true
	}
	return false
}

func main() {
	
	logs.SetUsePrefix(false)
	if Help() {
		return
	}
	
	fileName, sourceName, searchType, err := HandleArg()
	if err != nil {
		return
	}
	
	targetArray, err := file.ReadFile(fileName)
	if err != nil {
		logs.Debug("open file(%s) error:%s", fileName, err.Error())
		return
	}
	
	sourceFile, err := os.Open(sourceName)
	if err != nil {
		logs.Debug("open file(%s) error: %s", sourceName, err.Error())
		return
	}
	defer sourceFile.Close()
	
	switch searchType {
		case SEARCH_NO_MATCH:
			SearchNoMatch(sourceFile, targetArray, sourceName)
		case SEARCH_MATCH:
			SearchMatch(sourceFile, targetArray, sourceName)
		case SEARCH_NO_MATCH_UP:
			SearchNoMatchUp(sourceFile, targetArray, sourceName)
		case SEARCH_MATCH_GREP_LINE:
			SearchMatchGrepLine(sourceFile, targetArray, sourceName)
		default:
			logs.Debug("no match search type")
	}
	
}