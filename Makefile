.PHONY: setup
setup:
	curl -sL https://deb.nodesource.com/setup_18.x | sudo -E bash -

	sudo apt install -y \
		nodejs \
		golang \
		postgresql