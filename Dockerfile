FROM golang:latest


LABEL Project="youbei-api"

RUN apt update &&\
    apt install tzdata 

## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

RUN git config --global --add url."https://$REPO_TOKEN@github.com".insteadOf "https://github.com"
RUN go mod vendor

RUN export CGO_CPPFLAGS="-I /usr/local/include"

ENV TZ=UTC

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags $CGO_CPPFLAGS -a -installsuffix cgo -o youbei-api cmd