FROM golang:1.21.1-alpine AS builder 

WORKDIR /app 

COPY go.mod go.sum ./  

RUN go mod download && go mod verify 

COPY . .

RUN pwd 
RUN ls -als /app 

# RUN go build -o /my-app  
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /out/app ./cmd/app/main.go

FROM alpine:latest  

COPY --from=builder /out/app /app
COPY --from=builder /app/docker.yml /

RUN pwd 
RUN ls -al /

ENTRYPOINT ["/app"]

