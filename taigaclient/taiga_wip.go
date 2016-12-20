package taigaclient

import (
	"fmt"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

//MapStoriesPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesPerUsers(status string) {
	t.StoriesPerUsers = make(map[string][]taiga.Userstory)
	for _, us := range t.Milestone.UserStoryList {
		if us.Assigne != 0 && us.Status == usStatusMap[status] {
			t.StoriesPerUsers[userList[us.Assigne]] = append(t.StoriesPerUsers[userList[us.Assigne]], *us)
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
			t.IssuesPerUsers[userList[issue.Assigne]] = append(t.IssuesPerUsers[userList[issue.Assigne]], *issue)
		}
	}
}
