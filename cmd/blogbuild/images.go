package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	coverWidth  = 1600
	coverHeight = 900

	cardWidth  = 1000
	cardHeight = 500
)

func processImages(
	postDir string,
) error {
	imagesDir := filepath.Join(
		postDir,
		"images",
	)

	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		ext := strings.ToLower(
			filepath.Ext(name),
		)

		if !supportedImageExtension(ext) {
			continue
		}

		baseName := strings.TrimSuffix(
			name,
			filepath.Ext(name),
		)

		sourcePath := filepath.Join(
			imagesDir,
			name,
		)

		/*
			Generated files should not
			be processed again.
		*/
		if name == "cover.webp" ||
			name == "card.webp" {
			continue
		}

		/*
			Cover source.

			Supported examples:

			cover.png
			cover.jpg
			cover.jpeg
			cover.svg

			cover-source.png
			cover-source.jpg
			cover-source.svg
			cover-source.webp
		*/
		if baseName == "cover" ||
			baseName == "cover-source" {

			if err := processCover(
				sourcePath,
				imagesDir,
			); err != nil {
				return fmt.Errorf(
					"process cover %s: %w",
					sourcePath,
					err,
				)
			}

			continue
		}

		/*
			Existing WebP article images
			require no conversion.
		*/
		if ext == ".webp" {
			continue
		}

		outputPath := filepath.Join(
			imagesDir,
			baseName+".webp",
		)

		if err := convertArticleImage(
			sourcePath,
			outputPath,
			ext,
		); err != nil {
			return fmt.Errorf(
				"convert image %s: %w",
				sourcePath,
				err,
			)
		}
	}

	return nil
}

func supportedImageExtension(
	ext string,
) bool {
	switch ext {
	case ".png",
		".jpg",
		".jpeg",
		".svg",
		".webp":
		return true
	default:
		return false
	}
}

func convertArticleImage(
	sourcePath string,
	outputPath string,
	ext string,
) error {

	/*
		Article diagrams and screenshots are
		encoded losslessly.

		This prevents UI text, diagrams,
		lines and screenshots becoming blurry.
	*/
	args := []string{
		"-y",
		"-i",
		sourcePath,
		"-c:v",
		"libwebp",
		"-lossless",
		"1",
		"-compression_level",
		"6",
		outputPath,
	}

	return runFFmpeg(args)
}

func processCover(
	sourcePath string,
	imagesDir string,
) error {

	coverPath := filepath.Join(
		imagesDir,
		"cover.webp",
	)

	cardPath := filepath.Join(
		imagesDir,
		"card.webp",
	)

	/*
		Article hero.

		16:9
		1600 × 900
	*/
	coverFilter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d",
		coverWidth,
		coverHeight,
		coverWidth,
		coverHeight,
	)

	coverArgs := []string{
		"-y",
		"-i",
		sourcePath,
		"-vf",
		coverFilter,
		"-c:v",
		"libwebp",
		"-quality",
		"92",
		"-compression_level",
		"6",
		coverPath,
	}

	if err := runFFmpeg(
		coverArgs,
	); err != nil {
		return err
	}

	/*
		Blog card.

		2:1
		1000 × 500

		This is also used by the blurred
		card shadow image.
	*/
	cardFilter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d",
		cardWidth,
		cardHeight,
		cardWidth,
		cardHeight,
	)

	cardArgs := []string{
		"-y",
		"-i",
		sourcePath,
		"-vf",
		cardFilter,
		"-c:v",
		"libwebp",
		"-quality",
		"90",
		"-compression_level",
		"6",
		cardPath,
	}

	return runFFmpeg(
		cardArgs,
	)
}

func runFFmpeg(
	args []string,
) error {
	cmd := exec.Command(
		"ffmpeg",
		args...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func copyImages(
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
			if err := copyImages(
				sourcePath,
				targetPath,
			); err != nil {
				return err
			}

			continue
		}

		data, err := os.ReadFile(
			sourcePath,
		)
		if err != nil {
			return err
		}

		if err := os.WriteFile(
			targetPath,
			data,
			0o644,
		); err != nil {
			return err
		}
	}

	return nil
}