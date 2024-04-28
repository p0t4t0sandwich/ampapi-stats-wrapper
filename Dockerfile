FROM golang:1.22.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ampapi-stats-wrapper .

FROM alpine:edge AS release-stage

WORKDIR /app

COPY ./public ./public
COPY --from=build /app/ampapi-stats-wrapper .

CMD ["/app/ampapi-stats-wrapper"]
