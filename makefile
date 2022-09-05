build:
	go build -o out/proxy cmd/proxy/*.go
	go build -o out/server cmd/server/*.go
	go build -o out/client cmd/client/*.go

clean:
	go clean
	rm -rfi out/*
	rm -rfi log/*

