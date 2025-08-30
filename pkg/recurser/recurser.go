// Package recurser provides tools for identifying recursers on the HEAP
package recurser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Recurser struct {
	UserId  uint32
	HomeDir string
}

func (r *Recurser) RcId() uint32 {
	return r.UserId - 1000
}

func List() ([]Recurser, error) {
	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}

	defer passwdFile.Close()

	var recursers []Recurser
	scanner := bufio.NewScanner(passwdFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), ":")
		uidStr := splitLine[2]

		homeDir := splitLine[5]

		uid, err := strconv.ParseUint(uidStr, 10, 32)
		if err != nil {
			return nil, err
		}

		if uid <= 1000 || uid >= 30000 {
			continue
		}

		r := Recurser{
			UserId:  uint32(uid),
			HomeDir: homeDir,
		}

		recursers = append(recursers, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return recursers, nil
}
