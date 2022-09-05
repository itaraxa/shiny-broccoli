build:
	go build -o out/proxy cmd/proxy/*.go

clean:
	go clean
	rm -rfi out/*
	rm -rfi log/*

