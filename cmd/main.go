package main

import (
	"drawgraph"
	"drawgraph/draw"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Logger         = log.Default()
	defaultDirPath = draw_graph.AppDir() + "/data"

	inDir  = flag.String("d", "", "the directory to look for data files")
	inFile = flag.String("i", "", "input files, comma to separate multiple files")
	outDir = flag.String("o", "", "the dir to store picture")
)

func main() {
	flag.Parse()

	filePaths := inputFiles(*inDir)
	if *inFile != "" {
		filePaths = append(filePaths, strings.Split(*inFile, ",")...)
	}
	if len(filePaths) == 0 {
		//没有输入数据文件地址，使用测试数据
		files, err := ioutil.ReadDir(defaultDirPath)
		if err != nil {
			Logger.Fatal(err)
		}
		for _, file := range files {
			filePath := defaultDirPath + string(os.PathSeparator) + file.Name()
			filePaths = append(filePaths, filePath)
		}
	}

	if *outDir == "" {
		*outDir = defaultDirPath
	}
	fileInfo, err := os.Stat(*outDir)
	switch {
	case os.IsNotExist(err):
		if err := os.MkdirAll(*outDir, 0766); err != nil {
			Logger.Fatal(err)
		}
	case err == nil:
		if !fileInfo.IsDir() {
			Logger.Fatalf("out dir not dir")
		}
	default:
		Logger.Fatal(err)
	}

	for _, filePath := range filePaths {
		fileName := filePath[strings.LastIndex(filePath, "/")+1:]
		dotIdx := strings.LastIndex(fileName, ".")
		name := fileName[:dotIdx]
		suffix := fileName[dotIdx+1:]
		if suffix != "json" {
			continue
		}
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			Logger.Printf("filePath=%s, err=%v\n", filePath, err)
		}
		var v []draw_graph.Node
		if err := json.Unmarshal(bytes, &v); err != nil {
			Logger.Printf("filePath=%s, err=%v\n", filePath, err)
		}
		if err := draw.Draw(name, *outDir, v); err != nil {
			Logger.Printf("%s: %+v", name, err)
		}
	}
}

func inputFiles(dir string) []string {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil
	}
	var files []string

	jsonFiles, err := filepath.Glob(dir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files = append(files, jsonFiles...)
	return files
}
