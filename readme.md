## Installation


```bash
go run main.go
```

## Usage

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"a":2,"b":2}' \
  http://localhost:8080/calculator.div
```
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"a":2,"b":2}' \
  http://localhost:8080/calculator.sum
```
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"a":2,"b":2}' \
  http://localhost:8080/calculator.mul
```
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"a":2,"b":2}' \
  http://localhost:8080/calculator.sub
```