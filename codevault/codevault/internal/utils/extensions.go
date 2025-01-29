package utils

import "strings"

var langToExt = map[string]string{
	"GNU C++20":  "cpp",
	"GNU C++17":  "cpp",
	"Python 3":   "py",
	"Java 8":     "java",
	"JavaScript": "js",
	"TypeScript": "ts",
	"c++":        "cpp", // LeetCode's C++ label
	"python":     "py",
	"java":       "java",
}

func GetFileExtension(lang string) string {
	if ext, ok := langToExt[strings.ToLower(lang)]; ok {
		return ext
	}
	if ext, ok := langToExt[lang]; ok {
		return ext
	}
	return "txt" // Fallback
}
