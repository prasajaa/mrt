package station

import (
	"encoding/json"
	"mrt-schedule/common/client"
	"net/http"
	"time"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return response, err
	}

	var station []Station
	err = json.Unmarshal(byteResponse, &station)

	for _, item := range station {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return
}
