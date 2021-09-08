FROM golang:1.16.6

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/home24-page-analyser

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# ---> Expose
EXPOSE 8000

# Run the executable
CMD ["home24-page-analyser"]


