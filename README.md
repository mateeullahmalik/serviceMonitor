# Service Monitor
This is a service which monitors other sevices and keeps track of their response time.
It uses postgres database and exposes several endpoints for retrieveing reports related to response time & availability of services

## Libraries Used
* [Gin](https://github.com/gin-gonic/gin): For Rest APIs
* Go-pg: As an ORM for postgresql database
* Ginkgo & Gomega: For Test cases
