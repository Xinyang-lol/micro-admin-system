package model

type FileInfo struct {
	ID           int64
	OriginalName string
	StoredName   string
	Path         string
	Size         int64
	ContentType  string
	UploaderID   int64
	CreatedAt    string
}
