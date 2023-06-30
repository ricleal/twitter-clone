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

FROM alpine:3.18.2

COPY --from=builder /app/app /app/

EXPOSE 8888
ENTRYPOINT [ "/app/app" ]
