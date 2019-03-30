# FROM golang:latest
# ENV GOPATH /home/luoguibin/CompanyCode/go
# ENV APP_ROOT /home/luoguibin/CompanyCode/go/src/SghenApi
# WORKDIR ${APP_ROOT}
# COPY ./ ${APP_ROOT}
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
RUN mkdir conf && touch ./conf/app.conf
RUN mkdir -p ./data/
COPY ./main .
COPY ./conf/app.conf ./conf/app.conf
COPY ./data/sys-peotry-set.json ./data/sys-peotry-set.json
COPY ./data/sys-peotry.json ./data/sys-peotry.json
EXPOSE 8088
ENV SGHENENV prod
ENTRYPOINT ["/app/main"]