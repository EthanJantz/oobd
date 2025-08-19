package recurser

import "strconv"
import "strings"
import "os"
import "fmt"
import "bufio"

func Test(){
	fmt.Println("This is OOBD")
}

type Recurser struct {
	uid uint32
	homeDir string
}

func (r *Recurser) GetRcid() (uint32) {
	return r.uid - 1000
}

func  List() ([]Recurser, error) {
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

		r := Recurser{
			uid: uint32(uid), 
			homeDir: homeDir,
		}

		recursers = append(recursers, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return recursers, nil
}
