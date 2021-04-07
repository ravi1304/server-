FROM golang

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /home/koushik/code/tests2/my_servers
COPY ./goserver.go .


# Build the application
RUN go build -o goserver ./goserver.go

#Move to /dist directory as the place for resulting binary folder
# WORKDIR /home/koushik/code/tests2

# Copy binary from build to main folder
#RUN cp /home/koushik/code/tests2/goserver .

EXPOSE 5000

#ARG sname="server"

# Command to run the executable
#ENTRYPOINT ["/home/koushik/code/tests2/my_servers/goserver", ":5000", sname]
ENTRYPOINT ["/home/koushik/code/tests2/my_servers/goserver"]