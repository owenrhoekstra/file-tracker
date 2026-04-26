package ocr

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func OcrEndpoint(w http.ResponseWriter, r *http.Request) {
	println("OCR Endpoint Hit")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "bad multipart", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, err := mw.CreateFormFile("image", header.Filename)
	if err != nil {
		http.Error(w, "multipart build failed", http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		http.Error(w, "copy failed", http.StatusInternalServerError)
		return
	}
	mw.Close()

	pythonURL := os.Getenv("PYTHON_OCR_URL")
	resp, err := http.Post(pythonURL, mw.FormDataContentType(), &buf)
	if err != nil {
		fmt.Printf("[OCR] python unreachable: %v\n", err)
		http.Error(w, "ocr failed", http.StatusBadGateway)
		return
	}
	println("Py OCR Called")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[OCR] failed to read python response: %v\n", err)
		http.Error(w, "ocr failed", http.StatusInternalServerError)
		return
	}

	fmt.Printf("[OCR] status: %d, result: %s\n", resp.StatusCode, string(body))

	// Mirror the Python service's content type (JSON or plain error text)
	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	// Proxy status code + body straight through
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
