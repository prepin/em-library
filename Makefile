.PHONY: test

test:
	go test -v -count=1 -parallel=4 -coverprofile=cov.out ./...
	go tool cover -func=cov.out

coverage:
	go tool cover -html=cov.out
