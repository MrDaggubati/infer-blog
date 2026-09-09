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
	switch strings.ToLower(ext) {
	case ".svg",
		".png",
		".jpg",
		".jpeg",
		".webp":
		return true

	default:
		return false
	}
}

/*
	processImages copies images that live beside article.md.

	Source:

	  content/blog/my-post/
	    article.md
	    card.png
	    cover.jpg
	    diagram.svg

	Output:

	  public/blog/my-post/
	    article.html
	    card.png
	    cover.jpg
	    diagram.svg
*/
func processImages(
	postDir string,
	postOutputDir string,
	meta *BlogMeta,
) error {

	if err := copyPostImages(
		postDir,
		postOutputDir,
	); err != nil {
		return err
	}

	/*
		Validate front matter card image.
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
		Validate front matter cover image.
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

/*
	copyPostImages copies supported image files from
	the same directory as article.md.

	article.md itself is ignored.
*/
func copyPostImages(
	srcDir string,
	dstDir string,
) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(
			srcDir,
			entry.Name(),
		)

		targetPath := filepath.Join(
			dstDir,
			entry.Name(),
		)

		/*
			Recurse into subdirectories.

			This allows either:

			  article.md
			  card.png

			or:

			  article.md
			  images/card.png

			or arbitrarily nested image directories.
		*/
		if entry.IsDir() {
			if err := copyPostImages(
				sourcePath,
				targetPath,
			); err != nil {
				return err
			}

			continue
		}

		ext := strings.ToLower(
			filepath.Ext(entry.Name()),
		)

		/*
			Only publish supported image files.
			article.md and other source files are ignored.
		*/
		if !supportedImageExtension(ext) {
			continue
		}

		if err := copyFile(
			sourcePath,
			targetPath,
		); err != nil {
			return fmt.Errorf(
				"copy image %s: %w",
				sourcePath,
				err,
			)
		}
	}

	return nil
}
/*
	resolveImageSource validates an image path referenced
	from front matter.

	Examples:

	  image: card.png
	  cover: cover.jpg
*/
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

	/*
		Remote images can simply be left alone.
	*/
	if strings.HasPrefix(
		value,
		"http://",
	) ||
		strings.HasPrefix(
			value,
			"https://",
		) {
		return value, nil
	}

	clean :=
		filepath.Clean(
			filepath.FromSlash(value),
		)

	/*
		Do not allow article metadata to escape
		the article directory via ../ paths.
	*/
	if clean == ".." ||
		strings.HasPrefix(
			clean,
			".."+string(filepath.Separator),
		) {
		return "", fmt.Errorf(
			"image path escapes article directory: %s",
			value,
		)
	}

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