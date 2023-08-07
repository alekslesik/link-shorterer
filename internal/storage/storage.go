package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
    ErrURLNotFound = errors.New("url not found")
    ErrURLExists   = errors.New("url exists")
)

func CreateStorageFile(filePath string)  {
	const op = "storage.sqlite.createStorageFile"
	// Получаем директорию из полного пути
	dir := filepath.Dir(filePath)

	// Проверяем, существует ли папка. Если нет, создаем ее
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("%s: creating dir error %s\n",op, err)
			return
		}
	}

	// Открываем файл для записи. Если файл не существует, он будет создан.
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("%s: creating file error %s\n",op, err)
		return
	}
	defer file.Close() // По окончании работы с файлом, закрываем его

	fmt.Printf("%s file has created\n", file.Name())
}