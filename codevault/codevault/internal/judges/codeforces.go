package judges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Codeforces struct {
	Handle string
	Cookie string
}

type codeforcesSubmissionResponse struct {
	Status string `json:"status"`
	Result []struct {
		ID        int    `json:"id"`
		ProblemID int    `json:"contestId"`
		Index     string `json:"index"`
		Verdict   string `json:"verdict"`
		Language  string `json:"programmingLanguage"`
	} `json:"result"`
}

func (cf *Codeforces) GetSubmissions() ([]Submission, error) {
	apiURL := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", cf.Handle)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Cookie", cf.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var data codeforcesSubmissionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	var submissions []Submission
	for _, sub := range data.Result {
		if sub.Verdict == "OK" { // Only accepted solutions
			submissions = append(submissions, Submission{
				ID:        fmt.Sprintf("%d", sub.ID),
				ProblemID: fmt.Sprintf("%d%s", sub.ProblemID, sub.Index),
				Language:  sub.Language,
			})
		}
	}

	return submissions, nil
}

func (cf *Codeforces) FetchCode(submissionID string) (string, error) {
	submissionURL := fmt.Sprintf("https://codeforces.com/contest/ENTER_CONTEST_ID/submission/%s", submissionID)
	// ^ Replace ENTER_CONTEST_ID dynamically (needs contestId from submission data)

	req, _ := http.NewRequest("GET", submissionURL, nil)
	req.Header.Set("Cookie", cf.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("submission fetch failed: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("HTML parse error: %w", err)
	}

	code := doc.Find("#program-source-text").Text()
	if code == "" {
		return "", fmt.Errorf("code not found (session cookie expired?)")
	}

	return code, nil
}

type ProblemMetadata struct {
	Title      string
	Difficulty int
	Tags       []string
}

func (cf *Codeforces) GetProblemMetadata(problemID string) (ProblemMetadata, error) {
	// Split problemID like "123A" into contest ID and index
	contestID := problemID[:len(problemID)-1]
	index := problemID[len(problemID)-1:]

	apiURL := fmt.Sprintf("https://codeforces.com/api/problemset.problem?contestId=%s&index=%s",
		contestID, index)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ProblemMetadata{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ProblemMetadata{}, err
	}
	defer resp.Body.Close()

	// TODO: Parse JSON to fill ProblemMetadata fields correctly

	return ProblemMetadata{}, nil
}
