package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const githubAPIURL = "https://api.github.com"

type Repository struct {
	Name string `json:"name"`
}

type Commit struct {
	Sha string `json:"sha"`
}

type Event struct {
	Type    string `json:"type"`
	Payload struct {
		Commits []Commit `json:"commits"`
	} `json:"payload"`
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func fetchRepos(username string, token string) []Repository {
	url := fmt.Sprintf("%s/users/%s/repos", githubAPIURL, username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching repositories:", err)
		return nil
	}
	defer resp.Body.Close()

	var repos []Repository
	json.NewDecoder(resp.Body).Decode(&repos)
	return repos
}

func fetchCommits(owner, repo, token string) []Commit {
	url := fmt.Sprintf("%s/repos/%s/%s/commits", githubAPIURL, owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching commits:", err)
		return nil
	}
	defer resp.Body.Close()

	var commits []Commit
	json.NewDecoder(resp.Body).Decode(&commits)
	return commits
}

func fetchEvents(owner, repo, token string) []Event {
	url := fmt.Sprintf("%s/repos/%s/%s/events", githubAPIURL, owner, repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return nil
	}
	defer resp.Body.Close()

	var events []Event
	json.NewDecoder(resp.Body).Decode(&events)
	return events
}

func findDeletedCommits(commits []Commit, events []Event) []string {
	currentCommitShas := map[string]bool{}
	for _, commit := range commits {
		currentCommitShas[commit.Sha] = true
	}

	var deletedCommits []string
	for _, event := range events {
		if event.Type == "PushEvent" {
			for _, commit := range event.Payload.Commits {
				if !currentCommitShas[commit.Sha] {
					deletedCommits = append(deletedCommits, commit.Sha)
				}
			}
		}
	}

	return deletedCommits
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Please set your GitHub token in the GITHUB_TOKEN environment variable.")
		return
	}

	username := getUserInput("Enter the GitHub username: ")

	repos := fetchRepos(username, token)
	if len(repos) == 0 {
		fmt.Println("No repositories found or error fetching repositories.")
		return
	}

	fmt.Println("Repositories found:")
	for i, repo := range repos {
		fmt.Printf("[%d] %s\n", i+1, repo.Name)
	}

	repoIndexInput := getUserInput("Choose a repository (enter the number): ")
	repoIndex := 0
	fmt.Sscanf(repoIndexInput, "%d", &repoIndex)
	if repoIndex < 1 || repoIndex > len(repos) {
		fmt.Println("Invalid selection.")
		return
	}
	selectedRepo := repos[repoIndex-1].Name

	commits := fetchCommits(username, selectedRepo, token)
	events := fetchEvents(username, selectedRepo, token)

	deletedCommits := findDeletedCommits(commits, events)

	if len(deletedCommits) > 0 {
		fmt.Println("Deleted commits found:")
		for _, sha := range deletedCommits {
			fmt.Println(sha)
		}
	} else {
		fmt.Println("No deleted commits found.")
	}
}
