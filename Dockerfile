# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /hubitatedeviceevents

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage

WORKDIR /app

COPY --from=build-stage /hubitatedeviceevents /app/hubitatedeviceevents

EXPOSE 3000

#USER nonroot:nonroot

ENTRYPOINT ["/app/hubitatedeviceevents"]