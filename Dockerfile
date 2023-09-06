FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl
RUN curl -OL "https://golang.org/dl/go1.21.0.linux-amd64.tar.gz"
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
#RUN apt update && apt-get install -y golang-go sudo
RUN export PATH=$PATH:/usr/local/go/bin

WORKDIR /app_dir
COPY ./ /app_dir/

RUN go build -o main ./cmd/main/main.go
RUN chmod +x main

CMD ["/app_dir/main"]

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
