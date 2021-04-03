# Dockerfile Example
# https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888
# Based on this image: https:/hub.docker.com/_/golang/
FROM golang:latest as builder

RUN mkdir -p /app
WORKDIR /app

RUN git config --global url."https://bilalislam:e759a33264bc283b8f448135c3008fb62c6319bb@github.com".insteadOf "https://github.com"
# Force the go compiler to use modules
ENV GO111MODULE on
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

WORKDIR /app
RUN chmod +x /app
COPY . .
# Compile application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grafana-webhook



#Image Diff
#(Not Scratch) 1.23GB
#(Scratch    ) 34.3MB
# <- Second step to build minimal image
FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app .
# Execite application when container is started

ENV ENV_FILE=dev
EXPOSE 80/tcp
CMD ["./grafana-webhook"]

EXPOSE 8080