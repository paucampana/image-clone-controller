# build stage 
FROM golang:1.14

COPY . /src/
WORKDIR /src/
ENV CGO_ENABLED=0
RUN go build -o kubernetes-controller-backup ./src


# final stage  
FROM gliderlabs/alpine:3.8
COPY --from=0 /src/kubernetes-controller-backup /usr/bin/kubernetes-controller-backup
ENTRYPOINT ["/usr/bin/kubernetes-controller-backup"]