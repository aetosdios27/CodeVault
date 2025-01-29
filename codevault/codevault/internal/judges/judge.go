package judges

type Submission struct {
	ID        string
	ProblemID string
	Language  string
	Code      string
}

type Judge interface {
	GetSubmissions() ([]Submission, error)
	FetchCode(string) (string, error)
}