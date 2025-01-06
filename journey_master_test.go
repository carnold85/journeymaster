package main

import (
	api "journeymaster/core"
	"journeymaster/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var analyzeVehicleDataProvider = []struct {
	Name           string
	RequestData    models.VehiclePush
	ExpectedResult models.VehiclePushAnalysis
	ExpectedError  string
}{
	{
		Name: "noValidDataElement",
		RequestData: models.VehiclePush{
			BreakThreshold: 1800,
			Data:           nil,
			GasTankSize:    100,
			Vin:            "ABC00000000000001",
		},
		ExpectedResult: models.VehiclePushAnalysis{},
		ExpectedError:  "Data not valid",
	},
	{
		Name: "oneValidDataElement",
		RequestData: models.VehiclePush{
			BreakThreshold: 1800,
			Data: []*models.VehiclePushDataPoint{
				{
					FuelLevel:    func(i int32) *int32 { return &i }(100),
					Odometer:     0,
					PositionLat:  49.013297,
					PositionLong: 8.404205,
					Timestamp:    time.Now().Unix() - 3600,
				},
			},
			GasTankSize: 100,
			Vin:         "ABC00000000000001",
		},
		ExpectedResult: models.VehiclePushAnalysis{
			Breaks:      nil,
			Consumption: 0.0,
			Departure:   "Karlsruhe",
			Destination: "Unknown City",
			RefuelStops: nil,
			Vin:         "ABC00000000000001",
		},
		ExpectedError: "",
	},
	{
		Name: "validTwoElements",
		RequestData: models.VehiclePush{
			BreakThreshold: 1800,
			Data: []*models.VehiclePushDataPoint{
				{
					FuelLevel:    func(i int32) *int32 { return &i }(100),
					Odometer:     0,
					PositionLat:  49.013297,
					PositionLong: 8.404205,
					Timestamp:    time.Now().Unix() - 3600,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(90),
					Odometer:     100,
					PositionLat:  48.885614,
					PositionLong: 8.692087,
					Timestamp:    time.Now().Unix() - 1800,
				},
			},
			GasTankSize: 100,
			Vin:         "ABC00000000000001",
		},
		ExpectedResult: models.VehiclePushAnalysis{
			Breaks:      nil,
			Consumption: 10,
			Departure:   "Karlsruhe",
			Destination: "Pforzheim",
			RefuelStops: nil,
			Vin:         "ABC00000000000001",
		},
		ExpectedError: "",
	},
	{
		Name: "validFourElementsWithBreak",
		RequestData: models.VehiclePush{
			BreakThreshold: 1800,
			Data: []*models.VehiclePushDataPoint{
				{
					FuelLevel:    func(i int32) *int32 { return &i }(100),
					Odometer:     0,
					PositionLat:  49.013297,
					PositionLong: 8.404205,
					Timestamp:    time.Now().Unix() - 7200,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(90),
					Odometer:     100,
					PositionLat:  48.885614,
					PositionLong: 8.692087,
					Timestamp:    time.Now().Unix() - 5400,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(90),
					Odometer:     100,
					PositionLat:  48.885614,
					PositionLong: 8.692087,
					Timestamp:    time.Now().Unix() - 3600,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(80),
					Odometer:     200,
					PositionLat:  48.774822,
					PositionLong: 9.166767,
					Timestamp:    time.Now().Unix(),
				},
			},
			GasTankSize: 100,
			Vin:         "ABC00000000000001",
		},
		ExpectedResult: models.VehiclePushAnalysis{
			Breaks: []*models.Break{
				{
					EndTimestamp:   time.Now().Unix() - 3600,
					PositionLat:    48.885614,
					PositionLong:   8.692087,
					StartTimestamp: time.Now().Unix() - 5400,
				},
			},
			Consumption: 10,
			Departure:   "Karlsruhe",
			Destination: "Stuttgart",
			RefuelStops: nil,
			Vin:         "ABC00000000000001",
		},
		ExpectedError: "",
	},
	{
		Name: "validFourElementsWithRefuel",
		RequestData: models.VehiclePush{
			BreakThreshold: 1800,
			Data: []*models.VehiclePushDataPoint{
				{
					FuelLevel:    func(i int32) *int32 { return &i }(100),
					Odometer:     0,
					PositionLat:  49.013297,
					PositionLong: 8.404205,
					Timestamp:    time.Now().Unix() - 7200,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(90),
					Odometer:     100,
					PositionLat:  48.885614,
					PositionLong: 8.692087,
					Timestamp:    time.Now().Unix() - 5400,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(100),
					Odometer:     100,
					PositionLat:  48.885614,
					PositionLong: 8.692087,
					Timestamp:    time.Now().Unix() - 5000,
				},
				{
					FuelLevel:    func(i int32) *int32 { return &i }(90),
					Odometer:     200,
					PositionLat:  48.774822,
					PositionLong: 9.166767,
					Timestamp:    time.Now().Unix() - 3600,
				},
			},
			GasTankSize: 100,
			Vin:         "ABC00000000000001",
		},
		ExpectedResult: models.VehiclePushAnalysis{
			Breaks:      nil,
			Consumption: 10,
			Departure:   "Karlsruhe",
			Destination: "Stuttgart",
			RefuelStops: []*models.Break{
				{
					EndTimestamp:   time.Now().Unix() - 5400,
					PositionLat:    48.885614,
					PositionLong:   8.692087,
					StartTimestamp: time.Now().Unix() - 5400,
				},
			},
			Vin: "ABC00000000000001",
		},
		ExpectedError: "",
	},
}

func TestAnalyzeVehicleTrip(t *testing.T) {
	vehicle := api.JourneyMasterAPI{}
	var requestData models.VehiclePush

	for _, tt := range analyzeVehicleDataProvider {
		errorOccured := false
		requestData = tt.RequestData
		expectedResult := tt.ExpectedResult
		res, err := vehicle.AnalyzeVehicleTrip(requestData)

		if err != nil {
			if len(tt.ExpectedError) == 0 {
				t.Errorf("foo '%v'", err)
				return
			}
			errorOccured = true
		}

		if errorOccured {
			assert.Contains(t, err.Error(), tt.ExpectedError)
		}

		if !errorOccured {
			assert.Equal(t, expectedResult.Breaks, res.Breaks)
			assert.Equal(t, expectedResult.Consumption, res.Consumption)
			assert.Equal(t, expectedResult.Departure, res.Departure)
			assert.Equal(t, expectedResult.Destination, res.Destination)
			assert.Equal(t, expectedResult.Vin, res.Vin)
		}
	}
}
