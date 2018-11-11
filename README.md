# transformer-go
Simple HTTP service with image operations 

# Usage
```
make 
make run-server # server will be listing on 0.0.0.0:8082
```

## Docker
```apple js
make docker
docker run aldor007/transformer-go

```

## Example request
```
 curl -X POST http://localhost:8082/accept  -d @data.json  -H "content-type: application/json"  >zip.zip
```


# Test
```
make test
```