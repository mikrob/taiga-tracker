package taigaclient

import (
	"fmt"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers(status string) {
	t.StoriesPerUsers = make(map[string][]taiga.Userstory)
	for _, us := range t.Milestone.UserStoryList {
		if us.Status == usStatusMap[status] {
			assign := ""
			if us.Assigne == 0 || userList[us.Assigne] == "" {
				assign = "Not Assigned"
			} else {
				assign = userList[us.Assigne]
			}
			t.StoriesPerUsers[assign] = append(t.StoriesPerUsers[assign], *us)
		}
	}
}

//MapIssuesPerUsers retrieve issue in progress and map them per users
func (t *TaigaManager) MapIssuesPerUsers(status string) {
	t.IssuesPerUsers = make(map[string][]taiga.Issue)
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		if issue.Status == issuesStatusMap[status] {
			assign := ""
			if issue.Assigne == 0 || userList[issue.Assigne] == "" {
				assign = "Not Assigned"
			} else {
				assign = userList[issue.Assigne]
			}
			t.IssuesPerUsers[assign] = append(t.IssuesPerUsers[assign], *issue)
		}
	}
}
