// Package rcapi provides an interface for accessing the RC API
package rcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type rcAPI struct {
	Stints []struct {
		EndDate *string `json:"end_date"`
	} `json:"stints"`
}

func IsInBatch(rcid uint32) (bool, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://www.recurse.com/api/v1/people/%d?access_token=%s",
		rcid, os.Getenv("RECURSE_TOKEN")))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("Unexpected StatusCode: " + resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var inBatch rcAPI
	if err := dec.Decode(&inBatch); err != nil {
		return false, err
	}

	now := time.Now()
	for _, stint := range inBatch.Stints {
		// current employee?
		if stint.EndDate == nil {
			return true, nil
		}
		parsed, err := time.Parse("2006-01-02", *stint.EndDate)
		if err != nil {
			return false, err
		}
		if parsed.After(now) {
			return true, nil
		}
	}

	return false, nil
}
