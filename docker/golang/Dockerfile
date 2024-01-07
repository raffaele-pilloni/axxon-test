FROM golang:1.21-alpine

LABEL maintainer="raffaele.pilloni@hotmail.it"

ARG UID=1000
ARG GID=1000

RUN apk update  \
    && apk --no-cache add build-base curl

RUN addgroup -S -g ${GID} axxontest \
    && adduser -S -D -H -u ${UID} axxontest -G axxontest -h /app \
    && mkdir /app \
    && chown -R axxontest:axxontest /app \
    && chmod -R 755 /app

USER axxontest

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.55.2

RUN wget https://github.com/vektra/mockery/releases/download/v2.37.1/mockery_2.37.1_Linux_x86_64.tar.gz && \
    tar -xf mockery_2.37.1_Linux_x86_64.tar.gz && \
    cp ./mockery ${GOPATH}/bin && \
    rm mockery_2.37.1_Linux_x86_64.tar.gz

EXPOSE 8080
CMD ["go","run", "./src/cmd/http-server"]