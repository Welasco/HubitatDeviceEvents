# Build the application from source
FROM golang:1.21 AS build-go

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /hubitatedeviceevents
#######################################################################

# Build fronend from source
FROM node:21-bookworm AS build-frontend
WORKDIR /app
COPY . .
RUN cd /app/frontend \
    && npm install \
    && npm run build

# Deploy the application binary into a lean image
FROM alpine AS build-release

WORKDIR /app

COPY --from=build-go /hubitatedeviceevents /app/hubitatedeviceevents
COPY --from=build-frontend /app/frontend/dist /app/public

EXPOSE 3000

#USER nonroot:nonroot

ENTRYPOINT ["/app/hubitatedeviceevents"]