PROJECT_NAME = logdiff
MAIN = main.go
BUILD_DIR =_build
E2E_TESTS = $(wildcard test/e2e/*_test.go)

DIRECTORIES = $(BUILD_DIR)

GO = go
BUILD_CMD = build -o $(BUILD_DIR)/$(PROJECT_NAME) $(PROJECT_NAME)/$(MAIN)
E2E_CMD = test -v

MKDIR = mkdir
MKDIR_FLAGS = -pv

RM = rm
RM_FLAGS = -rvf

.PHONY: make_directories

all: make_directories

make_directories:
	$(MKDIR) $(MKDIR_FLAGS) $(DIRECTORIES)  

build: make_directories
	$(GO) $(BUILD_CMD)

e2e: build
	$(GO) $(E2E_CMD) $(E2E_TESTS)

clean:
	$(RM) $(RM_FLAGS) $(DIRECTORIES)
