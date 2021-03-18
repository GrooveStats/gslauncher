package unlocks

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type DownloadInfo struct {
	TotalSize  int64
	Downloaded int64
	Error      error
}

type Download struct {
	Url      string
	Filename string
	Progress chan DownloadInfo

	cancel bool
	stream io.ReadCloser
}

func (download *Download) Cancel() {
	if download.stream != nil {
		download.stream.Close()
	}

	download.cancel = true
}

func Fetch(url, filename string) *Download {
	download := &Download{
		Url:      url,
		Filename: filename,
		Progress: make(chan DownloadInfo),
	}

	go fetch(download)

	return download
}

func fetch(download *Download) {
	filename := download.Filename + ".part"

	defer close(download.Progress)

	info := DownloadInfo{
		TotalSize:  -1,
		Downloaded: 0,
	}
	download.Progress <- info

	outFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		info.Error = err
		download.Progress <- info
		return
	}
	defer func() {
		outFile.Close()
		os.Remove(filename)
	}()

	resp, err := http.Get(download.Url)
	if err != nil {
		info.Error = err
		download.Progress <- info
		return
	}
	defer func() {
		download.stream = nil
		resp.Body.Close()
	}()

	download.stream = resp.Body

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		info.Error = fmt.Errorf("HTTP status code %s", resp.Status)
		download.Progress <- info
		return
	}

	info.TotalSize = resp.ContentLength
	download.Progress <- info

	for {
		written, err := io.CopyN(outFile, resp.Body, 32*1024)
		info.Downloaded += written
		download.Progress <- info

		if download.cancel {
			err = fmt.Errorf("Download cancelled")
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			info.Error = err
			download.Progress <- info
			return
		}
	}

	outFile.Close()

	err = os.Rename(filename, download.Filename)
	if err != nil {
		info.Error = err
		download.Progress <- info
		return
	}
}
