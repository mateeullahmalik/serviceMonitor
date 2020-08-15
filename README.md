# Service Monitor
This is a service which monitors other sevices and keeps track of their response time.
It uses postgres database and exposes several endpoints for retrieveing reports related to response time & availability of services

## Libraries Used
* [Gin](https://github.com/gin-gonic/gin) for Rest APIs
* [Go-pg](https://github.com/go-pg/pg) as an ORM for postgresql database. Although [Gorm](https://github.com/jinzhu/gorm) could have been another great choice as well.
* [Ginkgo](https://github.com/onsi/ginkgo) & [Gomega](https://github.com/onsi/gomega) for behaviour-driven Test suite development
* [go-httpstat](https://github.com/tcnksm/go-httpstat) for retrieving http request statistics.

## Getting Started
1. Right now, it makes use of the default postgres database
2. Make sure to create a database 'test' before executing test cases
3. Main.go defines a list of services which will be monitored by it.
4. Make sure to add your postgres user password in main.go Line # 30

```
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "",
		Addr:     "localhost:5432",
	})
```
## Exposed APIs
nbb
* `/reports/responsetime?time_from=100&time_to=150` //Returns average response time of all the services within given time period. Paramters `time_from` & `time_to` are optional and accept Timestamp in UNIX format (BigInt).

* `/reports/availability?time_from=100&time_to=150` //Returns the availability (Total Hits, Successful & Unsuccessful Hits) within given time period. Paramters `time_from` & `time_to` are optional and accept Timestamp in UNIX format (BigInt).

### Possible Future Improvements

* Define a seperate table for service (ID, Name, URL) and  ServicePing (Takes ServiceID as foreign Key). It working abosluetly fine but this would be a better design decision.
* Add Migrations to automate the task of creating a 'test' named database for test cases
