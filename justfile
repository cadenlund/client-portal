# Run tests
# ./... Means this directory and all subdirectories recursively
test: 
    cd apps/api && go test ./...

test-v: 
    cd apps/api && go test -v ./...

# Codegen for sqlc queries
sqlc:
    cd apps/api && sqlc generate

tidy: 
    cd apps/api && go mod tidy