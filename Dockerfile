FROM golang
RUN go get github.com/gorilla/websocket
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/go-sql-driver/mysql
ADD ./*.go /go/src/web-bomberman/
RUN go install web-bomberman/
ADD ./frontend/ /go/bin/frontend/
ADD ./images/ /go/bin/images/

WORKDIR /go/bin/
RUN chmod -R -v u+rwx frontend
ENTRYPOINT /go/bin/web-bomberman  
EXPOSE 2100

