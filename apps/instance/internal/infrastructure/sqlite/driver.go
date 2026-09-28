package sqlite

// Registers the "sqlite" database/sql driver: pure Go, no cgo, so the Instance
// cross-compiles and runs in a scratch or distroless container.
import _ "modernc.org/sqlite"
