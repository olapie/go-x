package xmobile

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.olapie.com/x/xmobile/nomobile"
	"go.olapie.com/x/xtest"
)

type SortFieldType = int

const (
	SortByModTime SortFieldType = iota
	SortByName
)

type DirInfo interface {
	FileInfo
	NumFile() int
	File(i int) FileInfo

	InsertAt(f FileInfo, index int)
	Remove(fileID string) bool

	Move(fileID, dirID string) bool
	FileByID(id string, recursive bool) FileInfo
	FileByName(name string, recursive bool) FileInfo

	Sort(field SortFieldType, asc bool)

	ReadFiles(typ FileType) *FileInfoList
}

type FileInfo interface {
	// GetID is friendly to swift syntax
	GetID() string
	Parent() DirInfo
	Name() string
	Size() int64
	ModTime() int64
	MIMEType() string
	AsDir() DirInfo
}

var _ FileInfo = (*FileTreeNode)(nil)

var _ nomobile.FileEntry = (*virtualEntry)(nil)

type virtualEntry struct {
	ID          string
	EntryName   string
	SubEntryIDs []string
}

func (v *virtualEntry) GetID() string {
	return v.ID
}

func (v *virtualEntry) Name() string {
	return v.EntryName
}

func (v *virtualEntry) IsDir() bool {
	return true
}

func (v *virtualEntry) Size() int64 {
	return 0
}

func (v *virtualEntry) ModTime() int64 {
	return 0
}

func (v *virtualEntry) MIMEType() string {
	return ""
}

func (v *virtualEntry) SubIDs() []string {
	return v.SubEntryIDs
}

type FileTreeNode struct {
	entry  nomobile.FileEntry
	parent *FileTreeNode
	files  []*FileTreeNode
}

func (f *FileTreeNode) Entry() nomobile.FileEntry {
	return f.entry
}

func (f *FileTreeNode) Parent() DirInfo {
	if f.parent == nil {
		return nil
	}
	return f.parent.AsDir()
}

func (f *FileTreeNode) GetID() string {
	return f.entry.GetID()
}

func (f *FileTreeNode) Name() string {
	return f.entry.Name()
}

func (f *FileTreeNode) AsDir() DirInfo {
	if f.entry.IsDir() {
		return f
	}
	return nil
}

func (f *FileTreeNode) Size() int64 {
	return f.entry.Size()
}

func (f *FileTreeNode) ModTime() int64 {
	return f.entry.ModTime()
}

func (f *FileTreeNode) MIMEType() string {
	return f.entry.MIMEType()
}

func (f *FileTreeNode) NumFile() int {
	return len(f.files)
}

func (f *FileTreeNode) File(i int) FileInfo {
	return f.files[i]
}

func (f *FileTreeNode) ReadFiles(typ FileType) *FileInfoList {
	l := new(FileInfoList)
	for _, sub := range f.files {
		if GetFileType(sub)&typ != 0 {
			l.List = append(l.List, sub)
		}
	}
	return l
}

func (f *FileTreeNode) InsertAt(sub FileInfo, index int) {
	node := sub.(*FileTreeNode)
	if sub.GetID() != "" && node.FileByID(f.GetID(), true) != nil {
		panic("recycle file reference")
	}

	if node.parent != nil {
		node.parent.Remove(node.GetID())
		node.parent = nil
	}

	if index <= 0 {
		f.files = append([]*FileTreeNode{node}, f.files...)
	} else if index >= len(f.files) {
		f.files = append(f.files, node)
	} else {
		f.files = append(f.files, node)
		copy(f.files[index+1:], f.files[index:len(f.files)-1])
		f.files[index] = node
	}
	node.parent = f
}

func (f *FileTreeNode) Remove(id string) bool {
	for i, v := range f.files {
		if v.GetID() == id {
			v.parent = nil
			f.files = append(f.files[:i], f.files[i+1:]...)
			return true
		}

		if dir := v.AsDir(); dir != nil {
			if dir.Remove(id) {
				return true
			}
		}
	}
	return false
}

