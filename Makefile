
COMMIT = `git rev-parse HEAD`
BUILDTIME = `date +'%Y-%m-%d_%T'`
GOVER = `go env GOVERSION`
utils = "github.com/yzimhao/trading_engine/utils/app"

version ?= "0.0.0"

mainname = "haotrader"
distdir = "./dist"
exedir = "$(distdir)/$(mainname)"



test:
	go test -v ./...
dist:
	mkdir -p $(exedir)
clean:
	rm -rf $(distdir)


define build_haomatch
	@echo "Building for haotrader $1 $2"
	CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=$4 go build -ldflags="-s -w -X $(utils).Version=${version} -X $(utils).Commit=$(COMMIT) -X $(utils).Build=$(BUILDTIME) -X $(utils).Goversion=$(GOVER)" -o $(exedir)/haomatch$3 cmd/haomatch/main.go
	upx -9 $(exedir)/haomatch$3
endef


define build_haoquote
	@echo "Building for haoquote $1 $2"
	CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=$4 go build -ldflags="-s -w -X $(utils).Version=${version} -X $(utils).Commit=$(COMMIT) -X $(utils).Build=$(BUILDTIME) -X $(utils).Goversion=$(GOVER)" -o $(exedir)/haoquote$3 cmd/haoquote/main.go
	upx -9 $(exedir)/haoquote$3
endef


define build_haobase
	@echo "Building for haobase $1 $2"
	CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=$4 go build -ldflags="-s -w -X $(utils).Version=${version} -X $(utils).Commit=$(COMMIT) -X $(utils).Build=$(BUILDTIME) -X $(utils).Goversion=$(GOVER)" -o $(exedir)/haobase$3 cmd/haobase/main.go
	upx -9 $(exedir)/haobase$3
endef


copy_doc:
	cp README.md $(exedir)/
	cp -rf cmd/config.toml $(exedir)/config.toml_sample

define zipfile
	# tar
	cd $(distdir) && tar czvf $(mainname).$(version).$1-$2.tar.gz `basename $(exedir)`
	# zip
	cd $(distdir) && zip -r -m $(mainname).$(version).$1-$2.zip `basename $(exedir)` -x "*/\.*"
endef

build_linux_amd64:
	@make copy_doc
	$(call build_haobase,linux,amd64,'',x86_64-unknown-linux-gnu-gcc)
	$(call build_haomatch,linux,amd64,'',x86_64-unknown-linux-gnu-gcc)
	$(call build_haoquote,linux,amd64,'',x86_64-unknown-linux-gnu-gcc)
	
	$(call zipfile,linux,amd64)
	

build_darwin_amd64:
	@make copy_doc
	$(call build_haobase,darwin,amd64,'','')
	$(call build_haomatch,darwin,amd64,'','')
	$(call build_haoquote,darwin,amd64,'','')
	
	$(call zipfile,darwin,amd64)


release: clean dist
	@make build_linux_amd64
	@make build_darwin_amd64
	


upload_example:
	mkdir -p $(distdir)/trading_engine_example
	cd example && GOOS=linux GOARCH=amd64 go build -o ../$(distdir)/trading_engine_example/example example.go
	upx -9 $(distdir)/trading_engine_example/example
	cp -rf example/statics $(distdir)/trading_engine_example/
	cp -rf example/demo.html $(distdir)/trading_engine_example/
	scp -r $(distdir)/trading_engine_example/ demo:~/
	


upload_all:
	@make clean
	@make dist
	@make build_linux_amd64
	@make upload_example
	scp $(distdir)/haotrader.$(version).linux-amd64.tar.gz demo:~/
	ssh demo "tar xzvf haotrader.$(version).linux-amd64.tar.gz"
	ssh demo 'rm -f haotrader.$(version).linux-amd64.tar.gz'
	@make example_reload

example_start:

	ssh demo 'cd haotrader/ && ./haobase -d'
	ssh demo 'cd haotrader/ && ./haomatch -d'
	ssh demo 'cd haotrader/ && ./haoquote -d'
	ssh demo 'cd trading_engine_example/ && ./example -d'


example_stop:
   	
	ssh demo 'pgrep haobase | xargs kill'
	ssh demo 'pgrep haomatch | xargs kill'
	ssh demo 'pgrep haoquote | xargs kill'
	ssh demo 'pgrep example | xargs kill'
	


example_reload:
	@make example_stop
	@make example_start

example_clean:

	ssh demo 'cd haotrader/ && rm -f ./*.db'
	@make example_stop




require:
	brew tap messense/macos-cross-toolchains
	# install x86_64-unknown-linux-gnu toolchain
	brew install x86_64-unknown-linux-gnu
	
	brew install upx
	


tag:
	git tag -a $(version)
	git push --tags

.PHONY: dist clean build release
	