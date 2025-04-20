package zypher

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	ZypherDir     = ".zypher"
	ObjectsDir    = "objects"
	RefsDir       = "refs"
	HeadsDir      = "heads"
	HeadFile      = "HEAD"
	DefaultBranch = "main"
	BlobPrefix    = "blob"
)

type Commit struct {
	Hash    string            `json:"hash"`
	Message string            `json:"message"`
	Time    time.Time         `json:"time"`
	Files   map[string]string `json:"files"`
	Parent  string            `json:"parent,omitempty"`
	Branch  string            `json:"branch"`
}

type Branch struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type FileStatus struct {
	Path   string
	Status string // "modified", "added", "deleted"
	Hash   string
}

func InitRepository() error {
	// Check if repository is already initialized
	if isZypherRepository() {
		return fmt.Errorf("zypher repository is already initialized")
	}

	PrintLogo()

	dirs := []string{
		ZypherDir,
		filepath.Join(ZypherDir, ObjectsDir),
		filepath.Join(ZypherDir, RefsDir),
		filepath.Join(ZypherDir, RefsDir, HeadsDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize HEAD file to point to main branch
	headContent := fmt.Sprintf("ref: refs/heads/%s", DefaultBranch)
	headPath := filepath.Join(ZypherDir, HeadFile)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	// Create initial branch
	if err := createBranch(DefaultBranch); err != nil {
		return err
	}

	fmt.Println("Initialized empty Zypher repository")
	return nil
}

func createBranch(name string) error {
	branchPath := filepath.Join(ZypherDir, RefsDir, HeadsDir, name)
	return os.WriteFile(branchPath, []byte(""), 0644)
}

func getCurrentBranch() (string, error) {
	headPath := filepath.Join(ZypherDir, HeadFile)
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	// Format: ref: refs/heads/branch-name
	if len(data) < 16 { // Minimum length for "ref: refs/heads/"
		return "", fmt.Errorf("invalid HEAD file format")
	}

	ref := string(data)
	if ref[:16] != "ref: refs/heads/" {
		return "", fmt.Errorf("invalid HEAD reference")
	}

	return ref[16:], nil
}

func getBranchHead(branch string) (string, error) {
	branchPath := filepath.Join(ZypherDir, RefsDir, HeadsDir, branch)
	data, err := os.ReadFile(branchPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func updateBranchHead(branch, hash string) error {
	branchPath := filepath.Join(ZypherDir, RefsDir, HeadsDir, branch)
	return os.WriteFile(branchPath, []byte(hash), 0644)
}

func ShowStatus() error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	PrintLogo()

	// Get current branch
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Get current branch head
	headHash, err := getBranchHead(branch)
	if err != nil {
		return err
	}

	// Get all files in the current directory
	files, err := getAllFiles(".")
	if err != nil {
		return err
	}

	// Get status for each file
	var statuses []FileStatus
	for _, file := range files {
		if file == ZypherDir {
			continue
		}

		status, err := getFileStatus(file, headHash)
		if err != nil {
			return err
		}
		if status != "" {
			statuses = append(statuses, FileStatus{
				Path:   file,
				Status: status,
			})
		}
	}

	// Print status
	if len(statuses) == 0 {
		fmt.Println("No changes to commit")
		return nil
	}

	fmt.Println("Changes to be committed:")
	for _, status := range statuses {
		fmt.Printf("  %s: %s\n", status.Status, status.Path)
	}

	return nil
}

func CreateCommit(message string) error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	// Get current branch
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Get current branch head
	headHash, err := getBranchHead(branch)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Get all files in the current directory
	files, err := getAllFiles(".")
	if err != nil {
		return err
	}

	// Create commit object
	commit := Commit{
		Message: message,
		Time:    time.Now(),
		Files:   make(map[string]string),
		Parent:  headHash,
		Branch:  branch,
	}

	// Calculate hash for each file and store it
	for _, file := range files {
		if file == ZypherDir {
			continue
		}

		hash, err := calculateFileHash(file)
		if err != nil {
			return err
		}
		commit.Files[file] = hash

		// Store file content in objects
		if err := storeFileContent(file, hash); err != nil {
			return err
		}
	}

	// Calculate commit hash
	commitData, err := json.Marshal(commit)
	if err != nil {
		return err
	}
	commit.Hash = calculateHash(commitData)

	// Store commit object
	if err := storeCommit(commit); err != nil {
		return err
	}

	// Update branch head
	if err := updateBranchHead(branch, commit.Hash); err != nil {
		return err
	}

	fmt.Printf("Created commit %s on branch %s\n", commit.Hash, branch)
	return nil
}

func RevertToCommit(hash string) error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	// Load commit
	commit, err := loadCommit(hash)
	if err != nil {
		return err
	}

	// Restore files from commit
	for file, fileHash := range commit.Files {
		if err := restoreFileFromHash(file, fileHash); err != nil {
			return err
		}
	}

	// Update branch head
	if err := updateBranchHead(commit.Branch, hash); err != nil {
		return err
	}

	fmt.Printf("Reverted to commit %s on branch %s\n", hash, commit.Branch)
	return nil
}

func ShowCommitHistory() error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	// Get current branch
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	// Get current branch head
	headHash, err := getBranchHead(branch)
	if err != nil {
		return err
	}

	if headHash == "" {
		fmt.Println("No commits yet")
		return nil
	}

	// Start from HEAD and traverse through parents
	currentHash := headHash
	for currentHash != "" {
		commit, err := loadCommit(currentHash)
		if err != nil {
			return err
		}

		// Print commit information
		fmt.Printf("commit %s (branch: %s)\n", commit.Hash, commit.Branch)
		fmt.Printf("Author: %s\n", "Zypher User")
		fmt.Printf("Date:   %s\n", commit.Time.Format("Mon Jan 2 15:04:05 2006 -0700"))
		fmt.Printf("\n    %s\n\n", commit.Message)

		// Move to parent commit
		currentHash = commit.Parent
	}

	return nil
}

