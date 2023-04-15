## RusProfile gRPC wrapper

is a CLI tool & HTTP service (via grpc-gateway) to search company info on RusProfile by company INN

### usage (grpc):
![Image alt](https://github.com/SubochevaValeriya/gRPC-task/blob/dev/tools/logo/grpc.png)

```

- cmd/client/main.go [flags] INNs (you can input several INNs divided by backspaces)
- cmd/client/main.go [flags] file name.ext 
```

### usage (grpc-gateway):
![Image alt](https://github.com/SubochevaValeriya/gRPC-task/blob/dev/tools/logo/http.png)

### usage examples:
```
go run cmd/client/main.go 5008042065
go run cmd/client/main.go file companies.txt
```

### commands:

``` file name.ext ```

### flags:
```
--help     Show help message
```  

### To run a server:

```
make build && make run
```

**Used:** *gRPC, grpc-gateway, swagger, docker.*
