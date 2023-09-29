init-project:
	go mod init job-runner-app

run:
	go run -race ./cmd/app

health:
	curl localhost:4000/health

test_service:
	go test -race -v ./test/unit/service

WORKER-ID := "1b48a324-da9f-40cb-b22a-569ce37a4e5b"

create_worker:
	curl -XPUT localhost:4000/workers/ -d '{"name": "John", "description": "line leader"}'

get_worker:
	curl localhost:4000/workers/${WORKER-ID}

delete_worker:
	curl -XDELETE localhost:4000/workers/${WORKER-ID}
