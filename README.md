# Infer Origins Blog

Independent blog content repository for `${SITE_NAME}`.

This repository contains the authored Markdown, metadata, source images, Go static-site builder, image-processing pipeline, and GitHub Pages deployment output used by the Infer Origins blog.

- Blog content service: `${BLOG_BASE_URL}`
- Main website: `${MAIN_SITE_URL}`

The main website consumes the generated blog index and article HTML from this repository.

---

## Configuration

Keep deployment-specific and locally changeable values in environment variables rather than hard-coding domains, CNAMEs, repository locations, or build paths.

Use a repository-safe `.env.example` such as:

```dotenv
SITE_NAME=Infer Origins

MAIN_SITE_SCHEME=https
MAIN_SITE_HOST=www.inferorigins.com
MAIN_SITE_URL=https://www.inferorigins.com

BLOG_SCHEME=https
BLOG_CNAME=blog.inferorigins.com
BLOG_BASE_URL=https://blog.inferorigins.com

REPO_SSH_URL=git@github.com:MrDaggubati/infer-blog.git
REPO_DIR=infer-blog

CONTENT_DIR=content/blog
PUBLIC_DIR=public
CACHE_DIR=.cache

LOCAL_HOST=localhost
LOCAL_PORT=8080
```

Create a local `.env` or export the equivalent variables in your shell/CI environment. Add `.env` to `.gitignore`.

The `Makefile` and Go builder should use these environment variables, with sensible defaults where appropriate. Generated URLs should be composed from `MAIN_SITE_URL` and `BLOG_BASE_URL`.

### GitHub Pages CNAME

Generate the GitHub Pages `CNAME` from `BLOG_CNAME`:

```bash
printf '%s\n' "${BLOG_CNAME}" > "${PUBLIC_DIR}/CNAME"
```

Do not hard-code the production hostname in the generated file or deployment workflow.

### GitHub Actions

Store non-sensitive deployment configuration as GitHub repository/environment variables and expose it to the build:

```yaml
env:
  SITE_NAME: ${{ vars.SITE_NAME }}
  MAIN_SITE_HOST: ${{ vars.MAIN_SITE_HOST }}
  MAIN_SITE_URL: ${{ vars.MAIN_SITE_URL }}
  BLOG_CNAME: ${{ vars.BLOG_CNAME }}
  BLOG_BASE_URL: ${{ vars.BLOG_BASE_URL }}
  CONTENT_DIR: ${{ vars.CONTENT_DIR }}
  PUBLIC_DIR: ${{ vars.PUBLIC_DIR }}
  CACHE_DIR: ${{ vars.CACHE_DIR }}
```

Use GitHub **Secrets** instead of **Variables** for sensitive values.

---

## Repository Structure

```text
infer-blog/
├── cmd/
│   └── blogbuild/
│       ├── main.go
│       ├── content.go
│       └── images.go
├── content/
│   └── blog/
│       └── <slug>/
│           ├── meta.json
│           ├── article.md
│           └── images/
├── public/
├── .cache/
├── .github/
│   └── workflows/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

Configured locations:

- Authored content: `${CONTENT_DIR}/`
- Generated deployment output: `${PUBLIC_DIR}/`
- Generated blog content: `${PUBLIC_DIR}/blog/`
- Image cache: `${CACHE_DIR}/blogbuild/`

`PUBLIC_DIR` and `CACHE_DIR` are generated/disposable. Do not edit them manually.

---

## Dependencies

Dependency   
- Go           Blog builder and Markdown processing 
- Git          Source control 
Make         Build commands Python 3     
- Optional local HTTP server for `make serve`

Check/install supported dependencies:

```bash
make deps
```

Verify without installing:

```bash
make check-deps
```

The dependency helper supports package managers such as `apt-get`, `dnf`, `yum`, `pacman`, and `brew`. Go should already be installed; if it is missing, install it from the official Go distribution and rerun `make deps`.

---

## Initial Setup

```bash
git clone "${REPO_SSH_URL}"
cd "${REPO_DIR}"

