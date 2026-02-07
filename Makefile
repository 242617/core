


.PHONY: test
test:
	@go clean -testcache
	@go test -race -v \
		./...

test\:integration:
	@mkdir -p .tmp/ && \
		go test -v -tags=integration ./...

doc:
	@echo > doc.txt && \
		go list ./... | \
			grep -v /cmd/ | \
			grep -v /mocks/ | \
			grep -v /vendor/ | \
			xargs -n1 -I {} sh -c 'echo "=== {} ===" >> doc.txt && go doc -all {} >> doc.txt'
