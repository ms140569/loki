FROM golang:1.11
RUN apt-get update
RUN apt install -y pandoc
# RUN apt install -y protobuf-compiler
# RUN apt install -y golang-goprotobuf-dev
ADD . /usr/local/loki
RUN echo 'alias l="ls -ltr"' >> ~/.bashrc