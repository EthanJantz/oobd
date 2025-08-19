package recurser

import "strconv"
import "strings"
import "os"
import "fmt"
import "bufio"

func Test(){
	fmt.Println("This is recurser")
}

func List() ([]uint32, error) {
	// read etc passwd
	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}

	defer passwdFile.Close()

	// iterate over each line |> extract uid
	var recursers []uint32
	scanner := bufio.NewScanner(passwdFile)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), ":")
		uidStr := splitLine[2]

		uid, err := strconv.ParseUint(uidStr, 10, 32)
		if err != nil {
			return nil, err
		}

		recursers = append(recursers, uint32(uid))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// return a list of users
	return recursers, nil
}

func Info(){
	
}
