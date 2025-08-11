

sign:
	mkdir -p "out"
	go build -ldflags "-s -w"  -o ./out/ ./cmd/sign

parser:
	mkdir -p "out"
	go build -ldflags "-s -w"  -o ./out/ ./cmd/parser



dev:
	wails dev -s -skipbindings -extensions refresh -tags "debug"


generate:
	wails generate module

.PHONY: build
build:
	cp res/favicon.ico build/windows/icon.ico
	cp res/favicon.ico build/darwin/icon.ico
	cp res/appicon.png build/appicon.png
	cp res/favicon.ico frontend/public/favicon.ico

	wails build  -trimpath -webview2 embed   -ldflags "-checklinkname=0"
