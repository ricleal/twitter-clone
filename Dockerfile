# syntax=docker/dockerfile:1

###############
# Build stage #
###############
FROM golang:1.20.5-bullseye as builder

WORKDIR /app

# Add go module files
COPY go.mod go.sum ./

# Add source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN go build -v -o /app/app ./cmd/twitter


#################
# Runtime stage #
#################

FROM ubuntu:22.04

ARG API_PORT=8888

COPY --from=builder /app/app /app/

EXPOSE ${API_PORT}

ENTRYPOINT /app/app -port ${API_PORT}
