package taigaclient

import (
	"fmt"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

//MapStoriesDemos allow to make map of data with stories mapped per users
func (t *TaigaManager) MapStoriesDemos() {
	t.StoriesDemos = make([]taiga.Userstory, 0)
	for _, us := range t.Milestone.UserStoryList {
		if us.Status == usStatusMap[ReadyToTest] {
			assign := ""
			if us.Assigne == 0 || userList[us.Assigne] == "" {
				assign = NotAssigned
			} else {
				assign = userList[us.Assigne]
			}
			us.AssignedUser = assign
			t.StoriesDemos = append(t.StoriesDemos, *us)
		}
	}
}

//MapIssuesDemos retrieve issue in progress and map them per users
func (t *TaigaManager) MapIssuesDemos() {
	t.IssuesDemos = make([]taiga.Issue, 0)
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		if issue.Status == issuesStatusMap[ReadyToTest] {
			assign := ""
			if issue.Assigne == 0 || userList[issue.Assigne] == "" {
				assign = NotAssigned
			} else {
				assign = userList[issue.Assigne]
			}
			issue.AssignedUser = assign
			t.IssuesDemos = append(t.IssuesDemos, *issue)
		}
	}
}
