package jtt808

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	
	"github.com/flash520/jtt808/protocol"
)

// 分包文件
type partFile struct {
	f    *os.File
	seq  uint16
	path string
}

type partFileSlice []partFile

func (p partFileSlice) Len() int           { return len(p) }
func (p partFileSlice) Less(i, j int) bool { return p[i].seq < p[j].seq }
func (p partFileSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// 多分包文件
type MultipartFile struct {
	IccID uint64
	MsgID protocol.MsgID
	Tag   uint32
	Sum   uint16
}

// 是否完整
func (m MultipartFile) IsFull() bool {
	set := make(map[uint16]struct{})
	for i := 1; i <= int(m.Sum); i++ {
		set[uint16(i)] = struct{}{}
	}
	
	m.walkParts(func(path string, seq uint16) error {
		delete(set, seq)
		return nil
	})
	return len(set) == 0
}

// 合并文件
func (m MultipartFile) Merge() (io.ReadCloser, error) {
	inputs := make(partFileSlice, 0)
	err := m.walkParts(func(path string, seq uint16) error {
		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		inputs = append(inputs, partFile{f: f, seq: seq, path: path})
		return nil
	})
	
	defer func() {
		for _, item := range inputs {
			item.f.Close()
		}
	}()
	
	if err != nil {
		return nil, err
	}
	
	sep := filepath.Separator
	full := fmt.Sprintf("%s%c%d-%d.full", m.rootDir(), sep, m.MsgID, m.Tag)
	w, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	
	sort.Sort(inputs)
	for _, part := range inputs {
		if _, err = io.Copy(w, part.f); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(full, os.O_RDONLY, 0644)
}

// 获取根目录
func (m MultipartFile) rootDir() string {
	sep := filepath.Separator
	dir := fmt.Sprintf("%s%cgo808%c%d", os.TempDir(), sep, sep, m.IccID)
	return dir
}

// 写入分包文件
func (m MultipartFile) Write(seq uint16, data []byte) error {
	dir := m.rootDir()
	if err := os.MkdirAll(dir, 0644); err != nil {
		return err
	}
	sep := filepath.Separator
	name := fmt.Sprintf("%s%c%d-%d-%d.part", dir, sep, m.MsgID, m.Tag, seq)
	return ioutil.WriteFile(name, data, 0644)
}

// 遍历分包文件
func (m MultipartFile) walkParts(walkFn func(path string, seq uint16) error) error {
	dir := m.rootDir()
	id := strconv.Itoa(int(m.MsgID))
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		
		if filepath.Ext(info.Name()) != ".part" {
			return nil
		}
		
		slice := strings.Split(strings.Split(info.Name(), ".")[0], "-")
		if len(slice) != 3 || slice[0] != id {
			return nil
		}
		
		tag, _ := strconv.Atoi(slice[1])
		if uint32(tag) != m.Tag {
			return nil
		}
		
		seq, _ := strconv.Atoi(slice[2])
		return walkFn(path, uint16(seq))
	})
}
