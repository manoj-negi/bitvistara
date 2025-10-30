# bitVistara - Go server with Gorilla Mux

This repository serves your existing static HTML files through a Go server using Gorilla Mux and Go's `html/template` engine.

## Prerequisites
- Go 1.22+

## Install dependencies
Run inside the project root:

```bash
go mod tidy
```

## Run the server
```bash
go run .
```

The server will start on `http://localhost:8080`.

## Routes
- `/` → `index.html`
- `/about-us` → `about-us.html`
- `/services` → `our-services.html`
- `/training` → `training.html`
- `/blog` → `bloglisting.html`
- `/blog/{slug}` → `blogDetails.html` (receives `Slug` in template data)
- `/contact` → `contact_us.html`
- `/server` → `server.html`

## Static assets
Files in `public/` are served at `/public/`.

Example: `public/images/screen.png` → `http://localhost:8080/public/images/screen.png`

## Notes
- Templates are rendered file-by-file without a layout; this matches the current project structure. If you later want a shared layout, we can refactor to use a base template and `{{define}}` blocks.

