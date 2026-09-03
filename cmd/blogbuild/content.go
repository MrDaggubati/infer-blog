package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

const (
	sourceDir = "content/blog"
	outputDir = "public/blog"

	publicBaseURL = "https://blog.inferorigins.com"
	siteBaseURL   = "https://www.inferorigins.com"
)

var markdown = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		extension.Typographer,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

type BlogMeta struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Date        string   `json:"date"`
	Author      string   `json:"author"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Featured    bool     `json:"featured"`

	// Card image.
	Image string `json:"image,omitempty"`

	// Article hero image.
	Cover string `json:"cover,omitempty"`

	// Generated URLs.
	Article string `json:"article"`
	URL     string `json:"url"`
}


func buildBlog() error {
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}

	if err := os.MkdirAll(
		outputDir,
		0o755,
	); err != nil {
		return err
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	var posts []BlogMeta

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		postDir := filepath.Join(
			sourceDir,
			entry.Name(),
		)

		meta, err := readMeta(
			filepath.Join(
				postDir,
				"meta.json",
			),
		)
		if err != nil {
			return fmt.Errorf(
				"%s: %w",
				entry.Name(),
				err,
			)
		}

		if err := validateMeta(meta); err != nil {
			return fmt.Errorf(
				"%s: %w",
				entry.Name(),
				err,
			)
		}

		body, err := renderMarkdownFile(
			filepath.Join(
				postDir,
				"article.md",
			),
		)
		if err != nil {
			return fmt.Errorf(
				"%s: %w",
				meta.Slug,
				err,
			)
		}

		postOutputDir := filepath.Join(
			outputDir,
			meta.Slug,
		)

		if err := os.MkdirAll(
			postOutputDir,
			0o755,
		); err != nil {
			return err
		}

		/*
			Process all article images first.

			This generates:
			  cover.webp
			  card.webp
			  *.webp
		*/
		if err := processImages(
			postDir,
		); err != nil {
			return fmt.Errorf(
				"process images for %s: %w",
				meta.Slug,
				err,
			)
		}

		/*
			Copy article image directory
			into public output.
		*/
		if err := copyImages(
			filepath.Join(
				postDir,
				"images",
			),
			filepath.Join(
				postOutputDir,
				"images",
			),
		); err != nil {
			return fmt.Errorf(
				"copy images for %s: %w",
				meta.Slug,
				err,
			)
		}

		/*
			Write generated Markdown HTML fragment.
		*/
		if err := os.WriteFile(
			filepath.Join(
				postOutputDir,
				"article.html",
			),
			[]byte(body),
			0o644,
		); err != nil {
			return err
		}

		/*
			Generated public URLs.
		*/
		meta.Article =
			publicBaseURL +
				"/blog/" +
				meta.Slug +
				"/article.html"

		meta.URL =
			siteBaseURL +
				"/blog/" +
				meta.Slug

		meta.Image = articleAssetURL(
			meta.Slug,
			meta.Image,
		)

		meta.Cover = articleAssetURL(
			meta.Slug,
			meta.Cover,
		)

		posts = append(
			posts,
			meta,
		)
	}

	sortPosts(posts)

	return writeIndex(posts)
}

func readMeta(
	path string,
) (BlogMeta, error) {
	var meta BlogMeta

	data, err := os.ReadFile(path)
	if err != nil {
		return meta, err
	}

	if err := json.Unmarshal(
		data,
		&meta,
	); err != nil {
		return meta, err
	}

	return meta, nil
}

func validateMeta(
	meta BlogMeta,
) error {
	if strings.TrimSpace(meta.Title) == "" {
		return fmt.Errorf(
			"title is required",
		)
	}

	if strings.TrimSpace(meta.Slug) == "" {
		return fmt.Errorf(
			"slug is required",
		)
	}

	if strings.TrimSpace(meta.Date) == "" {
		return fmt.Errorf(
			"date is required",
		)
	}

	if _, err := time.Parse(
		"2006-01-02",
		meta.Date,
	); err != nil {
		return fmt.Errorf(
			"invalid date %q: %w",
			meta.Date,
			err,
		)
	}

	return nil
}

func renderMarkdownFile(
	path string,
) (template.HTML, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer

	if err := markdown.Convert(
		source,
		&out,
	); err != nil {
		return "",
			fmt.Errorf(
				"render markdown: %w",
				err,
			)
	}

	return template.HTML(
		out.String(),
	), nil
}

func articleAssetURL(
	slug string,
	value string,
) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	if strings.HasPrefix(
		value,
		"http://",
	) ||
		strings.HasPrefix(
			value,
			"https://",
		) {
		return value
	}

	return publicBaseURL +
		"/blog/" +
		slug +
		"/" +
		strings.TrimPrefix(
			value,
			"/",
		)
}

func sortPosts(
	posts []BlogMeta,
) {
	sort.SliceStable(
		posts,
		func(i, j int) bool {
			left, leftErr := time.Parse(
				"2006-01-02",
				posts[i].Date,
			)

			right, rightErr := time.Parse(
				"2006-01-02",
				posts[j].Date,
			)

			if leftErr != nil ||
				rightErr != nil {
				return posts[i].Date >
					posts[j].Date
			}

			return left.After(right)
		},
	)
}

func writeIndex(
	posts []BlogMeta,
) error {
	data, err := json.MarshalIndent(
		posts,
		"",
		"  ",
	)
	if err != nil {
		return err
	}

	data = append(
		data,
		'\n',
	)

	return os.WriteFile(
		filepath.Join(
			outputDir,
			"index.json",
		),
		data,
		0o644,
	)
}