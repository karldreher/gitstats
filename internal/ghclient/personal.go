package ghclient

import "time"

const personalQuery = `
query($login: String!, $from: DateTime!, $to: DateTime!) {
  user(login: $login) {
    contributionsCollection(from: $from, to: $to) {
      commitContributionsByRepository(maxRepositories: 100) {
        repository { nameWithOwner }
        contributions(first: 100) {
          nodes {
            commit {
              message
              committedDate
            }
          }
          pageInfo { hasNextPage endCursor }
        }
      }
    }
  }
}`

type personalResult struct {
	User struct {
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Repository struct {
					NameWithOwner string
				}
				Contributions struct {
					Nodes []struct {
						Commit struct {
							Message       string
							CommittedDate string
						}
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
				}
			}
		}
	}
}

func (c *Client) fetchPersonalCommits(since time.Time) ([]CommitData, error) {
	var result personalResult
	err := c.query(personalQuery, map[string]any{
		"login": c.user,
		"from":  since.UTC().Format(time.RFC3339),
		"to":    time.Now().UTC().Format(time.RFC3339),
	}, &result)
	if err != nil {
		return nil, err
	}

	var commits []CommitData
	for _, rc := range result.User.ContributionsCollection.CommitContributionsByRepository {
		repo := rc.Repository.NameWithOwner
		for _, node := range rc.Contributions.Nodes {
			commits = append(commits, CommitData{
				Message: node.Commit.Message,
				Repo:    repo,
				Author:  c.user,
			})
		}
		// TODO: paginate repos with >100 commits in the window (rare for 15-min incremental polls)
	}
	return commits, nil
}
