package memorystorage

import (
	"io"
    "runtime"
    "sync"
    //"fmt"

    "github.com/anacrolix/torrent/metainfo"
    "github.com/hashicorp/golang-lru"
)

var maxCount = 16 // Default element count is 16 for 4 MByte (1 << 22) piece length if max memory is 64 MB
                  // Minimum element count is 8 for 8 MByte (1 << 23) piece length if max memory is 64 MB

var lruStorage, _ = lru.NewWithEvict(maxCount, onEvicted)

var needToDeleteKey = -1
var neetToDeleteValue = []byte{}

var maxMemorySize int64 = 64 // Max memory size in MByte

func SetMaxMemorySize(size int64) {
	maxMemorySize = size
}

func GetMaxMemorySize() int64 {
    return maxMemorySize
}

func onEvicted(key interface{}, value interface{}) {
	needToDeleteKey = key.(int)
    neetToDeleteValue = value.([]byte)
    //runtime.GC()
	//fmt.Printf("Evicted key: %d\n", needToDeleteKey)
}

// Restricting all I/O through a single mutex, which would stop simultanious read/writes.
func storageWriteAt(mt *memoryTorrent, key int, b []byte, off int64) (int, error) {
    mt.storageMutex.Lock()
    defer mt.storageMutex.Unlock()

    elementCount := maxCount
    if mt.pl <= (1 << 23) {
    	if mt.pl >= (1 << 13) {
    		elementCount = int((maxMemorySize * (1 << 20)) / mt.pl)
    	} else {
    		elementCount = (1 << 13)
    	}
    }

    if maxCount != elementCount {
    	lruStorage.Resize(elementCount)
    	maxCount = elementCount
    	//fmt.Printf("Memory size / Piece lenght = Element count\n%d / %d = %d\n", maxMemorySize * (1 << 20), mt.pl, elementCount)
    }

    data := []byte{}
    dataInterface, present := lruStorage.Get(key)
    if present == true {
    	data = dataInterface.([]byte)
    }
    

    ioff := int(off)
	iend := ioff + len(b)
	if len(data) < iend {
		if len(data) == ioff {
            if lruStorage.Add(key, append(data, b...)) == true {
                if needToDeleteKey == 0 {
                    if lruStorage.Add(0, neetToDeleteValue) == true && needToDeleteKey > -1 {
                        mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
                    }
                } else if needToDeleteKey > -1 {
                    mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
                }
            }
			return len(b), nil
		}
		zero := make([]byte, iend-len(data))
        if lruStorage.Add(key, append(data, zero...)) == true {
            if needToDeleteKey == 0 {
                if lruStorage.Add(0, neetToDeleteValue) == true && needToDeleteKey > -1 {
                    mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
                }
            } else if needToDeleteKey > -1 {
                mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
            }
        }
	}

	data = []byte{}
    dataInterface, present = lruStorage.Get(key)
    if present == true {
    	data = dataInterface.([]byte)
    }

	copy(data[ioff:], b)
	if lruStorage.Add(key, data) == true {
        if needToDeleteKey == 0 {
            if lruStorage.Add(0, neetToDeleteValue) == true && needToDeleteKey > -1 {
                mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
            }
        } else if needToDeleteKey > -1 {
            mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
        }
    }

	return len(b), nil
}

func storageReadAt(mu *sync.Mutex, key int, b []byte, off int64) (int, error) {
    /*mu.Lock()
    defer mu.Unlock()*/

    data := []byte{}
    dataInterface, present := lruStorage.Get(key)
    if present == true {
    	data = dataInterface.([]byte)
    }

    ioff := int(off)
	if len(data) <= ioff {
		return 0, io.EOF
	}

	n := copy(b, data[ioff:])
	if n != len(b) {
		return n, io.EOF
	}

	return len(b), nil
}

func storageDelete(mu *sync.Mutex) {
	mu.Lock()
    defer mu.Unlock()

    lruStorage.Purge()
    runtime.GC()
}