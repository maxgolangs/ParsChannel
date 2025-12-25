package parser

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"pars/models"
)

type FileWriters struct {
	UsernamesWriter *bufio.Writer
	IDsWriter       *bufio.Writer
	CSVWriter       *csv.Writer
	UsernamesFile   *os.File
	IDsFile         *os.File
	CSVFile         *os.File
}

func initOutputFiles(outputDir string) (*FileWriters, *models.ParseResult, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("create output dir: %w", err)
	}

	usernamesPath := filepath.Join(outputDir, "username.txt")
	idsPath := filepath.Join(outputDir, "id.txt")
	csvPath := filepath.Join(outputDir, "full.csv")

	usernamesFile, err := os.Create(usernamesPath)
	if err != nil {
		return nil, nil, fmt.Errorf("create username.txt: %w", err)
	}

	idsFile, err := os.Create(idsPath)
	if err != nil {
		usernamesFile.Close()
		return nil, nil, fmt.Errorf("create id.txt: %w", err)
	}

	csvFile, err := os.Create(csvPath)
	if err != nil {
		usernamesFile.Close()
		idsFile.Close()
		return nil, nil, fmt.Errorf("create full.csv: %w", err)
	}

	usernamesWriter := bufio.NewWriter(usernamesFile)
	idsWriter := bufio.NewWriter(idsFile)
	csvWriter := csv.NewWriter(csvFile)

	csvWriter.Write([]string{"ID", "FirstName", "LastName", "Username", "IsPremium", "IsBot", "IsRestricted"})

	result := &models.ParseResult{
		UsernamesFile: usernamesPath,
		IDsFile:       idsPath,
		FullCSVFile:   csvPath,
	}

	return &FileWriters{
		UsernamesWriter: usernamesWriter,
		IDsWriter:       idsWriter,
		CSVWriter:       csvWriter,
		UsernamesFile:   usernamesFile,
		IDsFile:         idsFile,
		CSVFile:         csvFile,
	}, result, nil
}

func (fw *FileWriters) Close() {
	fw.UsernamesWriter.Flush()
	fw.IDsWriter.Flush()
	fw.CSVWriter.Flush()
	fw.UsernamesFile.Close()
	fw.IDsFile.Close()
	fw.CSVFile.Close()
}

