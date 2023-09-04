package controllers

import "github.com/bonus2k/go-metrics-tpl/cmd/server/repositories"

var MemStorage repositories.MemStorage

func init() {
	MemStorage = repositories.NewMemStorage()
}
