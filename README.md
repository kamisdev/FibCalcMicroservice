# Simple Microservice for Fibonacci and Mathematical Calculation

## How to run

### Running Server
```
go run cmd/run_machine_server.go
```

### Running Client
```
go run client/machine.go
```

## How to test

### Server general test
```
go test server/machine.go server/machine_test.go
```

### Server general live test
```
go test server/machine.go server/machine_live_test.go
```

### Client mock test
```
go test mock_machine/machine_mock.go mock_machine/machine_mock_test.go
```