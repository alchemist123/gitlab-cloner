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
const (
	gitlabToken = "<personal access token>"
	gitlabUrl   = "https://gitlab.com"
	groupName   = "Flaviant"
)

var outputDir = filepath.Join("./", groupName)

// Repository struct to unmarshal JSON
type Repository struct {
	Name         string `json:"name"`
	SshUrlToRepo string `json:"ssh_url_to_repo"`
}

// Ensure output directory exists
func ensureOutputDir() {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}
}

// Clone repository
func cloneRepo(repoUrl string, repoName string) {
	fmt.Printf("Cloning %s...\n", repoName)
	cmd := exec.Command("git", "clone", repoUrl, filepath.Join(outputDir, repoName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error cloning %s: %v\n", repoName, err)
	} else {
		fmt.Printf("Cloned %s\n", repoName)
	}
}

// Fetch repositories from a specific group
func fetchRepositories() {
	page := 1
	headers := map[string]string{
		"PRIVATE-TOKEN": gitlabToken,
	}

	var repos []Repository

	fmt.Printf("Fetching repositories from GitLab group: %s...\n", groupName)

	for {
		url := fmt.Sprintf("%s/api/v4/groups/%s/projects?per_page=100&page=%d", gitlabUrl, groupName, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error fetching repositories: %d\n", resp.StatusCode)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		var tempRepos []Repository
		json.Unmarshal(body, &tempRepos)

		if len(tempRepos) == 0 {
			break
		}

		repos = append(repos, tempRepos...)
		page++
	}

	fmt.Printf("Found %d repositories in group %s.\n", len(repos), groupName)

	for _, repo := range repos {
		cloneRepo(repo.SshUrlToRepo, repo.Name)
	}
}

func main() {
	ensureOutputDir()
	fetchRepositories()
}
