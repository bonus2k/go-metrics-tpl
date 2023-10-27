package repositories

import (
	"context"
	"encoding/json"
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	path2 "path"
	"time"
)

var fileService MemStorageService

type MemStorageService struct {
	interval int
	lastSave time.Time
	file     *os.File
	encoder  *json.Encoder
	mem      *Storage
}

func NewMemStorageService(interval int, path string, restore bool, mem *Storage) (*MemStorageService, error) {
	if fileService.file == nil {
		dir, _ := path2.Split(path)

		err := os.Mkdir(path2.Dir(dir), 0644)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			return nil, errors.Wrap(err, "can't create directory")
		}
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "can't open file")
		}

		ms := &MemStorageService{
			interval: interval,
			lastSave: time.Now(),
			file:     file,
			encoder:  json.NewEncoder(file),
			mem:      mem,
		}

		if restore {
			ms.loadMem()
		}
		fileService = *ms
	}
	return &fileService, nil
}

func (ms MemStorageService) Save() error {
	run := ms.lastSave.Add(time.Duration(ms.interval) * time.Second)
	metrics, err := mem.GetAllMetrics(context.TODO())
	if err != nil {
		return err
	}
	if run.After(time.Now()) || len(metrics) == 0 {
		return nil
	}
	ms.file.Truncate(0)
	ms.file.Seek(0, 0)
	err = ms.encoder.Encode(mem)
	if err != nil {
		return err
	}
	ms.lastSave = time.Now()
	logger.Log.Info("save mem storage")
	return nil
}

func (ms MemStorageService) loadMem() {
	decoder := json.NewDecoder(ms.file)
	decoder.Decode(ms.mem)
	logger.Log.Info("load mem storage")
}
