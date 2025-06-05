package station

import (
	"encoding/json"
	"errors"
	"mrt-schedule/common/client"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []StationResponse, err error)
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

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return response, err
	}

	var schedule []Schedule
	err = json.Unmarshal(byteResponse, &schedule)

	if err != nil {
		return response, err
	}

	var scheduleSelected Schedule
	err = json.Unmarshal(byteResponse, &scheduleSelected)
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("Station not found")
		return
	}

	response, err = ConvertDataToResponse(scheduleSelected)

	if err != nil {
		return
	}

	return
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	ScheduleLebakBulus := schedule.ScheduleLebakBulus
	ScheduleBundaranHI := schedule.ScheduleBunderanHI

	ScheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(ScheduleLebakBulus)
	if err != nil {
		return
	}
	ScheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(ScheduleBundaranHI)
	if err != nil {
		return
	}

	//convert to response
	for _, item := range ScheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range ScheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}
		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			err = errors.New("Error parsing time" + trimmedTime)
		}
		response = append(response, parsedTime)
	}
	return
}
