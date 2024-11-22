all: run

build: pipeline

graph: pipeline.png

run: build
	export PROJECT=$$( gcloud config get core/project ) &&\
	./pipeline -project $$PROJECT

pipeline: pipeline.go
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o pipeline

pipeline.dot: pipeline
	./pipeline -runner dot -dot_file pipeline.dot

pipeline.png: pipeline.dot
	dot -Tpng -opipeline.png pipeline.dot

.PHONY: all build graph
