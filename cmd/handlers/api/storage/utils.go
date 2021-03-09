package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func getChecksums(file string) checksums {
	if buf, err := ioutil.ReadFile(file); err != nil {
		fmt.Printf("WARN  error reading file for checksum calculation: %v\n", err)
		return checksums{}
	} else {
		sha1 := sha1.Sum(buf)
		md5 := md5.Sum(buf)
		sha256 := sha256.Sum256(buf)
		return checksums{
			Sha1:   hex.EncodeToString(sha1[:]),
			Md5:    hex.EncodeToString(md5[:]),
			Sha256: hex.EncodeToString(sha256[:]),
		}
	}
}
