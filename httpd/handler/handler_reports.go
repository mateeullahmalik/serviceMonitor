package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mateeullahmalik/serviceMonitor/httpd/metadata"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
)

func SetupRouter(h *storage.StorageHandler) *gin.Engine {
	router := gin.Default()
	//This endpoint returns Average Response Time report of services b/w a time period (Unix Time Format)
	router.GET("/reports/responsetime", averageResponseGET(h))
	//This endpoint returns Availability report of services b/w a time period (Unix Time Format)
	router.GET("/reports/availability", availabilityGET(h))

	return router
}

//averageResponseGET returns average response time of services
func averageResponseGET(h *storage.StorageHandler) gin.HandlerFunc {

	return func(c *gin.Context) {
		result := []metadata.ServiceAverageResponse{}

		timeFrom, timeTo, err := getTimeDuration(c)
		if err != nil {
			c.String(http.StatusInternalServerError, " -- %s", err.Error())
		}

		err = h.PullAverageResponseTime(&result, timeFrom, timeTo)
		if err != nil {
			c.String(http.StatusInternalServerError, " -- %s ", err.Error())
		}

		c.JSON(http.StatusOK, result)
	}

}

//availabilityGET returns services availability report
func availabilityGET(h *storage.StorageHandler) gin.HandlerFunc {

	return func(c *gin.Context) {
		result := []metadata.ServiceHits{}

		timeFrom, timeTo, err := getTimeDuration(c)
		if err != nil {
			c.String(http.StatusInternalServerError, " -- %s", err.Error())
		}

		err = h.PullServicesAvailability(&result, timeFrom, timeTo)
		if err != nil {
			c.String(http.StatusInternalServerError, " -- %s", err.Error())
		}

		c.JSON(http.StatusOK, result)
	}

}

//Helper Method that gets 'time from' & 'time to' parameters
func getTimeDuration(c *gin.Context) (timeFrom int64, timeTo int64, err error) {
	timeFromParam := strings.TrimSpace(c.Query("time_from"))
	timeToParam := strings.TrimSpace(c.Query("time_to"))

	parseTimeParam := func(param string, time *int64, defaultVal int64) error {
		if param == "" {
			*time = defaultVal
		} else {
			timeVal, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			}

			*time = timeVal
		}
		return nil
	}

	if err := parseTimeParam(timeFromParam, &timeFrom, 0); err != nil {
		return 0, 0, errors.New("From: Parse Time From Parameter: " + err.Error())
	}

	if err := parseTimeParam(timeToParam, &timeTo, time.Now().Unix()); err != nil {
		return 0, 0, errors.New("From: Parse Time To Parameter: " + err.Error())
	}

	return timeFrom, timeTo, nil
}
