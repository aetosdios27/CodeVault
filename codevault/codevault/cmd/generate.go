package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/schollz/progressbar/v3"
    "codevault/internal/utils"
    "strings"
)

// Inside the Codeforces processing block:
cf := &judges.Codeforces{ /* ... */ }
subs, err := cf.GetSubmissions()
if err != nil {
	fmt.Printf("⚠️ Failed to fetch Codeforces submissions: %v\n", err)
	continue
}

bar := progressbar.Default(int64(len(subs)), "Processing Codeforces")
for _, sub := range subs {
	code, err := cf.FetchCode(sub.ID)
	if err != nil {
		fmt.Printf("⚠️ Skipping %s: %v\n", sub.ProblemID, err)
		continue
	}

	ext := utils.GetFileExtension(sub.Language)
	filePath := fmt.Sprintf("%s/codeforces/%s/solution.%s", outputPath, sub.ProblemID, ext)

	// Save code to filePath (create directories if needed)
	// Then commit using git.CommitFile()

	bar.Add(1)
}

func saveFile(path, content string) (bool, error) {
	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return false, err
	}

	// Check if file exists and has same content
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == content {
			return false, nil // Skip duplicate
		}
	}

	return true, os.WriteFile(path, []byte(content), 0644)
}

func saveMetadata(path string, meta ProblemMetadata) error {
	readmePath := filepath.Join(filepath.Dir(path), "README.md")
	content := fmt.Sprintf("# %s\n- Difficulty: %s\n- Tags: %s",
		meta.Title,
		meta.Difficulty,
		strings.Join(meta.Tags, ", "),
	)
	return os.WriteFile(readmePath, []byte(content), 0644)
}

// In processing loop:
filePath := fmt.Sprintf("%s/%s/%s/solution.%s", 
	outputPath, 
	"codeforces", 
	sub.ProblemID, 
	utils.GetFileExtension(sub.Language),
)

created, err := saveFile(filePath, code)
if err != nil {
	fmt.Printf("⚠️ Failed to save %s: %v\n", filePath, err)
	continue
}

if created {
	git.CommitFile(repo, filePath, fmt.Sprintf("Add %s solution", sub.ProblemID))
}

meta, err := judge.GetProblemMetadata(sub.ProblemID)
if err == nil {
	if err := saveMetadata(filePath, meta); err != nil {
		log.Printf("Failed to save metadata: %v", err)
	}
}