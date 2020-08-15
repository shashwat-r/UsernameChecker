package main

import (
	"net/http"
	"strings"
)

// This file contains elementary hacks to identify user names

// getCodeforcesResult checks if the username exists on codeforces by
// checking for the "/submissions/<username>" string in the html response
func getCodeforcesResult(name string) (bool, error) {
	req, err := http.NewRequest("GET", "https://codeforces.com/profile/"+name, nil)
	if err != nil {
		return false, err
	}
	resp, err := makeRequest(req)
	if err != nil {
		return false, err
	}

	identificationStr := "\"/submissions/" + name + "\""
	return strings.Contains(resp, identificationStr), nil
}

// getCodechefResult checks if the username exists on codeforces by
// checking for the "/users/<username>/teams" string in the html response
func getCodechefResult(name string) (bool, error) {
	req, err := http.NewRequest("GET", "https://www.codechef.com/users/"+name, nil)
	if err != nil {
		return false, err
	}
	resp, err := makeRequest(req)
	if err != nil {
		return false, err
	}

	identificationStr := "\"/users/" + name + "/teams\""
	return strings.Contains(resp, identificationStr), nil
}

// getKaggleResult checks if the username exists on codeforces by
// checking for the "/users/<username>/teams" string in the html response
func getKaggleResult(name string) (bool, error) {
	req, err := http.NewRequest("GET", "https://www.kaggle.com/"+name, nil)
	if err != nil {
		return false, err
	}
	resp, err := makeRequest(req)
	if err != nil {
		return false, err
	}

	identificationStr := "\"/" + name + "/activity.json\""
	return strings.Contains(resp, identificationStr), nil
}

func getMediumResult(name string) (bool, error) {
	req, err := http.NewRequest("GET", "https://www.medium.com/@"+name, nil)
	if err != nil {
		return false, err
	}
	resp, err := makeRequest(req)
	if err != nil {
		return false, err
	}

	identificationStr := ">Follow</button>"
	return strings.Contains(resp, identificationStr), nil
}
