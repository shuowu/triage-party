package provider

import (
	"context"
	"fmt"
	"github.com/google/triage-party/pkg/constants"
	"github.com/google/triage-party/pkg/models"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

type Provider interface {
	IssuesListByRepo(sp models.SearchParams) ([]*models.Issue, *models.Response, error)
	IssuesListComments(sp models.SearchParams) ([]*models.IssueComment, *models.Response, error)
	IssuesListIssueTimeline(sp models.SearchParams) ([]*models.Timeline, *models.Response, error)
	PullRequestsList(sp models.SearchParams) ([]*models.PullRequest, *models.Response, error)
	PullRequestsGet(sp models.SearchParams) (*models.PullRequest, *models.Response, error)
	PullRequestsListComments(sp models.SearchParams) ([]*models.PullRequestComment, *models.Response, error)
	PullRequestsListReviews(sp models.SearchParams) ([]*models.PullRequestReview, *models.Response, error)
}

var (
	githubProvider *GithubProvider
	gitlabProvider *GitlabProvider
)

type Config struct {
	GithubAPIRawURL *string
	GithubTokenFile *string
	GitlabTokenFile *string
}

func InitProviders(ctx context.Context, c Config) {
	initGithub(ctx, c)
	initGitlab(c)
}

func ResolveProviderByHost(providerHost string) Provider {
	switch providerHost {
	case constants.GithubProviderHost:
		return githubProvider
	case constants.GitlabProviderHost:
		return gitlabProvider
	}
	fmt.Println("not existing provider")
	return nil
}

func mustReadToken(path string, env string) string {
	token := os.Getenv(env)
	if path != "" {
		t, err := ioutil.ReadFile(path)
		if err != nil {
			klog.Exitf("unable to read token file: %v", err)
		}
		token = string(t)
		klog.Infof("loaded %d byte github/gitlab token from %s", len(token), path)
	} else {
		klog.Infof("loaded %d byte github/gitlab token from %s", len(token), env)
	}

	token = strings.TrimSpace(token)
	if len(token) < 8 {
		klog.Exitf("github/gitlab token impossibly small: %q", token)
	}
	return token
}
