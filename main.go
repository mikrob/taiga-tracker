package main

import (
	"flag"
	"fmt"
	"net/http"
	"taiga_tracker/taigaclient"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	taigaUsername = flag.String("u", "api", "taiga username")
	taigaPassword = flag.String("p", "botsunit8075", "taiga password")
	taigaManager  *taigaclient.TaigaManager
)

func refreshDatas(c *gin.Context) {
	start := time.Now()
	ch := make(chan bool)
	taigaManager.GetMilestoneWithDetails("0.5", "Ufancyme", ch)
	// ready := <-ch
	// if ready {
	taigaManager.MapStoriesPerUsers()
	taigaManager.MapIssuesPerUsers()
	// }
	elapsed := time.Since(start)
	fmt.Printf("Took %s to run", elapsed)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":       "BotsUnit Taiga Tracker",
		"userStories": taigaManager.StoriesPerUsers,
		"issues":      taigaManager.IssuesPerUsers,
		"time":        elapsed,
	})
}

func main() {
	flag.Parse()
	taigaManager = (&taigaclient.TaigaManager{}).NewTaigaManager(taigaUsername, taigaPassword)

	// for user, userStories := range taigaManager.StoriesPerUsers {
	// 	fmt.Println("================================================================================================")
	// 	fmt.Println("User : ", user)
	// 	fmt.Println("Stories : ")
	// 	for _, us := range userStories {
	// 		fmt.Println(us.Subject)
	// 	}
	// 	fmt.Println("================================================================================================")
	// }
	//
	// for user, issueList := range taigaManager.IssuesPerUsers {
	// 	fmt.Println("================================================================================================")
	// 	fmt.Println("User:", user)
	// 	for _, issue := range issueList {
	// 		fmt.Println(issue.Subject)
	// 	}
	//
	// }
	// fmt.Println("Milestone content : ")
	// for _, us := range milestone.UserStoryList {
	// 	fmt.Println(us.Subject)
	// }
	router := gin.Default()
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/fonts", "./fonts")
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin":  "password",
		"mikrob": "password",
		"benj":   "password",
	}))

	router.LoadHTMLGlob("templates/*")

	authorized.GET("/wip", refreshDatas)
	router.Run(":8080")

}
