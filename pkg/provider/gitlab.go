package provider

import (
	"encoding/json"
	"fmt"
	"github.com/google/triage-party/pkg/constants"
	"github.com/google/triage-party/pkg/models"
	"github.com/xanzy/go-gitlab"
	"log"
	"strconv"
	"time"
)

type GitlabProvider struct {
	client *gitlab.Client
}

func initGitlab(c Config) {
	cl, err := gitlab.NewClient(mustReadToken(*c.GithubTokenFile, constants.GitlabTokenEnvVar))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	gitlabProvider = &GitlabProvider{
		client: cl,
	}
}

func (p *GitlabProvider) getListProjectIssuesOptions(sp models.SearchParams) *gitlab.ListProjectIssuesOptions {
	return &gitlab.ListProjectIssuesOptions{
		ListOptions:  p.getListOptions(sp.IssueListByRepoOptions.ListOptions),
		State:        &sp.IssueListByRepoOptions.State,
		CreatedAfter: &sp.IssueListByRepoOptions.Since,
	}
}

func (p *GitlabProvider) getListOptions(m models.ListOptions) gitlab.ListOptions {
	return gitlab.ListOptions{
		Page:    m.Page,
		PerPage: m.PerPage,
	}
}

func (p *GitlabProvider) getIssues(i []*gitlab.Issue) []*models.Issue {
	r := make([]*models.Issue, len(i))
	for k, v := range i {
		m := models.Issue{}
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			fmt.Println(err)
		}
		r[k] = &m
	}
	return r
}

func (p *GitlabProvider) getRate(i *gitlab.Response) models.Rate {
	l, err := strconv.Atoi(i.Header.Get(constants.GitlabRateLimitHeader))
	if err != nil {
		fmt.Println(err)
	}
	r, err := strconv.Atoi(i.Header.Get(constants.GitlabRateLimitRemainingHeader))
	if err != nil {
		fmt.Println(err)
	}
	rs, err := strconv.Atoi(i.Header.Get(constants.GitlabRateLimitResetHeader))
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Unix(int64(rs), 0)
	return models.Rate{
		Limit:     l,
		Remaining: r,
		Reset:     models.Timestamp{tm},
	}
}

func (p *GitlabProvider) getResponse(i *gitlab.Response) *models.Response {
	r := models.Response{
		NextPage: i.NextPage,
		Rate:     p.getRate(i),
	}
	return &r
}

// https://docs.gitlab.com/ee/api/issues.html#list-project-issues
func (p *GitlabProvider) IssuesListByRepo(sp models.SearchParams) (i []*models.Issue, r *models.Response, err error) {
	opt := p.getListProjectIssuesOptions(sp)
	is, gr, err := p.client.Issues.ListProjectIssues(sp.Repo.Project, opt)
	i = p.getIssues(is)
	r = p.getResponse(gr)
	return
}

func (p *GitlabProvider) IssuesListComments(sp models.SearchParams) ([]*models.IssueComment, *models.Response, error) {

}

func (p *GitlabProvider) IssuesListIssueTimeline(sp models.SearchParams) ([]*models.Timeline, *models.Response, error) {

}

func (p *GitlabProvider) PullRequestsList(sp models.SearchParams) ([]*models.PullRequest, *models.Response, error) {

}

func (p *GitlabProvider) PullRequestsGet(sp models.SearchParams) (*models.PullRequest, *models.Response, error) {

}

func (p *GitlabProvider) PullRequestsListComments(sp models.SearchParams) ([]*models.PullRequestComment, *models.Response, error) {

}

func (p *GitlabProvider) PullRequestsListReviews(sp models.SearchParams) ([]*models.PullRequestReview, *models.Response, error) {

}
