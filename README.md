# ASCII Art Web - Export File

A Go web application and API that converts plain text into ASCII-art banners and allows users to export the generated output in multiple formats (TXT, HTML, JSON) with correct HTTP headers.

## Features

- **ASCII Generation Engine:** Supports `standard`, `shadow`, and `thinkertoy` banner fonts.
- **Multiple Export Formats:**
  - **TXT (`.txt`):** Raw ASCII text. Served as `text/plain`.
  - **HTML (`.html`):** Monospace pre-formatted template. Served as `text/html`.
  - **JSON (`.json`):** Structured JSON representation with `text`, `banner`, and `ascii` properties. Served as `application/json`.
- **Correct HTTP Headers:**
  - `Content-Type`: Set dynamically with character encoding (e.g., `text/plain; charset=utf-8`).
  - `Content-Length`: Set dynamically to the exact size of the payload body in bytes.
  - `Content-Disposition`: Set to `attachment; filename="..."` to force browser download.
- **REST API & Web UI:** Includes a responsive futuristic grid UI with instant format selection and a backend export API.

## Quick Start

### 1. Run the Server
Ensure Go is installed, and run:
```bash
go run .
```
Visit the application at:
```text
http://localhost:8080
```

To configure another port:
```bash
PORT=8081 go run .
```

## How To Export
1. Enter printable ASCII text in the input box.
2. Select your desired banner style.
3. Choose the export format from the dropdown: `.TXT`, `.HTML`, or `.JSON`.
4. Click **Export**. The file will automatically download to your device with standard user read/write permissions.

## Routes & API Endpoints

| Method | Path | Format / Purpose |
| --- | --- | --- |
| `GET` | `/` | Renders the main web interface. |
| `POST` | `/api/ascii-art` | Generates ASCII art and returns it in JSON. |
| `POST` | `/export` | Serves export file via form POST. |
| `GET` | `/export` | Serves export file via query parameters. |
| `GET` | `/static/*` | Serves stylesheets and scripts. |

### API Curl Example (Direct File Download)
```bash
curl -i -X POST -d "text=HELLO&banner=standard&format=json" http://localhost:8080/export
```

Response Headers:
```http
HTTP/1.1 200 OK
Content-Disposition: attachment; filename="ascii-art-standard.json"
Content-Length: 337
Content-Type: application/json; charset=utf-8
```

## Security & Best Practices
- **No Third-Party Packages:** The project compiles using only the Go standard library.
- **Security Headers:** Hardened with CSP, Content-Type Options (`nosniff`), X-Frame-Options (`DENY`), and Referrer Policy.
- **Request Limiting:** Maximum body and text payload length checks are enforced.

## Development & Tests
Run the tests:
```bash
go test ./...
```
