FROM golang:1.24.2-alpine AS build

WORKDIR /app

RUN adduser -D scratchuser

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/*.go ./

RUN CGO_ENABLED=0 go build -o /aleff-challenge-responder -trimpath -ldflags "-s -w"

FROM scratch

WORKDIR /

USER scratchuser

COPY --from=0 /etc/passwd /etc/passwd
COPY --from=build /aleff-challenge-responder /aleff-challenge-responder

ENTRYPOINT ["/aleff-challenge-responder"]

