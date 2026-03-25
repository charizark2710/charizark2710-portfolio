package storage

import (
	"embed"
	"encoding/csv"
	"io"
	"strings"

	"popfolio/internal/model"
)

var DataFS embed.FS

// indexMap builds a map of header name -> column index from a header row
func indexMap(headers []string) map[string]int {
	m := make(map[string]int, len(headers))
	for i, h := range headers {
		m[strings.TrimSpace(h)] = i
	}
	return m
}

// safeGet returns record[index] if index is valid, otherwise ""
func safeGet(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}
	return record[index]
}

func LoadCSVData() (*model.PortfolioData, error) {
	file, err := DataFS.Open("data/portfolio.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	data := &model.PortfolioData{}

	isFirstLine := true
	isFirstLineWorkHeaders := true
	isFirstLineEducationHeaders := true
	isFirstLineTechHeaders := true

	currentSection := ""

	// Header index maps — populated when each section's header row is read
	var headerIdx map[string]int
	var workIdx map[string]int
	var eduIdx map[string]int
	var techIdx map[string]int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip blank lines
		if len(record) == 0 {
			continue
		}

		// Very first line is the top-level headers row
		if isFirstLine {
			isFirstLine = false
			data.Headers = record
			headerIdx = indexMap(record)
			currentSection = "HEADER"
			continue
		}

		section := strings.TrimSpace(record[0])

		// Section markers
		if section == "WORK" || section == "[WORK]" {
			currentSection = "WORK"
			continue
		}
		if section == "EDUCATION" || section == "[EDUCATION]" {
			currentSection = "EDUCATION"
			continue
		}

		if section == "TECHNICAL SKILLS" || section == "[TECHNICAL SKILLS]" {
			currentSection = "TECHNICAL SKILLS"
			continue
		}

		switch currentSection {
		case "HEADER":
			data.Name = safeGet(record, headerIdx["Name"])
			data.Title = safeGet(record, headerIdx["Title"])
			data.About = safeGet(record, headerIdx["About"])
			data.Email = safeGet(record, headerIdx["Email"])
			data.GitHub = safeGet(record, headerIdx["GitHub"])
			data.LinkedIn = safeGet(record, headerIdx["LinkedIn"])

		case "WORK":
			if isFirstLineWorkHeaders {
				data.WorkExp.Headers = append(data.WorkExp.Headers, record...)
				workIdx = indexMap(record)
				isFirstLineWorkHeaders = false
				continue
			}
			data.WorkExp.Details = append(data.WorkExp.Details, model.WorkExperienceDetail{
				Company:     safeGet(record, workIdx["Company"]),
				Position:    safeGet(record, workIdx["Position"]),
				Duration:    safeGet(record, workIdx["Duration"]),
				Summary:     safeGet(record, workIdx["Summary"]),
				Description: safeGet(record, workIdx["Description"]),
				PreviewFile: safeGet(record, workIdx["PreviewFile"]),
				GitHub:      safeGet(record, workIdx["GitHub"]),
			})

		case "EDUCATION":
			if isFirstLineEducationHeaders {
				data.Education.Headers = append(data.Education.Headers, record...)
				eduIdx = indexMap(record)
				isFirstLineEducationHeaders = false
				continue
			}
			data.Education.Details = append(data.Education.Details, model.EducationDetail{
				Institution: safeGet(record, eduIdx["Institution"]),
				Degree:      safeGet(record, eduIdx["Degree"]),
				Field:       safeGet(record, eduIdx["Field"]),
				Process:     safeGet(record, eduIdx["Process"]),
			})

		case "TECHNICAL SKILLS":
			if isFirstLineTechHeaders {
				data.TechnicalSkills.Headers = append(data.TechnicalSkills.Headers, record...)
				techIdx = indexMap(record)
				isFirstLineTechHeaders = false
				continue
			}
			data.TechnicalSkills.Details = append(data.TechnicalSkills.Details, model.TechnicalSkillDetail{
				ProgrammingLanguages: strings.Split(safeGet(record, techIdx["ProgrammingLanguages"]), ";"),
				Frameworks:           strings.Split(safeGet(record, techIdx["Frameworks"]), ";"),
				Databases:            strings.Split(safeGet(record, techIdx["Databases"]), ";"),
				Others:               strings.Split(safeGet(record, techIdx["Others"]), ";"),
			})
		}
	}

	return data, nil
}
