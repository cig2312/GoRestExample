FROM varikin/golang-glide-alpine

WORKDIR /go/src/app

# copy source code to docker container
COPY / /go/src/gorestexample

# Make this as working directory
WORKDIR /go/src/gorestexample

#install vendor dependencies using glide install
RUN glide install

# build go binaries
RUN go build 


EXPOSE 8080