func ListBranches() error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	branchesDir := filepath.Join(ZypherDir, RefsDir, HeadsDir)
	entries, err := os.ReadDir(branchesDir)
	if err != nil {
		return err
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	fmt.Println("Branches:")
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		branchName := entry.Name()
		if branchName == currentBranch {
			fmt.Printf("* %s\n", branchName)
		} else {
			fmt.Printf("  %s\n", branchName)
		}
	}

	return nil
}

func CreateNewBranch(name string) error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	// Check if branch already exists
	branchPath := filepath.Join(ZypherDir, RefsDir, HeadsDir, name)
	if _, err := os.Stat(branchPath); !os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' already exists", name)
	}

	// Get current branch head
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	headHash, err := getBranchHead(currentBranch)
	if err != nil {
		return err
	}

	// Create new branch pointing to the same commit
	if err := os.WriteFile(branchPath, []byte(headHash), 0644); err != nil {
		return err
	}

	fmt.Printf("Created branch '%s'\n", name)
	return nil
}

func CheckoutBranch(checkoutBranch string) error {
	if !isZypherRepository() {
		return fmt.Errorf("not a zypher repository")
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	if currentBranch == checkoutBranch {
		return fmt.Errorf("already on branch '%s'", checkoutBranch)
	}

	// Check if branch exists
	branchPath := filepath.Join(ZypherDir, RefsDir, HeadsDir, checkoutBranch)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", checkoutBranch)
	}

	// Update HEAD to point to the new branch
	headContent := fmt.Sprintf("ref: refs/heads/%s", checkoutBranch)
	headPath := filepath.Join(ZypherDir, HeadFile)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return err
	}

	// Get the commit hash for the branch
	commitHash, err := getBranchHead(checkoutBranch)
	if err != nil {
		return err
	}

	// Restore files to the state of the branch's head commit
	if commitHash == "" {
		return nil
	}

	commit, err := loadCommit(commitHash)
	if err != nil {
		return err
	}

	for file, fileHash := range commit.Files {
		if err := restoreFileFromHash(file, fileHash); err != nil {
			return err
		}
	}

	// remove files that are not in the branch's head commit
	files, err := getAllFiles(".")
	if err != nil {
		return err
	}

	for _, file := range files {
		if _, ok := commit.Files[file]; !ok {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Switched to branch '%s'\n", checkoutBranch)
	return nil
}

// Helper functions

func isZypherRepository() bool {
	_, err := os.Stat(ZypherDir)
	return !os.IsNotExist(err)
}

func getAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip .zypher directory and its contents
		if info.IsDir() && info.Name() == ZypherDir {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func getFileStatus(file string, headHash string) (string, error) {
	// Get current file hash
	currentHash, err := calculateFileHash(file)
	if err != nil {
		return "", err
	}

	// If no HEAD commit, all files are new
	if headHash == "" {
		return "added", nil
	}

	// Load HEAD commit
	headCommit, err := loadCommit(headHash)
	if err != nil {
		return "", err
	}

	// Check if file exists in HEAD commit
	oldHash, exists := headCommit.Files[file]
	if !exists {
		return "added", nil
	}

	if oldHash != currentHash {
		return "modified", nil
	}

	return "", nil
}

func applyChanges(oldHash string, newHash string) error {
	oldCommit, err := loadCommit(oldHash)
	if err != nil {
		return err
	}

	newCommit, err := loadCommit(newHash)
	if err != nil {
		return err
	}

	for file, fileHash := range newCommit.Files {
		fmt.Println("file--//\n", file)
		fmt.Println("fileHash--//\n", fileHash)
	}

	for file, fileHash := range oldCommit.Files {
		fmt.Println("file--//\n", file)
		fmt.Println("fileHash--//\n", fileHash)
	}
	return nil

}
func calculateFileHash(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return calculateHash(data), nil
}

func storeFileContent(file string, hash string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Create blob content: "blob <size>\0<content>"
	header := fmt.Sprintf("%s %d\000", BlobPrefix, len(data))
	blobContent := append([]byte(header), data...)

	// Create object directory if it doesn't exist
	objectDir := filepath.Join(ZypherDir, ObjectsDir, hash[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return err
	}

	// Store blob in objects directory
	objectPath := filepath.Join(objectDir, hash[2:])
	return os.WriteFile(objectPath, blobContent, 0644)
}

func storeCommit(commit Commit) error {
	data, err := json.Marshal(commit)
	if err != nil {
		return err
	}

	// Create commit content: "commit <size>\0<content>"
	header := fmt.Sprintf("commit %d\000", len(data))
	commitContent := append([]byte(header), data...)

	// Create object directory if it doesn't exist
	objectDir := filepath.Join(ZypherDir, ObjectsDir, commit.Hash[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return err
	}

	// Store commit in objects directory
	objectPath := filepath.Join(objectDir, commit.Hash[2:])
	return os.WriteFile(objectPath, commitContent, 0644)
}

func loadCommit(hash string) (Commit, error) {
	// Read commit from objects directory
	objectPath := filepath.Join(ZypherDir, ObjectsDir, hash[:2], hash[2:])

	data, err := os.ReadFile(objectPath)

	if err != nil {
		return Commit{}, err
	}

	// Parse commit header
	headerEnd := 0
	for i, b := range data {
		if b == 0 { // Find null byte
			headerEnd = i
			break
		}
	}
	if headerEnd == 0 {
		return Commit{}, fmt.Errorf("invalid commit format")
	}

	// Extract content after header
	content := data[headerEnd+1:]

	var commit Commit
	if err := json.Unmarshal(content, &commit); err != nil {
		return Commit{}, err
	}

	return commit, nil
}

func restoreFileFromHash(file string, hash string) error {
	// Read blob from objects directory
	objectPath := filepath.Join(ZypherDir, ObjectsDir, hash[:2], hash[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return err
	}

	// Parse blob header
	headerEnd := 0
	for i, b := range data {
		if b == 0 { // Find null byte
			headerEnd = i
			break
		}
	}
	if headerEnd == 0 {
		return fmt.Errorf("invalid blob format")
	}

	// Extract content after header
	content := data[headerEnd+1:]

	// Create directory if it doesn't exist
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(file, content, 0644)
}

func calculateHash(content []byte) string {
	h := sha1.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}
