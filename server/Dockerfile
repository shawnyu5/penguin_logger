FROM golang:1.18.3-alpine3.16 AS build

WORKDIR /server

# copy go config files
# COPY ./go.* ./
# # copy go server file
# COPY ./server.go ./

# COPY ./utils/ ./utils/
# COPY ./coin_products/ ./coin_products/
COPY . .

# RUN go get github.com/joho/godotenv
# RUN go get github.com/patrickmn/go-cache
RUN go build -o server

FROM golang:1.18.3-alpine3.16 AS production

LABEL maintainer="Shawn Yu"
LABEL description="A golang web server for penguin bot"

WORKDIR /server
COPY --from=build /server/server ./server

CMD ["./server"]

EXPOSE 8080
