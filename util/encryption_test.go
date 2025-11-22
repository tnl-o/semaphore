package util

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePrivateKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "generate valid RSA key pair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var privateKeyBuf bytes.Buffer
			publicKey, err := GeneratePrivateKey(&privateKeyBuf)

			require.NoError(t, err)
			assert.NotEmpty(t, publicKey)
			assert.NotEmpty(t, privateKeyBuf.Bytes())

			// Verify private key format
			privateKeyPEM := privateKeyBuf.Bytes()
			block, _ := pem.Decode(privateKeyPEM)
			require.NotNil(t, block)
			assert.Equal(t, "RSA PRIVATE KEY", block.Type)

			// Verify private key can be parsed
			privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			require.NoError(t, err)
			assert.NotNil(t, privateKey)
			assert.Equal(t, 2048, privateKey.N.BitLen())

			// Verify public key format
			publicKeyBlock, _ := pem.Decode([]byte(publicKey))
			require.NotNil(t, publicKeyBlock)
			assert.Equal(t, "PUBLIC KEY", publicKeyBlock.Type)

			// Verify public key can be parsed
			publicKeyParsed, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
			require.NoError(t, err)
			assert.NotNil(t, publicKeyParsed)

			// Verify public key matches private key
			assert.Equal(t, privateKey.PublicKey.N, publicKeyParsed.N)
			assert.Equal(t, privateKey.PublicKey.E, publicKeyParsed.E)
		})
	}
}

func TestGeneratePrivateKeyMultipleTimes(t *testing.T) {
	// Generate multiple keys and verify they are different
	keys := make(map[string]bool)
	for i := 0; i < 5; i++ {
		var privateKeyBuf bytes.Buffer
		publicKey, err := GeneratePrivateKey(&privateKeyBuf)
		require.NoError(t, err)

		// Verify each key is unique
		assert.False(t, keys[publicKey], "Generated duplicate public key")
		keys[publicKey] = true
	}
}

func TestGenerateRecoveryCode(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "generate valid recovery code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, hash, err := GenerateRecoveryCode()

			require.NoError(t, err)
			assert.NotEmpty(t, code)
			assert.NotEmpty(t, hash)
			assert.GreaterOrEqual(t, len(code), 10)
			assert.LessOrEqual(t, len(code), 20)

			// Verify code can be verified
			assert.True(t, VerifyRecoveryCode(code, hash))
		})
	}
}

func TestVerifyRecoveryCode(t *testing.T) {
	tests := []struct {
		name        string
		setupCode   func() (string, string)
		inputCode   string
		expected    bool
		description string
	}{
		{
			name: "valid recovery code",
			setupCode: func() (string, string) {
				code, hash, _ := GenerateRecoveryCode()
				return code, hash
			},
			inputCode:   "", // Will be set in test
			expected:    true,
			description: "should verify correct recovery code",
		},
		{
			name: "invalid recovery code",
			setupCode: func() (string, string) {
				code, hash, _ := GenerateRecoveryCode()
				return code, hash
			},
			inputCode:   "INVALID_CODE_12345",
			expected:    false,
			description: "should reject incorrect recovery code",
		},
		{
			name: "empty code",
			setupCode: func() (string, string) {
				_, hash, _ := GenerateRecoveryCode()
				return "", hash
			},
			inputCode:   "",
			expected:    false,
			description: "should reject empty code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCode, hash := tt.setupCode()
			inputCode := tt.inputCode
			if inputCode == "" && originalCode != "" {
				inputCode = originalCode
			}

			result := VerifyRecoveryCode(inputCode, hash)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGenerateRecoveryCodeUniqueness(t *testing.T) {
	// Generate multiple recovery codes and verify they are unique
	codes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		code, _, err := GenerateRecoveryCode()
		require.NoError(t, err)

		// Verify each code is unique
		assert.False(t, codes[code], "Generated duplicate recovery code")
		codes[code] = true
	}
}

func TestRecoveryCodeFormat(t *testing.T) {
	code, hash, err := GenerateRecoveryCode()
	require.NoError(t, err)

	// Recovery code should be base32 encoded (alphanumeric, uppercase)
	for _, char := range code {
		assert.True(t, strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", char),
			"Recovery code should only contain base32 characters")
	}

	// Hash should be a bcrypt hash (starts with $2a$, $2b$, or $2y$)
	assert.True(t, strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$"),
		"Hash should be a valid bcrypt hash")
}

func TestPrivateKeyPEMFormat(t *testing.T) {
	var privateKeyBuf bytes.Buffer
	_, err := GeneratePrivateKey(&privateKeyBuf)
	require.NoError(t, err)

	privateKeyPEM := privateKeyBuf.String()

	// Verify PEM format
	assert.Contains(t, privateKeyPEM, "-----BEGIN RSA PRIVATE KEY-----")
	assert.Contains(t, privateKeyPEM, "-----END RSA PRIVATE KEY-----")

	// Verify it can be decoded
	block, rest := pem.Decode([]byte(privateKeyPEM))
	require.NotNil(t, block)
	assert.Empty(t, rest) // Should consume entire input
}

func TestPublicKeyPEMFormat(t *testing.T) {
	var privateKeyBuf bytes.Buffer
	publicKey, err := GeneratePrivateKey(&privateKeyBuf)
	require.NoError(t, err)

	// Verify PEM format
	assert.Contains(t, publicKey, "-----BEGIN PUBLIC KEY-----")
	assert.Contains(t, publicKey, "-----END PUBLIC KEY-----")

	// Verify it can be decoded
	block, rest := pem.Decode([]byte(publicKey))
	require.NotNil(t, block)
	assert.Empty(t, rest) // Should consume entire input
}

func TestPrivateKeySize(t *testing.T) {
	var privateKeyBuf bytes.Buffer
	_, err := GeneratePrivateKey(&privateKeyBuf)
	require.NoError(t, err)

	privateKeyPEM := privateKeyBuf.Bytes()
	block, _ := pem.Decode(privateKeyPEM)
	require.NotNil(t, block)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	require.NoError(t, err)

	// Verify key size is 2048 bits
	assert.Equal(t, 2048, privateKey.N.BitLen())
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	var privateKeyBuf bytes.Buffer
	publicKeyStr, err := GeneratePrivateKey(&privateKeyBuf)
	require.NoError(t, err)

	// Parse private key
	privateKeyPEM := privateKeyBuf.Bytes()
	privateBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	require.NoError(t, err)

	// Parse public key
	publicBlock, _ := pem.Decode([]byte(publicKeyStr))
	publicKey, err := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	require.NoError(t, err)

	// Verify public key matches private key's public key
	assert.Equal(t, privateKey.PublicKey.N, publicKey.N)
	assert.Equal(t, privateKey.PublicKey.E, publicKey.E)
}

// Benchmark tests
func BenchmarkGeneratePrivateKey(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, err := GeneratePrivateKey(&buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateRecoveryCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := GenerateRecoveryCode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyRecoveryCode(b *testing.B) {
	code, hash, err := GenerateRecoveryCode()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyRecoveryCode(code, hash)
	}
}
