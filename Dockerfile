FROM golang:latest

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/akulinski/api-token-manager
ADD . /go/src/github.com/akulinski/api-token-manager
WORKDIR /go/src/github.com/akulinski/api-token-manager
RUN dep ensure
RUN go build -o api_token_manager .
EXPOSE 8080
CMD ["/go/src/github.com/akulinski/api-token-manager/api_token_manager"]