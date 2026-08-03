generate name:
  go run . generate migrations {{name}}

sdk:
  go run . generate sdk

# Install UI deps (run once / after dependency changes).
ui-install:
  cd ui && yarn install

# Build the frontend into ui/dist so it gets embedded by `go build`.
ui:
  cd ui && yarn build

# Build the single binary with the frontend embedded (UI built first).
build: ui
  go build -o bin/waggle .

release: ui
    goreleaser build --snapshot --clean

# Cut a release: bump VERSION, regenerate the spec/SDKs so they carry the new
# number, then tag. Order matters -- the artifacts must be committed BEFORE the
# tag, or the tagged tree ships a spec claiming the previous version.
# Usage: just release-prep 0.2.1
release-prep version:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "{{version}}" > VERSION
    just sdk
    just terraform
    echo
    echo "VERSION is now {{version}} and artifacts are regenerated."
    echo "Next: review the diff, commit, merge to main, then tag on the merge commit:"
    echo "    git tag v{{version}} && git push origin v{{version}}"
    echo "The Docker Publish workflow verifies the tag matches VERSION."

up:
  docker-compose up -d
  migrate up

down *args:
  docker-compose down {{args}}

reset:
  down -v
  up

terraform:
  go run . generate terraform openapi-generator

readme:
  npx repomix --include "CLAUDE.md,AGENTS.md,Justfile,Dockerfile,docker-compose*.yml,go.mod,agent/go.mod,ui/package.json,.goreleaser.yml,cmd/**/*.go,internal/handlers/*.go,internal/models/**/*.go,ui/src/routes/index.tsx,mprocs.yml" -o repomix-readme.md
  cat repomix-readme.md | claude --dangerously-skip-permissions -p "You are a technical writer. Using the codebase snapshot from stdin, generate a comprehensive README.md for this project. Include: project overview, features list, tech stack, prerequisites, setup/installation, development workflow, architecture overview, API documentation summary (reference /docs for full OpenAPI spec), deployment notes, and contributing guidelines. Keep it professional and concise. Do NOT include the repomix snapshot in the output — just the README content." --model haiku > README.md
  rm -f repomix-readme.md
  @echo "README.md generated successfully."

migrate *args:
  go run . migrate {{args}}
