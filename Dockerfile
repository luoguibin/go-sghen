FROM golang:latest as builder
ENV GOPATH /var/lib/jenkins/workspace-go
ENV APP_ROOT /var/lib/jenkins/workspace-go/src/go-sghen
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest as final
ENV APP_ROOT /var/lib/jenkins/workspace-go/src/go-sghen
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
RUN mkdir conf && touch ./conf/app.conf
RUN mkdir -p ./data/
COPY --from=builder ${APP_ROOT}/main .
COPY --from=builder ${APP_ROOT}/conf/app.conf ./conf/app.conf
EXPOSE 8085
ENV SGHENENV prod
ENTRYPOINT ["/app/main"]