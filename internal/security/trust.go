// Package security manages trust decisions and persistence for repository config files.
package security

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TrustStatus represents the outcome of a trust check on a file.
type TrustStatus int

const (
	// TrustStatusTrusted indicates the file matches a known hash.
	TrustStatusTrusted TrustStatus = iota
	// TrustStatusUntrusted means the file either changed or has not been trusted yet.
	TrustStatusUntrusted
	// TrustStatusNotFound is returned when the file does not exist.
	TrustStatusNotFound
)

func getTrustDBPath() string {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "lazyworktree", "trusted.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "lazyworktree", "trusted.json")
}

// TrustManager stores trusted hashes and enforces TOFU (Trust On First Use).
type TrustManager struct {
	dbPath        string
	trustedHashes map[string]string // Map absolute path -> sha256 hash
}

// NewTrustManager creates and loads the persisted trust database.
func NewTrustManager() *TrustManager {
	tm := &TrustManager{
		dbPath:        getTrustDBPath(),
		trustedHashes: make(map[string]string),
	}
	tm.load()
	return tm
}

func (tm *TrustManager) load() {
	if _, err := os.Stat(tm.dbPath); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(tm.dbPath)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &tm.trustedHashes); err != nil {
		// If corrupt, start fresh for safety
		tm.trustedHashes = make(map[string]string)
	}
}

func (tm *TrustManager) save() error {
	dir := filepath.Dir(tm.dbPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tm.trustedHashes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tm.dbPath, data, 0o600)
}

// calculateHash calculates SHA256 of the file content
func (tm *TrustManager) calculateHash(filePath string) string {
	resolvedPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	}

	// #nosec G304 -- resolvedPath is an absolute path derived from trusted input
	file, err := os.Open(resolvedPath)
	if err != nil {
		return ""
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	buf := make([]byte, 65536)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return ""
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// CheckTrust validates the given path against the trust database.
func (tm *TrustManager) CheckTrust(filePath string) TrustStatus {
	resolvedPath, err := filepath.Abs(filePath)
	if err != nil {
		return TrustStatusNotFound
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return TrustStatusNotFound
	}

	currentHash := tm.calculateHash(resolvedPath)
	if currentHash == "" {
		return TrustStatusUntrusted
	}

	storedHash, exists := tm.trustedHashes[resolvedPath]
	if !exists {
		return TrustStatusUntrusted
	}

	if storedHash == currentHash {
		return TrustStatusTrusted
	}

	return TrustStatusUntrusted
}

// TrustFile records the current hash of a file as trusted.
func (tm *TrustManager) TrustFile(filePath string) error {
	resolvedPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", resolvedPath)
	}

	currentHash := tm.calculateHash(resolvedPath)
	if currentHash == "" {
		return fmt.Errorf("failed to calculate hash for: %s", resolvedPath)
	}

	tm.trustedHashes[resolvedPath] = currentHash
	return tm.save()
}
