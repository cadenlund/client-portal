# Run tests
# ./... Means this directory and all subdirectories recursively
test: 
    cd apps/api && go test ./...

test-v: 
    cd apps/api && go test -v ./...