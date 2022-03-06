# Go_web_scrapping

Web scrapping Using golang

## Technology Used

1. Golang
2. Postgresql (database)
3. Docker

## Features of the project

1. Create User
2. Login User
3. User and Scrape multiple website by passing multiple url through POST method (using GO Routine and Channel)
4. User can view all their Scrapped Content with pagination
5. Any one can view all Scrapped Content with pagination
6. Search by url
7. Filter by DateTime

## Package Used

1. gin [ framework ] ( url :- https://github.com/gin-gonic/gin )
2. swagger [ api documentation ] ( url :- https://github.com/swaggo/gin-swagger )
3. ginCors [ cors ] ( url :- https://github.com/gin-contrib/cors )
4. sqlc [ orm ] ( url :- https://sqlc.dev/ )
5. pesto [ token ] ( url :- https://github.com/o1egl/paseto )
6. viper [ to read config file ] ( url :- https://github.com/spf13/viper )

## Doumentation

http://localhost:8080/swagger/index.html#/

## Steps to run

1. install docker on your machine
2. clone repo " git clone https://github.com/amallick86/Go_web_scrapping.git "
3. OPEN the project in vs code
4. RUN command " make dcup " OR " docker-compose up " on vs code teminal for up your composer and wait for it
5. Hit http://localhost:8080/swagger/index.html#/ url in chrome for documentation
6. You can use above documentation to test all REST API except SEARCH API
7. IN above documentation click on " Authorize " button and paste your token that you get from login api, paste toke as " Bearer your_token "
8. For search open your postman
9. in url " http://localhost:8080/search?q=https://www.facebook.com " with get method

## Steps to close the project

1. RUN command " make dcdown " OR " docker-compose down " for down your composer
2. RUN command " make drmi " OR " docker rmi go_web_scrapping_api " for remove image from your computer

