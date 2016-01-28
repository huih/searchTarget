package main

import (
	"strconv"
	"io"
	"strings"
	"bufio"
	"os"
	"github.com/gotools/logs"
	"github.com/gotools/files"
)

const (
	SEARCH_NO_MATCH = 1
	SEARCH_MATCH = 2
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

func main() {
	searchType := SEARCH_MATCH
	fileName := "target.txt"
	targetArray, err := file.ReadFile(fileName)
	if err != nil {
		logs.Debug("open file(%s) error:%s", fileName, err.Error())
		return
	}
	
	sourceName := "source.txt"
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
		default:
			logs.Debug("no match search type")
	}
}