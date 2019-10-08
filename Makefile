help:
	@echo "You can perform the following:"
	@echo ""
	@echo "  check         Format, lint, vet, and test Go code"
	@echo "  cover         Show test coverage in html"
	@echo "  deploy        Deploy to IBM Cloud Functions"
	@echo "  prep          Prepare to develop"
	@echo "  list          List APIs"

check:
	@echo 'Formatting, linting, vetting, and testing Go code'
	go fmt ./...
	golint ./...
	go vet ./...
	go test ./... -cover

cover:
	@echo 'Test coverage in html'
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

deploy:
	@echo 'Deploy to IBM Cloud Functions'
	ibmcloud fn deploy

prep:
	@echo 'Prepare for development by setting resource group'
	ibmcloud target -g TodoBackend
	ibmcloud fn property set --namespace todobackend

list:
	ibmcloud fn api list
