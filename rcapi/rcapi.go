package rcapi
import "fmt"
import "io"
import "errors"
import "net/http"
import "encoding/json"
type rcApi struct {
	Stint []struct {
		EndDate string `json:"end_date"`
	} `json:"stint"`
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

	fmt.Printf("%s", InBatch)
	return true, nil
}
