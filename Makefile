GO111MODULE=on

# -----------------------------
# Perform all actions and serve
# -----------------------------
# Run validations, generate all recipes, generate toc and start up the service.
all: ## Run validations, generate all recipes, generate toc and start up the service.
	yamale -s recipe_schema.yaml pages/recipes
	go run tools/mdgen/cmd/main.go generate
	go run tools/mdgen/cmd/main.go toc
	docker-compose up -d

# -------------
# Main Commands
# -------------

# Run main.go script to actually generate all the markdown files based on existing yaml files and also to create the respective toc on main README file.
generate_all: ## Run main.go script to actually generate all the markdown files based on existing yaml files and also to create the respective toc on main README file.
	go run tools/mdgen/cmd/main.go generate
	go run tools/mdgen/cmd/main.go toc

# Run main.go script to generate all the markdown files based on the existing yaml files.
generate: ## Run main.go script to generate all the markdown files based on the existing yaml files.
	go run tools/mdgen/cmd/main.go generate

# Run main.go script to create the toc file in the README.md file on the pages folder.
generate_toc: ## Run main.go script to create the toc file in the README.md file on the pages folder.
	go run tools/mdgen/cmd/main.go toc

# Validate that the created yaml files respect the schema defined.
validate: ## Validate that the created yaml files respect the schema defined.
	yamale -s recipe_schema.yaml pages/recipes

# Run main.go script to watch the yaml files for changes and once saved, generates the respective markdown files.
watch: ## Run main.go script to watch the yaml files for changes and once saved, generates the respective markdown files.
	go run tools/mdgen/cmd/main.go watch

# -------------------------------
# Deploy to Github Pages Commands
# -------------------------------

# Deploy mkdocs to github pages.
deploy: ## Deploy mkdocs to github pages.
	mkdocs gh-deploy

# -----------------------
# Docker Related Commands
# -----------------------

# Rebuild image and startup the docker server to display mkdocs with all the data.
build_serve: ## Rebuild image and startup the docker server to display mkdocs with all the data.
	docker-compose up --build -d

# Startup the docker server to display mkdocs with all the data.
serve: ## Startup the docker server to display mkdocs with all the data.
	docker-compose up -d

# Kills the running container of mkdos.
stop: ## Kills the running container of mkdos.
	docker-compose down

# --------------------
# Setup/Infra Commands
# --------------------

# Download all the dependencies.
deps: ## Download all the dependencies.
	go get -u github.com/urfave/cli/v2
	go get -u gopkg.in/yaml.v2
	brew install pyenv
	pyenv install 3.7.3
	pyenv global 3.7.3
	pip install --upgrade pip

# ------------
# Help Command
# ------------

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
