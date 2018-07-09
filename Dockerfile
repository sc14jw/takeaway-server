FROM golang
ADD . /go/src/takeaway/takeaway-server

RUN go get "github.com/facebookgo/inject"
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/globalsign/mgo"
RUN go get "gopkg.in/mgo.v2/bson"

RUN go install takeaway/takeaway-server

ENTRYPOINT [ "/go/bin/takeaway-server", "-mongoHost", "mongo", "-mongoPort", "8081"]
EXPOSE 8080