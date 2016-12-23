package taigaclient

import (
	"fmt"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

//MapStoriesWipPerUsers allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesWipPerUsers() {
	t.StoriesPerUsers = make(map[string][]taiga.Userstory)
	for _, us := range t.Milestone.UserStoryList {
		if us.Status == usStatusMap[InProgress] {
			assign := ""
			if us.Assigne == 0 || userList[us.Assigne] == "" {
				assign = NotAssigned
			} else {
				assign = userList[us.Assigne]
			}
			t.StoriesPerUsers[assign] = append(t.StoriesPerUsers[assign], *us)
		}
	}
}

//MapIssuesWipPerUsers retrieve issue in progress and map them per users
func (t *TaigaManager) MapIssuesWipPerUsers() {
	t.IssuesPerUsers = make(map[string][]taiga.Issue)
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	if err != nil {
		fmt.Println("Error while retrieving issue list", err.Error())
	}
	for _, issue := range issueList {
		if issue.Status == issuesStatusMap[InProgress] {
			assign := ""
			if issue.Assigne == 0 || userList[issue.Assigne] == "" {
				assign = NotAssigned
			} else {
				assign = userList[issue.Assigne]
			}
			t.IssuesPerUsers[assign] = append(t.IssuesPerUsers[assign], *issue)
		}
	}
}
