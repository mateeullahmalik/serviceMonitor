package storage

import (
	"errors"

	orm "github.com/go-pg/pg/v9/orm"
	"github.com/mateeullahmalik/serviceMonitor/httpd/metadata"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/models"
)

//CreateServicePingsTable creates a table 'service_ping' if it doesn't already exist
func (h *StorageHandler) CreateServicePingsTable() error {
	options := &orm.CreateTableOptions{
		IfNotExists: true,
	}

	err := h.Db.CreateTable(&models.ServicePing{}, options)
	if err != nil {
		return errors.New("Error Creating Service Pings Table: " + err.Error())
	}

	return nil
}

//PullAverageResponseTime Fetches average response time of all the services
func (h *StorageHandler) PullAverageResponseTime(res *[]metadata.ServiceAverageResponse, timeFrom int64, timeTo int64) error {

	return h.Db.Model((*models.ServicePing)(nil)).
		Column("name").
		Where("is_available IS NOT NULL").
		ColumnExpr("AVG(response_time) AS avrg_response_time").
		Where("ping_timestamp > ?", timeFrom).
		Where("ping_timestamp < ?", timeTo).
		Group("name").
		Order("avrg_response_time").
		Select(res)
}

//PullServicesAvailability Fetches Availability report of all the services
func (h *StorageHandler) PullServicesAvailability(res *[]metadata.ServiceHits, timeFrom int64, timeTo int64) error {

	return h.Db.Model((*models.ServicePing)(nil)).
		Column("name").
		ColumnExpr("COUNT(name) AS total_hits").
		ColumnExpr("Count(is_available) AS successful_hits").
		ColumnExpr("COUNT(name) - Count(is_available) as unsuccessful_hits").
		Where("ping_timestamp > ?", timeFrom).
		Where("ping_timestamp < ?", timeTo).
		Group("name").
		Order("successful_hits").
		Select(res)
}
