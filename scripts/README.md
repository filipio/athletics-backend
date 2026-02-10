# Scripts

This directory contains utility scripts for the athletics-backend project.

## generate-mocks.sh

Regenerates all mocks in the `mocks/` directory.

### Usage

```bash
./scripts/generate-mocks.sh
```

### How It Works

The script runs `go generate ./...` which executes all `//go:generate` directives found in the codebase. Mocks are defined using these directives in their respective package files.

### Adding a New Mock

1. Define an interface in your package (e.g., `email/interface.go`)
2. Add a `//go:generate` directive above the interface:
   ```go
   //go:generate mockgen -destination=../mocks/mock_your_interface.go -package=mocks github.com/filipio/athletics-backend/your_package YourInterface
   ```
3. Run `./scripts/generate-mocks.sh` to generate the mock

### Example

In `email/interface.go`:
```go
//go:generate mockgen -destination=../mocks/mock_email_sender.go -package=mocks github.com/filipio/athletics-backend/email EmailSender

type EmailSender interface {
    SendVerificationEmail(ctx context.Context, params VerificationEmailParams) error
}
```

After adding the directive, run the script to generate `mocks/mock_email_sender.go`.

### Git Configuration

Mock files in `mocks/*.go` are gitignored to avoid committing generated code. Ensure the `mocks/` directory exists before running the script.
