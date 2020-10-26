package fat32storage

import (
    "log"
    "fmt"
    "sync"
    "path/filepath"
    "github.com/anacrolix/torrent/storage"
    "github.com/anacrolix/torrent/metainfo"
)

type fat32Client struct {
    filePath        string

    pc              storage.PieceCompletion
}

type fat32Torrent struct {
    cl              *fat32Client
    pl              int64
    ih              metainfo.Hash

    storagePath     string
    storageMutex    *sync.Mutex
}

type fat32Piece struct {
    trt     *fat32Torrent

    p       metainfo.Piece
}

func NewFat32Storage(filePath string) storage.ClientImpl {
    ret := &fat32Client{
        filePath:   filePath,
        pc:         storage.NewMapPieceCompletion(),
    }

    return ret
}

func (me *fat32Client) Close() error {
    return nil
}

func (me *fat32Client) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
    return &fat32Torrent{
        cl:             me,
        pl:             info.PieceLength,
        ih:             infoHash,

        storagePath:    filepath.Join(me.filePath, infoHash.String() + ".diskv"),
        storageMutex:   &sync.Mutex{},
    }, nil
}

func (me *fat32Torrent) Piece(p metainfo.Piece) storage.PieceImpl {
    return &fat32Piece{
        trt:        me,
        p:          p,
    }
}

func (me *fat32Torrent) Close() error {
    storageDeletePath(me.storageMutex, me.storagePath)

    return nil
}

func (sp *fat32Piece) pieceKey() metainfo.PieceKey {
    return metainfo.PieceKey{sp.trt.ih, sp.p.Index()}
}

func (sp *fat32Piece) chunkKey(index int) (string) {
    return fmt.Sprintf("%d", index)
}

func (sp *fat32Piece) Completion() storage.Completion {
    ret, _ := sp.trt.cl.pc.Get(sp.pieceKey())
    return ret
}

func (sp *fat32Piece) MarkComplete() error {
    sp.trt.cl.pc.Set(sp.pieceKey(), true)
    return nil
}

func (sp *fat32Piece) MarkNotComplete() error {
    sp.trt.cl.pc.Set(sp.pieceKey(), false)
    return nil
}

func (sp *fat32Piece) ReadAt(b []byte, off int64) (n int, err error) {
    ci := sp.p.Index()
    bToRead := sp.trt.pl - off
    //log.Printf("Got read for chunk (%d) offset (%d).", ci, off)
    for len(b) != 0 {
        ck := sp.chunkKey(int(ci))
        var rLen int
        if len(b) < int(bToRead) {
            rLen = len(b)
        } else {
            rLen = int(bToRead)
        }
        _b := make([]byte, rLen)
        i, rerr := storageReadAt(sp.trt.storageMutex, sp.trt.storagePath, string(ck[:]), _b, off)
        //log.Printf("Doing read for chunk (%d) offset (%d) for (%d) bytes and got (%d) bytes.", ci, off, rLen, i)
        n1 := copy(b, _b[:i])
        off = 0
        ci++
        b = b[n1:]
        n += n1
        if rerr != nil {
            log.Printf("Error Reading During Read: %s", rerr)
            err = rerr
            return
        }
    }
    return
}

func (sp *fat32Piece) WriteAt(b []byte, off int64) (n int, err error) {
    ci := sp.p.Index()
    //log.Printf("At chunk (%d), got bytes (%d) and offset (%d).", ci, len(b), off)
    bToWrite := sp.trt.pl - off
    var btw int
    for len(b) != 0 {
        if len(b) > int(bToWrite) {
            btw = int(bToWrite)
        } else {
            btw = len(b)
        }
        ck := sp.chunkKey(int(ci))
        n1, werr := storageWriteAt(sp.trt.storageMutex, sp.trt.storagePath, string(ck[:]), b[:btw], off)
        //log.Printf("Writing (%d) bytes [confirm %d] to chunk (%d) - total bytes (%d) - offset (%d) - written (%d) bytes.", btw, len(b[:btw]), ci, len(b), off, n1)
        if werr != nil {
            log.Printf("Error Writing During Write: %s", werr)
            err = werr
            return
        }
        if n1 > len(b) {
            break
        }
        b = b[n1:]
        off = 0
        bToWrite = sp.trt.pl - off
        ci++
        n += n1
    }
    return
}
