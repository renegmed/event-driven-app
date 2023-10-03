init-project:
	go mod init job-runner-app

run:
	go run -race ./cmd/app

health:
	curl localhost:4000/health

test_service_worker:
	go test -race -v -run=Worker -count=1 ./test/unit/service

test_service_job:
	go test -race -v -run=Job -count=1 ./test/unit/service

test_repository:
	go test -race -v -run=Repository -count=1 ./test/integration/repository 

test_e2e:
	go test -race -v -count=1 ./test/e2e 

	
#---------------------------------------

WORKER-ID := "bd72aa18-ae9a-4e83-b902-182658e16cba"

WORKER-ID-2 :="9fdd7d7e-7085-478a-b432-9903921e6309"

create_worker:
	curl -XPUT localhost:4001/workers/ -d '{"name": "John", "description": "line leader"}'

create_worker_2:
	curl -XPUT localhost:4001/workers/ -d '{"name": "Peter", "description": "line supervisor"}'


get_worker:
	curl localhost:4001/workers/${WORKER-ID}

get_worker_2:
	curl localhost:4001/workers/${WORKER-ID-2}


delete_worker:
	curl -XDELETE localhost:4001/workers/${WORKER-ID}

#--------------------------------------

launch_job:
	curl -XPOST localhost:4001/jobs/ -d '{"worker_id":${WORKER-ID}, "input": {"test1":"Test 1", "test2":"Test 2"}}'

JOB-ID := "f66c41b7-991f-48d8-98e0-f670d34db2f6"

get_job_status:
	curl -XGET localhost:4001/jobs/${JOB-ID}/statuses

