# auth-service

## Using:

### Docker Compose:

- docker-compose -f docker-compose.auth.yaml --env-file local.env up --build

### Kubernetes:

kubectl apply -f kubernetes/secrets.yaml -f kubernetes/postgres-deployment.yaml -f kubernetes/auth-service-deployment.yaml 
kubectl delete -f kubernetes/secrets.yaml -f kubernetes/postgres-deployment.yaml -f kubernetes/auth-service-deployment.yaml


### Установите плагины для Go
Для работы с Go вам нужны плагины для генерации кода. Установите их с помощью go install:
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

### Генерация gRPC кода
Теперь выполните команду для генерации:
    protoc --go_out=. --go-grpc_out=. ../../proto/auth.proto