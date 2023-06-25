# go-microservice
Example of a Go microservice with API points, DB connection and modular architecture focused in using my test framework to test APIs.


## To run:
Install the dependencies with "go get ./..."
Run the microservice with "go run ./..."
Tests can be run using the makefile or individually with 
"go test -timeout 30s -run ^TestService_Mongo$ github.com/Martin-Jast/go-microservice/server"



When testing try to use a docker that will put the mongoDB up so it has a clean database when it starts or remember to drop the database between tests ( to make sure no database is wrongly dropped the code don't do that )