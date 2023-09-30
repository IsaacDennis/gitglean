package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/google/go-github/v50/github"
)

type Format string

type Info struct {
	Repos         []*github.Repository
	Contributions []*github.Event
	Format        Format
}

const (
	Markdown Format = "md"
	Org      Format = "org"
)

func Eq[A comparable](a *A, b A) bool {
	return *a == b
}

func Deref[T any](p *T) T {
	return *p
}

var funcMap template.FuncMap = template.FuncMap{
	"humanizeEvent": humanizeEvent,
	"eqstr":         Eq[string],
	"deref":         Deref[string],
}

func createMDLink(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func createOrgLink(text, url string) string {
	return fmt.Sprintf("[[%s][%s]]", url, text)
}

func humanizeEvent(e *github.Event, format Format) string {
	repoName := *e.Repo.Name
	repoURL := "https://github.com/" + repoName
	var link string
	var createLink func(text, url string) string
	switch format {
	case Markdown:
		createLink = createMDLink
	case Org:
		createLink = createOrgLink
	}
	link = createLink(repoName, repoURL)
	payload, _ := e.ParsePayload()
	switch payload := payload.(type) {
	case *github.PushEvent:
		return fmt.Sprintf("⤴️ Pushed %d commits to %s", len(payload.Commits), link)
	case *github.WatchEvent:
		return fmt.Sprintf("⭐ Starred %s", link)
	case *github.PublicEvent:
		return fmt.Sprintf("🔓 Made %s public", link)
	case *github.CommitCommentEvent:
		return fmt.Sprintf("📝 Created a commit comment in %s", link)
	case *github.CreateEvent:
		return fmt.Sprintf("➕ Created a branch/tag in %s", link)
	case *github.DeleteEvent:
		return fmt.Sprintf("➖ Deleted a branch/tag in %s", link)
	case *github.ForkEvent:
		return fmt.Sprintf("🍴 Forked %s", link)
	case *github.GollumEvent:
		return fmt.Sprintf("📃 Created/updated a wiki page in %s", link)
	case *github.PullRequestEvent:
		pr := payload.PullRequest
		number := *pr.Number
		title := *pr.Title
		requestLink := createLink(fmt.Sprintf("#%d %s", number, title), *pr.HTMLURL)
		return fmt.Sprintf("⛙ Opened %s in %s", requestLink, link)

	}
	return fmt.Sprintf(":bangbang: %s (not implemented)", *e.Type)
}

func ListEventsPerformedByUser(gh *github.Client, user string, page, perPage int) ([]*github.Event, *github.Response, error) {
	listOptions := github.ListOptions{
		Page:    page,
		PerPage: perPage,
	}
	return gh.Activity.ListEventsPerformedByUser(context.Background(), user, true, &listOptions)
}

func ListRecentRepositories(gh *github.Client, user string, page, perPage int) ([]*github.Repository, *github.Response, error) {
	repoOptions := github.RepositoryListOptions{
		Visibility: "public",
		Sort:       "created",
		Direction:  "desc",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	return gh.Repositories.List(context.Background(), user, &repoOptions)
}
func main() {
	var name, templatePath, outputPath, format string
	var page, perPage int
	flag.StringVar(&name, "name", "", "username to use in GitHub API requests")
	flag.StringVar(&templatePath, "template", "", "path to template file")
	flag.StringVar(&outputPath, "output", "README", "path to output file")
	flag.StringVar(&format, "format", "md", "export format (md, org)")
	flag.IntVar(&page, "page", 1, "Page of results to retrieve")
	flag.IntVar(&perPage, "perPage", 1, "Number of results to include per page")
	flag.Parse()

	templ, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error while loading template file: %s\n", err.Error())
		os.Exit(1)
	}

	var readme = template.Must(template.New("readme").Funcs(funcMap).Parse(string(templ)))
	gh := github.NewClient(nil)
	events, _, _ := ListEventsPerformedByUser(gh, name, page, perPage)
	repos, _, _ := ListRecentRepositories(gh, name, page, perPage)
	info := Info{repos, events, Format(format)}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error while creating output file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := readme.Execute(outputFile, info); err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}
