package rcapi
import "fmt"
import "time"
import "io"
import "errors"
import "net/http"
import "encoding/json"
type rcApi struct {
	Stints []struct {
		EndDate string `json:"end_date"`
	} `json:"stints"`
}

func IsInBatch(rcid uint32) (bool, error) {
	resp, err := http.Get(BASE_URL)
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.Status)
	if resp.StatusCode > 299 {
		return false, errors.New("Unexpected StatusCode: " + resp.Status)
	}
	if err != nil {
		return false, err
	}

	var InBatch rcApi
	err = json.Unmarshal(body, &InBatch)
	if err != nil {
		fmt.Println("error: ", err)
		return false, err
	}
	
	LastEndDate := InBatch.Stints[len(InBatch.Stints)-1].EndDate
	fmt.Println(LastEndDate)

	const YYYYMMDD = "2006-01-02"
	t, err := time.Parse(YYYYMMDD, LastEndDate)
	if err != nil {
		fmt.Println("error: ", err)
		return false, err
	}

	curDate := time.Now()
	OutOfBatch := t.Before(curDate)

	return OutOfBatch, nil
}

