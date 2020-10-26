package fat32storage

import (
    "os"
    "fmt"
    "sync"
    "path/filepath"
)

const transformBlockSize = 4

const storageFilePerm = 0666
const storagePathPerm = 0777

/*
* The Mutex
*******************************
* You may have noticed we are restricting all I/O
* through a single mutex, which would stop simultanious
* read/writes and multiple open file handles.
*
* While on a very good file system this would make things very slow,
* on FAT32 and slow flash media, when writing to lots of different locations
* there are very often lock ups in the OS and the performance is awful.
*
* Using this locking method, the performance was more predictable
* and actually gave faster throughput than trying to do multiple I/Os at once.
* Although, fragmentation seems to make quite a big difference.
* On Windows at-least, the performance between a freshly formatted drive
* and a previously used one was very big, tested with 1 low end USB drive.
*
* One thing to investigate in future might be 2 mutexes, one for reads and one for writes.
*/

func storageWriteAt(mu *sync.Mutex, basePath, key string, b []byte, off int64) (int, error) {
    mu.Lock()
    defer mu.Unlock()

	if err := storageEnsurePath(basePath, key); err != nil {
		return 0, fmt.Errorf("ensure path: %s", err)
	}

	mode := os.O_WRONLY | os.O_CREATE
	f, err := os.OpenFile(storageCompleteFilename(basePath, key), mode, storageFilePerm)
	if err != nil {
		return 0, fmt.Errorf("open file: %s", err)
	}

	var i int
	if i, err = f.WriteAt(b, off); err != nil {
		f.Close() // error deliberately ignored
		return 0, fmt.Errorf("i/o copy: %s", err)
	}

	if err := f.Close(); err != nil {
		return 0, fmt.Errorf("file close: %s", err)
	}

	return i, nil
}

func storageReadAt(mu *sync.Mutex, basePath, key string, b []byte, off int64) (int, error) {
    mu.Lock()
    defer mu.Unlock()

	filename := storageCompleteFilename(basePath, key)

	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, os.ErrNotExist
	}

	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.ReadAt(b, off)
}

func storageDeletePath(mu *sync.Mutex, path string) {
    mu.Lock()
    defer mu.Unlock()

    os.RemoveAll(path)
}

func storagePathFor(basePath, key string) string {
	return filepath.Join(basePath, filepath.Join(storageTransform(key)...))
}

func storageEnsurePath(basePath, key string) error {
	return os.MkdirAll(storagePathFor(basePath, key), storagePathPerm)
}

func storageCompleteFilename(basePath, key string) string {
	return filepath.Join(storagePathFor(basePath, key), key)
}

func storageTransform(s string) []string {
	var (
		sliceSize = len(s) / transformBlockSize
		pathSlice = make([]string, sliceSize)
	)
	for i := 0; i < sliceSize; i++ {
		from, to := i*transformBlockSize, (i*transformBlockSize)+transformBlockSize
		pathSlice[i] = s[from:to]
	}
	return pathSlice
}
