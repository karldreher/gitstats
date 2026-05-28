package github

import (
	"strings"
	"time"
)

const orgReposQuery = `
query($org: String!, $since: GitTimestamp!, $cursor: String) {
  organization(login: $org) {
    repositories(
      first: 50
      after: $cursor
      isArchived: false
      orderBy: {field: PUSHED_AT, direction: DESC}
    ) {
      nodes {
        nameWithOwner
        defaultBranchRef {
          target {
            ... on Commit {
              history(first: 100, since: $since) {
                nodes {
                  message
                  author {
                    name
                    user { login }
                  }
                  committedDate
                }
                pageInfo { hasNextPage endCursor }
              }
            }
          }
        }
      }
      pageInfo { hasNextPage endCursor }
    }
  }
}`

type orgResult struct {
	Organization struct {
		Repositories struct {
			Nodes []struct {
				NameWithOwner    string
				DefaultBranchRef *struct {
					Target struct {
						History struct {
							Nodes []struct {
								Message string
								Author  struct {
									Name string
									User *struct {
										Login string
									}
								}
								CommittedDate string
							}
							PageInfo struct {
								HasNextPage bool
								EndCursor   string
							}
						}
					}
				}
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		}
	}
}

func (c *Client) FetchOrgCommits(since time.Time) ([]CommitData, error) {
	sinceStr := since.UTC().Format(time.RFC3339)
	var allCommits []CommitData
	var cursor *string

	for {
		vars := map[string]any{
			"org":   c.org,
			"since": sinceStr,
		}
		if cursor != nil {
			vars["cursor"] = *cursor
		}

		var result orgResult
		if err := c.query(orgReposQuery, vars, &result); err != nil {
			return nil, err
		}

		for _, repo := range result.Organization.Repositories.Nodes {
			if repo.DefaultBranchRef == nil {
				continue
			}
			for _, node := range repo.DefaultBranchRef.Target.History.Nodes {
				author := node.Author.Name
				if node.Author.User != nil {
					login := node.Author.User.Login
					if strings.Contains(login, "[bot]") {
						continue
					}
					author = login
				}
				allCommits = append(allCommits, CommitData{
					Message: node.Message,
					Repo:    repo.NameWithOwner,
					Author:  author,
				})
			}
			// TODO: paginate repo history if DefaultBranchRef.Target.History.PageInfo.HasNextPage
		}

		if !result.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		next := result.Organization.Repositories.PageInfo.EndCursor
		cursor = &next
	}

	return allCommits, nil
}
