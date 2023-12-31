PLUGIN_NAME = t1d333/ftp-driver


all: clean rootfs create

test:
	@go test -v -coverpkg ./... -coverprofile=profile.cov ./...
	@cat profile.cov | grep -v mocks > profile.filtred.cov
	
lint:
	@golangci-lint run ./...

coverage: test
	@go tool cover -func  profile.filtred.cov
	
clean:
	@rm -rf ./plugin || true
	@rm profile.cov || true

rootfs:
	@docker build -t ${PLUGIN_NAME}:rootfs .
	@mkdir -p ./plugin/rootfs
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C ./plugin/rootfs
	@cp config.json ./plugin
	@docker rm -vf tmp 
    
create:
	@docker plugin rm -f ${PLUGIN_NAME} || true
	@docker plugin create ${PLUGIN_NAME} ./plugin
    
enable:
	@docker plugin enable ${PLUGIN_NAME}