// FileByID searches a descendant node
func (f *FileTreeNode) FileByID(id string, recursive bool) FileInfo {
	if id != "" && f.GetID() == id {
		return f
	}

	for _, fi := range f.files {
		if fi.GetID() == id {
			return fi
		}

		if dir := fi.AsDir(); dir != nil && recursive {
			if sub := dir.FileByID(id, recursive); sub != nil {
				return sub
			}
		}
	}

	if f.GetID() == id {
		return f
	}

	return nil
}

func (f *FileTreeNode) FileByName(name string, recursive bool) FileInfo {
	for _, fi := range f.files {
		if fi.Name() == name {
			return fi
		}
		if dir := fi.AsDir(); dir != nil && recursive {
			if sub := dir.FileByName(name, recursive); sub != nil {
				return sub
			}
		}
	}
	return nil
}

func (f *FileTreeNode) DirInfo() DirInfo {
	if f.entry.IsDir() {
		return f
	}
	return nil
}

func (f *FileTreeNode) Move(fileID, dirID string) bool {
	if fileID == dirID {
		log.Printf("Move: same file %s, %s\n", fileID, dirID)
		return false
	}

	fi := f.FileByID(fileID, true)
	if fi == nil {
		log.Println("Move: no file", fileID)
		return false
	}

	if fi.Parent() != nil && fi.Parent().GetID() == dirID {
		log.Println("Move: already in dir", dirID)
		return true
	}

	if fiDir := fi.AsDir(); fiDir != nil {
		if fiDir.FileByID(dirID, true) != nil {
			log.Printf("Move: %s is under %s\n", dirID, fileID)
			return false
		}
	}

	dirFile := f.FileByID(dirID, true)
	if dirFile == nil {
		log.Println("Move: no dir", dirID)
		return false
	}

	dir := dirFile.AsDir()
	if dir == nil {
		log.Println("Move: cannot convert to dir", dirID)
		return false
	}
	dir.InsertAt(fi, dir.NumFile()+1)
	return true
}

func (f *FileTreeNode) Sort(typ SortFieldType, asc bool) {
	switch typ {
	case SortByName:
		f.sortSubsByName(asc)
	case SortByModTime:
		f.sortSubsByModTime(asc)
	default:
		log.Println("unsupported type", typ)
		break
	}
}

func (f *FileTreeNode) sortSubsByModTime(asc bool) {
	for _, fi := range f.files {
		if fi.AsDir() != nil {
			fi.sortSubsByModTime(asc)
		}
	}

	sort.Slice(f.files, func(i, j int) bool {
		fi, fj := f.files[i], f.files[j]
		if fi.ModTime() == fj.ModTime() {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		}
		return asc == (fi.ModTime() < fj.ModTime())
	})
}

func (f *FileTreeNode) sortSubsByName(asc bool) {
	for _, fi := range f.files {
		if fi.AsDir() != nil {
			fi.sortSubsByName(asc)
		}
	}

	sort.Slice(f.files, func(i, j int) bool {
		fi, fj := f.files[i], f.files[j]
		if fi.Name() == fj.Name() {
			return asc == (fi.ModTime() < fj.ModTime())
		}
		return asc == (fi.Name() == fj.Name())
	})
}

func NewVirtualDir(id, name string) DirInfo {
	return &FileTreeNode{
		entry: &virtualEntry{
			ID:        id,
			EntryName: name,
		},
	}
}

func FileInfoFromEntry(entry nomobile.FileEntry) FileInfo {
	return &FileTreeNode{
		entry: entry,
	}
}

