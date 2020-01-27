PLATS = linux windows
.PHONY : none $(PLATS)

CGO_ENABLED := 0
GOOS := windows
GOARCH := amd64

linux : GOOS := linux

windows linux :
	$(MAKE) all CGO_ENABLED="$(CGO_ENABLED)" GOOS="$(GOOS)"

