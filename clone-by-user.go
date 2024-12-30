package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Configuration
const gitlabToken = "<personal access token>" // Use valid token
const gitlabUrl = "https://gitlab.com"        // No trailing slash
const outputDir = "gitlab-repos"

// Repo represents a GitLab project structure
type Repo struct {
	Name         string `json:"name"`
	SSHURLToRepo string `json:"ssh_url_to_repo"`
}

func main() {
	// Ensure the output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	fetchRepositories()
}

func fetchRepositories() {
	page := 1
	repos := []Repo{}
	fmt.Println("Fetching repositories from GitLab...")

	for {
		url := fmt.Sprintf("%s/api/v4/projects?membership=true&per_page=100&page=%d", gitlabUrl, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("PRIVATE-TOKEN", gitlabToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch repos. Status: %d\n", resp.StatusCode)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		var result []Repo
		json.Unmarshal(body, &result)

		if len(result) == 0 {
			break
		}

		repos = append(repos, result...)
		page++
	}

	fmt.Printf("Found %d repositories.\n", len(repos))
	for _, repo := range repos {
		cloneRepo(repo.SSHURLToRepo, repo.Name)
	}
}

func cloneRepo(repoUrl, repoName string) {
	fmt.Printf("Cloning %s...\n", repoName)
	cloneCmd := exec.Command("git", "clone", repoUrl, filepath.Join(outputDir, repoName))
	output, err := cloneCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error cloning %s: %v\nOutput: %s\n", repoName, err, output)
	} else {
		fmt.Printf("Cloned %s\n", repoName)
	}
}
