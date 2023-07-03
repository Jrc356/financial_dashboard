.PHONY: setup
setup:
	curl -sL https://deb.nodesource.com/setup_18.x | sudo -E bash -

	sudo apt install -y \
		nodejs \
		golang \
		postgresql

.PHONY: dev
dev:
	./scripts/run-dev.sh

.PHONY: build
build: clean build-client build-backend

.PHONY: build-client
build-client:
	cd client && npm run build

.PHONY: build-backend
build-backend:
	go build -o app .

.PHONY: clean
clean:
	rm app
	rm -r client/build
