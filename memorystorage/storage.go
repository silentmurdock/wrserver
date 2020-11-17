package memorystorage

import (
    //"fmt"
	"io"
    //"log"
    "runtime"
    "sync"

    "github.com/anacrolix/torrent/metainfo"
    "github.com/hashicorp/golang-lru"
)

var maxCount = 16 // Default element count is 16 for 4 MByte (1 << 22) piece length if max memory is 64 MB

var lruStorage, _ = lru.NewWithEvict(maxCount, onEvicted)

var needToDeleteKey = -1

var maxMemorySize int64 = 64 // Maximum memory size in MByte

var memStats runtime.MemStats

var setMaxCount = true

func SetMaxMemorySize(size int64) {
	maxMemorySize = size
}

func GetMaxMemorySize() int64 {
    return maxMemorySize
}

func onEvicted(key interface{}, value interface{}) {
	needToDeleteKey = key.(int)
	//log.Printf("Evicted key: %d\n", needToDeleteKey)
}

// Restricting all I/O through a single mutex, which would stop simultanious read/writes.
func storageWriteAt(mt *memoryTorrent, key int, b []byte, off int64) (int, error) {
    mt.storageMutex.Lock()
    defer mt.storageMutex.Unlock()

    if setMaxCount == true {
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
        	//log.Printf("Memory size / Piece lenght = Element count\n%d / %d = %d\n", maxMemorySize * (1 << 20), mt.pl, elementCount)
        }
        setMaxCount = false
    }

    dataInterface, present := lruStorage.Get(key)
    if present == false {
        dataInterface = []byte{}
    }
    
    ioff := int(off)
	iend := ioff + len(b)
	if len(dataInterface.([]byte)) < iend {
		if len(dataInterface.([]byte)) == ioff {
            if lruStorage.Add(key, append(dataInterface.([]byte), b...)) == true {
                if needToDeleteKey > -1 {
                    mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
                }
            }
			return len(b), nil
		}
		// Add zero bytes to the end of data
        if lruStorage.Add(key, append(dataInterface.([]byte), make([]byte, iend-len(dataInterface.([]byte)))...)) == true {
            if needToDeleteKey > -1 {
                mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
            }
        }
	}

    dataInterface, present = lruStorage.Get(key)
    if present == false {
        dataInterface = []byte{}
    }

	copy(dataInterface.([]byte)[ioff:], b)
	if lruStorage.Add(key, dataInterface.([]byte)) == true {
        if needToDeleteKey > -1 {
            mt.cl.pc.Set(metainfo.PieceKey { mt.ih, needToDeleteKey }, false)
        }
    }

    // Before return check if need to free up some memory
    FreeMemoryPercent(mt, uint64(maxMemorySize), 15)

	return len(b), nil
}

func storageReadAt(mu *sync.Mutex, key int, b []byte, off int64) (int, error) {
    dataInterface, present := lruStorage.Get(key)
    if present == false {
    	dataInterface = []byte{}
    }

    ioff := int(off)
	if len(dataInterface.([]byte)) <= ioff {
		return 0, io.EOF
	}

	n := copy(b, dataInterface.([]byte)[ioff:])
	if n != len(b) {
		return n, io.EOF
	}

	return len(b), nil
}

func storageDelete(mu *sync.Mutex) {
	mu.Lock()
    defer mu.Unlock()

    setMaxCount = true

    lruStorage.Purge()

    needToDeleteKey = -1

    runtime.GC()
}

func FreeMemoryPercent(mt *memoryTorrent, threshold uint64, percent int) {
    runtime.ReadMemStats(&memStats)

    if memStats.Alloc / (1 << 20) > threshold + ((threshold * uint64(percent)) / 100) {
        //log.Printf("Alloc = %v MiB, NumGC = %v\n", memStats.Alloc / (1 << 20), memStats.NumGC)
        var deleteCount = (maxCount * percent) / 100

        if deleteCount == 0 {
            deleteCount++
        }

        for i := 0; i < deleteCount; i++ {
            key, _, ok := lruStorage.RemoveOldest()
            if ok == true {
                if needToDeleteKey > -1 {
                    mt.cl.pc.Set(metainfo.PieceKey { mt.ih, key.(int) }, false)
                }
            }
        }
        //log.Println(lruStorage.Len())
        needToDeleteKey = -1

        runtime.GC()
    }
}