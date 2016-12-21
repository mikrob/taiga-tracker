package taigaclient

import (
	"fmt"
	"time"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

// Sort history entries : filter the one that are status modification only, and take the latest modified today
func getLatestHistoryEntryWithStatusModification(historyEntries []*taiga.HistoryEntry) *taiga.HistoryEntry {
	var historyEntryResult *taiga.HistoryEntry
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, historyEntry := range historyEntries {
		if historyEntry.Diff.Status != nil && len(historyEntry.Diff.Status) == 2 {
			historyModificationYear, historyModificationMonth, historyModificationDay := historyEntry.CreatedAt.Date()
			if historyModificationYear == nowYear && historyModificationMonth == nowMonth && historyModificationDay == nowDay {
				if historyEntryResult == nil {
					historyEntryResult = historyEntry
				} else {
					if historyEntry.CreatedAt.After(historyEntryResult.CreatedAt) {
						historyEntryResult = historyEntry
					}
				}
			}
		}
	}
	return historyEntryResult
}

func (t *TaigaManager) retrieveUserStoryHistory(us taiga.Userstory) (string, string) {
	historyEntries, _, err := t.taigaClient.Userstories.GetUserStoryHistory(us.ID)
	if err != nil {
		fmt.Println("Error while retrieving history", err.Error())
	}
	latestHistoryEntry := getLatestHistoryEntryWithStatusModification(historyEntries)
	var fromStatus, toStatus string
	if latestHistoryEntry != nil {
		fromStatus = latestHistoryEntry.HistoryValueList.Status[0]
		toStatus = latestHistoryEntry.HistoryValueList.Status[1]
		return "", ""
	}
	return fromStatus, toStatus
}

func (t *TaigaManager) retrieveIssueHistory(issue taiga.Issue) (string, string) {
	historyEntries, _, err := t.taigaClient.Userstories.GetIssueHistory(issue.ID)
	if err != nil {
		fmt.Println("Error while retrieving history", err.Error())
	}
	latestHistoryEntry := getLatestHistoryEntryWithStatusModification(historyEntries)
	if latestHistoryEntry != nil && len(latestHistoryEntry.HistoryValueList.Status) > 0 {
		fromStatus := latestHistoryEntry.HistoryValueList.Status[0]
		toStatus := latestHistoryEntry.HistoryValueList.Status[1]
		return fromStatus, toStatus
	}
	return "", ""
}

//MapStoriesDonePerUsers make data to storie that have been done today
func (t *TaigaManager) MapStoriesDonePerUsers() {
	var storiesDones []taiga.Userstory
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, us := range t.Milestone.UserStoryList {
		year, month, day := us.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveUserStoryHistory(*us)
			if fromStatus == "Ready for test" && toStatus == "Done" {
				us.AssignedUser = userList[us.Assigne]
				storiesDones = append(storiesDones, *us)
			}
		}
	}
	t.StoriesDonePerUsers = storiesDones
}

//MapStoriesRejectedPerUsers make data to storie that have been done today
func (t *TaigaManager) MapStoriesRejectedPerUsers() {
	var storiesRejected []taiga.Userstory
	nowYear, nowMonth, nowDay := time.Now().Date()
	for _, us := range t.Milestone.UserStoryList {
		year, month, day := us.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveUserStoryHistory(*us)
			fmt.Println(fmt.Sprintf("US : %s, fromStatus : %s, toStatus :%s", us.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "In progress" {
				us.AssignedUser = userList[us.Assigne]
				storiesRejected = append(storiesRejected, *us)
			}
		}
	}
	t.StoriesRejectedPerUsers = storiesRejected
}

//MapIssuesDonePerUsers map issue done per users
func (t *TaigaManager) MapIssuesDonePerUsers() {
	var issuesDone []taiga.Issue
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	nowYear, nowMonth, nowDay := time.Now().Date()

	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		year, month, day := issue.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveIssueHistory(*issue)
			//fmt.Println(fmt.Sprintf("Issue : %s, FromStatus : %s, toStatus : %s", issue.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "Closed" {
				issue.AssignedUser = userList[issue.Assigne]
				issuesDone = append(issuesDone, *issue)
			}
		}
	}
	t.IssuesDonePerUsers = issuesDone
}

//MapIssuesRejectedPerUsers map issue rejected per users
func (t *TaigaManager) MapIssuesRejectedPerUsers() {
	var issuesRejected []taiga.Issue
	issueList, _, err := t.taigaClient.Issues.ListIssues()
	nowYear, nowMonth, nowDay := time.Now().Date()

	if err != nil {
		fmt.Println(fmt.Errorf("Error while retrieving issue list"))
	}
	for _, issue := range issueList {
		year, month, day := issue.LastModified.Date()
		if nowYear == year && nowMonth == month && nowDay == day {
			fromStatus, toStatus := t.retrieveIssueHistory(*issue)
			fmt.Println(fmt.Sprintf("Issue : %s, FromStatus : %s, toStatus : %s", issue.Subject, fromStatus, toStatus))
			if fromStatus == "Ready for test" && toStatus == "In progress" {
				issue.AssignedUser = userList[issue.Assigne]
				issuesRejected = append(issuesRejected, *issue)
			}
		}
	}
	t.IssuesRejectedPerUsers = issuesRejected
}
