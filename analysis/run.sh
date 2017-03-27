# Create the go data volume
# docker create -v /go/src/github.com/ --name go golang:1.7
# Mount it into the Go 1.7 container and get all of the git repos.
docker run -it --rm -v $(pwd)/download:/download/ --volumes-from go golang:1.7 go run /download/download.go -in /download/github.csv
# Then run the ver container, against the repositories and export the JSON.