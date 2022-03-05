container:
	sudo docker run --name webscrape -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root123 -d postgres

createdb:
	sudo docker exec -it webscrape createdb --username=root --owner=root webscrapping

dropdb:
	sudo docker exec -it webscrape dropdb webscrapping

gooseup: 
		goose -dir ./db/migration/ -v postgres "postgres://root:root123@localhost:5432/webscrapping?sslmode=disable" up

goosedown: 
		goose -dir ./db/migration/ -v postgres "postgres://root:root123@localhost:5432/webscrapping?sslmode=disable" down

sqlc:
		sqlc generate

dockerImage:
		sudo docker images

dockerbuild:
		sudo docker build -t go_web_scrapping:latest .

dockerremove:
		sudo docker rm go_web_scrapping 
dockerrm:
		sudo docker rmi go_web_scrapping_api
startcontainer:
		sudo docker run --name go_web_scrapping -p 8080:8080  go_web_scrapping:latest
coumposeup:
		sudo docker-compose up
coumposedown:
		sudo docker-compose down

.PHONY:gooseup  goosedown sqlc dockerImage dockerbuild