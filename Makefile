.PHONY: setup
setup:
	curl -sL https://deb.nodesource.com/setup_18.x | sudo -E bash -

	sudo apt install -y \
		nodejs \
		golang \
		postgresql

.PHONY: dev
dev:
	docker compose up --remove-orphans

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
