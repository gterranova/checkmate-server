# Specifies a parent image
FROM golang:1.20 AS builder
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY ./lib ./lib
COPY ./server ./server 

RUN export PATH=/bin:$PATH

WORKDIR /app/server

# Installs Go dependencies
RUN go mod download
 
# Builds your app with optional configuration
RUN go build -o server server.go

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

#COPY ./documents /app/documents
COPY ./www /www 
COPY /server/.env /app/.env

COPY --from=builder /app/server/server /app/server
COPY --from=builder /app/server/fiber.db /app/fiber.db
#COPY --from=builder /bin/sleep /bin/sleep

EXPOSE 4300

USER nonroot:nonroot

# Tells Docker which network port your container listens on
EXPOSE 4300

WORKDIR /app

# Specifies the executable command that runs when the container starts
#CMD [ "/bin/sleep", "infinity" ]
CMD [ "/app/server" ]