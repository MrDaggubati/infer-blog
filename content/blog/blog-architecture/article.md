---
title: Building a Portable Web Platform with Go, HTMX and Static Output
slug: building-portable-web-platform-go-htmx-static
date: 2026-09-04
author: Sudhakar Daggubati
summary: >
  How we assembled the Infer Origins portal and technical blog around Go,
  HTMX, Markdown, generated static output, and Cloudflare while preserving
  the option to run the same platform as a server-rendered application.
tags:
  - Go
  - HTMX
  - Cloudflare
  - Static Sites
  - Platform Engineering
featured: true
---


Websites often become application platforms before they really need to.

A blog introduces a content system. Search introduces APIs. Deployment introduces containers, runtime configuration, health checks, and another service that has to stay alive.

For the Infer Origins portal, the goal was different:

**Keep the site deployable as static files, while preserving the ability to run it as a Go application when server-side behaviour becomes useful.**

The architecture therefore starts with one content and rendering model, not two separate systems.

```text
                  +----------------------+
                  |  Infer Origins Web   |
                  +----------+-----------+
                             |
               +-------------+-------------+
               |                           |
               v                           v
        +-------------+             +-------------+
        | Static mode |             | Server mode |
        +------+------+             +------+------+
               |                           |
               v                           v
        HTML/CSS/JS                  Go HTTP + HTMX
        JSON/images                  rendered HTML
               |                           |
               +-------------+-------------+
                             |
                             v
                     Same content model
```

## Go as the assembly layer

Go owns the durable structure of the site:

```text
Go
|
+-- content models
+-- templates
+-- build commands
+-- Markdown processing
+-- metadata validation
+-- asset handling
+-- static generation
+-- optional HTTP serving
```

The browser receives normal HTML.

There is no requirement to bootstrap a client-side application before the page becomes useful.

```text
                    +----------------+
                    |   Go models    |
                    +--------+-------+
                             |
             +---------------+---------------+
             |                               |
             v                               v
      +--------------+                +--------------+
      | Templates    |                | Content      |
      |              |                |              |
      | base         |                | products     |
      | products     |                | services     |
      | case studies |                | blog         |
      | policies     |                | pages        |
      +------+-------+                +------+-------+
             |                               |
             +---------------+---------------+
                             |
                             v
                       Go assembly
                             |
               +-------------+-------------+
               |                           |
               v                           v
          Static files                HTTP responses
```

That gives the portal one rendering system regardless of deployment mode.

## HTMX without turning the site into an SPA

HTMX fits naturally because the server can still own the HTML.

A traditional client-heavy model often looks like:

```text
Browser
   |
   v
JavaScript application
   |
   +-- routing
   +-- state
   +-- rendering
   +-- API client
   |
   v
JSON API
```

With Go and HTMX:

```text
Browser
   |
   | request
   v
Go
   |
   v
HTML fragment
   |
   v
Browser
```

HTMX becomes an enhancement to ordinary HTTP instead of the foundation of the whole site.

That keeps page rendering understandable and leaves the static mode viable.

## One portal, two deployment modes

The same project can generate a deployable directory:

```text
dist/
|
+-- index.html
+-- products/
+-- services/
+-- case-studies/
+-- about/
+-- blog/
+-- static/
+-- sitemap.xml
+-- robots.txt
```

and those files can be hosted by almost any static platform.

```text
              +----------------+
              |   Go builder   |
              +--------+-------+
                       |
                       v
              +----------------+
              | Static output  |
              +--------+-------+
                       |
          +------------+------------+
          |            |            |
          v            v            v
     GitHub Pages   Cloudflare    object/static
                    delivery        hosting
```

But the same content and templates can later sit behind a Go process:

```text
Internet
   |
   v
Cloudflare
   |
   v
Go server
   |
   +-- full HTML pages
   +-- HTMX fragments
   +-- live data
   +-- authenticated routes
```

That means the portal can start operationally simple without locking the architecture into static-only behaviour.

## Static-first, not static-only

The design principle is:

```text
              Can this be generated?
                      |
             +--------+--------+
             |                 |
            yes                no
             |                 |
             v                 v
        static output      runtime path
```

Most content can be generated.

Only behaviour that genuinely needs a runtime should introduce one.

This makes static output an operational optimisation rather than an ideological constraint.

## Cloudflare as the edge layer

Cloudflare sits around the platform rather than becoming the application itself.

```text
                     +----------------+
                     |   Cloudflare   |
                     +--------+-------+
                              |
         +--------------------+--------------------+
         |                    |                    |
         v                    v                    v
       DNS/CDN             Workers             Security
         |                    |                    |
         v                    v                    v
    static portal      small dynamic paths   edge controls
    static blog
```

This keeps responsibilities narrow.

The portal owns content and presentation.

Cloudflare owns delivery and edge execution.

External systems can remain external rather than being pulled into the core application.

## The blog as a publishing pipeline

The blog follows the same philosophy.

Originally the article model separated metadata and content:

```text
article/
|
+-- meta.json
+-- article.md
+-- images/
```

That meant maintaining two source files for one article.

The simplified model uses Markdown front matter:

```text
article.md
|
+-- YAML front matter
|     |
|     +-- title
|     +-- slug
|     +-- date
|     +-- author
|     +-- summary
|     +-- tags
|     +-- image
|     +-- cover
|
+-- Markdown body
```

The builder splits the document:

```text
                         article.md
                             |
                             v
                      +-------------+
                      | readArticle |
                      +------+------+
                             |
                 +-----------+-----------+
                 |                       |
                 v                       v
          YAML front matter         Markdown body
                 |                       |
                 v                       v
             BlogMeta                 Goldmark
                 |                       |
                 v                       v
             index.json            article.html
```

