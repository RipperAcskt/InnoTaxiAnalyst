# InnoTaxiDriver API

This is a microservice for analyse service.

## Installation

To install and run the service, follow these steps:

Install Go on your system if you haven't already done so.
Clone the repository to your local machine.
Run the following command in the root directory of the project:

    go run ./cmd/main.go

Also you can run project using docker-compose.

## Code Description

# Project structure

The code for the microservice is organized into several packages:

- cmd/main.go contains the main function for the service.
- internal/app/app.go contains functions which sets up the API routes and starts the server.
- internal/models/ contains the data models for the application. In this case, there is only one model - Driver.
- internal/repo/ contains the repository implementation for working with the databases.
- internal/services/ contains the business logic services for the application.
- internal/handlers/ contains the API request handlers for the application. Service provides handlers for registartion and auth driver, also handlers for working with driver's profile.

---

- Business logic services are located in the internal/services/ package. Service uses repository layer to get data.

### Conclusion

This microservice demonstrates a simple way to analyse and follow after service using Go. The code is organized into packages, making it easy to maintain and extend. The `AnalystService` and `AnalystHandler` objects provide the business logic and API endpoints, respectively.