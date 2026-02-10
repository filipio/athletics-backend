#!/bin/bash

# Generate all mocks using mockgen directives
# This script regenerates mocks defined via //go:generate comments in the codebase

set -e

echo "Generating mocks..."
go generate ./...

echo "✓ Mocks generated successfully"
