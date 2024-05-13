[![CI](https://github.com/meysamhadeli/shop-golang-microservices/actions/workflows/ci.yml/badge.svg?branch=main&style=flat-square)](https://github.com/meysamhadeli/shop-golang-microservices/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.21-61CFDD.svg?style=flat-square)
# 📍 Shop Golang Microservices

> **A practical and imaginary microservices for implementing an infrastructure for up and running distributed system with the latest technology and architecture like Vertical Slice Architecture, Event Driven Architecture, CQRS, Postgres, RabbitMq in Golang.** 🚀

> 💡 **This project is not business-oriented and most of my focus was in the thechnical part for implement a distributed system with a sample project. In this project I implemented some concept in microservices like Messaging, Tracing, Event Driven Architecture, Vertical Slice Architecture, CQRS.**

<a href="https://gitpod.io/#https://github.com/meysamhadeli/shop-golang-microservices"><img alt="Open in Gitpod" src="https://gitpod.io/button/open-in-gitpod.svg"/></a>

# Table of Contents

- [The Goals of This Project](#the-goals-of-this-project)
- [Plan](#plan)
- [Technologies - Libraries](#technologies---libraries)
- [The Domain and Bounded Context - Service Boundary](#the-domain-and-bounded-context---service-boundary)
- [Structure of Project](#structure-of-project)
- [How to Run](#how-to-run)
    - [Docker-Compose](#docker-compose)
    - [Build](#build)
    - [Run](#run)
    - [Test](#test)
- [Documentation Apis](#documentation-apis)  
- [Support](#support)
- [Contribution](#contribution)

## The Goals of This Project

- :sparkle: Using `Vertical Slice Architecture` for `architecture level`.
- :sparkle: Using `Rabbitmq` for `Event Driven Architecture` between our microservices with `streadway/amqp` library.
- :sparkle: Using `gRPC` for `internal communication` between our microservices with `grpc/grpc-go` library.
- :sparkle: Using `CQRS` implementation with `mehdihadeli/Go-MediatR` library.
- :sparkle: Using `Postgres` for `database` in our microservices with `go-gorm/gorm` library.
- :sparkle: Using `go-playground/validator` for `validating input` data in the REST calls.
- :sparkle: Using `OpenTelemetry` for `distributed tracing` top of `Jaeger`.
- :sparkle: Using `OAuth2` for implementation `authentication` and `authorization` with `go-oauth2/oauth2` library.
- :sparkle: Using `Echo framework` for `RESTFul api`.
- :sparkle: Using `Swagger` with `swaggo/swag` library for api documentation.
- :sparkle: Using `uber-go/fx` library for `dependency injection`.
- :sparkle: Using `Viper` for `configuration management`.
- :sparkle: Using `logrus` as a `structured logger`.
- :sparkle: Using `Unit Testing`,`Integration Testing` and `End To End Testing` for testing level.
- :sparkle: Using `Docker-Compose` for our `deployment` mechanism.
- :construction: Using `OpenTelemetry` for `monitoring` top of `Prometteuse` and `Grafana`
- :construction: Using `MongoDB` for read side with `mongo-driver`.
- :construction: Using `Domain Driven Design` (DDD) to implement all `business` processes in microservices.
- :construction: Using `Inbox Pattern` for ensuring message idempotency for receiver and `Exactly once Delivery`.
- :construction: Using `Outbox Pattern` for ensuring no message is lost and there is at `At Least One Delivery`.
  
## Plan

> 🌀This project is a work in progress, new features will be added over time.🌀

I will try to register future goals and additions in the [Issues](https://github.com/meysamhadeli/shop-golang-microservices/issues) section of this repository.

## Technologies - Libraries

- ✔️ **[`labstack/echo`](https://github.com/labstack/echo)** - High performance, minimalist Go web framework
- ✔️ **[`go-gorm/gorm`](https://github.com/go-gorm/gorm)** - The fantastic ORM library for Go, aims to be developer friendly
- ✔️ **[`sirupsen/logrus`](https://github.com/sirupsen/logrus)** - Logrus is a structured logger for Go
- ✔️ **[`streadway/amqp`](https://github.com/streadway/amqp)** - Go RabbitMQ Client Library
- ✔️ **[`spf13/viper`](https://github.com/spf13/viper)** - Go configuration with fangs
- ✔️ **[`swaggo/echo-swagger`](https://github.com/swaggo/echo-swagger)** - Echo middleware to automatically generate RESTful API documentation
- ✔️ **[`mehdihadeli/Go-MediatR`](https://github.com/mehdihadeli/Go-MediatR)** - This package is a Mediator Pattern implementation in Go
- ✔️ **[`go-playground/validator`](https://github.com/go-playground/validator)** - Implements value validations for structs and individual fields based on tags
- ✔️ **[`open-telemetry/opentelemetry-go`](https://github.com/open-telemetry/opentelemetry-go)** - Implementation of OpenTelemetry in Go for distributed-tracing
- ✔️ **[`meysamhadeli/problem-details`](https://github.com/meysamhadeli/problem-details)** - Error Handler for mapping our error to standardized error payload to client
- ✔️ **[`go-resty/resty`](https://github.com/go-resty/resty)** - Simple HTTP and REST client library for Go (inspired by Ruby rest-client)
- ✔️ **[`grpc/grpc-go`](https://github.com/grpc/grpc-go)** - The Go language implementation of gRPC. HTTP/2 based RPC
- ✔️ **[`go-oauth2/oauth2`](https://github.com/go-oauth2/oauth2)** - An open protocol to allow secure authorization in a simple and standard method
- ✔️ **[`stretchr/testify`](https://github.com/stretchr/testify)** - A toolkit with common assertions and mocks that plays nicely with the standard library
- ✔️ **[`uber-go/fx`](https://github.com/uber-go/fx)** - Fx is a dependency injection system for Go
- ✔️ **[`cenkalti/backoff`](https://github.com/cenkalti/backoff)** - This is a Go port of the exponential backoff algorithm
- ✔️ **[`stretchr/testify`](https://github.com/stretchr/testify)** - A toolkit with common assertions and mocks that plays nicely with the standard library
- ✔️ **[`testcontainers/testcontainers-go`](https://github.com/testcontainers/testcontainers-go)** - it's a package to create and clean up container for automated integration/smoke tests
- ✔️ **[`avast/retry-go`](https://github.com/avast/retry-go)** - Simple golang library for retry mechanism
- ✔️ **[`ahmetb/go-linq`](https://github.com/ahmetb/go-linq)** - .NET LINQ capabilities in Go

## The Domain And Bounded Context - Service Boundary

![](./assets/shop-golang-microservices.png)

## Structure of Project

In this project I used [vertical slice architecture](https://jimmybogard.com/vertical-slice-architecture/) and [feature folder structure](http://www.kamilgrzybek.com/design/feature-folders/) to structure my files.

I used [RabbitMQ](https://github.com/rabbitmq) as my MessageBroker for async communication between microservices using the eventual consistency mechanism. 

Microservices are `event based` which means they can publish and/or subscribe to any events occurring in the setup. By using this approach for communicating between services, each microservice does not need to know about the other services or handle errors occurred in other microservices.

I treat each request as a distinct use case or slice, encapsulating and grouping all concerns from front-end to back.
When adding or changing a feature in an application in n-tire architecture, we are typically touching many "layers" in an application. We are changing the user interface, adding fields to models, modifying validation, and so on. Instead of coupling across a layer, we couple vertically along a slice. We `minimize coupling` `between slices`, and `maximize coupling` `in a slice`.

With this approach, each of our vertical slices can decide for itself how to best fulfill the request. New features only add code, we're not changing shared code and worrying about side effects.

<div align="center">
  <img src="./assets/vertical-slice-architecture.png" />
</div>

Instead of grouping related action methods in one endpoint, I used the [REPR pattern](https://deviq.com/design-patterns/repr-design-pattern). Each action gets its own small endpoint, and for communication between our endpoint and handlers, I use [Go-MediatR](https://github.com/mehdihadeli/Go-MediatR) for decouple our endpoint to handlers directly, and it gives use some pipeline behavior for logging, caching, validation and... easily.

The use of the [mediator pattern](https://golangbyexample.com/mediator-design-pattern-golang/) in my endpoints creates clean and thin endpoint. By separating action logic into individual handlers we support the [Single Responsibility Principle](https://en.wikipedia.org/wiki/Single_responsibility_principle) and [Don't Repeat Yourself principles](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself), this is because traditional controllers tend to become bloated with large action methods and several injected `Services` only being used by a few methods.

I used CQRS to decompose my features into small parts that makes our application:

- Maximize performance, scalability and simplicity.
- Easy to maintain and add features to. Changes only affect one command or query, avoiding breaking changes or creating side effects.
- It gives us better separation of concerns and cross-cutting concern (with help of mediatr behavior pipelines), instead of bloated service classes doing many things.

Using the CQRS pattern, we cut each business functionality into vertical slices, for each of these slices we group classes (see [technical folders structure](http://www.kamilgrzybek.com/design/feature-folders)) specific to that feature together (command, handlers, infrastructure, repository, controllers, etc). In our CQRS pattern each command/query handler is a separate slice. This is where you can reduce coupling between layers. Each handler can be a separated code unit, even copy/pasted. Thanks to that, we can tune down the specific method to not follow general conventions (e.g. use custom postgresql query or even different storage). In a traditional layered architecture, when we change the core generic mechanism in one layer, it can impact all methods.


## How to Run

> ### Docker-Compose

Use the command below to run our `infrastructure` with `docker` using the [infrastructure.yaml](./deployments/docker-compose/infrastructure.yaml) file at the `root` of the app:

```bash
docker-compose -f ./deployments/docker-compose/infrastructure.yaml up -d
```
##### Todo
I will add `docker-compsoe` for up and running whole app here in the next...


> ### Build
To `build` each microservice, run this command in the root directory of each microservice where the `go.mod` file is located:

```bash
go build -v ./...
```

> ### Run
To `run` each microservice, run this command in the root of the microservice where `go.mod` is located:

```bash
go run -v ./...
```

> ### Test
To `test` each microservice, run this command in the root directory of the microservice where the `go.mod` file is located:

```bash
go test -v ./...
```

> ### Documentation Apis


Each microservice has a `Swagger OpenAPI`. Browse to `/swagger/index.html` for a list of endpoints.

> Note: For generate Swagger OpenAPI, we need to install `swag cli` with this command below:
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.3
```

As part of API testing, I created the [shop.rest](./shop.rest) file which can be run with the [REST Client](https://github.com/Huachao/vscode-restclient) `VSCode plugin`.


# Support

If you like my work, feel free to:

- ⭐ this repository. And we will be happy together :)

Thanks a bunch for supporting me!

## Contribution

Thanks to all [contributors](https://github.com/meysamhadeli/shop-golang-microservices/graphs/contributors), you're awesome and this wouldn't be possible without you! The goal is to build a categorized community-driven collection of very well-known resources.

Please follow this [contribution guideline](./CONTRIBUTION.md) to submit a pull request or create the issue.

## Project References & Credits

- [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [https://github.com/AleksK1NG/Go-Elasticsearch-RabbitMQ](https://github.com/AleksK1NG/Go-Elasticsearch-RabbitMQ)
- [https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [https://github.com/kgrzybek/modular-monolith-with-ddd](https://github.com/kgrzybek/modular-monolith-with-ddd)
- [https://github.com/jbogard/ContosoUniversityDotNetCore-Pages](https://github.com/jbogard/ContosoUniversityDotNetCore-Pages)

## License
This project is made available under the MIT license. See [LICENSE](https://github.com/meysamhadeli/shop-golang-microservices/blob/main/LICENSE) for details.
