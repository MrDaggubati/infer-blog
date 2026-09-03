package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)
func supportedImageExtension(
	ext string,
) bool {
	switch ext {
	case ".svg",
		".png",
		".jpg",
		".jpeg":
		return true

	default:
		return false
	}
}

func copyReferencedImage(
	postDir string,
	postOutputDir string,
	value string,
) error {

	source, err :=
		resolveImageSource(
			postDir,
			value,
		)
	if err != nil {
		return err
	}

	target := filepath.Join(
		postOutputDir,
		filepath.FromSlash(value),
	)

	return copyFile(
		source,
		target,
	)
}

func processImages(
	postDir string,
	postOutputDir string,
	meta *BlogMeta,
) error {
	sourceImagesDir := filepath.Join(
		postDir,
		"images",
	)

	outputImagesDir := filepath.Join(
		postOutputDir,
		"images",
	)

	/*
		Copy every authored image.

		This covers:
		- front matter image
		- front matter cover
		- images referenced inside article.md
	*/
	if err := copyOriginalImages(
		sourceImagesDir,
		outputImagesDir,
	); err != nil {
		return err
	}

	/*
		Validate front matter image reference.
	*/
	if strings.TrimSpace(meta.Image) != "" {
		if _, err := resolveImageSource(
			postDir,
			meta.Image,
		); err != nil {
			return fmt.Errorf(
				"invalid card image: %w",
				err,
			)
		}
	}

	/*
		Validate front matter cover reference.
	*/
	if strings.TrimSpace(meta.Cover) != "" {
		if _, err := resolveImageSource(
			postDir,
			meta.Cover,
		); err != nil {
			return fmt.Errorf(
				"invalid cover image: %w",
				err,
			)
		}
	}

	return nil
}

func resolveImageSource(
	postDir string,
	value string,
) (string, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", fmt.Errorf(
			"image path is empty",
		)
	}

	if strings.HasPrefix(
		value,
		"http://",
	) ||
		strings.HasPrefix(
			value,
			"https://",
		) {
		return "", fmt.Errorf(
			"remote image sources are not supported: %s",
			value,
		)
	}

	clean := filepath.Clean(
		filepath.FromSlash(value),
	)

	path := filepath.Join(
		postDir,
		clean,
	)

	info, err := os.Stat(path)
	if err != nil {
		return "",
			fmt.Errorf(
				"source image %q: %w",
				path,
				err,
			)
	}

	if info.IsDir() {
		return "", fmt.Errorf(
			"%s is a directory",
			path,
		)
	}

	ext := strings.ToLower(
		filepath.Ext(path),
	)

	if !supportedImageExtension(ext) {
		return "", fmt.Errorf(
			"unsupported image format %q",
			ext,
		)
	}

	return path, nil
}


func copyOriginalImages(
	src string,
	dst string,
) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := os.MkdirAll(
		dst,
		0o755,
	); err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(
			src,
			entry.Name(),
		)

		targetPath := filepath.Join(
			dst,
			entry.Name(),
		)

		if entry.IsDir() {
			if err := copyOriginalImages(
				sourcePath,
				targetPath,
			); err != nil {
				return err
			}

			continue
		}

		if err := copyFile(
			sourcePath,
			targetPath,
		); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(
	src string,
	dst string,
) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := os.MkdirAll(
		filepath.Dir(dst),
		0o755,
	); err != nil {
		return err
	}

	target, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(
		target,
		source,
	); err != nil {
		target.Close()
		return err
	}

	return target.Close()
}