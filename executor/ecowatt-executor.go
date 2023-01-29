package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rte-etl-routine/authentication"
	"rte-etl-routine/config"
	"time"
)

type EcoWatExecutor struct{}

type EcoWatt struct {
	Signals []struct {
		GenerationFichier time.Time `json:"GenerationFichier"`
		Jour              time.Time `json:"jour"`
		Dvalue            int       `json:"dvalue"`
		Message           string    `json:"message"`
		Values            []struct {
			Pas    int `json:"pas"`
			Hvalue int `json:"hvalue"`
		} `json:"values"`
	} `json:"signals"`
}

func (e EcoWatExecutor) Execute() error {
	client := &http.Client{}
	token, err := authentication.New().GetToken()
	if err != nil {
		fmt.Println("Unable to get authentication token", err)
		return err
	}

	req, err := http.NewRequest("GET", config.GetEnv().Ecowatt.URI, nil)
	if err != nil {
		fmt.Println("Unable to create request for ecowatt", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Unable to request ecowatt", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Unable to close Body reader ", err)
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body) // response body is []byte
	if err != nil {
		return err
	}

	var ecowatt EcoWatt
	if err := json.Unmarshal(body, &ecowatt); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON from authentication")
	}

	log.Println("response: ", ecowatt.Signals[0].Message)
	return nil
}
