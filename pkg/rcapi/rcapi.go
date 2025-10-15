// Package rcapi provides an interface for accessing the RC API
package rcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Users struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	ZulipId       int    `json:"zulip_id"`
	CurrentlyAtRc bool   `json:"currently_at_rc"`
}

func GetRecursers() ([]Users, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://www.recurse.com/api/public/cluster_accounts?secret=%s",
		os.Getenv("ACCOUNTS_API_SECRET")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unexpected StatusCode: " + resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var userList []Users
	if err := dec.Decode(&userList); err != nil {
		return nil, err
	}

	return userList, nil
}
