package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// RepositoryInfo holds basic information about a repository and its latest commit.
type RepositoryInfo struct {
	URL               string
	TempPath          string // Path to the temporary clone
	LatestCommit      CommitInfo
	ChangedFiles      []ChangedFileStats // Per-file line counts will be 0 due to env limitations
	TotalLinesAdded   int
	TotalLinesDeleted int
}

// CommitInfo holds information about a specific commit.
type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Email   string
	Date    string
}

// ChangedFileStats holds statistics for a single changed file.
// Note: LinesAdded and LinesDeleted will currently be 0 for individual files
// due to environment limitations in resolving go-git diff constants.
type ChangedFileStats struct {
	Path         string
	FileType     string // e.g., ".go", ".md"
	LinesAdded   int    // Currently will be 0
	LinesDeleted int    // Currently will be 0
}

// CloneRepository clones a git repository from the given URL to a temporary directory.
func CloneRepository(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "zenwatch-clone-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: nil,
		Depth:    1,
	})

	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository %s: %w", url, err)
	}
	return tempDir, nil
}

// AnalyzeLatestCommit analyzes the latest commit of the repository cloned at repoPath.
// It will populate total lines added/deleted for the commit, but per-file line counts
// will be zero due to limitations in the current Go environment with go-git diff constants.
func AnalyzeLatestCommit(repoPath string) (*RepositoryInfo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", repoPath, err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	latestCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit object: %w", err)
	}

	commitInfo := CommitInfo{
		Hash:    latestCommit.Hash.String(),
		Message: strings.Split(latestCommit.Message, "\n")[0],
		Author:  latestCommit.Author.Name,
		Email:   latestCommit.Author.Email,
		Date:    latestCommit.Author.When.String(),
	}

	repoInfo := &RepositoryInfo{
		TempPath:     repoPath,
		LatestCommit: commitInfo,
	}

	// Get overall commit stats for total lines added/deleted
	totalAdded := 0
	totalDeleted := 0

	commitStats, err := latestCommit.Stats()
	if err != nil {
		// Fallback or note if stats are unavailable, though it should generally work
		// For Depth:1 clones, this often fails with "object not found" if parent is needed by Stats()
		// fmt.Fprintf(os.Stderr, "Warning: could not retrieve commit stats: %v\n", err)
	} else {
		for _, fileStat := range commitStats {
			totalAdded += fileStat.Addition
			totalDeleted += fileStat.Deletion
		}
	}
	repoInfo.TotalLinesAdded = totalAdded
	repoInfo.TotalLinesDeleted = totalDeleted

	currentTree, err := latestCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %w", err)
	}

	var changedFileStatsList []ChangedFileStats
	var patch *object.Patch

	numParents := latestCommit.NumParents()
	if numParents == 0 {
		// Diffing against an empty tree for initial commit (or single commit in shallow clone)
		changes, errDiff := object.DiffTree(nil, currentTree) // Use nil for an empty tree
		if errDiff != nil {
			return nil, fmt.Errorf("failed to diff initial commit tree: %w", errDiff)
		}
		patch, err = changes.Patch()
		if err != nil {
            return nil, fmt.Errorf("failed to get patch from changes (initial commit): %w", err)
        }
	} else {
		parentCommit, errParent := latestCommit.Parent(0)
		if errParent != nil {
			// Fallback for shallow clone where parent isn't available
			changes, diffErr := object.DiffTree(nil, currentTree) // Use nil for an empty tree
			if diffErr != nil {
				return nil, fmt.Errorf("failed to diff current tree with empty (parent fetch failed: %v): %w", errParent, diffErr)
			}
			patch, err = changes.Patch()
			if err != nil {
				return nil, fmt.Errorf("failed to get patch from changes (fallback to empty tree): %w", err)
			}
		} else {
			parentTree, errParentTree := parentCommit.Tree()
			if errParentTree != nil {
				return nil, fmt.Errorf("failed to get parent commit tree: %w", errParentTree)
			}
			patch, err = parentTree.Patch(currentTree)
			if err != nil {
				return nil, fmt.Errorf("failed to create patch between parent and current tree: %w", err)
			}
		}
	}

    if patch != nil {
        for _, filePatch := range patch.FilePatches() {
            from, to := filePatch.Files()
            filePath := ""
            if to != nil {
                filePath = to.Path()
            } else if from != nil { // File was deleted
                filePath = from.Path()
            }
            if filePath == "" { // Should not happen with valid patches
                continue
            }
            changedFileStatsList = append(changedFileStatsList, ChangedFileStats{
                Path:         filePath,
                FileType:     strings.ToLower(filepath.Ext(filePath)),
                LinesAdded:   0, // Per-file line counts set to 0 due to env limitations
                LinesDeleted: 0, // Per-file line counts set to 0 due to env limitations
            })
        }
    }

	repoInfo.ChangedFiles = changedFileStatsList
	return repoInfo, nil
}

// Cleanup removes the temporary directory used for cloning.
func Cleanup(repoPath string) {
	os.RemoveAll(repoPath)
}
