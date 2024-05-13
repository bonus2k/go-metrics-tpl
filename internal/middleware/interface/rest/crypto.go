package rest

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
)

type Encrypt struct {
	key      *rsa.PublicKey
	label    []byte
	hash     hash.Hash
	isActive bool
}

func NewEncrypt(file string) (*Encrypt, error) {
	if file == "" {
		return &Encrypt{isActive: false}, nil
	}

	cert, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(cert)
	if err != nil {
		return nil, err
	}

	blocks, _ := pem.Decode(all)
	if blocks == nil {
		return nil, errors.New("can't read certificate")
	}

	certificate, err := x509.ParseCertificate(blocks.Bytes)
	if err != nil {
		return nil, err
	}

	return &Encrypt{key: certificate.PublicKey.(*rsa.PublicKey), label: []byte(""), hash: sha256.New(), isActive: true}, nil
}

func (crypto *Encrypt) EncryptRequest(c *resty.Client, r *http.Request) error {
	if !crypto.isActive {
		return nil
	}
	if r.Body == nil {
		return nil
	}
	body, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			logger.Log.Error("EncryptRequest", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("can't read body %w", err)
	}

	chunkSize := crypto.key.N.BitLen()/8 - 2*crypto.hash.Size() - 2
	result := []byte("")
	chunks := chunkBy[byte](body, chunkSize)

	for _, chunk := range chunks {
		ciphertext, err := rsa.EncryptOAEP(crypto.hash, rand.Reader, crypto.key, chunk, crypto.label)
		if err != nil {
			return fmt.Errorf("can't encrypt body %w", err)
		}
		result = append(result, ciphertext...)
	}
	r.Header.Add(m.KeyContentEncrypt, m.TypeEncryptContent)
	r.Body = io.NopCloser(bytes.NewReader(result))
	return nil
}

type Decrypt struct {
	key      *rsa.PrivateKey
	label    []byte
	hash     hash.Hash
	isActive bool
}

func NewDecrypt(file string) (*Decrypt, error) {
	if file == "" {
		return &Decrypt{isActive: false}, nil
	}

	cert, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(cert)
	if err != nil {
		return nil, err
	}

	blocks, _ := pem.Decode(all)
	if blocks == nil {
		return nil, errors.New("can't read certificate")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(blocks.Bytes)
	if err != nil {
		return nil, err
	}

	return &Decrypt{key: privateKey, label: []byte(""), hash: sha256.New(), isActive: true}, nil
}

func (crypto *Decrypt) DecryptRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !crypto.isActive {
			h.ServeHTTP(w, r)
			return
		}
		if !strings.Contains(r.Header.Get(m.KeyContentEncrypt), m.TypeEncryptContent) {
			h.ServeHTTP(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer func() {
			err = r.Body.Close()
			if err != nil {
				logger.Log.Error("can't close body", err)
			}
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var result []byte
		for _, chnk := range chunkBy[byte](body, crypto.key.N.BitLen()/8) {
			plaintext, err := rsa.DecryptOAEP(crypto.hash, rand.Reader, crypto.key, chnk, crypto.label)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Log.Error("can't decrypt body", err)
			}
			result = append(result, plaintext...)
		}
		r.Body = io.NopCloser(bytes.NewReader(result))
		h.ServeHTTP(w, r)
	})
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
