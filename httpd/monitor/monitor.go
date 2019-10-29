package monitor

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/models"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
	"github.com/tcnksm/go-httpstat"
)

const timeInterval = 50 * time.Second

//MonitorServices monitors a list of services
func MonitorServices(h *storage.StorageHandler, services ...string) {
	for _, service := range services {
		go monitor(h, service)
	}
}

func monitor(h *storage.StorageHandler, service string) {
	for {
		err := PingService(h, service)
		if err != nil {
			log.Println(err.Error())
		}

		time.Sleep(timeInterval)
	}

}

//PingService pings a service & stores Service Ping record in the table
func PingService(h *storage.StorageHandler, service string) error {
	isAvailable := true
	// Create a new HTTP request
	req, err := http.NewRequest("GET", service, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Error Gnerating Request to %s :%s", service, err.Error()))
	}
	// Create a httpstat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by default HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		errors.New(fmt.Sprintf("Error Connecting to %s :%s", service, err.Error()))
	}

	if res == nil || res.StatusCode == 503 || res.StatusCode == 502 {
		isAvailable = false
		saveResult(h, service, time.Now().Unix(), 0, isAvailable)
		return errors.New(fmt.Sprintf("%s: Service Unavailable", service))
	}

	result.End(time.Now())

	responseTime := int((result.TCPConnection + result.TLSHandshake + result.ServerProcessing) / time.Millisecond)
	log.Printf("%s Response Time: %v ms", service, responseTime)

	if err := saveResult(h, service, time.Now().Unix(), responseTime, isAvailable); err != nil {
		return errors.New(fmt.Sprintf("%s: Unable to Save Ping Record into DB: %s", service, err.Error()))
	}

	return nil
}

func saveResult(h *storage.StorageHandler, serviceName string, timeStamp int64, responseDuration int, availablity bool) error {
	servicePing := models.ServicePing{
		Name:         serviceName,
		Timestamp:    timeStamp,
		ResponseTime: responseDuration,
		IsAvailable:  availablity,
	}

	return h.Save(&servicePing)
}
