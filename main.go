package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/google/go-github/v50/github"
)

type Info struct {
	Repos         []*github.Repository
	Contributions []*github.Event
}

//go:embed readme.template
var templ string
var readme = template.Must(template.New("readme").Parse(templ))

func ListEventsPerformedByUser(gh *github.Client, user string) ([]*github.Event, *github.Response, error){
	listOptions := github.ListOptions {1, 5}
	return gh.Activity.ListEventsPerformedByUser(context.Background(), user, true, &listOptions)
}

func ListRecentRepositories(gh *github.Client, user string) ([]*github.Repository, *github.Response, error) {
	repoOptions := github.RepositoryListOptions{
		Visibility: "public",
		Sort: "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{1, 5},
	}
	return gh.Repositories.List(context.Background(), user, &repoOptions)
}
func main() {
	var name = "YOUR_NAME"
	gh := github.NewClient(nil)
	events, _, _ := ListEventsPerformedByUser(gh, name)
	repos, _, _ := ListRecentRepositories(gh, name)
	info := Info{repos, events}
	if err := readme.Execute(os.Stdout, info); err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}
