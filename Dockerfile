FROM golang
ADD . /go/src/takeaway/takeaway-server

RUN go get "github.com/facebookgo/inject"
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/globalsign/mgo"
RUN go get "gopkg.in/mgo.v2/bson"
RUN go get github.com/rs/cors

RUN go install takeaway/takeaway-server

ARG Username="test"
ARG Password="test"
ARG Host="mongo"

RUN echo "Password = " ${Password}
RUN echo "Username = " ${Username}

ENTRYPOINT /go/bin/takeaway-server -mongoHost "${Host}" -mongoUsername "${Username}" -mongoPassword "${Password}"
EXPOSE 8080