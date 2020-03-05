MIGRATIONDIR := pkg/store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
EXECUTABLE := api

bindata:
	cd ${MIGRATIONDIR} && go-bindata -pkg migrations .

run:
	go run cmd/restapi/*

	# docker run --publish 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
	# go get -u -d github.com/golang-migrate/migrate/v4/cmd/migrate
	# ./migrate.linux-amd64 create -ext sql -dir pkg/store/postgres/migrations -seq create_users_table