package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Storage interface {
	IdentifyExists()
	Save([]byte) error
	InitStorageFolderStructure(string)
}

type DataStorage struct {
	files   map[string]struct{}
	path    string
	fileExt string
}

func NewDataStorage(path, fileExt string) *DataStorage {
	d := &DataStorage{
		path:    path,
		fileExt: fileExt,
		files:   make(map[string]struct{}),
	}
	d.InitStorageFolderStructure(fileExt)
	d.IdentifyExists()
	return d
}

func (d *DataStorage) IdentifyExists() {
	currentPath := filepath.Join(".", d.path, d.fileExt)
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		f, err := os.Open(filepath.Join(currentPath, f.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Println(err)
			continue
		}
		err = f.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		name := md5ToString(h.Sum(nil))
		d.files[name] = struct{}{}
	}
}

func (d *DataStorage) Save(yourData []byte) error {
	md5Sum := md5.Sum(yourData)
	md5SumString := md5ToString(md5Sum[:16])
	_, ok := d.files[md5SumString]
	if ok {
		return &DontUniqueError{}
	}
	fullName := md5SumString + "." + d.fileExt
	currentPath := filepath.Join(".", d.path, d.fileExt, fullName)
	fileErr := ioutil.WriteFile(currentPath, yourData, os.ModePerm)
	if fileErr != nil {
		return fmt.Errorf("Cant write to file: %s with error %v\n", currentPath, fileErr)
	}
	d.files[md5SumString] = struct{}{}
	return nil
}

func (d *DataStorage) InitStorageFolderStructure(ext string) {
	currentPath := filepath.Join(".", d.path, ext)
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		err := os.MkdirAll(currentPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func md5ToString(md5Sum []byte) string {
	return hex.EncodeToString(md5Sum)
}
