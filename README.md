
## _Nicolae Gherman_ / _Николай_ 

###  Runner: 
```
go run .
``` 

## Task 1 

- The SQL queries are in /SQL_queirs 
- Modify SQL connection config in main.go

## Task 2 

- The logic is shown in /http-controller + /repository 
- benchmark test is in benchmark_test.go

##### Running the test:
```
go test -bench=BenchmarkGetCampaginsPerSource ./http-controller/
``` 

## Task 3

- The logic is shown in /repository 
- The unit test is in ./http-controller/controller_test.go

##### Running the test:
```
go test ./http-controller/
``` 


## Task 4 

[Click me](task4/readme.md)
