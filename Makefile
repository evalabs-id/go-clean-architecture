BINARY_NAME=app
BUILD_DIR=build
CMD_DIR=cmd

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -v ./$(CMD_DIR)
