package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"lakpahana/write-my-naita/internal/llm"
)

const resourceDocxFile = "resources/IT.docx"
const projectDirName = "write-my-naita"
const modulePath = "internal/docx"

func InsertWeeklyTimelineToDocx(report llm.WeeklyTrainingReport, outputDocx string) error {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	resourceDocxFile := string(rootPath) + `/` + modulePath + `/` + resourceDocxFile
	zipReader, err := zip.OpenReader(resourceDocxFile)
	if err != nil {
		return fmt.Errorf("failed to open resource DOCX file: %w", err)
	}
	defer zipReader.Close()

	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	var documentXMLContent []byte
	for _, file := range zipReader.File {
		readCloser, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in DOCX: %w", err)
		}
		defer readCloser.Close()

		if file.Name == "word/document.xml" {
			documentXMLContent, err = ioutil.ReadAll(readCloser)
			if err != nil {
				return fmt.Errorf("failed to read document.xml: %w", err)
			}
			continue
		}

		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			return fmt.Errorf("failed to create file in new DOCX: %w", err)
		}
		if _, err := io.Copy(writer, readCloser); err != nil {
			return fmt.Errorf("failed to copy file to new DOCX: %w", err)
		}
	}

	if documentXMLContent == nil {
		return errors.New("document.xml not found in the resource DOCX")
	}

	modifiedDocumentXML, err := insertWeeklyTimelineToXML(documentXMLContent, report)
	if err != nil {
		return fmt.Errorf("failed to modify document.xml: %w", err)
	}

	writer, err := zipWriter.Create("word/document.xml")
	if err != nil {
		return fmt.Errorf("failed to create modified document.xml in new DOCX: %w", err)
	}
	if _, err := writer.Write(modifiedDocumentXML); err != nil {
		return fmt.Errorf("failed to write modified document.xml: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to finalize new DOCX: %w", err)
	}

	if err := ioutil.WriteFile(outputDocx, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write output DOCX: %w", err)
	}

	return nil
}

func insertWeeklyTimelineToXML(xmlContent []byte, report llm.WeeklyTrainingReport) ([]byte, error) {

	xmlStr := string(xmlContent)
	xmlStr = strings.ReplaceAll(xmlStr, "!---DAY1---!", formatDayContent(report.DailyProgress.Day1))
	xmlStr = strings.ReplaceAll(xmlStr, "!---DAY2---!", formatDayContent(report.DailyProgress.Day2))
	xmlStr = strings.ReplaceAll(xmlStr, "!---DAY3---!", formatDayContent(report.DailyProgress.Day3))
	xmlStr = strings.ReplaceAll(xmlStr, "!---DAY4---!", formatDayContent(report.DailyProgress.Day4))
	xmlStr = strings.ReplaceAll(xmlStr, "!---DAY5---!", formatDayContent(report.DailyProgress.Day5))
	xmlStr = strings.ReplaceAll(xmlStr, "!---LEARNINGS---!", escapeXML(report.KeyLearnings))

	var doc xml.TokenReader
	if err := xml.Unmarshal([]byte(xmlStr), &doc); err != nil {
		return nil, fmt.Errorf("invalid XML after modification: %w", err)
	}

	return []byte(xmlStr), nil
}

func formatDayContent(points []string) string {
	if len(points) == 0 {
		return "No activity documented."
	}
	return strings.Join(points, ", ")
}

func escapeXML(input string) string {

	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"'", "&apos;",
		"\"", "&quot;",
	)
	return replacer.Replace(input)
}
