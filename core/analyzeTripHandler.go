package api

import (
	"errors"
	"journeymaster/libs"
	"journeymaster/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (api *JourneyMasterAPI) analyzeTripHandler(c *gin.Context) {
	// incoming data
	requestBody := models.VehiclePush{}

	// outgoing data
	responseBody := models.VehiclePushAnalysis{}

	// unmarshal check
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		api.SendError(c, "VTA0001", string(err.Error()), http.StatusBadRequest)
		return
	}

	responseBody, err := api.AnalyzeVehicleTrip(requestBody)
	if err != nil {
		api.SendError(c, "VTA0002", err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	api.SendResponse(c, responseBody)
}

func (api *JourneyMasterAPI) AnalyzeVehicleTrip(requestBody models.VehiclePush) (models.VehiclePushAnalysis, error) {
	// break element
	var breaks []*models.Break

	// refuel stops
	var refuelStops []*models.Break

	// city vars
	startCity := string("Unknown City")
	endCity := string("Unknown City")

	// consumption
	latestFuelFillLevel := float32(0.0)
	actualFuelFillLevel := float32(0.0)
	var consumption []float32

	// latest VehiclePushDataPoint element
	latestData := models.VehiclePushDataPoint{}

	if requestBody.BreakThreshold < 0 || requestBody.BreakThreshold > 36000 {
		return models.VehiclePushAnalysis{}, errors.New("BreakThreshold not valid")
	}
	if requestBody.GasTankSize <= 0 || requestBody.GasTankSize > 120 {
		return models.VehiclePushAnalysis{}, errors.New("GasTankSize not valid")
	}
	if len(requestBody.Vin) != 17 {
		return models.VehiclePushAnalysis{}, errors.New("Vin not valid")
	}

	if len(requestBody.Data) == 0 {
		return models.VehiclePushAnalysis{}, errors.New("Data not valid")
	}

	for counter, tripData := range requestBody.Data {
		err := validateTripData(*tripData)
		if err != nil {
			return models.VehiclePushAnalysis{}, err
		}
		if counter == 0 {
			// calculate start
			latestFuelFillLevel = float32(requestBody.GasTankSize) * float32(*tripData.FuelLevel) / 100.0
			// TODO: use go routine for fastness
			startCity, _ = libs.FindCity(float64(tripData.PositionLat), float64(tripData.PositionLong))
		}

		// check if latestData is filled aka more than one VehiclePushDataPoint is in request
		if (models.VehiclePushDataPoint{}) != latestData {
			actualFuelFillLevel = float32(requestBody.GasTankSize) * float32(*tripData.FuelLevel) / 100.0

			// Validate fuel consumption consistency with movement
			if tripData.Odometer > latestData.Odometer && actualFuelFillLevel >= latestFuelFillLevel {
				return models.VehiclePushAnalysis{}, errors.New("Fuel consumption is inconsistent with movement")
			}

			// check if there is a refuel (Adoption: only valid between two data points, if requested expandable for multiple data points)
			if tripData.Odometer == latestData.Odometer {
				if *tripData.FuelLevel > *latestData.FuelLevel {
					refuelBreak := models.Break{
						PositionLat:    latestData.PositionLat,
						PositionLong:   latestData.PositionLong,
						StartTimestamp: latestData.Timestamp,
						EndTimestamp:   tripData.Timestamp,
					}
					refuelStops = append(refuelStops, &refuelBreak)

					latestFuelFillLevel = actualFuelFillLevel
				}
			}

			// check if there is a break (Adoption: only valid between two data points, if requested expandable for multiple data points)
			if tripData.Odometer == latestData.Odometer {
				if tripData.Timestamp-latestData.Timestamp >= int64(requestBody.BreakThreshold) {
					breaks = append(breaks, &models.Break{
						PositionLat:    latestData.PositionLat,
						PositionLong:   latestData.PositionLong,
						StartTimestamp: latestData.Timestamp,
						EndTimestamp:   tripData.Timestamp,
					})
				}
			}

			// insert new consumption point
			if tripData.Odometer > latestData.Odometer {
				consumption = append(consumption, (latestFuelFillLevel-actualFuelFillLevel)/float32(tripData.Odometer-latestData.Odometer)*100)
				latestFuelFillLevel = actualFuelFillLevel
			}
		}

		// store element for next loop
		latestData = *tripData
	}

	// look up destination city
	if len(requestBody.Data) > 1 {
		// TODO: use go routine for fastness
		endCity, _ = libs.FindCity(float64(latestData.PositionLat), float64(latestData.PositionLong))
	}

	// return complete VehiclePushAnalysis element
	return models.VehiclePushAnalysis{
		Vin:         requestBody.Vin,
		Consumption: libs.Average(consumption),
		Departure:   startCity,
		Destination: endCity,
		Breaks:      breaks,
		RefuelStops: refuelStops,
	}, nil
}

func validateTripData(tripData models.VehiclePushDataPoint) error {
	if *tripData.FuelLevel < 0 || *tripData.FuelLevel > 100 {
		return errors.New("FuelLevel not valid")
	}
	if tripData.Odometer < 0 {
		return errors.New("Odometer not valid")
	}
	if tripData.PositionLong < -180 || tripData.PositionLong > 180 {
		return errors.New("PositionLong not valid")
	}
	if tripData.PositionLat < -90 || tripData.PositionLat > 90 {
		return errors.New("PositionLat not valid")
	}
	if tripData.Timestamp > time.Now().Unix() {
		return errors.New("Timestamp not valid")
	}
	return nil
}