make deps
make tidy
make build
```

For day-to-day work, use `make rebuild` followed by `make serve`.

---

## Creating an Article

Each article lives in its own directory:

```text
${CONTENT_DIR}/<slug>/
├── meta.json
├── article.md
└── images/
```

Example:

```bash
mkdir -p "${CONTENT_DIR}/platform-engineering/images"
```

### Article Metadata

Each article requires `meta.json`:

```json
{
  "title": "Platform Engineering",
  "slug": "platform-engineering",
  "date": "2026-09-03",
  "author": "Infer Origins",
  "summary": "A short description used by blog listing cards.",
  "description": "A longer description suitable for article and SEO metadata.",
  "tags": [
    "engineering",
    "platform"
  ],
  "featured": false,
  "image": "images/Infer-PaaS-foundation.svg",
  "cover": "images/Infer-PaaS-foundation.svg"
}
```

Fields:

----------------------------------------------------------------------- Field                               Purpose ----------------------------------- ----------------------------------- `title`                             Article title

`slug`                              URL-safe article identifier

`date`                              Publication date in `YYYY-MM-DD` format

`author`                            Article author

`summary`                           Short text for cards/listings

`description`                       Longer description for metadata, SEO, and article context

`tags`                              Article classifications

`featured`                          Whether the article is featured

`image`                             Source image used to generate the card image

`cover`                             Source image used to generate the hero/cover image -----------------------------------------------------------------------

The resulting public article URL is:

```text
${MAIN_SITE_URL}/blog/<slug>
```

`image` and `cover` may reference the same source file or different files. Source filenames are not part of the public card/cover API.

### Writing the Article

Write article content in `article.md`:

``` markdown
# Platform Engineering

Platform engineering is the practice of building internal systems that provide
developers with reliable paved roads for delivering software.

## Platform Foundations

A production platform typically combines:

- infrastructure automation
- deployment workflows
- observability
- security controls
- developer interfaces
```

The builder converts Markdown to an HTML fragment using Goldmark with GitHub Flavored Markdown, footnotes, typographer extensions, and automatic heading IDs.

Generated article HTML:

```text
${PUBLIC_DIR}/blog/<slug>/article.html
```

---

## Images

Source images belong under:

```text
${CONTENT_DIR}/<slug>/images/
```

Supported source formats:

```text
.png
.jpg
.jpeg
.svg
.webp
```

The builder copies source images to the generated article directory and creates WebP derivatives where appropriate.

For metadata:

```json
{
  "image": "images/Infer-PaaS-foundation.svg",
  "cover": "images/Infer-PaaS-foundation.svg"
}
```

the builder generates:

-------------------------------------------------------------------------------------------------- Output                  Size                    Public path ----------------------- ----------------------- -------------------------------------------------- Card                    1000 × 500              `${BLOG_BASE_URL}/blog/<slug>/images/card.webp`

Cover                   1600 × 900              `${BLOG_BASE_URL}/blog/<slug>/images/cover.webp` --------------------------------------------------------------------------------------------------

### Image Cache

Image transformations use a content-addressed cache under:

```text
${CACHE_DIR}/blogbuild/images/
```

Cache keys are based on the source image contents plus the transformation. Unchanged images reuse cached derivatives; changed source images or transformation versions generate new entries automatically.

To force image regeneration:

```bash
make clean-cache
make build
```

To remove both generated output and cache:

```bash
make clean-all
make build
```

---

## Generated Blog API

Every build creates:

```text
${PUBLIC_DIR}/blog/index.json
```

Example:

```json
[
  {
    "title": "Platform Engineering",
    "slug": "platform-engineering",
    "date": "2026-09-03",
    "author": "Infer Origins",
    "summary": "A short description used by blog listing cards.",
    "description": "A longer description suitable for article and SEO metadata.",
    "tags": [
      "engineering",
      "platform"
    ],
    "featured": false,
    "image": "${BLOG_BASE_URL}/blog/platform-engineering/images/card.webp",
    "cover": "${BLOG_BASE_URL}/blog/platform-engineering/images/cover.webp",
    "article": "${BLOG_BASE_URL}/blog/platform-engineering/article.html",
    "url": "${MAIN_SITE_URL}/blog/platform-engineering"
  }
]
```

`index.json` is generated automatically and should not be maintained manually.

Public URL patterns:

Resource               URL ---------------------- -------------------------------------------------- Blog content service   `${BLOG_BASE_URL}` Blog index             `${BLOG_BASE_URL}/blog/index.json` Article HTML           `${BLOG_BASE_URL}/blog/<slug>/article.html` Card image             `${BLOG_BASE_URL}/blog/<slug>/images/card.webp` Cover image            `${BLOG_BASE_URL}/blog/<slug>/images/cover.webp` Public article page    `${MAIN_SITE_URL}/blog/<slug>`

---

## Author Workflow

### Add a New Article

```bash
git pull
mkdir -p "${CONTENT_DIR}/my-new-article/images"

