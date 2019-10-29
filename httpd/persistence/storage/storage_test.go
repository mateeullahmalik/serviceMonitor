package storage_test

import (
	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mateeullahmalik/serviceMonitor/httpd/metadata"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/models"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
	. "github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
)

const (
	totalHitsIndex        = 0
	successfulHitsIndex   = 1
	unsuccessfulHitsIndex = 2
)

var _ = Describe("Storage Test", func() {
	var (
		h  *StorageHandler
		db *pg.DB
	)

	BeforeEach(func() {
		db = pg.Connect(&pg.Options{
			Database: "test", //Make sure Test Database exists
			User:     "postgres",
		})

		h = &storage.StorageHandler{
			Db: db,
		}

		Expect(h.CreateServicePingsTable()).To(Succeed())
	})

	It("Tests the creation of Service Ping Row", func() {
		servicePing := models.ServicePing{
			ID:           171,
			Name:         "https://example.com",
			Timestamp:    1572332712,
			IsAvailable:  true,
			ResponseTime: 500,
		}

		Expect(h.Save(&servicePing)).To(Succeed())
	})

	It("Tests the Average Response Time & Availability Report", func() {

		responseResults := map[string]float64{
			"A": 500,
			"B": 1000,
			"C": 250,
		}

		availabilityResults := map[string][]int{
			"A": []int{6, 4, 2}, //Total hits, Successful Hits, unsuccessful hits
			"B": []int{3, 2, 1},
			"C": []int{2, 2, 0},
		}

		rows := []models.ServicePing{
			//Service 'A'
			{10, "A", 1000, true, 500},
			{10, "A", 1000, false, 0},
			{10, "A", 1250, false, 0},
			{10, "A", 1200, true, 500},
			{10, "A", 1220, true, 1000},
			{10, "A", 1250, true, 200},
			{10, "A", 1270, true, 300},
			{10, "A", 1370, true, 300},
			{10, "A", 1270, false, 0},
			//Service 'B'
			{10, "B", 1000, true, 500},
			{10, "B", 1200, true, 1500},
			{10, "B", 1250, true, 500},
			{10, "B", 1250, false, 0},
			//Service 'C'
			{10, "C", 1200, true, 125},
			{10, "C", 1250, true, 375},
		}

		for i := 0; i < len(rows); i++ {
			servicePing := models.ServicePing{
				ID:           rows[i].ID + i,
				Name:         rows[i].Name,
				Timestamp:    rows[i].Timestamp,
				IsAvailable:  rows[i].IsAvailable,
				ResponseTime: rows[i].ResponseTime,
			}

			Expect(h.Save(&servicePing)).To(Succeed())
		}

		responsedata := []metadata.ServiceAverageResponse{}
		Expect(h.PullAverageResponseTime(&responsedata, 1150, 1350)).To(Succeed())

		for _, averageResponse := range responsedata {
			switch averageResponse.Name {
			case "A":
				Expect(averageResponse.AvrgResponseTime).To(Equal(responseResults["A"]))
			case "B":
				Expect(averageResponse.AvrgResponseTime).To(Equal(responseResults["B"]))
			case "C":
				Expect(averageResponse.AvrgResponseTime).To(Equal(responseResults["C"]))
			}
		}

		availabilitydata := []metadata.ServiceHits{}
		Expect(h.PullServicesAvailability(&availabilitydata, 1150, 1350)).To(Succeed())

		for _, serviceHit := range availabilitydata {
			switch serviceHit.Name {
			case "A":
				Expect(serviceHit.TotalHits).To(Equal(availabilityResults["A"][totalHitsIndex]))
				Expect(serviceHit.SuccessfulHits).To(Equal(availabilityResults["A"][successfulHitsIndex]))
				Expect(serviceHit.UnsuccessfulHits).To(Equal(availabilityResults["A"][unsuccessfulHitsIndex]))
			case "B":
				Expect(serviceHit.TotalHits).To(Equal(availabilityResults["B"][totalHitsIndex]))
				Expect(serviceHit.SuccessfulHits).To(Equal(availabilityResults["B"][successfulHitsIndex]))
				Expect(serviceHit.UnsuccessfulHits).To(Equal(availabilityResults["B"][unsuccessfulHitsIndex]))
			case "C":
				Expect(serviceHit.TotalHits).To(Equal(availabilityResults["C"][totalHitsIndex]))
				Expect(serviceHit.SuccessfulHits).To(Equal(availabilityResults["C"][successfulHitsIndex]))
				Expect(serviceHit.UnsuccessfulHits).To(Equal(availabilityResults["C"][unsuccessfulHitsIndex]))
			}
		}

	})

	AfterEach(func() {
		Expect(db.DropTable(&models.ServicePing{}, &orm.DropTableOptions{
			IfExists: true,
		})).To(Succeed())
		Expect(db.Close()).To(Succeed())
	})
})
