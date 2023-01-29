package executor

import (
	"context"
	"encoding/json"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"io"
	"net/http"
	"rte-etl-routine/authentication"
	"rte-etl-routine/config"
	influxdb_client "rte-etl-routine/influxdb-client"
	"time"
)

type EcoWatExecutor struct{}

type EcoWatt struct {
	Signals []struct {
		GenerationDate time.Time `json:"GenerationFichier"`
		Day            time.Time `json:"jour"`
		ScoreDay       int8      `json:"dvalue"`
		Message        string    `json:"message"`
		HoursValues    []struct {
			Hour      int8 `json:"pas"`
			ScoreHour int8 `json:"hvalue"`
		} `json:"values"`
	} `json:"signals"`
}

func (e EcoWatExecutor) Execute() error {
	influxClient := influxdb_client.NewClient()
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

	var point *write.Point
	writeAPIBlocking := influxClient.WriteAPIBlocking(config.GetEnv().Influxdb.Org, "rte-ecowatt")
	for _, signal := range ecowatt.Signals {
		hoursValues := make(map[string]interface{})
		hoursValues["score_day"] = signal.ScoreDay
		for _, hourValue := range signal.HoursValues {
			hoursValues["hour"] = hourValue.Hour
			hoursValues["score_hour"] = hourValue.ScoreHour
		}
		point = influxdb2.NewPoint(signal.GenerationDate.String(),
			map[string]string{
				"day_date": signal.Day.String(),
				"message":  signal.Message,
			},
			hoursValues,
			time.Now())
		err = writeAPIBlocking.WritePoint(context.Background(), point)
		if err != nil {
			fmt.Println("Unable to write in influx db", err)
			return err
		}
	}
	influxClient.Close()
	fmt.Println("response: ", ecowatt.Signals[0].Message)
	return nil
}
