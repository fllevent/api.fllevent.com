FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/joho/godotenv
RUN go build -o main . 

EXPOSE 8000:8000

CMD ["/app/main"]