`article.md` becomes the single authored source of truth.

## Generate the blog index

The portal needs a lightweight way to list and search articles without loading every article body.

The builder therefore generates `index.json`.

```text
                 +------------------+
                 |  Blog sources    |
                 +--------+---------+
                          |
                          v
                    Go builder
                          |
             +------------+------------+
             |                         |
             v                         v
       article.html                index.json
             |                         |
             |              +----------+----------+
             |              |          |          |
             |              v          v          v
             |            cards      search      tags
             |
             +-------------------------------> article
```

The index only contains discovery metadata.

```json
{
  "title": "Linux isolation",
  "slug": "linux-isolation",
  "date": "2026-07-27",
  "author": "Sudhakar Daggubati",
  "summary": "How modern workloads can combine...",
  "tags": [
    "Security",
    "Virtualization",
    "Linux"
  ],
  "featured": false,
  "image": "...",
  "cover": "...",
  "article": "..."
}
```

The article body stays separate.

## Article assets stay with the article

An article can reference images from front matter:

```yaml
image: images/DNSControl-flow.png
cover: images/DNSControl-flow.png
```

and from Markdown:

```markdown
![DNS Flow](images/dns-flow.png)
```

Those assets are part of the article content and need to survive the build.

```text
SOURCE

automating-dns/
|
+-- article.md
|
+-- images/
    +-- dns-flow.png
    +-- DNSControl-flow.png


             build
               |
               v


PUBLIC

blog/Automating-DNS/
|
+-- article.html
|
+-- images/
    +-- dns-flow.png
    +-- DNSControl-flow.png
```

The author controls filenames.

The builder controls publication paths.

## Markdown rendering and article styling

Goldmark converts Markdown into semantic HTML:

```text
Markdown                 HTML

# Heading          --->  <h1>
## Heading         --->  <h2>
paragraph          --->  <p>
```bash          --->  <pre><code>
| table |          --->  <table>
```

But technical articles need different typography from editorial landing pages.

The portal uses large headings intentionally.

Those same rules should not control long-form technical content.

```text
                  Shared design system
                         |
             +-----------+-----------+
             |                       |
             v                       v
      Editorial pages          Article content
             |                       |
       large headings             h1-h3
       hero layouts               paragraphs
       visual sections            tables
                                  code blocks
```

The solution is scoped CSS:

```css
.article-content h2 {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-size: clamp(2rem, 3vw, 3rem);
  line-height: 1.1;
}

.article-content p {
  margin: 0 0 1.2rem;
  font-size: 1.125rem;
  line-height: 1.7;
}
```

Code blocks get their own presentation:

```css
.article-content pre {
  margin: 1.75rem 0;
  padding: 1.25rem 1.5rem;
  overflow-x: auto;
  border: 1px solid var(--line);
  border-radius: 14px;
}
```

The HTML stays semantic while the visual language adapts to the reading context.

## Build caches are disposable

As the number of articles grows, parsing everything on every build becomes unnecessary.

A cache can retain metadata for unchanged articles.

```text
                    article.md
                        |
                        v
                  calculate hash
                        |
             +----------+----------+
             |                     |
          unchanged               changed
             |                     |
             v                     v
        use cache            parse + render
             |                     |
             +----------+----------+
                        |
                        v
                  update cache
                        |
                        v
                  generate index
```

The important rule is:

```text
article.md       = source of truth
cache            = optimisation
index.json       = generated
article.html     = generated
```

Deleting the cache must never destroy information.

## Build and rebuild mean different things

Normal build:

```text
infer-blog build
|
+-- load cache
+-- inspect source
+-- rebuild changed articles
+-- reuse unchanged metadata
+-- regenerate index
```

Clean rebuild:

```text
infer-blog rebuild
|
+-- remove generated output
+-- remove cache
+-- parse every article
+-- render every article
+-- copy assets
+-- regenerate index
```

The useful invariant is:

```text
incremental build
       |
       | should produce
       | the same result as
       v
clean rebuild
```

If they differ, the cache has become a source of truth by accident.

## The overall architecture

The portal and blog now fit into one model:

```text
                           SOURCE
                             |
          +------------------+------------------+
          |                  |                  |
          v                  v                  v
       Go code            Templates          Content
          |                  |                  |
          +------------------+------------------+
                             |
                             v
                        Go builder
                             |
             +---------------+---------------+
             |                               |
             v                               v
       Static deployment               Server deployment
             |                               |
             v                               v
        CDN / files                      Go + HTMX
             |
             v
         Cloudflare
```

The blog sits inside the same pipeline:

```text
article.md
   |
   +-- front matter
   |
   +-- Markdown
   |
   v
Go builder
   |
   +-- article.html
   +-- index.json
   +-- copied assets
```

That keeps the platform small while preserving room to grow.

## The main architectural principle

The result is not purely static.

It is not purely server-rendered.

It is not an SPA.

It is a platform that keeps the durable parts close to source and chooses runtime only where runtime adds value.

```text
          source
            |
            v
        generate first
            |
      +-----+-----+
      |           |
      v           v
   static      runtime
   output      only when needed
```

The practical benefit is portability.

The portal can remain a collection of static files today and become a Go + HTMX application tomorrow without throwing away its content model, templates, blog pipeline, or design system.

That is the property we wanted most:

> **Build once around durable content and simple interfaces, then choose the deployment model that fits the problem rather than letting the deployment model dictate the architecture.**