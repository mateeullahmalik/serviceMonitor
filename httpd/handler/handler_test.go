package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mateeullahmalik/serviceMonitor/httpd/handler"
	"github.com/mateeullahmalik/serviceMonitor/httpd/metadata"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/models"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
)

var _ = Describe("Handler", func() {
	var (
		h  *storage.StorageHandler
		db *pg.DB
	)

	BeforeEach(func() {
		db = pg.Connect(&pg.Options{
			Database: "test", //Make sure 'test' Database exists.
			User:     "postgres",
			Password: "iamrobot",
			Addr:     "localhost:5432",
		})

		h = &storage.StorageHandler{
			Db: db,
		}

		Expect(h.CreateServicePingsTable()).To(Succeed())
		servicePing := models.ServicePing{
			ID:           171,
			Name:         "A",
			Timestamp:    1572332712,
			IsAvailable:  true,
			ResponseTime: 500,
		}

		Expect(h.Save(&servicePing)).To(Succeed())
	})

	It("Tests the Response Time Report End Point", func() {

		router := SetupRouter(h)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/reports/responsetime?time_from=1572332710", nil)
		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))

		response := []metadata.ServiceAverageResponse{}
		Expect(json.Unmarshal([]byte(w.Body.String()), &response)).To(Succeed())
		for _, data := range response {
			if data.Name == "A" {
				Expect(data.AvrgResponseTime).To(Equal(500.00))
			}
		}

	})

	It("Tests the Availability Report End Point", func() {
		router := SetupRouter(h)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/reports/availability", nil)
		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))

		response := []metadata.ServiceHits{}
		Expect(json.Unmarshal([]byte(w.Body.String()), &response)).To(Succeed())

		for _, data := range response {
			if data.Name == "A" {
				Expect(data.TotalHits).To(Equal(1))
				Expect(data.SuccessfulHits).To(Equal(1))
				Expect(data.UnsuccessfulHits).To(Equal(0))
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
