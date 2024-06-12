FROM golang:1.22-alpine as builder

RUN apk --update add build-base

WORKDIR .

ADD go.mod .

RUN go mod download

ADD . .

RUN go run .

FROM alpine

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /bin/

# Copying binaries
COPY --from=builder ./bin/app .

CMD /bin/app