func BuildFileTree(entries []nomobile.FileEntry) DirInfo {
	root := NewVirtualDir("", "").(*FileTreeNode)
	if len(entries) == 0 {
		return root
	}

	idToEntry := make(map[string]nomobile.FileEntry)
	for _, f := range entries {
		idToEntry[f.GetID()] = f
	}

	idToNode := make(map[string]*FileTreeNode)
	for _, e := range entries {
		buildFileTreeNode(nil, e.GetID(), idToEntry, idToNode)
	}

	for _, node := range idToNode {
		if node.parent == nil {
			root.files = append(root.files, node)
			node.parent = root
		}
	}

	root.Sort(SortByModTime, false)
	return root
}

func buildFileTreeNode(parent *FileTreeNode, id string, idToEntry map[string]nomobile.FileEntry, result map[string]*FileTreeNode) {
	node := result[id]
	if node != nil {
		if parent != nil {
			parent.InsertAt(node, parent.NumFile()+1)
		}
		return
	}

	entry, exists := idToEntry[id]
	if !exists {
		log.Println("Warn: no entry for file id", id)
		return
	}
	delete(idToEntry, id)

	node = &FileTreeNode{
		entry: entry,
	}
	if parent != nil {
		parent.InsertAt(node, parent.NumFile()+1)
	}
	result[node.GetID()] = node
	if node.AsDir() == nil {
		return
	}

	for _, subID := range entry.SubIDs() {
		buildFileTreeNode(node, subID, idToEntry, result)
	}
}

var _ nomobile.FileEntry = (*mockFileEntry)(nil)

type mockFileEntry struct {
	id       string
	name     string
	isDir    bool
	size     int64
	modTime  int64
	mimeType string
	subIDs   []string
}

func (m *mockFileEntry) GetID() string {
	return m.id
}

func (m *mockFileEntry) Name() string {
	return m.name
}

func (m *mockFileEntry) IsDir() bool {
	return m.isDir
}

func (m *mockFileEntry) Size() int64 {
	return m.size
}

func (m *mockFileEntry) ModTime() int64 {
	return m.modTime
}

func (m *mockFileEntry) MIMEType() string {
	return m.mimeType
}

func (m *mockFileEntry) SubIDs() []string {
	return m.subIDs
}

func NewMockFileInfo(isDir bool) FileInfo {
	if isDir {
		return &FileTreeNode{
			entry: &mockFileEntry{
				id:      uuid.NewString(),
				name:    "dir" + xtest.RandomString(10),
				isDir:   true,
				modTime: time.Now().Unix(),
				subIDs:  []string{uuid.NewString(), uuid.NewString()},
			},
		}
	}
	return &FileTreeNode{
		entry: &mockFileEntry{
			id:      uuid.NewString(),
			name:    "file" + xtest.RandomString(10),
			modTime: time.Now().Unix(),
		},
	}
}

type FileInfoList struct {
	List []FileInfo
}

func NewFileInfoList() *FileInfoList {
	return new(FileInfoList)
}

func (l *FileInfoList) Len() int {
	return len(l.List)
}

func (l *FileInfoList) Get(i int) FileInfo {
	return l.List[i]
}

func (l *FileInfoList) Delete(i int) {
	l.List = append(l.List[:i], l.List[i+1:]...)
}

type FileType = int64

const (
	FileTypeDir FileType = 1 << iota
	FileTypeText
	FileTypeAudio
	FileTypeVideo
	FileTypeImage
	FileTypeUnknown

	FileTypeNotDir = FileTypeText | FileTypeAudio | FileTypeVideo | FileTypeImage | FileTypeUnknown
)

func GetFileType(f FileInfo) FileType {
	if f.AsDir() != nil {
		return FileTypeDir
	}

	if IsMIMEAudio(f.MIMEType()) {
		return FileTypeAudio
	}

	if IsMIMEVideo(f.MIMEType()) {
		return FileTypeVideo
	}

	if IsMIMEImage(f.MIMEType()) {
		return FileTypeImage
	}

	if IsMIMEText(f.MIMEType()) {
		return FileTypeText
	}

	return FileTypeUnknown
}
