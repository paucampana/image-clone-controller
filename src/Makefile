.PHONY: clean image push-image

IMAGE := paucampana/backupoperator

all: backupoperator

backupoperator:
	 go build -o kubernetes-controller-backup ./src

clean:
	go clean ./...


image:
	docker build -t $(IMAGE) .

push-image:
	docker push $(IMAGE)
