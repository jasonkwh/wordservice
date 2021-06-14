# wordservice
Simple RESTful API gateway for word searching & updating operations. Use RAM for temporary word storage.

### Directory Structure

- `api`: Front-end REST gateway, forwarding requests onto service(s).
- `wordservice`: Back-end grpc word service.

```
wordservice/
├─ api/
│  ├─ proto/
│     ├─ wordservice/
│        ├─ wordservice.proto
│  ├─ main.go
├─ wordservice/
│  ├─ proto/
│     ├─ wordservice/
│        ├─ wordservice.proto
│  ├─ service/
│     ├─ wordservice.go
│     ├─ wordservice_test.go
│     ├─ storage.go
│  ├─ main.go
├─ README.md
```

### Getting Started

- Build & start Front-end API service
```bash
cd ./api
go build && ./deltatre_api
➜ 2021/06/14 11:12:40 API server listening on: localhost:8000
```

- Test, Build, and Start Back-end service
```bash
cd ./wordservice
go test ./service && go build && ./deltatre_grpc
➜ 2021/06/14 11:12:41 gRPC server listening on: localhost:9000
```

### API Usage

- Add new word
```bash
curl -X "POST" "http://localhost:8000/v1/words" \
     -H 'Content-Type: application/json' \
     -d $'{
  "values": ["{new_word_1}","{new_word_2}"]
}'
```

- Use keyword to search words
```bash
curl -X "GET" "http://localhost:8000/v1/words/{keyword}" \
     -H 'Content-Type: application/json'
```
PS: Use single keyword

- Update searchable words
```bash
curl -X "PATCH" "http://localhost:8000/v1/words" \
     -H 'Content-Type: application/json' \
     -d $'{
  "orig_value": "{original_word}","new_value":"{new_word}"
}'
```

- Return the top 5 search keywords and how many times each has been searched
```bash
curl -X "GET" "http://localhost:8000/v1/trends" \
     -H 'Content-Type: application/json'
```
