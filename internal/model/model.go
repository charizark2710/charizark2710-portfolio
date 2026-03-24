package model

type WorkExperience struct {
	Headers []string
	Details []WorkExperienceDetail
}

type TechnicalSkills struct {
	Headers []string
	Details []TechnicalSkillDetail
}

type TechnicalSkillDetail struct {
	ProgrammingLanguages []string
	Frameworks           []string
	Databases            []string
	Others               []string
}

type WorkExperienceDetail struct {
	Company     string
	Position    string
	Duration    string
	Summary     string
	Description string
	PreviewFile string
}

type Education struct {
	Headers []string
	Details []EducationDetail
}

type EducationDetail struct {
	Institution string
	Degree      string
	Field       string
	Process     string
}

type PortfolioData struct {
	Headers         []string
	Title           string
	Name            string
	About           string
	WorkExp         WorkExperience
	Education       Education
	TechnicalSkills TechnicalSkills
	Email           string
	GitHub          string
	LinkedIn        string
}
