package scm

import (
	"context"

	"github.com/google/go-github/v74/github"
	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/rs/zerolog/log"
)

type GithubScm struct {
	*ScmConfig

	token  string
	client *github.Client
}

func NewGithubScm(config *ScmConfig, token string) (*GithubScm, error) {
	var err error

	client := github.NewClient(nil).WithAuthToken(token)

	if config.Server == "" || config.Server == "github.com" {
		config.Server = "github.com"
	} else {
		client, err = client.WithEnterpriseURLs("https://"+config.Server, config.Server)
		if err != nil {
			return nil, err
		}
	}

	return &GithubScm{ScmConfig: config, token: token, client: client}, nil
}

func (gs GithubScm) GitRepoUrl(repository string) string {
	return git.RepoUrl(repository, gs.GitProtocol, gs.Server, gs.token)
}

func (gs GithubScm) ListRepos(ctx context.Context) ([]string, error) {
	repos := make([]string, 0)

	opts := &github.RepositoryListByOrgOptions{
		Type:        gs.Visiblity.String(),
		Sort:        "full_name",
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {
		ghRepos, resp, err := gs.client.Repositories.ListByOrg(ctx, gs.Org, opts)
		if err != nil {
			return nil, err
		}
		for _, repo := range ghRepos {
			if gs.IncludeArchived || !*repo.Archived {
				// Skip empty repositories
				if *repo.Size != 0 {
					repos = append(repos, *repo.FullName)
				}
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	log.Info().
		Int("num", len(repos)).
		Msg("found repositories to scan")

	return repos, nil
}
