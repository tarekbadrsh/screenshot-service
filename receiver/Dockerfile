##################################
# STEP 1 build executable binary #
##################################
FROM golang:latest as builder

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

# install dependencies
RUN go mod download

# copy the content 
COPY . .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Run the binary.
# CMD ["/app/main"]

##############################
# STEP 2 build a final image #
##############################
FROM scratch

# Copy the static executable.
COPY --from=builder /app/main /app/main

# Run the binary.
CMD ["/app/main"]