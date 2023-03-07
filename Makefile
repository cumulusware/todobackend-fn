help:
	@echo "You can perform the following:"
	@echo ""
	@echo "  check         Format, vet, and test Go code"
	@echo "  cover         Show test coverage in html"
	@echo "  deploy        Deploy to IBM Cloud Functions"
	@echo "  lint          Lint Go code"
	@echo "  list          List APIs"

check:
	@echo 'Formatting, vetting, and testing Go code'
	go fmt ./...
	go vet ./...
	go test ./... -cover

cover:
	@echo 'Test coverage in html'
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	@echo 'Linting code using staticcheck'
	staticcheck -f stylish ./...

todos:
	cd resources/todos/create/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-create-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-create-src.zip >todos-create-bin.zip
	cd resources/todos/readall/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-readall-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-readall-src.zip >todos-readall-bin.zip
	cd resources/todos/read/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-read-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-read-src.zip >todos-read-bin.zip
	cd resources/todos/delete/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-delete-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-delete-src.zip >todos-delete-bin.zip
	cd resources/todos/deleteall/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-deleteall-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-deleteall-src.zip >todos-deleteall-bin.zip
	cd resources/todos/update/src/main; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
		zip -r ../todos-update-src.zip *; \
		cd ..; \
		docker run -i openwhisk/actionloop-base -compile main <todos-update-src.zip >todos-update-bin.zip

deploy: todos
	@echo 'Deploy to IBM Cloud Functions'
	ibmcloud target --cf -o TodoBackendOrg -s dev
	ibmcloud fn property unset --namespace
	source ./set_env.sh; ibmcloud fn deploy

list:
	ibmcloud fn api list
