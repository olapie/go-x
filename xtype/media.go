package xtype

type Audio struct {
	URL      string `json:"url,omitempty"`
	Format   string `json:"format,omitempty"`
	Duration int32  `json:"duration,omitempty"`
	Size     int32  `json:"size,omitempty"`
	Name     string `json:"name,omitempty"`
	Data     []byte `json:"data,omitempty"`
}

type Image struct {
	URL       string `json:"url,omitempty"`
	Width     int32  `json:"width,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Format    string `json:"format,omitempty"`
	Size      int32  `json:"size,omitempty"`
	Name      string `json:"name,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Data      []byte `json:"data,omitempty"`
}

type Video struct {
	URL      string `json:"url,omitempty"`
	Format   string `json:"format,omitempty"`
	Duration int32  `json:"duration,omitempty"`
	Size     int32  `json:"size,omitempty"`
	Image    *Image `json:"image,omitempty"`
	Name     string `json:"name,omitempty"`
	Data     []byte `json:"data,omitempty"`
}

type PhotoID struct {
	Type       string `json:"type,omitempty"`
	Front      string `json:"front,omitempty"`
	Back       string `json:"back,omitempty"`
	Number     string `json:"number,omitempty"`
	IssueTime  int64  `json:"issue_time,omitempty"`
	ExpireTime int64  `json:"expire_time,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
}

type File struct {
	URL    string `json:"url,omitempty"`
	Name   string `json:"name,omitempty"`
	Size   int32  `json:"size,omitempty"`
	Format string `json:"format,omitempty"`
	Data   []byte `json:"data,omitempty"`
}

type WebPage struct {
	Title   string `json:"title,omitempty"`
	Summary string `json:"summary,omitempty"`
	Image   *Image `json:"image,omitempty"`
	URL     string `json:"url,omitempty"`
}
