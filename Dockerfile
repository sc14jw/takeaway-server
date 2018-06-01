FROM golang
ADD . /go/src/takeaway/takeaway-server

RUN go get "github.com/facebookgo/inject"
RUN go get "github.com/gorilla/mux"

RUN go install takeaway/takeaway-server

ENTRYPOINT [ "/go/bin/takeaway-server" ]
EXPOSE 8080