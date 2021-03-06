package lib

import "gitlab.botsunit.com/infra/taiga-gitlab/taiga"

//AllPredicate allow to filter a collection with a given function predicate
func AllPredicate(vs []bool, f func(bool) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

//AllTaskSameStatus allow to verify with a given func that all task have same status
func AllTaskSameStatus(taskList []*taiga.Task, f func(*taiga.Task) bool) bool {
	for _, v := range taskList {
		if !f(v) {
			return false
		}
	}
	return true
}

//Index Returns the first index of the target string t, or -1 if no match is found.
func Index(vs []interface{}, t interface{}) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include Returns true if the target string t is in the slice.
func Include(vs []interface{}, t interface{}) bool {
	return Index(vs, t) >= 0
}

//Any Returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

//All Returns true if all of the strings in the slice satisfy the predicate f.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

//FilterMilestone Returns a new slice containing all strings in the slice that satisfy the predicate f.
func FilterMilestone(vs []*taiga.Milestone, f func(*taiga.Milestone) bool) []*taiga.Milestone {
	vsf := make([]*taiga.Milestone, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Map Returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
