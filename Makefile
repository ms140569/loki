SHELL=/bin/bash
BIN_NAME=loki
TARGET_DIR=bin
AGENT_NAME=loki-agentd
AGENT_DIR=agentd
PACKAGE_NAME=loki
PKG_VERSION=$(shell VERSION="$$(cat VERSION)"; if [ "$${VERSION}" == "latest" ]; then echo "$$(git rev-parse HEAD)"; else echo "$${VERSION}"; fi)
PKG_DIR_NAME=${PACKAGE_NAME}_${PKG_VERSION}
CONF_DIR=/etc
BIN_DIR=/usr/bin
BIN_LIST=$(TARGET_DIR)/${AGENT_NAME} $(TARGET_DIR)/${BIN_NAME}
COMPLETION_FILE_DIR=/usr/share/bash-completion/completions
VERSION_TARGET=config/version.go
MAN_DIR=/usr/share/man/man1
MAN_BASE=man
MAN_PAGE=loki.1.gz
OS=$(shell uname -s)
PROTOBUF_DEVS=storage master
MAC_BIN_PATH=/usr/local/bin
MAC_MAN_PATH=/usr/local/share/man/man1
DOCKER_IMAGE=lokidev

ifeq ($(OS), Linux)
	install_target=install_linux	
else
	install_target=install_mac
endif

run: build
	@$(TARGET_DIR)/$(BIN_NAME)
build: proto buildversion man
	@go build -o $(TARGET_DIR)/$(BIN_NAME)
	@cd $(AGENT_DIR); go build -o $(CURDIR)/$(TARGET_DIR)/$(AGENT_NAME)
run2: proto
	@go run *.go data/single/matthias.loki

proto: $(addprefix storage/,$(PROTOBUF_DEVS:=.pb.go)) 

.PHONY: $(PROTOBUF_DEVS:=.proto)

$(addprefix storage/,$(PROTOBUF_DEVS:=.pb.go)): %.pb.go: %.proto
	@protoc --go_out=. $<

clean: undeb
	@rm -f storage/*.pb.go
	@rm -f $(TARGET_DIR)/*
	@rm -f debug
	@rm -f $(VERSION_TARGET)
	@rm -f $(MAN_PAGE)

.PHONY: stat
stat: clean
	@find . -type f -name \*.go |xargs wc -l

# The go test command has a cache, even for os.Getenv. See
# https://github.com/kelseyhightower/envconfig/issues/107
.PHONY: test
test: build
	go test -count=1 -v loki/utils
	go test -count=1 -v loki/cmd

.PHONY: man
man:
	@pandoc --standalone --to man $(MAN_BASE)/loki.1.md -o $(MAN_BASE)/loki.1
	@gzip -c $(MAN_BASE)/loki.1 > $(MAN_PAGE)
	@rm -f $(MAN_BASE)/loki.1

lint:
	golint `find . -depth 1 -type d -not -name '.git' -not -name '.vscode'`

godoc:
	godoc -http=:6060

version:
	@echo $(PKG_VERSION)

buildversion:
	@echo "package config" > $(VERSION_TARGET)
	@echo "" >> $(VERSION_TARGET)
	@echo "// SoftwareVersion generated from VERSION file." >> $(VERSION_TARGET)
	@echo "const SoftwareVersion = \"$(PKG_VERSION)\"" >> $(VERSION_TARGET)

deb: build undeb
	@echo ${PKG_VERSION}
	@mkdir -p ${PKG_DIR_NAME}${BIN_DIR}
	@mkdir -p ${PKG_DIR_NAME}${CONF_DIR}
	@mkdir -p ${PKG_DIR_NAME}${COMPLETION_FILE_DIR}
	@mkdir -p ${PKG_DIR_NAME}${MAN_DIR}
	@cp -r ${BIN_LIST} ${PKG_DIR_NAME}${BIN_DIR}
	@cp bash/completion.bash ${PKG_DIR_NAME}${COMPLETION_FILE_DIR}/loki
	@cp $(MAN_PAGE) ${PKG_DIR_NAME}${MAN_DIR}
	@debian/create-control-file ${PACKAGE_NAME} ${PKG_VERSION}
	@cp debian/postinst ${PACKAGE_NAME}_${PKG_VERSION}/DEBIAN
	@dpkg-deb --build ${PKG_DIR_NAME}

undeb:
	@rm -rf ${PKG_DIR_NAME}
	@rm -rf ${PKG_DIR_NAME}.deb

install: test $(install_target)

install_linux: deb
	@echo "Installing on Linux. Use: ${PKG_DIR_NAME}.deb"

install_mac: build
	@echo "Installing on Mac"
	@cp -r ${BIN_LIST} ${MAC_BIN_PATH}
	@cp -r ${MAN_PAGE} ${MAC_MAN_PATH}
	@cp bash/completion.bash /usr/local/etc/bash_completion.d/loki

# The host system gotta have a decent protoc compiler.
# protoc 3.0.0 won't do, 3.6.0 will
# no idea about the versions inbetween
docker: clean proto
	@docker build -t ${DOCKER_IMAGE} .

jump:
	@docker run -it ${DOCKER_IMAGE} bash
