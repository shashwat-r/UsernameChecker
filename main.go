package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

const (
	Yes   = "yes"
	No    = "no"
	Maybe = "maybe"
	Start = "start"
)

var (
	router = gin.Default()
	client = http.Client{
		Timeout: time.Second * 10,
	}
	results = []result{
		{
			Name:     "codeforces",
			function: getCodeforcesResult,
		},
		{
			Name:     "codechef",
			function: getCodechefResult,
		},
		{
			Name:     "kaggle",
			function: getKaggleResult,
		},
		{
			Name:     "medium",
			function: getMediumResult,
		},
	}
)

type result struct {
	Name     string
	Status   string
	function func(string) (bool, error)
}

func main() {
	var err error
	err = initRoutes()
	if err != nil {
		panic(err)
	}
	err = router.Run()
	if err != nil {
		panic(err)
	}
}

func initRoutes() error {
	router.LoadHTMLGlob("templates/*")
	router.Use(favicon.New("./resources/favicon.ico"))
	path, err := filepath.Abs("./resources")
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if err != nil {
		return err
	}
	router.Static("resources", path)

	router.GET("/", func(c *gin.Context) {
		initialiseResults()
		vals := c.Request.URL.Query()
		username := vals.Get("username")
		if username == "" {
			c.HTML(http.StatusOK, "index.html", gin.H{"title": "Username Checker", "username": "<blank>", "results": results})
			return
		}
		results = getResults(username, results)
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "Username Checker", "username": username, "results": results})
	})
	return nil
}

func initialiseResults() {
	for i, _ := range results {
		results[i].Status = Start
	}
}

func getResults(name string, results []result) []result {
	var resultMap sync.Map
	var wg sync.WaitGroup

	for _, resStruct := range results {
		res := resStruct
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := res.function(name)
			var status string
			if err != nil {
				status = Maybe
			} else if val {
				status = No
			} else {
				status = Yes
			}
			resultMap.Store(res.Name, status)
		}()
	}
	wg.Wait()

	for i := range results {
		resType, ok := resultMap.Load(results[i].Name)
		if !ok {
			results[i].Status = Maybe
		} else {
			results[i].Status = resType.(string)
		}
	}
	return results
}

func makeRequest(req *http.Request) (string, error) {
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	respStr := string(body)
	return respStr, err
}
