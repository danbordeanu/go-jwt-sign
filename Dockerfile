FROM golang:1.18 as build
WORKDIR /go-jwt-sign//src
ADD go-jwt-sign /go-jwt-sign
ADD go-shared-noversion /go-shared-noversion
RUN go get -d -v ./... \
    && go build -o /jwt-sign
FROM gcr.io/distroless/base
USER 1000
EXPOSE 8080 53835
ENTRYPOINT [ "/jwt-sign" ]
COPY --from=build /jwt-sign /
ADD go-jwt-sign/src/swagger.yaml /swagger.yaml

