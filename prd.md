# PRD – ascii-art-web

## 1. Problem Statement
The ascii-art-web project involves creating a Go-based HTTP server that provides a graphical user interface (GUI) for the previous ascii-art project. It allows users to input text and select a banner style via a web interface, and displays the generated ASCII art on the page.

---

## 2. User / Use Case

### Primary Users
- Users wanting to generate ASCII-art visually without using the command line.
- Web developers learning Go HTTP server implementation and HTML templates.

### Use Cases
- Access a web page to input a string of text.
- Select from multiple ASCII-art banners (shadow, standard, thinkertoy).
- Submit the data and view the resulting ASCII-art directly in the browser.

---

## 3. Web Interface Requirements

### Main Page Features
- A text input field for the string to convert.
- Radio buttons or a select object to switch between banners (`shadow`, `standard`, `thinkertoy`).
- A submit button to send the request and generate the result.
- The generated ASCII-art should be displayed either on the home page or redirected to the `/ascii-art` route after completion.

---

## 4. Technical Requirements

### 4.1 Endpoints
- **`GET /`**:
  - Sends the main HTML page.
  - Uses Go `html/template` package to render and display data.
- **`POST /ascii-art`**:
  - Receives the text and selected banner from the user.
  - Generates the ASCII-art using the logic from the previous project.
  - Displays the result.

### 4.2 Error Handling and HTTP Status Codes
Endpoints must return appropriate HTTP status codes:
- **`200 OK`**: Everything went without errors.
- **`400 Bad Request`**: Incorrect requests (e.g., malformed POST data).
- **`404 Not Found`**: Resources not found (e.g., missing templates or banner files, invalid URL routes).
- **`500 Internal Server Error`**: Unhandled server errors.

### 4.3 Architecture
- **Language**: Go
- **Packages**: Only standard Go packages are allowed (no external libraries/frameworks).
- **Templates**: HTML templates must be stored in the `templates` directory at the project root.
- **Best Practices**: The code must respect good Go programming practices.

---

## 5. Non-Goals
- No external Go packages or routers (like Gin or Echo).

---

## 6. Acceptance Criteria
- Running the server and visiting the correct port displays the form.
- Submitting the form with a text and banner choice successfully displays the ASCII art.
- Missing banner files or invalid endpoints correctly trigger 404/500/400 errors.
