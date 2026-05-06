.PHONY: build download codegen refresh check clean extension

build:
	go build -o webex .

download:
	cd codegen && python3 download_collections.py

codegen:
	cd codegen && python3 extract_api_spec.py
	cd codegen && python3 generate_cli.py
	cd codegen && python3 generate_skills.py

refresh: download codegen build

check: build
	go vet ./...

extension: build
	cd extension && zip -r ../webex-mcp.dxt manifest.json icon.png
	@echo "Built webex-mcp.dxt"

clean:
	rm -f webex webex-cli webex-mcp.dxt
	rm -f codegen/api_spec.json codegen/api_outline.json
