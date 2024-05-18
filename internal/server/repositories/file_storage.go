package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	path2 "path"
	"time"

	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	"github.com/bonus2k/go-metrics-tpl/internal/utils"
	"github.com/pkg/errors"
)

var fileService MemStorageService

type MemStorageService struct {
	interval time.Duration
	lastSave time.Time
	file     *os.File
	encoder  *json.Encoder
	mem      *Storage
}

func NewMemStorageService(interval time.Duration, path string, restore bool, mem *Storage) (*MemStorageService, error) {
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
			err := ms.loadMem()
			if err != nil {
				logger.Log.Error("can't load mem storage", err)
			}
		}
		fileService = *ms
	}
	return &fileService, nil
}

func (ms MemStorageService) Save() error {
	run := ms.lastSave.Add(ms.interval)
	metrics, err := mem.GetAllMetrics(context.TODO())
	if err != nil {
		return err
	}
	if run.After(time.Now()) || len(metrics) == 0 {
		return nil
	}
	f := func() error {
		err := ms.file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = ms.file.Seek(0, 0)
		if err != nil {
			return err
		}
		err = ms.encoder.Encode(mem)
		if err != nil {
			return err
		}
		return nil
	}
	err = utils.RetryAfterError(f)
	if err != nil {
		return fmt.Errorf("can't save file %w", err)
	}

	ms.lastSave = time.Now()
	logger.Log.Info("save mem storage")
	return nil
}

func (ms MemStorageService) loadMem() error {
	decoder := json.NewDecoder(ms.file)
	f := func() error {
		err := decoder.Decode(ms.mem)
		if err != nil {
			return err
		}
		return nil
	}
	err := utils.RetryAfterError(f)
	if err != nil {
		return fmt.Errorf("can't load mem storage %w", err)
	}
	logger.Log.Info("load mem storage")
	return nil
}
