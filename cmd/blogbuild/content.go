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
	"gopkg.in/yaml.v3"	
)

var (
	publicRootDir = envOrDefault("PUBLIC_DIR", "public")
	sourceDir     = envOrDefault("CONTENT_DIR", "content/blog")
	outputDir     = filepath.Join(publicRootDir, "blog")

	publicBaseURL = strings.TrimRight(
		envOrDefault("BLOG_BASE_URL", "https://blog.inferorigins.com"),
		"/",
	)
	siteBaseURL = strings.TrimRight(
		envOrDefault("MAIN_SITE_URL", "https://www.inferorigins.com"),
		"/",
	)
	blogCNAME = envOrDefault("BLOG_CNAME", "blog.inferorigins.com")
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
	Title string `json:"title" yaml:"title"`
	Slug string `json:"slug" yaml:"slug"`
	Date string `json:"date" yaml:"date"`
	Author string `json:"author" yaml:"author"`
	Summary string `json:"summary" yaml:"summary"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Tags []string `json:"tags" yaml:"tags"`
	Featured bool `json:"featured" yaml:"featured"`

	Image string `json:"image,omitempty" yaml:"image,omitempty"`
	Cover string `json:"cover,omitempty" yaml:"cover,omitempty"`

	/*
		Generated fields.
		These do NOT come from article.md.
	*/
	Article string `json:"article"`

	URL string `json:"url"`
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

		meta, body, err := readArticle(
			filepath.Join(
				postDir,
				"article.md",
			),
		)
		if err != nil {
			return fmt.Errorf(
				"%s: %w",
				entry.Name(),
				err,
			)
		}
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
			Process/copy images.

			Important:
			processImages receives &meta and rewrites:

			  meta.Image -> images/card.webp
			  meta.Cover -> images/cover.webp
		*/
		if err := processImages(
			postDir,
			postOutputDir,
			&meta,
		); err != nil {
			return fmt.Errorf(
				"process images for %s: %w",
				meta.Slug,
				err,
			)
		}

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

		meta.Article =
			publicBaseURL +
				"/blog/" +
				meta.Slug +
				"/article.html"

		meta.URL =
			siteBaseURL +
				"/blog/" +
				meta.Slug +
				"/article.html"
		/*
			At this point processImages has already changed
			Image/Cover to their generated WebP filenames.
		*/
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

	if err := writeIndex(posts); err != nil {
		return err
	}

	if err := writeRootIndex(posts); err != nil {
		return err
	}

	if err := writeCNAME(); err != nil {
		return err
	}

	return nil
}

func readArticle(
	path string,
) (BlogMeta, template.HTML, error) {
	var meta BlogMeta

	source, err := os.ReadFile(path)
	if err != nil {
		return meta, "", err
	}

	/*
		Normalize line endings so front matter parsing
		works consistently across Linux/Windows.
	*/
	source = bytes.ReplaceAll(
		source,
		[]byte("\r\n"),
		[]byte("\n"),
	)

	if !bytes.HasPrefix(
		source,
		[]byte("---\n"),
	) {
		return meta, "",
			fmt.Errorf(
				"article is missing YAML front matter",
			)
	}

	/*
		Remove opening --- and find closing ---.
	*/
	remaining :=
		source[len("---\n"):]

	end :=
		bytes.Index(
			remaining,
			[]byte("\n---\n"),
		)

	if end == -1 {
		return meta, "",
			fmt.Errorf(
				"article front matter is not closed",
			)
	}

	frontMatter :=
		remaining[:end]

	body :=
		remaining[
			end+len("\n---\n"):
		]

	if err := yaml.Unmarshal(
		frontMatter,
		&meta,
	); err != nil {
		return meta, "",
			fmt.Errorf(
				"parse front matter: %w",
				err,
			)
	}

	if err := validateMeta(
		meta,
	); err != nil {
		return meta, "", err
	}

	var out bytes.Buffer

	if err := markdown.Convert(
		body,
		&out,
	); err != nil {
		return meta, "",
			fmt.Errorf(
				"render markdown: %w",
				err,
			)
	}

	return meta,
		template.HTML(
			out.String(),
		),
		nil
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

func writeRootIndex(
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
			publicRootDir,
			"index.json",
		),
		data,
		0o644,
	)
}

func writeCNAME() error {
	if strings.TrimSpace(blogCNAME) == "" {
		return nil
	}

	return os.WriteFile(
		filepath.Join(
			publicRootDir,
			"CNAME",
		),
		[]byte(strings.TrimSpace(blogCNAME)+"\n"),
		0o644,
	)
}
