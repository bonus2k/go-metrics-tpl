package rest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"hash"
	"io"
	"net/http"
)

type SignSHA256 struct {
	password []byte
	isActive bool
}

func NewSignSHA256(password string) *SignSHA256 {
	return &SignSHA256{password: []byte(password), isActive: password != ""}
}

func (sign *SignSHA256) AddSignToReq(c *resty.Client, r *http.Request) error {
	if !sign.isActive {
		return nil
	}
	if r.Body == nil {
		return nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read body %w", err)
	}
	sum := getSign(body, sign.password)
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.Header.Add(m.KeyHashSHA256, hex.EncodeToString(sum))
	return nil
}

func (sign *SignSHA256) AddSignToRes(h http.Handler) http.Handler {
	signFunc := func(w http.ResponseWriter, r *http.Request) {
		if !sign.isActive {
			h.ServeHTTP(w, r)
			return
		}
		writer := signResponseWriter{
			ResponseWriter: w,
			hash:           hmac.New(sha256.New, sign.password),
		}
		h.ServeHTTP(&writer, r)
	}
	return http.HandlerFunc(signFunc)
}

func (sign *SignSHA256) CheckSignReq(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sign.isActive {
			h.ServeHTTP(w, r)
			return
		}
		signR := r.Header.Get(m.KeyHashSHA256)
		if signR == "" {
			h.ServeHTTP(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		signB := getSign(body, sign.password)
		decodeSignR, _ := hex.DecodeString(signR)
		if !hmac.Equal(decodeSignR, signB) {
			http.Error(w, "signature is not valid", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		h.ServeHTTP(w, r)
	})
}

func getSign(body []byte, pass []byte) []byte {
	hash := hmac.New(sha256.New, pass)
	hash.Write(body)
	sum := hash.Sum(nil)
	return sum
}

type signResponseWriter struct {
	http.ResponseWriter
	hash hash.Hash
}

func (w signResponseWriter) Write(b []byte) (int, error) {
	w.hash.Write(b)
	w.Header().Set(m.KeyHashSHA256, hex.EncodeToString(w.hash.Sum(nil)))
	return w.ResponseWriter.Write(b)
}
