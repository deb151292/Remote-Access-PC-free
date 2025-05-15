package main

import (
	"archive/zip"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Define defaultBaseDir in a platform-agnostic way
var defaultBaseDir = filepath.FromSlash(getDefaultBaseDir())

// getDefaultBaseDir returns a platform-appropriate default directory
func getDefaultBaseDir() string {
	// Check the operating system at runtime
	switch sysOS := strings.ToLower(os.Getenv("OS")); {
	case strings.Contains(sysOS, "windows"):
		return "D:/" // Windows default
	default:
		// For Linux/macOS, use the user's home directory
		homeDir, err := os.UserHomeDir() // Now correctly refers to the os package
		if err != nil {
			log.Printf("Failed to get user home directory: %v, falling back to /tmp", err)
			return "/tmp" // Fallback for Linux/macOS
		}
		return filepath.Join(homeDir, "Documents") // e.g., /home/user/Documents or /Users/user/Documents
	}
}

func main() {
	// Define HTTP handlers
	http.HandleFunc("/", fileBrowser)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/delete", deleteFile)
	http.HandleFunc("/create-folder", createFolder)

	log.Println("Serving folder GUI on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fileBrowser(w http.ResponseWriter, r *http.Request) {
	// Get the current path from query parameter, default to defaultBaseDir
	currentPath := r.URL.Query().Get("path")
	if currentPath == "" {
		log.Printf("No path in query, using defaultBaseDir: %s", defaultBaseDir)
		currentPath = defaultBaseDir
	} else {
		log.Printf("Path from query: %s", currentPath)
	}

	// Sanitize and normalize the path
	currentPath = filepath.Clean(currentPath)
	log.Printf("Sanitized currentPath: %s", currentPath)

	// Ensure the path doesn't go above defaultBaseDir
	absDefaultBaseDir, err := filepath.Abs(defaultBaseDir)
	if err != nil {
		log.Printf("Failed to get absolute path for defaultBaseDir %s: %v", defaultBaseDir, err)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}
	absCurrentPath, err := filepath.Abs(currentPath)
	if err != nil {
		log.Printf("Failed to get absolute path for currentPath %s: %v", currentPath, err)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(absCurrentPath, absDefaultBaseDir) {
		log.Printf("Path %s is above defaultBaseDir %s, redirecting to defaultBaseDir", absCurrentPath, absDefaultBaseDir)
		redirectURL := "/?path=" + url.QueryEscape(defaultBaseDir)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Check if the path is safe
	if !isPathSafe(currentPath) {
		log.Printf("Path is not safe: %s, redirecting with error", currentPath)
		http.Error(w, "Invalid path or access denied", http.StatusBadRequest)
		return
	}

	// Read directory contents
	files, err := os.ReadDir(currentPath)
	if err != nil {
		log.Printf("Failed to read directory %s: %v", currentPath, err)
		http.Error(w, "Unable to read directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare file list for template
	type FileInfo struct {
		Name     string
		Type     string
		Size     string
		Modified string
	}

	var list []FileInfo
	for _, f := range files {
		ftype := "file"
		if f.IsDir() {
			ftype = "folder"
		}

		// Get file details
		info, err := f.Info()
		if err != nil {
			log.Printf("Failed to get info for %s: %v", f.Name(), err)
			continue // Skip files we can't get info for
		}

		// Format file size
		var sizeStr string
		if !f.IsDir() {
			size := info.Size()
			switch {
			case size >= 1<<30: // GB
				sizeStr = fmt.Sprintf("%.2f GB", float64(size)/float64(1<<30))
			case size >= 1<<20: // MB
				sizeStr = fmt.Sprintf("%.2f MB", float64(size)/float64(1<<20))
			case size >= 1<<10: // KB
				sizeStr = fmt.Sprintf("%.2f KB", float64(size)/float64(1<<10))
			default:
				sizeStr = fmt.Sprintf("%d bytes", size)
			}
		} else {
			sizeStr = "-"
		}

		// Format modified time
		modified := info.ModTime().Format("2006-01-02 15:04:05")

		list = append(list, FileInfo{
			Name:     f.Name(),
			Type:     ftype,
			Size:     sizeStr,
			Modified: modified,
		})
	}

	// Template data
	data := struct {
		Files          []FileInfo
		CurrentPath    string
		DefaultBaseDir string
		Separator      string
		Error          string
		Success        string
	}{
		Files:          list,
		CurrentPath:    currentPath,
		DefaultBaseDir: defaultBaseDir,
		Separator:      string(filepath.Separator),
		Error:          r.URL.Query().Get("error"),
		Success:        r.URL.Query().Get("success"),
	}

	log.Printf("Rendering template with CurrentPath: %s", currentPath)

	// Define custom functions for the template
	funcMap := template.FuncMap{
		"split": func(s, sep string) []string {
			// Normalize separators to match the OS
			s = strings.ReplaceAll(s, "/", string(filepath.Separator))
			s = strings.ReplaceAll(s, "\\", string(filepath.Separator))
			return strings.Split(s, sep)
		},
	}

	// Parse the template with the custom functions
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	currentPath := r.URL.Query().Get("path")
	name := r.URL.Query().Get("file")
	if name == "" || currentPath == "" {
		http.Error(w, "Missing file or path parameter", http.StatusBadRequest)
		return
	}

	// Sanitize paths
	currentPath = filepath.Clean(currentPath)
	if !isPathSafe(currentPath) {
		http.Error(w, "Invalid path or access denied", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(currentPath, filepath.Clean(name))
	if !isPathSafeParent(filePath, currentPath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if the path is a file or directory
	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File or folder not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		// Handle folder download as a ZIP file
		zipFileName := name + ".zip"
		w.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
		w.Header().Set("Content-Type", "application/zip")

		// Create a ZIP writer
		zw := zip.NewWriter(w)
		defer zw.Close()

		// Walk the directory and add files to the ZIP
		err = filepath.Walk(filePath, func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Create a header for the file in the ZIP
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// Adjust the name to be relative to the folder being zipped
			header.Name, err = filepath.Rel(filePath, file)
			if err != nil {
				return err
			}

			if info.IsDir() {
				header.Name += string(filepath.Separator)
			} else {
				header.Method = zip.Deflate
			}

			// Create the file in the ZIP
			writer, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				f, err := os.Open(file)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.Copy(writer, f)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			http.Error(w, "Error creating ZIP file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = zw.Close()
		if err != nil {
			http.Error(w, "Error finalizing ZIP file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Handle file download
		f, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, "Error downloading file", http.StatusInternalServerError)
			return
		}
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Log all form values for debugging
	log.Printf("Form values: %v", r.Form)

	// Get the current path from form data
	currentPath := r.FormValue("path")
	if currentPath == "" {
		log.Printf("Form path is empty, falling back to defaultBaseDir: %s", defaultBaseDir)
		currentPath = defaultBaseDir
	} else {
		log.Printf("Received path from form: %s", currentPath)
	}

	// Sanitize path
	currentPath = filepath.Clean(currentPath)
	log.Printf("Sanitized currentPath: %s", currentPath)

	// Check if the path is safe
	if !isPathSafe(currentPath) {
		log.Printf("Path is not safe: %s", currentPath)
		// Try to create the directory if it doesn't exist
		err := os.MkdirAll(currentPath, os.ModePerm)
		if err != nil {
			log.Printf("Failed to create directory %s: %v", currentPath, err)
			redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape(fmt.Sprintf("Cannot access or create directory %s: %v (try running the application as administrator)", currentPath, err))
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		log.Printf("Created directory: %s", currentPath)
		// Re-check if the path is safe after creating it
		if !isPathSafe(currentPath) {
			log.Printf("Path is still not safe after creation: %s", currentPath)
			redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Invalid path or access denied (try running the application as administrator)")
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving uploaded file: %v", err)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Upload error: "+err.Error())
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	defer file.Close()

	// Sanitize file path
	dstPath := filepath.Join(currentPath, filepath.Base(handler.Filename))
	log.Printf("Constructed destination path: %s", dstPath)
	if !isPathSafeParent(dstPath, currentPath) {
		log.Printf("Destination path is not safe: %s", dstPath)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Invalid file path")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Check if the file already exists and is locked
	if _, err := os.Stat(dstPath); err == nil {
		// File exists, try to remove it first
		err = os.Remove(dstPath)
		if err != nil {
			log.Printf("Failed to remove existing file %s: %v", dstPath, err)
			redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape(fmt.Sprintf("File %s already exists and cannot be overwritten: %v (ensure the file is not open in another program)", dstPath, err))
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		log.Printf("Removed existing file: %s", dstPath)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Failed to create file at %s: %v", dstPath, err)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape(fmt.Sprintf("Failed to save file to %s: %v (try running the application as administrator or check directory permissions)", dstPath, err))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file to %s: %v", dstPath, err)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape(fmt.Sprintf("Error saving file to %s: %v", dstPath, err))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Redirect to the current path to refresh the page with success message
	successMsg := fmt.Sprintf("File '%s' uploaded successfully to %s", handler.Filename, currentPath)
	log.Printf("Upload successful: %s", successMsg)
	redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&success=" + url.QueryEscape(successMsg)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	currentPath := r.URL.Query().Get("path")
	name := r.URL.Query().Get("file")
	if name == "" || currentPath == "" {
		http.Error(w, "Missing file or path parameter", http.StatusBadRequest)
		return
	}

	// Sanitize paths
	currentPath = filepath.Clean(currentPath)
	if !isPathSafe(currentPath) {
		http.Error(w, "Invalid path or access denied", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(currentPath, filepath.Clean(name))
	if !isPathSafeParent(filePath, currentPath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if the path exists
	_, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File or folder not found", http.StatusNotFound)
		return
	}

	// Delete the file or folder
	err = os.RemoveAll(filePath)
	if err != nil {
		http.Error(w, "Failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted successfully"))
}

func createFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form in createFolder: %v", err)
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get path and folderName from form data
	currentPath := r.FormValue("path")
	folderName := r.FormValue("folderName")
	if currentPath == "" {
		log.Printf("Form path is empty in createFolder, falling back to defaultBaseDir: %s", defaultBaseDir)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Missing path")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Sanitize paths
	currentPath = filepath.Clean(currentPath)
	if !isPathSafe(currentPath) {
		log.Printf("Path is not safe in createFolder: %s", currentPath)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Invalid path or access denied")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// If folderName is empty, use default name "New Folder"
	if folderName == "" {
		folderName = "New Folder"
		// Check if "New Folder" exists, and if so, append a number
		baseName := folderName
		count := 1
		folderPath := filepath.Join(currentPath, filepath.Clean(folderName))
		for {
			_, err := os.Stat(folderPath)
			if os.IsNotExist(err) {
				break
			}
			folderName = fmt.Sprintf("%s (%d)", baseName, count)
			folderPath = filepath.Join(currentPath, filepath.Clean(folderName))
			count++
		}
	}

	// Construct the folder path
	folderPath := filepath.Join(currentPath, filepath.Clean(folderName))

	// Ensure the new path is a child of currentPath (prevent directory traversal)
	if !isPathSafeParent(folderPath, currentPath) {
		log.Printf("Folder path is not safe in createFolder: %s", folderPath)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Invalid folder path")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Create the folder
	err = os.Mkdir(folderPath, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder at %s: %v", folderPath, err)
		redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&error=" + url.QueryEscape("Failed to create folder: "+err.Error())
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Redirect to the current path with success message
	successMsg := fmt.Sprintf("Folder '%s' created successfully", folderName)
	redirectURL := "/?path=" + url.QueryEscape(currentPath) + "&success=" + url.QueryEscape(successMsg)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// isPathSafe ensures the path exists and is safe (no directory traversal)
func isPathSafe(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Failed to get absolute path for %s: %v", path, err)
		return false
	}
	// Ensure the path is a valid directory
	_, err = os.Stat(absPath)
	if err != nil {
		log.Printf("Path does not exist or is inaccessible: %s, error: %v", absPath, err)
		return false
	}
	log.Printf("Path is safe: %s", absPath)
	return true
}

// isPathSafeParent ensures the path is a child of the parent path (for new paths)
func isPathSafeParent(path, parent string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Failed to get absolute path for %s: %v", path, err)
		return false
	}
	absParent, err := filepath.Abs(parent)
	if err != nil {
		log.Printf("Failed to get absolute path for parent %s: %v", parent, err)
		return false
	}
	// Check if absPath starts with absParent
	isSafe := strings.HasPrefix(absPath, absParent)
	if !isSafe {
		log.Printf("Path %s is not a child of %s", absPath, absParent)
	} else {
		log.Printf("Path %s is a safe child of %s", absPath, absParent)
	}
	return isSafe
}
