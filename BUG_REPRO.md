# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	stickerchallenge/cmd/stickerctl	[no test files]
ok  	stickerchallenge/internal/config	0.002s
ok  	stickerchallenge/internal/domain	0.002s
ok  	stickerchallenge/internal/httpapi	0.008s
ok  	stickerchallenge/internal/importer	0.002s
ok  	stickerchallenge/internal/report	0.001s
--- FAIL: TestBusiness05Regression (0.01s)
    regression_test.go:32: stale export retained old number: [{"id":"r1","batch_id":"2116-05","number":22,"divisors":[2,11],"result":"pass","confirmed":true,"updated_by":"operator"}]
FAIL
FAIL	stickerchallenge/internal/service	0.017s
ok  	stickerchallenge/internal/store	0.013s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/stickerctl): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/stickerctl): exit `0`
