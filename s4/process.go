package s4

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

var letters = []rune("abcdefghijklmnopqrstuvwxyz ")

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
		return fmt.Errorf("unable to select files.")
	}
	log.Printf("[v] files selection: %v\n", selection)
	p.selection = selection
	return nil
}

func (p *process) encode() error {
	log.Println("[?] encoding...")

	size := len(p.selection)
	messages := s1.Split(p.config.Message, size)
	if len(messages) != size {
		return fmt.Errorf("unable to split the message.")
	}
	log.Printf("[v] split message: %s\n", strings.Join(messages, "$"))

	elements := map[string]string{}
	for i := range messages {
		elements[p.selection[i]] = messages[i]
	}

	var wg sync.WaitGroup
	errs := make(chan error)
	wgDone := make(chan bool)
	err := run(p.config.Input, func(file fs.FileInfo) {
		wg.Add(1)
		message, ok := elements[file.Name()]
		go func(intput string, apply bool) {
			defer wg.Done()
			value := intput
			key := p.config.Key
			if !apply {
				key = randString(len(p.config.Key))
				value = randString(len(p.config.Message) / len(p.selection))
			}
			text, err := s2.Encrypt(value, key)
			if err != nil {
				errs <- err
				return
			}
			err = s3.Encrypt(filepath.Join(p.config.Input, file.Name()), filepath.Join(p.config.Output, file.Name()), text)
			if err != nil {
				errs <- err
				return
			}
		}(message, ok)
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
		log.Println("[v] successful encoding!")
		return p.check()
	case err := <-errs:
		close(errs)
		return err
	}
}

func (p *process) decode() (string, error) {
	log.Println("[?] decoding...")

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
	log.Println("[v] successful decoding!")
	return strings.Join(result, ""), nil
}

func (p *process) check() error {
	log.Println("[?] checking...")

	inputFiles, err := ioutil.ReadDir(p.config.Input)
	if err != nil {
		return err
	}
	outputFile, err := ioutil.ReadDir(p.config.Output)
	if err != nil {
		return err
	}

	if len(inputFiles) != len(outputFile) {
		return fmt.Errorf("missing output files")
	}

	for _, file := range outputFile {
		if time.Since(file.ModTime()) > time.Minute {
			return fmt.Errorf("modification date mismatch.")
		}
	}
	log.Println("[v] successful checking!")
	return nil
}
