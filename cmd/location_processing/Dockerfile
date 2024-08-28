FROM golang:1.22.0-alpine as builder
WORKDIR /app
RUN apk update && apk add --no-cache gcc musl-dev git
COPY go.* ./
RUN go mod download
COPY . .
WORKDIR /app
RUN cd cmd/location_processing && go build -ldflags '-w -s' -a -o app ./main.go

# Deployment environment
# ----------------------
FROM alpine:3.19.1
WORKDIR /app
RUN chown nobody:nobody /app
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody ./app/cmd/location_processing/app .
COPY --from=builder --chown=nobody:nobody ./app/cmd/location_processing/run.sh .

ENTRYPOINT sh run.sh
