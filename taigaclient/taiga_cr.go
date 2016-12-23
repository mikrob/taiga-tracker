package taigaclient

import (
	"fmt"
	"time"

	"gitlab.botsunit.com/infra/taiga-gitlab/taiga"
)

const (
	//StatusDoneUS represent the done status
	StatusDoneUS = "Done"
	//StatusDoneIssue closed
	StatusDoneIssue = "Closed"
	//StatusReadyUS represent statusReady
	StatusReadyUS = "Ready for test"
	//StatusReadyIssue represent statusReady
	StatusReadyIssue = "Ready for test"
	//StatusInProgressUS in progress
	StatusInProgressUS = "In progress"
	//StatusInProgressIssue in progress
	StatusInProgressIssue = "In progress"
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
		return fromStatus, toStatus
	}
	return "", ""
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
			if fromStatus == StatusReadyUS && toStatus == StatusDoneUS {
				assign := ""
				if us.Assigne == 0 || userList[us.Assigne] == "" {
					assign = NotAssigned
				} else {
					assign = userList[us.Assigne]
				}
				us.AssignedUser = assign
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
			if fromStatus == StatusReadyUS && toStatus == StatusInProgressUS {
				fmt.Println("Rejected US :", us.Subject)
				assign := ""
				if us.Assigne == 0 || userList[us.Assigne] == "" {
					assign = NotAssigned
				} else {
					assign = userList[us.Assigne]
				}
				us.AssignedUser = assign
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
			if fromStatus == StatusReadyIssue && toStatus == StatusDoneIssue {
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
			if fromStatus == StatusReadyIssue && toStatus == StatusInProgressIssue {
				issue.AssignedUser = userList[issue.Assigne]
				issuesRejected = append(issuesRejected, *issue)
			}
		}
	}
	t.IssuesRejectedPerUsers = issuesRejected
}
