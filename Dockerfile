FROM golang:1.10-alpine
RUN apk add --update ca-certificates \
    && apk add curl git coreutils make \
    && rm /var/cache/apk/*
COPY azure-nanny /usr/local/bin/azure-nanny
ENTRYPOINT [ "azure-nanny" ]
