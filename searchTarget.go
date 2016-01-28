package main

import (
	"strconv"
	"io"
	"strings"
	"bufio"
	"os"
	"github.com/gotools/logs"
	"github.com/gotools/files"
	"errors"
)

const (
	SEARCH_NO_MATCH = 1
	SEARCH_MATCH = 2
	SEARCH_NO_MATCH_UP = 3
)

func SearchMatch(f *os.File, targetArray []string) {
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
			line = strconv.Itoa(lineNo) + ": " + line
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

func SearchNoMatch(f *os.File, targetArray []string) {
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
		
		line = strconv.Itoa(lineNo) + ": " + line
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

func SearchNoMatchUp(f *os.File, targetArray []string) {
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
		
		line = strconv.Itoa(lineNo) + ": " + line
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
	logs.Debug("targetFileName:%s, checkFileName:%s, searchType:%s", targetFile, checkFile, typeStr)	
	return targetFile, checkFile, searchType, nil
}

func Help() bool {
	if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "help") {
		logs.Debug("Usage: searchTarget [targetFileName checkFileName searchType]")
		logs.Debug("searchType have two values 1: SEARCH_NO_MATCH 2:SEARCH_MATCH 3:SEARCH_NO_MATCH_UP")
		return true
	}
	return false
}

func main() {
	
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
			SearchNoMatch(sourceFile, targetArray)
		case SEARCH_MATCH:
			SearchMatch(sourceFile, targetArray)
		case SEARCH_NO_MATCH_UP:
			SearchNoMatchUp(sourceFile, targetArray)
		default:
			logs.Debug("no match search type")
	}
}