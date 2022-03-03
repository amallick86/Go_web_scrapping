gooseup: 
		goose -dir ./db/migration/ -v postgresql "root:root123/webscrapping?sslmode=disable" up

sqlc:
		sqlc generate

goosedown: 
		goose -dir ./db/migration/ -v postgresql "root:root123/webscrapping?sslmode=disable" down

.PHONY:gooseup  goosedown sqlc