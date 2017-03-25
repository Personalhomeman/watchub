package followers

import (
	"context"

	"github.com/google/go-github/github"
)

// Get the list of followers of a given user
func Get(client *github.Client) (result []*github.User, err error) {
	opt := &github.ListOptions{PerPage: 30}

	for {
		followers, nextPage, err := getPage(opt, client)
		if err != nil {
			return result, err
		}
		result = append(result, followers...)
		if opt.Page = nextPage; nextPage == 0 {
			break
		}
	}
	return result, nil
}

func getPage(
	opt *github.ListOptions, client *github.Client,
) (followers []*github.User, nextPage int, err error) {
	ctx := context.Background()
	followers, resp, err := client.Users.ListFollowers(ctx, "", opt)
	if err != nil {
		return followers, 0, err
	}
	return followers, resp.NextPage, err
}
