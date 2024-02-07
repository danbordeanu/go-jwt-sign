# JWT-SIGN Sample API

A go mock app capable of signing/decoding and validating a token

- verify-signature
- validate-jwt

# Building

After cloning the repo, you can:

## Build locally

__!!!NB!!!__ First, get the swagger/swag Go application

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

```shell
cd $GOPATH/src/go-jwt-sign/src
go get -d -v ./...
swag init && go build main.go && ./main -s -d -p 8080 -r local
```

# Running it locally

If you built it locally, then execute the binary and pass necessary command-line parameters.

```shell
./main [opts]
```

## Examples

### All params:

```shell
swag init --parseDependency && go build main.go && ./main -s -d -r local -p  8080
```

### To enable swagger:

```shell
./main -d -s -p 8080
```

(open browser: http://localhost:8080/swagger/index.html#/default/)

### To enable telemetry:

```shell
./main-r remote/local
```

(__!!!NB!!!__ user remote to print stuff at stdout)

If you built the docker image, you may do the same. Remember to expose the necessary ports.


# Command line parameters

You may specify a number of command-line parameters which change the behavior of the application

| Short | Long | Default | Usable in prod | Description |
|-----|-----|-----|-----|-----|
| -t | --timeout | 60 | Yes | Time to wait for graceful shutdown on SIGTERM/SIGINT in seconds |
| -p | --port | 8080 | Yes | TCP port for the HTTP listener to bind to |
| -s | --swagger | | No | Activate swagger. Do not use this in Production! |
| -d | --devel | | No | Start in development mode. Implies --swagger. Do not use this in Production! |
| -g | --gin-logger| | No | Activate Gin's logger, for debugging. **Warning**: This breaks structured logging. Do not use this in Production! |
| -r | --telemetry| | Yes | Enable telemetry. Values accepted: local (for local telemetry) remote(for jaeger telemetry)|


# Environment variables and options

## Telemetry env vars (jaeger)
For telemetry using jaeger app required jaeger endpoint (if not set, default local host will be used)

```
appConfig.JaegerEngine = utils.EnvOrDefault("JAEGER_ENGINE_NAME", "http://localhost:14268/api/traces")
```

```shell
export JAEGER_ENGINE_NAME=http://localhost:14268/api/traces
```

Starting local jaeger server

```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
```

(open browser http://localhost:16686/)


# API Docs

All endpoints are documented using [swagger](http://localhost:8080/swagger/index.html)

# Swagger for development

First, get the swagger/swag Go application

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

Now, every time you make a change to the Swagger headers, you will need to regenerate the docs

```shell
cd src
swag init
```

If you have a problem with your headers or mappings, you will get an error describing what's wrong. You **must** fix these before committing the code!

**Note:** The docs are regenerated automatically when building the docker image.


# API request sample

## Validate JWT

```shell
curl -X 'POST' \
  'http://localhost:8080/v1/validate-jwt' \
  -H 'accept: text/html' \
  -H 'Content-Type: application/json' \
  -d '{
  "answers": [
    "answer1",
    "answer2"
  ],
  "jwt": "your_jwt_here",
  "questions": [
    "question1",
    "question2"
  ]
}'
```

## Validate Signature

```shell
curl -X 'POST' \
  'http://localhost:8080/v1/verify-signature' \
  -H 'accept: text/html' \
  -H 'Content-Type: application/json' \
  -d '{
  "signature": "test-signature-JonnyBoy",
  "user": "JonnyBoy"
}'
```

