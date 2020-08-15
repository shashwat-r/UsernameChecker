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
	checkerMap = map[string]func(string) (bool, error){
		"codeforces": getCodeforcesResult,
		"codechef":   getCodechefResult,
	}
)

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
		results := getInitialResults()
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

type result struct {
	Name string
	Type string
}

func getInitialResults() []result {
	var results []result
	for k, _ := range checkerMap {
		results = append(results, result{
			Name: k,
			Type: Start,
		})
	}
	return results
}

func getResults(name string, results []result) []result {
	var resultMap sync.Map
	var wg sync.WaitGroup

	for k, fn := range checkerMap {
		wg.Add(1)
		go func(name string, k string) {
			defer wg.Done()
			val, err := fn(name)
			var res string
			if err != nil {
				res = Maybe
			} else if val {
				res = No
			} else {
				res = Yes
			}
			resultMap.Store(k, res)
		}(name, k)
	}
	wg.Wait()

	for i := range results {
		resType, ok := resultMap.Load(results[i].Name)
		if !ok {
			results[i].Type = Maybe
		} else {
			results[i].Type = resType.(string)
		}
	}
	return results
}

func makeRequest(req *http.Request) (string, error) {
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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
