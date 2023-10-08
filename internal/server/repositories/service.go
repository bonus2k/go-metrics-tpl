package repositories

import (
	"encoding/json"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"go.uber.org/zap"
	"os"
	path2 "path"
	"time"
)

var service MemStorageService

type MemStorageService struct {
	interval int
	lastSave time.Time
	file     *os.File
	encoder  *json.Encoder
}

func NewMemStorageService(interval int, path string, restore bool) *MemStorageService {
	if service.file == nil {
		dir, _ := path2.Split(path)
		os.Mkdir(path2.Dir(dir), 0644)
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			panic(err)
		}

		ms := &MemStorageService{
			interval: interval,
			lastSave: time.Now(),
			file:     file,
			encoder:  json.NewEncoder(file),
		}

		if restore {
			ms.loadMem()
		}
		service = *ms
	}
	return &service
}

func (ms MemStorageService) Save() error {
	run := ms.lastSave.Add(time.Duration(ms.interval) * time.Second)
	if run.After(time.Now()) || len(mem.GetAllMetrics()) == 0 {
		return nil
	}
	ms.file.Truncate(0)
	ms.file.Seek(0, 0)
	err := ms.encoder.Encode(mem)
	if err != nil {
		return err
	}
	ms.lastSave = time.Now()
	logger.Log.Info("save mem storage")
	return nil
}

func (ms MemStorageService) loadMem() {
	decoder := json.NewDecoder(ms.file)
	decoder.Decode(&mem)
	logger.Log.Info("load mem storage")
}
