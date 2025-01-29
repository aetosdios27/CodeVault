package judges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LeetCode struct {
	Cookie string
}

type leetcodeSubmissionResponse struct {
	Submissions []struct {
		ID        string `json:"id"`
		TitleSlug string `json:"title_slug"`
		Lang      string `json:"lang"`
	} `json:"submissions_dump"`
}

func (lc *LeetCode) GetSubmissions() ([]Submission, error) {
	req, _ := http.NewRequest("GET", "https://leetcode.com/api/submissions/", nil)
	req.Header.Set("Cookie", lc.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var data leetcodeSubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	var submissions []Submission
	for _, sub := range data.Submissions {
		submissions = append(submissions, Submission{
			ID:        sub.ID,
			ProblemID: sub.TitleSlug,
			Language:  sub.Lang,
		})
	}

	return submissions, nil
}

func (lc *LeetCode) FetchCode(submissionID string) (string, error) {
	apiURL := fmt.Sprintf("https://leetcode.com/submissions/detail/%s/", submissionID)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Cookie", lc.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("code fetch failed: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("HTML parse error: %w", err)
	}

	// Extract JSON from script tag
	script := doc.Find("#__NEXT_DATA__").First().Text()
	if script == "" {
		return "", fmt.Errorf("script tag not found")
	}

	var data struct {
		Props struct {
			PageProps struct {
				SubmissionDetails struct {
					Code string `json:"code"`
				} `json:"submissionDetails"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal([]byte(script), &data); err != nil {
		return "", fmt.Errorf("JSON unmarshal error: %w", err)
	}

	if data.Props.PageProps.SubmissionDetails.Code == "" {
		return "", fmt.Errorf("code not found (invalid session cookie?)")
	}

	return data.Props.PageProps.SubmissionDetails.Code, nil
}

func (lc *LeetCode) GetProblemMetadata(titleSlug string) (ProblemMetadata, error) {
	graphqlQuery := fmt.Sprintf(`{
		"query": "query { question(titleSlug: \"%s\") { title difficulty topicTags { name } }}"
	}`, titleSlug)

	req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", strings.NewReader(graphqlQuery))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", lc.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ProblemMetadata{}, err
	}
	defer resp.Body.Close()

	// TODO: Parse the JSON response body into ProblemMetadata

	return ProblemMetadata{}, nil // Add a return statement to fix the missing return error
}
