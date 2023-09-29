up:
	docker-compose up -d 

down:
	docker-compose down 

create_network:
	docker network create microservices

# --------------------------------------

POSTGRES_URL := "postgres://admin:password123@127.0.0.1:5432/db?sslmode=disable" 

remove_data:
	sudo rm -r postgres-data 

migrate-init:
	migrate create -ext sql -dir migrations/postgres -seq create_tables
	 
migrate-init-funcs:
	migrate create -ext sql -dir migrations/postgres -seq create_replace_functions 

migrate-list-funcs:
	migrate create -ext sql -dir migrations/postgres -seq create_list_functions 

migrate_up:
	migrate -database ${POSTGRES_URL} -path migrations/postgres up

migrate_down:
	migrate -database ${POSTGRES_URL} -path migrations/postgres down

migrate_drop:
	migrate -database ${POSTGRES_URL} -path migrations/postgres -verbose drop  

#-------------------------

