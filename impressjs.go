package main

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func getFile(url, fileName string) error {
	l.debug("getting file %s to %s", url, fileName)
	f, err := os.Create(fileName)
	if err != nil {
		l.error("error creating file : %s ; caused by : %s", fileName, err)
		return err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		l.error("cannot download %s", err)
		return err
	}
	defer resp.Body.Close()
	n, err := io.Copy(f, resp.Body)
	if err != nil {
		l.error("error writing file %s", err)
		return err
	}
	l.info("downloaded %s (%d bytes)", fileName, n)
	return nil
}

// expecting destFolder is valid abs path
func unzip(src, destFolder string) error {
	l.debug("unzip %s", src)
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		l.error("error reading %s", src)
		return err
	}
	defer zipReader.Close()

	var reader io.ReadCloser
	var writer io.WriteCloser
	defer func() {
		if reader != nil {
			reader.Close()
		}
		if writer != nil {
			writer.Close()
		}
	}()

	for _, f := range zipReader.Reader.File {
		info := f.FileInfo()
		mode := info.Mode()
		isdir := info.IsDir()
		fullName := f.Name
		l.debug("archive item %s (folder? %t; mode: %o)", fullName, isdir, mode)
		path := filepath.Join(destFolder, fullName)
		if isdir {
			l.debug("creating folder...")
			err := os.MkdirAll(path, mode)
			if err != nil {
				l.error("error extracting %s; caused by: %s", fullName, err)
				return err
			}
		} else {
			l.debug("creating file...%s", path)
			reader, err := f.Open()
			if err != nil {
				l.error("error reading archive item %s; caused by: %s", fullName, err)
				return err
			}
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, mode)
			if err != nil {
				l.error("error creating destination file %s; caused by: %S", fullName, err)
				return err
			}

			n, err := io.Copy(writer, reader)
			reader.Close()
			writer.Close()
			l.debug("extracted file %s(size %d bytes)", info.Name(), n)
		}
	}
	return nil
}
