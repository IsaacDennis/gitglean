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
var funcMap template.FuncMap = template.FuncMap{
	"humanizeEvent": humanizeEvent,
}
var readme = template.Must(template.New("readme").Funcs(funcMap).Parse(templ))

func createMDLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}
func humanizeEvent(e *github.Event) string {
	repoName := *e.Repo.Name
	repoURL := "https://github.com/" + repoName
	link := createMDLink(repoName, repoURL)
	payload, _ := e.ParsePayload()
	switch payload := payload.(type) {
	case *github.PushEvent:
		return fmt.Sprintf(":arrow_heading_up: Pushed %d commits to %s", len(payload.Commits), link)
	case *github.WatchEvent:
		return fmt.Sprintf(":star: Starred %s", link)
	case *github.PublicEvent:
		return fmt.Sprintf(":unlock: Made %s public", link)
	case *github.CommitCommentEvent:
		return fmt.Sprintf(":memo: Created a commit comment in %s", link)
	case *github.CreateEvent:
		return fmt.Sprintf(":heavy_plus_sign: Created a branch/tag in %s", link)
	case *github.DeleteEvent:
		return fmt.Sprintf(":heavy_minus_sign: Deleted a branch/tag in %s", link)
	case *github.ForkEvent:
		return fmt.Sprintf(":fork_and_knife: Forked %s", link)
	case *github.GollumEvent:
		return fmt.Sprintf(":page_with_curl: Created/updated a wiki page in %s", link)
	}
	return fmt.Sprintf(":bangbang: %s (not implemented)", *e.Type)
}

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
