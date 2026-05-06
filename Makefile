.PHONY: build download codegen refresh check clean extension skill

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

skill:
	@python3 -c "\
import zipfile, os; \
z = zipfile.ZipFile('webex.skill', 'w', zipfile.ZIP_DEFLATED); \
[z.write(os.path.join(r,f), 'webex-cli/' + os.path.relpath(os.path.join(r,f), 'skill')) \
 for r, _, files in os.walk('skill') for f in files if f.endswith('.md')]; \
z.close()"
	@echo "Built webex.skill"

clean:
	rm -f webex webex-cli webex-mcp.dxt webex.skill
	rm -f codegen/api_spec.json codegen/api_outline.json
