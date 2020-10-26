## FAT32 Storage Driver for [anacrolix/torrent](https://github.com/anacrolix/torrent)

This driver provides a FAT32 compatible storage driver for [anacrolix/torrent](https://github.com/anacrolix/torrent), by storing each piece as it's own file in a split folder structure.

### Performance

Even though this is designed the eliminate most the size and performance restrictions usually presented by the FAT32 file system, the performance still can't really match a real file system like ext4.

The performance of the underlying hardware is also really important, especially write throughput and latency.

On a low end USB 2.0 drive, I could only get 300 kB/sec max with 100 kB/sec avg and occasional lock ups and freezes.

In contrast, on a fairly mid range 32GB SanDisk MicroSD, I could get 3.5 MB/sec max with 3.0 MB/sec average and a very responsive experience.

One thing to bear in mind is that FAT32 always puts data at the beginning of the disk, which is pretty much the opposite of what you want for flash storage that wears out after an amount of write cycles, so it is probably best to go with a high endurance MLC flash if possible, which is widely available for not much more money in the MicroSD format.
