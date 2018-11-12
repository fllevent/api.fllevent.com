FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/joho/godotenv
RUN go get github.com/appleboy/gin-jwt
RUN go get golang.org/x/crypto/bcrypt
RUN go build -o main . 

EXPOSE 8000:8000

CMD ["/app/main"]