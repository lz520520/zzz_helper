

sign:
	mkdir -p "out"
	go build -ldflags "-s -w"  -o ./out/ ./cmd/sign

parser:
	mkdir -p "out"
	go build -ldflags "-s -w"  -o ./out/ ./cmd/parser
