fastwalk:
	-go get -u -v github.com/golang/tools/...
	cp -r $(GOPATH)/src/github.com/golang/tools/imports/fastwalk* fastwalk/
	sed -i 's/imports/fastwalk/' fastwalk/*
	sed -i 's/func fastWalk/func Walk/' fastwalk/fastwalk.go

.PHONY: fastwalk
