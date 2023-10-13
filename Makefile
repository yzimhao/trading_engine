
COMMIT = `git rev-parse HEAD`
BUILDTIME = `date +'%Y-%m-%d_%T'`
GOVER = `go env GOVERSION`
utils = "github.com/yzimhao/trading_engine/utils/app"

version ?= "1.0.0"

distdir = "./dist"



test:
	go test -v ./...
dist:
	mkdir -p $(distdir)
clean:
	rm -rf $(distdir)


define build_haotrader
	@echo "Building for haotrader $1 $2"
	mkdir -p $(distdir)/haotrader
	CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=$4 go build -ldflags="-s -w -X $(utils).Version=${version} -X $(utils).Commit=$(COMMIT) -X $(utils).Build=$(BUILDTIME) -X $(utils).Goversion=$(GOVER)" -o $(distdir)/haotrader/haotrader$3 cmd/haotrader/main.go
	cp README.md $(distdir)/haotrader/
	cp docs/haotrader.md $(distdir)/haotrader/
	upx -9 $(distdir)/haotrader/haotrader$3
	
	$(call build_haoquote,$1,$2,$3,$4)
	
	cp -rf cmd/config.toml $(distdir)/haotrader/config.toml_sample
	
	# tar
	cd $(distdir) && tar czvf haotrader.$(version).$1-$2.tar.gz `basename $(distdir)/haotrader`
	# zip
	cd $(distdir) && zip -r -m haotrader.$(version).$1-$2.zip `basename $(distdir)/haotrader` -x "*/\.*"
endef


define build_haoquote
	@echo "Building for haoquote $1 $2"
	CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=$4 go build -ldflags="-s -w -X $(utils).Version=${version} -X $(utils).Commit=$(COMMIT) -X $(utils).Build=$(BUILDTIME) -X $(utils).Goversion=$(GOVER)" -o $(distdir)/haotrader/haoquote$3 cmd/haoquote/main.go
	upx -9 $(distdir)/haotrader/haoquote$3
endef


build_linux_amd64:
	$(call build_haotrader,linux,amd64,'',x86_64-unknown-linux-gnu-gcc)

build_darwin_amd64:
	$(call build_haotrader,darwin,amd64,'','')

build_windows_amd64:
	$(call build_haotrader,windows,amd64,'.exe','')

release: clean dist
	@make build_linux_amd64
	@make build_darwin_amd64
	# @make build_windows_amd64


app_example:
	mkdir -p $(distdir)/trading_engine_example
	cd example && GOOS=linux GOARCH=amd64 go build -o ../$(distdir)/trading_engine_example/example example.go
	cp -rf example/statics $(distdir)/trading_engine_example/
	upx -9 $(distdir)/trading_engine_example/example
	cp -rf example/demo.html $(distdir)/trading_engine_example/
	scp -r $(distdir)/trading_engine_example/ demo:~/
	@make example_reload


pubdemo:
	@make clean
	@make dist
	@make build_linux_amd64
	@make app_example
	scp $(distdir)/haotrader.$(version).linux-amd64.tar.gz demo:~/
	ssh demo "tar xzvf haotrader.$(version).linux-amd64.tar.gz"
	ssh demo 'rm -f haotrader.$(version).linux-amd64.tar.gz'
	@make example_reload


example_reload:
   	
	ssh demo 'kill `cat haotrader/haotrader.pid`'
	ssh demo 'kill `cat haotrader/haoquote.pid`'
	ssh demo 'cd haotrader/ && ./haotrader -d'
	ssh demo 'cd haotrader/ && ./haoquote -d'
	ssh demo 'cd trading_engine_example/ && kill `cat run.pid`'
	ssh demo 'cd trading_engine_example/ && ./example -d -port=20001'

example_clean:

	ssh demo 'cd haotrader/ && rm -f ./*.db'
	@make example_reload




doc:
	swag init -g cmd/haoquote/main.go -o docs/api/


require:
	brew tap messense/macos-cross-toolchains
	# install x86_64-unknown-linux-gnu toolchain
	brew install x86_64-unknown-linux-gnu
	
	brew install upx
	go install github.com/swaggo/swag/cmd/swag@latest


tag:
	git tag -a $(version)
	git push --tags

.PHONY: dist clean build release
	