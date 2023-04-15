# syntax=docker/dockerfile:1

FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /drinkeeapp

FROM build-stage AS run-test-stage
RUN go test -v ./test/*_test.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /drinkeeapp /drinkeeapp

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/drinkeeapp"]