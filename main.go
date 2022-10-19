package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	zipFile := flag.String("z", "", "source zip file absolute path")
	outputDir := flag.String("o", "", "The output directory")

	flag.Parse()

	if *zipFile == "" {
		flag.PrintDefaults()
		return
	}

	if _, zipFileStatsErr := os.Stat(*zipFile); os.IsNotExist(zipFileStatsErr) {
		flag.PrintDefaults()
		log.Fatalln("Source zip file does not exist")
	}

	if *outputDir == "" {
		flag.PrintDefaults()
		return
	}

	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		flag.PrintDefaults()
		log.Fatalln("Output directory does not exist")
	}

	zipReader, err := zip.OpenReader(*zipFile)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
	}(zipReader)

	// unzip the quarterly zip file
	for _, file := range zipReader.File {
		unzipFileErr := unzipFile(file, *outputDir)

		if unzipFileErr != nil {
			log.Fatalln(unzipFileErr)
		}
	}

	zipFilePath := filepath.Dir(*zipFile)
	println(zipFilePath)

	zipFileName := filepath.Base(*zipFile)

	dbPath := filepath.Join(*outputDir, strings.TrimSuffix(zipFileName, ".zip")+".db")

	//err = RunSQLite3Command(dbPath, ".mode tabs")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	var files = []string{
		"sub",
		"num",
		"pre",
		"tag",
	}

	for _, file := range files {
		log.Printf("Loading %s.txt\n", file)
		txtPath := filepath.Join(*outputDir, fmt.Sprintf("%s.txt", file))

		err = RunSQLite3Command(dbPath, fmt.Sprintf(".import %s %s", txtPath, file))
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Loading %s.txt...DONE\n", file)
	}
}

func unzipFile(f *zip.File, outputDir string) error {
	data, fileErr := f.Open()

	if fileErr != nil {
		return fileErr
	}

	destination, destinationErr := os.Create(filepath.Join(outputDir, f.Name))

	if destinationErr != nil {
		return destinationErr
	}

	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)

	_, copyErr := io.Copy(destination, data)

	if copyErr != nil {
		return copyErr
	}

	return nil
}

func RunSQLite3Command(dbPath string, command string) error {
	cmd := exec.Command("sqlite3", dbPath, "-tabs", "-cmd", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
