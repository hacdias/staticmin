package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
)

func tusPostHandler() handleFunc {
	return withUser(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			if !d.user.Perm.Create || !d.Check(r.URL.Path) {
				return http.StatusForbidden, nil
			}

			dirPath := filepath.Dir(r.URL.Path)
			if _, statErr := d.user.Fs.Stat(dirPath); os.IsNotExist(statErr) {
				if mkdirErr := d.user.Fs.MkdirAll(dirPath, files.PermDir); mkdirErr != nil {
					return http.StatusInternalServerError, err
				}
			}
		case err != nil:
			return errToStatus(err), err
		}

		fileFlags := os.O_CREATE | os.O_WRONLY
		if r.URL.Query().Get("override") == "true" {
			fileFlags |= os.O_TRUNC
		}

		// if file exists
		if file != nil {
			if file.IsDir {
				return http.StatusBadRequest, fmt.Errorf("cannot upload to a directory %s", file.RealPath())
			}
		}

		openFile, err := d.user.Fs.OpenFile(r.URL.Path, fileFlags, files.PermFile)
		if err != nil {
			return errToStatus(err), err
		}
		if err := openFile.Close(); err != nil {
			return errToStatus(err), err
		}

		return http.StatusCreated, nil
	})
}

func tusHeadHandler() handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		w.Header().Set("Cache-Control", "no-store")
		if !d.Check(r.URL.Path) {
			return http.StatusForbidden, nil
		}

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})
		if err != nil {
			return errToStatus(err), err
		}

		w.Header().Set("Upload-Offset", strconv.FormatInt(file.Size, 10))
		w.Header().Set("Upload-Length", "-1")

		return http.StatusOK, nil
	})
}

func tusPatchHandler() handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if !d.user.Perm.Modify || !d.Check(r.URL.Path) {
			return http.StatusForbidden, nil
		}
		if r.Header.Get("Content-Type") != "application/offset+octet-stream" {
			return http.StatusUnsupportedMediaType, nil
		}

		uploadOffset, err := getUploadOffset(r)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("invalid upload offset: %w", err)
		}

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       r.URL.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})

		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			return http.StatusNotFound, nil
		case err != nil:
			return errToStatus(err), err
		}

		switch {
		case file.IsDir:
			return http.StatusBadRequest, fmt.Errorf("cannot upload to a directory %s", file.RealPath())
		case file.Size != uploadOffset:
			return http.StatusConflict, fmt.Errorf(
				"%s file size doesn't match the provided offset: %d",
				file.RealPath(),
				uploadOffset,
			)
		}

		if uploadOffset == 0 {
			err = d.RunBeforeHook("upload", r.URL.Path, "", d.user)
			if err != nil {
				return errToStatus(err), err
			}
		}

		openFile, err := d.user.Fs.OpenFile(r.URL.Path, os.O_WRONLY|os.O_APPEND, files.PermFile)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("could not open file: %w", err)
		}
		defer openFile.Close()

		_, err = openFile.Seek(uploadOffset, 0)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("could not seek file: %w", err)
		}

		defer r.Body.Close()
		bytesWritten, err := io.Copy(openFile, r.Body)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("could not write to file: %w", err)
		}

		w.Header().Set("Upload-Offset", strconv.FormatInt(uploadOffset+bytesWritten, 10))

		if bytesWritten < int64(d.Settings.Tus.ChunkSize) {
			err = d.RunAfterHook("upload", r.URL.Path, "", d.user)
			if err != nil {
				return errToStatus(err), err
			}
		}

		return http.StatusNoContent, nil
	})
}

func getUploadOffset(r *http.Request) (int64, error) {
	uploadOffset, err := strconv.ParseInt(r.Header.Get("Upload-Offset"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid upload offset: %w", err)
	}
	return uploadOffset, nil
}
