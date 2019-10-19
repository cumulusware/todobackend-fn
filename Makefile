help:
	@echo "You can perform the following:"
	@echo ""
	@echo "  check         Format, lint, vet, and test Go code"
	@echo "  cover         Show test coverage in html"
	@echo "  deploy        Deploy to IBM Cloud Functions"
	@echo "  deployiam     Deploy to IBM Cloud Functions to IAM namespace"
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
	ibmcloud target --cf -o TodoBackendOrg -s dev
	ibmcloud fn property unset --namespace
	ibmcloud fn deploy

deployiam:
	@echo 'Deploy to IBM Cloud Functions to IAM namespace'
	ibmcloud target -g TodoBackendRG
	ibmcloud fn property unset --namespace
	ibmcloud fn property set --namespace TodoBackendNamespace
	ibmcloud fn deploy

list:
	ibmcloud fn api list
