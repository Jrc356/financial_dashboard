.PHONY: setup
setup:
	which node >/dev/null || curl -sL https://deb.nodesource.com/setup_18.x | sudo -E bash -
	sudo apt install -y \
		nodejs \
		golang \
		postgresql
	which hadolint >/dev/null || (sudo wget -O /usr/local/bin/hadolint https://github.com/hadolint/hadolint/releases/download/v2.12.0/hadolint-Linux-x86_64 && chmod +x /usr/local/bin/hadolint)
	cp .githooks/* .git/hooks/

.PHONY: dev
dev:
	docker compose up --remove-orphans

.PHONY: test-go
test-go:
	docker build . --target test-go -t test-go
	docker run --rm -v ./coverage/go:/app/coverage test-go

.PHONY: test-node
test-node:
	docker build . --target test-node -t test-node
	docker run --rm -v ./coverage/node:/app/coverage test-node

.PHONY: test
test: test-go test-node

.PHONY: build
build: clean build-client build-backend

.PHONY: build-client
build-client:
	cd client && npm run build

.PHONY: build-backend
build-backend:
	cd server && go build -o app .

.PHONY: clean
clean:
	rm server/app
	rm -r client/build
