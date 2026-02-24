FROM golang:1.24.3-alpine AS builder

WORKDIR /department-management-service

RUN apk --no-cache add git bash make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=config/config.yaml
ENV CGO_ENABLED=0

RUN go build --ldflags="-w -s" -o department-management-service ./cmd/department-management-service

FROM alpine AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /department-management-service

COPY --from=builder /department-management-service/config/ /department-management-service/config/
COPY --from=builder /department-management-service/department-management-service /department-management-service/department-management-service

EXPOSE 8080

CMD [ "./department-management-service" ]