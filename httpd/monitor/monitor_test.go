package monitor_test

import (
	"log"
	"net/http"

	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mateeullahmalik/serviceMonitor/httpd/monitor"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/models"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
)

var _ = Describe("Monitor", func() {
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
	})

	It("Tests the service monitoring functionality", func() {
		testServices := []string{"https://www.google.com", "https://www.youtube.com", "https://www.amazon.com"} //99.99% availability

		for _, service := range testServices {
			resp, err := http.Get(service)
			if err != nil {
				log.Printf("%s: Service Unavailable", service)
				continue
			}
			if resp != nil && resp.StatusCode != 502 && resp.StatusCode != 503 {
				//Our core function 'Ping Service' should also be able to ping this service.
				Expect(PingService(h, service)).To(Succeed())
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