# Add meta.json, article.md, and source images.

make rebuild
make serve
```

Verify locally:

```text
http://${LOCAL_HOST}:${LOCAL_PORT}/blog/index.json
http://${LOCAL_HOST}:${LOCAL_PORT}/blog/my-new-article/article.html
```

Then publish:

```bash
git add .
git commit -m "Add my new article"
git push
```

GitHub Actions builds and publishes the content.

### Update an Article

Edit:

```text
${CONTENT_DIR}/<slug>/article.md
${CONTENT_DIR}/<slug>/meta.json
```

Then:

```bash
make rebuild
```

Unchanged images reuse their cached WebP derivatives.

### Delete an Article

```bash
rm -rf "${CONTENT_DIR}/<slug>"
make rebuild
```

Because the generated blog output is recreated during the build, the article is removed from both its generated directory and `index.json`.

---

## Make Targets

----------------------------------------------------------------------- Command                             Purpose ----------------------------------- ----------------------------------- `make help`                         Display available commands

`make deps`                         Check dependencies and install supported missing system packages

`make check-deps`                   Check dependencies without installing

`make tidy`                         Update Go module dependencies

`make fmt`                          Format Go builder source

`make build`                        Build generated content

`make rebuild`                      Clean/reformat/tidy/build while preserving image cache

`make clean`                        Remove generated public content

`make clean-cache`                  Remove image cache

`make clean-all`                    Remove generated output and image cache

`make serve`                        Build and serve generated content locally -----------------------------------------------------------------------

Local server:

```text
http://${LOCAL_HOST}:${LOCAL_PORT}
```

`make rebuild` preserves the image cache, so unchanged images do not require rebuild

---

## Deployment

The blog repository is deployed independently from the main website:

```text
infer-blog
    │
    │ push
    ▼
GitHub Actions
    │
    │ Go blog build
    ▼
${PUBLIC_DIR}/
    │
    ▼
GitHub Pages
    │
    ▼
${BLOG_CNAME}
```

The main website remains independently deployed at `${MAIN_SITE_HOST}`. This separation allows blog articles to be published without rebuilding the entire main website.

The source of truth is always `${CONTENT_DIR}/`. `${PUBLIC_DIR}/` contains deployment artifacts and `${CACHE_DIR}/` contains local build-acceleration artifacts; both can be reconstructed.

---

## `.gitignore`

Recommended entries:

```gitignore
/public/
/.cache/
.env
```

Generated deployment artifacts, local cache, and local environment configuration should not be committed.

---

## Quick Start

First-time setup:

```bash
make deps
make tidy
make build
```

Normal article development:

```bash
make rebuild
make serve
```

Publish:

```bash
git add .
git commit -m "Add blog article"
git push
```