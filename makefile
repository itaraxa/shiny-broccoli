build:
	go build -o out/server cmd/server/*.go
	go build -o out/client cmd/client/*.go
	go build -o out/proxy cmd/proxy/*.go

clean:
	go clean
	rm -rfi out/*
	rm -rfi log/*

