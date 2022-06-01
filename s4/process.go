package s4

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lgylgy/iscript/s1"
	"github.com/lgylgy/iscript/s2"
	"github.com/lgylgy/iscript/s3"
)

func Run(encode bool, config *Config) (string, error) {
	_, err := os.Stat(config.Input)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("input directory '%s' does not exist.", config.Input)
	}
	process := newProcess(config)
	err = process.selectFiles()
	if err != nil {
		return "", err
	}
	if encode {
		return "", process.encode()
	}
	return process.decode()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type process struct {
	config    *Config
	selection []string
}

func newProcess(config *Config) *process {
	return &process{
		config:    config,
		selection: []string{},
	}
}

func run(path string, fun func(fs.FileInfo)) error {
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range filesInfo {
		fun(file)
	}
	return nil
}

func (p *process) selectFiles() error {
	files := []string{}
	err := run(p.config.Input, func(file fs.FileInfo) {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	})
	if err != nil {
		return err
	}

	selection, err := s1.Select(files, p.config.Selection, p.config.Seed)
	if err != nil {
		return err
	}
	if len(selection) == 0 {
		return fmt.Errorf("internal error")
	}
	p.selection = selection
	return nil
}

func (p *process) encode() error {
	size := len(p.selection)
	messages := s1.Split(p.config.Message, size)
	if len(messages) != size {
		return fmt.Errorf("internal error")
	}

	elements := map[string]string{}
	for i := range messages {
		elements[p.selection[i]] = messages[i]
	}

	var wg sync.WaitGroup
	errs := make(chan error)
	wgDone := make(chan bool)
	err := run(p.config.Input, func(file fs.FileInfo) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := p.config.Key
			message, ok := elements[file.Name()]
			if !ok {
				key = randString(12)
				message = randString(len(p.selection[0]))
			}
			text, err := s2.Encrypt(message, key)
			if err != nil {
				errs <- err
				return
			}
			err = s3.Encrypt(filepath.Join(p.config.Input, file.Name()), filepath.Join(p.config.Output, file.Name()), text)
			if err != nil {
				errs <- err
				return
			}
		}()
	})
	if err != nil {
		return err
	}

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		return nil
	case err := <-errs:
		close(errs)
		return err
	}
}

func (p *process) decode() (string, error) {
	result := []string{}
	for _, file := range p.selection {
		value, err := s3.Decrypt(filepath.Join(p.config.Output, file))
		if err != nil {
			return "", err
		}
		text, err := s2.Decrypt(value, p.config.Key)
		if err != nil {
			return "", err
		}
		result = append(result, string(text))
	}
	return strings.Join(result, ""), nil
}
