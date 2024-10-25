package scm

import (
	"context"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
)

type GitlabScm struct {
	*ScmConfig

	token  string
	client *gitlab.Client
}

func NewGitlabScm(config *ScmConfig, token string) (*GitlabScm, error) {
	opts := make([]gitlab.ClientOptionFunc, 0, 1)

	if config.Server == "" || config.Server == "gitlab.com" {
		config.Server = "gitlab.com"
	} else {
		opts = append(opts, gitlab.WithBaseURL("https://"+config.Server+"/api/v4"))
	}

	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, err
	}

	return &GitlabScm{ScmConfig: config, token: token, client: client}, nil
}

func (gl GitlabScm) GitRepoUrl(repository string) string {
	return git.RepoUrl(repository, gl.GitProtocol, gl.Server, gl.token)
}

func (gl GitlabScm) ListRepos(ctx context.Context) ([]string, error) {
	repos := make([]string, 0)

	var archived *bool = nil
	if !gl.IncludeArchived {
		b := false
		archived = &b
	}
	includeSubgroups := true
	orderBy := "path"

	opts := &gitlab.ListGroupProjectsOptions{
		Visibility:       gl.Visiblity.Gitlab(),
		Archived:         archived,
		IncludeSubGroups: &includeSubgroups,
		OrderBy:          &orderBy,
	}

	for {
		glRepos, resp, err := gl.client.Groups.ListGroupProjects(gl.Org, opts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		for _, repo := range glRepos {
			repos = append(repos, repo.PathWithNamespace)
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
