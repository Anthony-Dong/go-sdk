package static

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/anthony-dong/go-sdk/commons"
)

type IDLInfo struct {
	Idl  map[string]string
	Main string
}

type IDLProvider interface {
	LoadIDL() (idl *IDLInfo, err error)
}

type loadIdlProvider struct {
	loadDir    string
	main       string
	concurrent bool
}

func NewLocalIdlProvider(local string, main string) IDLProvider {
	return &loadIdlProvider{
		loadDir:    local,
		main:       main,
		concurrent: false,
	}
}

func (l *loadIdlProvider) LoadIDL() (*IDLInfo, error) {
	absDir, err := filepath.Abs(l.loadDir)
	if err != nil {
		return nil, fmt.Errorf(`load abs dir: %s find err: %v`, l.loadDir, err)
	}
	l.loadDir = absDir
	absDir, err = filepath.Abs(l.main)
	if err != nil {
		return nil, fmt.Errorf(`load abs dir: %s find err: %v`, l.main, err)
	}
	l.main = absDir

	files, err := commons.GetAllFiles(l.loadDir, func(fileName string) bool {
		return strings.HasSuffix(fileName, ".proto")
	})
	if err != nil {
		return nil, err
	}

	mainRel, err := filepath.Rel(l.loadDir, l.main)
	if err != nil {
		return nil, fmt.Errorf(`load relative file name err, file: %s, dir: %s`, l.main, l.loadDir)
	}

	wg, _ := errgroup.WithContext(context.Background())
	idls := make(map[string]string)
	idlsLock := sync.Mutex{}
	for _, _filename := range files {
		filename := _filename
		loadIdl := func() error {
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return fmt.Errorf(`read file: %s find err: %v`, filename, err)
			}
			idlsLock.Lock()
			defer idlsLock.Unlock()
			rel, err := filepath.Rel(l.loadDir, filename)
			if err != nil {
				return fmt.Errorf(`load relative file name err, file: %s, dir: %s`, filename, l.loadDir)
			}
			idls[rel] = string(content)
			return nil
		}
		if l.concurrent {
			wg.Go(loadIdl)
		} else {
			if err := loadIdl(); err != nil {
				return nil, err
			}
		}
	}
	if l.concurrent {
		if err := wg.Wait(); err != nil {
			return nil, err
		}
	}
	return &IDLInfo{
		Idl:  idls,
		Main: mainRel,
	}, nil
}
