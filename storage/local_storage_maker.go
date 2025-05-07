package storage

import (
	"encoding/json"
	"fmt"
	"os"            // Import os package
	"path/filepath" // Use filepath for OS-independent paths
	"sort"
	"strings"
	"time"
)

// LocalStorageMaker implements the Maker interface for saving reports to the local filesystem.
type LocalStorageMaker struct {
	// baseReportDir stores the absolute path to the "chakara-report" folder within the project root.
	baseReportDir  string
	baseAnswersDir string // Added field for answers directory

}

// NewLocalStorageMaker creates a new local storage maker instance.
// It automatically determines the base directory as "chakara-report" within the current working directory.
func NewLocalStorageMaker() (Maker, error) {
	// Get the current working directory (assumed to be the project root)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Construct the base path for reports: <working_directory>/chakara-report
	baseReportDir := filepath.Join(wd, ReportFolderName)
	// Construct the base path for reports answers: <working_directory>/chakara-report-answers
	baseAnswersDir := filepath.Join(wd, ReportAnswerFolderName)

	// Optional: Log the determined path for verification
	fmt.Printf("Initialized LocalStorageMaker. Reports will be saved under: %s\n", baseReportDir)

	maker := &LocalStorageMaker{
		baseReportDir:  baseReportDir,
		baseAnswersDir: baseAnswersDir,
	}
	return maker, nil
}

// SaveChakaraReportAsText saves the chakra report content as a text file
// in the local filesystem under the automatically determined base report directory.
// Path: <project_root>/chakara-report/email/uniqueId-timestamp.txt
func (maker *LocalStorageMaker) SaveChakaraReportAsText(email string, uniqueId string, content string) error {
	// Get the current time and format it using the constant
	currentTime := time.Now()
	timestamp := currentTime.Format(TimeStampFormat)

	// Construct the specific user's directory path using filepath.Join
	// This creates: <project_root>/chakara-report/email
	dirPath := filepath.Join(maker.baseReportDir, email) // Use baseReportDir

	// Create the directory structure (including parent dirs) if it doesn't exist.
	err := os.MkdirAll(dirPath, FilePermission)
	if err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dirPath, err)
	}

	// Construct the full file path
	// Creates: uniqueId-timestamp.txt
	reportName := fmt.Sprintf("%s-%s%s", uniqueId, timestamp, ReportFileExtension)
	// Creates: <project_root>/chakara-report/email/uniqueId-timestamp.txt
	filePath := filepath.Join(dirPath, reportName)

	// Write the content to the file.
	// os.WriteFile creates the file if it doesn't exist, or truncates it if it does.
	// 0644 provides standard permissions (rw-r--r--).
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		// Use fmt.Errorf to wrap the original error for more context
		return fmt.Errorf("failed to write report file '%s': %w", filePath, err)
	}

	fmt.Printf("Successfully saved report locally to '%s'\n", filePath) // Log success
	return nil
}

// GetChakaraReportByUniqueCode retrieves the content of the chakra report file
// matching the uniqueCode for a given email from the local filesystem.
// If multiple matches exist, it returns the latest one based on the timestamp.
func (maker *LocalStorageMaker) GetChakaraReportByUniqueCode(email string, uniqueCode string) (string, error) {
	if uniqueCode == "" {
		return "", fmt.Errorf("uniqueCode cannot be empty")
	}
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	// Construct the specific user's directory path
	dirPath := filepath.Join(maker.baseReportDir, email)
	// Define the specific prefix for the unique code within the user's folder
	filePrefix := uniqueCode + "-"

	// Read the directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Handle cases where the directory doesn't exist
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no reports found for email: %s", email)
		}
		return "", fmt.Errorf("failed to read directory '%s': %w", dirPath, err)
	}

	// Filter and collect report filenames matching the uniqueCode prefix
	var matchingFiles []string
	for _, entry := range entries {
		// Skip directories and files that don't match the prefix and extension
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), filePrefix) && strings.HasSuffix(entry.Name(), ReportFileExtension) {
			matchingFiles = append(matchingFiles, entry.Name())
		}
	}

	// Check if any matching files were found
	if len(matchingFiles) == 0 {
		return "", fmt.Errorf("report with unique code '%s' not found for email '%s' locally", uniqueCode, email)
	}

	// If multiple matches, sort descending by filename to get the latest timestamp
	if len(matchingFiles) > 1 {
		sort.Sort(sort.Reverse(sort.StringSlice(matchingFiles)))
	}

	// Get the filename for the requested report (the first one after sorting)
	targetFilename := matchingFiles[0]
	filePath := filepath.Join(dirPath, targetFilename)

	// Read the file content
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read report file '%s': %w", filePath, err)
	}

	fmt.Printf("Successfully retrieved report locally from '%s'\n", filePath) // Log success
	return string(contentBytes), nil
}

// SaveChakaraReportAnswersAsText saves the chakra report answers map as a JSON file
// in the local filesystem under the base answers directory.
// Path: <project_root>/chakara-report-answers/email/uniqueId.txt
func (maker *LocalStorageMaker) SaveChakaraReportAnswersAsText(email string, uniqueId string, answers map[string]string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if uniqueId == "" {
		return fmt.Errorf("uniqueId cannot be empty")
	}
	if len(answers) == 0 {
		// Decide if saving an empty map is an error or should just do nothing
		// return fmt.Errorf("answers map cannot be empty")
		fmt.Printf("Warning: Attempting to save empty answers map for email '%s', uniqueId '%s'. Skipping.\n", email, uniqueId)
		return nil // Or return an error if empty answers are invalid
	}

	// Get the current time and format it using the constant
	currentTime := time.Now()
	timestamp := currentTime.Format(TimeStampFormat) // Added timestamp generation

	// Construct the specific user's directory path within the answers folder
	// This creates: <project_root>/chakara-report-answers/email
	dirPath := filepath.Join(maker.baseAnswersDir, email)

	// Create the directory structure (including parent dirs) if it doesn't exist.
	err := os.MkdirAll(dirPath, FilePermission) // Use the same permission constant
	if err != nil {
		return fmt.Errorf("failed to create answers directory '%s': %w", dirPath, err)
	}

	// Construct the full file path for the answers file
	// Creates: uniqueId.txt
	answersFilename := fmt.Sprintf("%s-%s%s", uniqueId, timestamp, ReportFileExtension) // Updated filename format
	// Creates: <project_root>/chakara-report-answers/email/uniqueId.txt
	filePath := filepath.Join(dirPath, answersFilename)

	// Marshal the answers map into JSON format
	jsonData, err := json.MarshalIndent(answers, "", "  ") // Use MarshalIndent for readability
	if err != nil {
		return fmt.Errorf("failed to marshal answers map to JSON for file '%s': %w", filePath, err)
	}

	// Write the JSON data to the file.
	// 0644 provides standard permissions (rw-r--r--).
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write answers file '%s': %w", filePath, err)
	}

	fmt.Printf("Successfully saved answers locally to '%s'\n", filePath) // Log success
	return nil
}
