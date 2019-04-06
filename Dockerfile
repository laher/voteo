FROM golang:1.12-alpine

RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /voteo

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN CGO_ENABLED=0 go get -d -v ./...

# Install the package
RUN CGO_ENABLED=0 go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["voteo"]
