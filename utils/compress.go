package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Zip(fileOrDir, destFile string) error {
	zipDest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer zipDest.Close()

	zipWriter := zip.NewWriter(zipDest)
	defer zipWriter.Close()

	err = filepath.Walk(fileOrDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(fileOrDir)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fh, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fh.Close()

			_, err = io.Copy(writer, fh)
		}
		return err
	}) // Walk
	return err
}

func UnZip(zipFile, destPath string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destPath, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			fh, err := f.Open()
			if err != nil {
				return err
			}
			defer fh.Close()

			ofh, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer ofh.Close()

			_, err = io.Copy(ofh, fh)
			if err != nil {
				return err
			}
		} // else
	} // for
	return nil
}

func TarGz(src, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	gzWriter := gzip.NewWriter(d)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	fh, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fh.Close()

	err =gz(fh, "", tarWriter)
	return err
}

func gz(f *os.File, prefix string, tw *tar.Writer) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		subDirInfos, err := f.Readdir(-1)
		if err != nil {
			return err
		}
		for _, subDirInfo := range subDirInfos {
			subDirFh, err := os.Open(path.Join(f.Name(), subDirInfo.Name()))
			if err != nil {
				return err
			}
			defer subDirFh.Close()

			err = gz(subDirFh, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = prefix + "/" + header.Name
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnTarGz(tarGz, dest string) error {
	fh, err := os.Open(tarGz)
	if err != nil {
		return err
	}
	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		hdr, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		filename := path.Join(dest, hdr.Name)
		f, err := EnsureOpenFile(filename)
		if err != nil {
			return err
		}
		defer f.Close() // TAKE CARE!!

		_, err = io.Copy(f, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnsureOpenFile(full string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(full)[0:strings.LastIndex(full, string(os.PathSeparator))]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(full)
}
