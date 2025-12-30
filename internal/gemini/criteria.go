package gemini

type Criteria struct {
	ForbiddenWords    []string
	ProfessionalCheck bool
	Tone              string
	ExcludePolitics   bool
}
