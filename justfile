set dotenv-load
set dotenv-required
set windows-shell := ["C:/Program Files/Git/bin/bash.exe", "-c"]

default: dev

# Generate wrangler.toml file
generate-config:
    envsubst < wrangler.template.toml > wrangler.toml

# Development server
dev: generate-config build
    wrangler dev

# Build the worker
build:
    go run github.com/syumai/workers/cmd/workers-assets-gen@v0.27.0
    tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

# Deploy the worker
deploy: generate-config build
    wrangler deploy
    wrangler secret bulk .env
