BLOG_COVER_SIZE := 1000x500
BLOG_COVER_QUALITY := 85

.PHONY: help tidy build clean rebuild blog-cover blog-cover-slug serve

help:
	@echo "Targets:"
	@echo "  make tidy"
	@echo "  make build"
	@echo "  make rebuild"
	@echo "  make clean"
	@echo "  make serve"
	@echo "  make blog-cover SRC=input.png OUT=content/blog/slug/images/cover.webp"
	@echo "  make blog-cover-slug SLUG=test-dynamic-blog SRC=input.png"

tidy:
	go mod tidy

build:
	go run ./cmd/blogbuild

clean:
	rm -rf public

rebuild: clean tidy build

blog-cover:
	@test -n "$(SRC)" || (echo "Usage: make blog-cover SRC=path/to/image.png OUT=path/to/cover.webp"; exit 1)
	@test -n "$(OUT)" || (echo "Usage: make blog-cover SRC=path/to/image.png OUT=path/to/cover.webp"; exit 1)
	@mkdir -p "$$(dirname "$(OUT)")"
	ffmpeg -y \
		-i "$(SRC)" \
		-vf "scale=$(BLOG_COVER_SIZE):force_original_aspect_ratio=increase,crop=$(BLOG_COVER_SIZE)" \
		-c:v libwebp \
		-quality $(BLOG_COVER_QUALITY) \
		"$(OUT)"

blog-cover-slug:
	@test -n "$(SLUG)" || (echo "Usage: make blog-cover-slug SLUG=test-dynamic-blog SRC=cover.png"; exit 1)
	@test -n "$(SRC)" || (echo "Usage: make blog-cover-slug SLUG=test-dynamic-blog SRC=cover.png"; exit 1)
	@mkdir -p "content/blog/$(SLUG)/images"
	ffmpeg -y \
		-i "$(SRC)" \
		-vf "scale=$(BLOG_COVER_SIZE):force_original_aspect_ratio=increase,crop=$(BLOG_COVER_SIZE)" \
		-c:v libwebp \
		-quality $(BLOG_COVER_QUALITY) \
		"content/blog/$(SLUG)/images/cover.webp"

serve: build
	@echo "Serving blog at http://localhost:8080"
	python3 -m http.server 8080 --directory public

