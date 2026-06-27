package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ── Config ────────────────────────────────────────────────────────────────────

const (
	sessionTTL    = 12 * time.Hour
	uploadMaxSize = 8 << 20 // 8 MB
)

// ── Data types ────────────────────────────────────────────────────────────────

type ContactSubmission struct {
	ID        string `json:"id"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Message   string `json:"message" binding:"required"`
	CreatedAt string `json:"created_at"`
	Read      bool   `json:"read"`
}

type BlogPost struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Slug    string `json:"slug"`
	Excerpt string `json:"excerpt"`
}

type BlogPostFull struct {
	BlogPost
	HTML string `json:"html"`
	Body string `json:"body,omitempty"` // raw markdown, admin only
}

type Session struct {
	Token     string
	ExpiresAt time.Time
}

// ── Session store (in-memory) ─────────────────────────────────────────────────

var (
	sessionsMu sync.Mutex
	sessions   = map[string]*Session{}
)

func newSession() string {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	sessionsMu.Lock()
	sessions[token] = &Session{Token: token, ExpiresAt: time.Now().Add(sessionTTL)}
	sessionsMu.Unlock()
	return token
}

func validSession(token string) bool {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(s.ExpiresAt) {
		delete(sessions, token)
		return false
	}
	return true
}

func deleteSession(token string) {
	sessionsMu.Lock()
	delete(sessions, token)
	sessionsMu.Unlock()
}

// ── Auth helpers ──────────────────────────────────────────────────────────────

func hashPassword(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

func adminPassword() string {
	pw := os.Getenv("ADMIN_PASSWORD")
	if pw == "" {
		pw = "admin123" // default — change via env var in production
	}
	return hashPassword(pw)
}

func adminUsername() string {
	u := os.Getenv("ADMIN_USERNAME")
	if u == "" {
		return "admin"
	}
	return u
}

// ── Auth middleware ───────────────────────────────────────────────────────────

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("admin_session")
		if err != nil || !validSession(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ── Markdown parser ───────────────────────────────────────────────────────────

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Table, extension.Footnote),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
)

func markdownToHTML(src []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ── Front-matter parser ───────────────────────────────────────────────────────

func parseFrontMatter(content string) (map[string]string, string) {
	meta := map[string]string{}
	if !strings.HasPrefix(content, "---") {
		return meta, content
	}
	end := strings.Index(content[3:], "---")
	if end == -1 {
		return meta, content
	}
	block := content[3 : end+3]
	body := content[end+6:]
	for _, line := range strings.Split(block, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			meta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return meta, body
}

func buildFrontMatter(title, date, excerpt string) string {
	return fmt.Sprintf("---\ntitle: %s\ndate: %s\nexcerpt: %s\n---\n\n", title, date, excerpt)
}

// ── Blog helpers ──────────────────────────────────────────────────────────────

func loadPosts() ([]BlogPost, error) {
	entries, err := os.ReadDir("data/posts")
	if err != nil {
		return nil, err
	}
	var posts []BlogPost
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("data/posts", e.Name()))
		if err != nil {
			continue
		}
		meta, _ := parseFrontMatter(string(raw))
		slug := strings.TrimSuffix(e.Name(), ".md")
		posts = append(posts, BlogPost{
			Title:   meta["title"],
			Date:    meta["date"],
			Slug:    slug,
			Excerpt: meta["excerpt"],
		})
	}
	return posts, nil
}

func loadPost(slug string, includeBody bool) (*BlogPostFull, error) {
	path := filepath.Join("data/posts", slug+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	meta, body := parseFrontMatter(string(raw))
	htmlContent, err := markdownToHTML([]byte(body))
	if err != nil {
		return nil, err
	}
	p := &BlogPostFull{
		BlogPost: BlogPost{
			Title:   meta["title"],
			Date:    meta["date"],
			Slug:    slug,
			Excerpt: meta["excerpt"],
		},
		HTML: htmlContent,
	}
	if includeBody {
		p.Body = strings.TrimSpace(body)
	}
	return p, nil
}

// ── JSON file helpers ─────────────────────────────────────────────────────────

func readJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ── CORS & security middleware ────────────────────────────────────────────────

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ── Public API handlers ───────────────────────────────────────────────────────

func handleAbout(c *gin.Context) {
	var data any
	if err := readJSON("data/about.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleSkills(c *gin.Context) {
	var data any
	if err := readJSON("data/skills.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleProjects(c *gin.Context) {
	var data any
	if err := readJSON("data/projects.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleBlogList(c *gin.Context) {
	posts, err := loadPosts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, posts)
}

func handleBlogPost(c *gin.Context) {
	post, err := loadPost(c.Param("slug"), false)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "post not found"}); return
		}
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, post)
}

func handleContact(c *gin.Context) {
	var sub ContactSubmission
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	sub.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	sub.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	sub.Read = false

	var contacts []ContactSubmission
	_ = readJSON("data/contacts.json", &contacts)
	contacts = append(contacts, sub)
	if err := writeJSON("data/contacts.json", contacts); err != nil {
		c.JSON(500, gin.H{"error": "could not save submission"}); return
	}
	c.JSON(200, gin.H{"message": "Message received — thank you!"})
}

// ── Admin: Auth ───────────────────────────────────────────────────────────────

func handleAdminLogin(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "username and password required"}); return
	}
	if body.Username != adminUsername() || hashPassword(body.Password) != adminPassword() {
		c.JSON(401, gin.H{"error": "invalid credentials"}); return
	}
	token := newSession()
	c.SetCookie("admin_session", token, int(sessionTTL.Seconds()), "/", "", false, true)
	c.JSON(200, gin.H{"message": "logged in", "token": token})
}

func handleAdminLogout(c *gin.Context) {
	token, _ := c.Cookie("admin_session")
	deleteSession(token)
	c.SetCookie("admin_session", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{"message": "logged out"})
}

func handleAdminMe(c *gin.Context) {
	c.JSON(200, gin.H{"username": adminUsername(), "authenticated": true})
}

// ── Admin: About ──────────────────────────────────────────────────────────────

func handleAdminGetAbout(c *gin.Context) {
	var data any
	if err := readJSON("data/about.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleAdminSaveAbout(c *gin.Context) {
	var data any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := writeJSON("data/about.json", data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "About updated"})
}

// ── Admin: Skills ─────────────────────────────────────────────────────────────

func handleAdminGetSkills(c *gin.Context) {
	var data any
	if err := readJSON("data/skills.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleAdminSaveSkills(c *gin.Context) {
	var data any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := writeJSON("data/skills.json", data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Skills updated"})
}

// ── Admin: Projects ───────────────────────────────────────────────────────────

func handleAdminGetProjects(c *gin.Context) {
	var data any
	if err := readJSON("data/projects.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleAdminSaveProjects(c *gin.Context) {
	var data any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := writeJSON("data/projects.json", data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Projects updated"})
}

// ── Admin: Blog ───────────────────────────────────────────────────────────────

func handleAdminGetPosts(c *gin.Context) {
	posts, err := loadPosts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, posts)
}

func handleAdminGetPost(c *gin.Context) {
	post, err := loadPost(c.Param("slug"), true)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "post not found"}); return
		}
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, post)
}

func handleAdminCreatePost(c *gin.Context) {
	var body struct {
		Title   string `json:"title" binding:"required"`
		Slug    string `json:"slug" binding:"required"`
		Date    string `json:"date" binding:"required"`
		Excerpt string `json:"excerpt" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	slug := slugify(body.Slug)
	path := filepath.Join("data/posts", slug+".md")
	if _, err := os.Stat(path); err == nil {
		c.JSON(409, gin.H{"error": "post with this slug already exists"}); return
	}
	content := buildFrontMatter(body.Title, body.Date, body.Excerpt) + body.Body
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Post created", "slug": slug})
}

func handleAdminUpdatePost(c *gin.Context) {
	slug := c.Param("slug")
	var body struct {
		Title   string `json:"title" binding:"required"`
		Date    string `json:"date" binding:"required"`
		Excerpt string `json:"excerpt" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	path := filepath.Join("data/posts", slug+".md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "post not found"}); return
	}
	content := buildFrontMatter(body.Title, body.Date, body.Excerpt) + body.Body
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Post updated"})
}

func handleAdminDeletePost(c *gin.Context) {
	slug := c.Param("slug")
	path := filepath.Join("data/posts", slug+".md")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "post not found"}); return
		}
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Post deleted"})
}

// ── Admin: Contacts ───────────────────────────────────────────────────────────

func handleAdminGetContacts(c *gin.Context) {
	var contacts []ContactSubmission
	if err := readJSON("data/contacts.json", &contacts); err != nil {
		contacts = []ContactSubmission{}
	}
	c.JSON(200, contacts)
}

func handleAdminMarkRead(c *gin.Context) {
	id := c.Param("id")
	var contacts []ContactSubmission
	if err := readJSON("data/contacts.json", &contacts); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	found := false
	for i := range contacts {
		if contacts[i].ID == id {
			contacts[i].Read = true
			found = true
			break
		}
	}
	if !found {
		c.JSON(404, gin.H{"error": "contact not found"}); return
	}
	writeJSON("data/contacts.json", contacts)
	c.JSON(200, gin.H{"message": "marked as read"})
}

func handleAdminDeleteContact(c *gin.Context) {
	id := c.Param("id")
	var contacts []ContactSubmission
	if err := readJSON("data/contacts.json", &contacts); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	newContacts := contacts[:0]
	for _, ct := range contacts {
		if ct.ID != id {
			newContacts = append(newContacts, ct)
		}
	}
	writeJSON("data/contacts.json", newContacts)
	c.JSON(200, gin.H{"message": "contact deleted"})
}

// ── Admin: Site Settings (CSS variables) ─────────────────────────────────────

type SiteSettings struct {
	// Dark theme colours
	DarkBg        string `json:"dark_bg"`
	DarkBgAlt     string `json:"dark_bg_alt"`
	DarkSurface   string `json:"dark_surface"`
	DarkAccent1   string `json:"dark_accent1"`
	DarkAccent2   string `json:"dark_accent2"`
	DarkAccent3   string `json:"dark_accent3"`
	DarkTextPrimary   string `json:"dark_text_primary"`
	DarkTextSecondary string `json:"dark_text_secondary"`

	// Light theme colours
	LightBg        string `json:"light_bg"`
	LightBgAlt     string `json:"light_bg_alt"`
	LightSurface   string `json:"light_surface"`
	LightAccent1   string `json:"light_accent1"`
	LightAccent2   string `json:"light_accent2"`
	LightTextPrimary   string `json:"light_text_primary"`
	LightTextSecondary string `json:"light_text_secondary"`

	// Typography
	FontBody    string `json:"font_body"`
	FontCode    string `json:"font_code"`

	// Layout & shape
	RadiusLg string `json:"radius_lg"`
	RadiusSm string `json:"radius_sm"`

	// Gradient direction
	GradientAngle string `json:"gradient_angle"`

	// WhatsApp number (syncs fab button)
	WhatsAppNumber string `json:"whatsapp_number"`

	// Navbar border colour
	NavBorderColor string `json:"nav_border_color"`

	// Animation speed
	TransitionSpeed string `json:"transition_speed"`
}

func defaultSettings() SiteSettings {
	return SiteSettings{
		DarkBg: "#0d0f1a", DarkBgAlt: "#111428", DarkSurface: "#161929",
		DarkAccent1: "#7c3aed", DarkAccent2: "#06b6d4", DarkAccent3: "#f0abfc",
		DarkTextPrimary: "#f0f2ff", DarkTextSecondary: "#8b92b8",
		LightBg: "#f8f9ff", LightBgAlt: "#eef0fb", LightSurface: "#ffffff",
		LightAccent1: "#7c3aed", LightAccent2: "#06b6d4",
		LightTextPrimary: "#0d0f1a", LightTextSecondary: "#4a5280",
		FontBody: "Inter", FontCode: "Fira Code",
		RadiusLg: "16px", RadiusSm: "8px",
		GradientAngle: "135deg",
		WhatsAppNumber: "+2348168642824",
		NavBorderColor: "rgba(230,14,91,0.51)",
		TransitionSpeed: "0.3s",
	}
}

func loadSettings() SiteSettings {
	var s SiteSettings
	if err := readJSON("data/settings.json", &s); err != nil {
		return defaultSettings()
	}
	return s
}

func handleAdminGetSettings(c *gin.Context) {
	c.JSON(200, loadSettings())
}

func handleAdminSaveSettings(c *gin.Context) {
	var s SiteSettings
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := writeJSON("data/settings.json", s); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	// Also write a generated CSS override file that the frontend loads
	css := generateSettingsCSS(s)
	if err := os.WriteFile("static/theme-override.css", []byte(css), 0644); err != nil {
		c.JSON(500, gin.H{"error": "saved settings but failed to write CSS: " + err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Settings saved and CSS rebuilt"})
}

func generateSettingsCSS(s SiteSettings) string {
	fontImport := ""
	if s.FontBody != "" && s.FontBody != "Inter" {
		fontImport = fmt.Sprintf("@import url('https://fonts.googleapis.com/css2?family=%s:wght@300;400;500;600;700;800&display=swap');\n", strings.ReplaceAll(s.FontBody, " ", "+"))
	}

	return fmt.Sprintf(`%s/* Auto-generated by Portfolio CMS — do not edit manually */
:root {
  --bg:           %s;
  --bg-alt:       %s;
  --surface:      %s;
  --accent-1:     %s;
  --accent-2:     %s;
  --accent-3:     %s;
  --text-primary:   %s;
  --text-secondary: %s;
  --grad:         linear-gradient(%s, var(--accent-1), var(--accent-2));
  --grad-text:    linear-gradient(90deg, var(--accent-1), var(--accent-2));
  --radius:       %s;
  --radius-sm:    %s;
  --transition:   %s cubic-bezier(0.4,0,0.2,1);
  --font-body:    '%s', sans-serif;
  --font-code:    '%s', monospace;
}
[data-theme="light"] {
  --bg:           %s;
  --bg-alt:       %s;
  --surface:      %s;
  --accent-1:     %s;
  --accent-2:     %s;
  --text-primary:   %s;
  --text-secondary: %s;
}
.navbar { border-bottom-color: %s; }
`,
		fontImport,
		s.DarkBg, s.DarkBgAlt, s.DarkSurface,
		s.DarkAccent1, s.DarkAccent2, s.DarkAccent3,
		s.DarkTextPrimary, s.DarkTextSecondary,
		s.GradientAngle,
		s.RadiusLg, s.RadiusSm, s.TransitionSpeed,
		s.FontBody, s.FontCode,
		s.LightBg, s.LightBgAlt, s.LightSurface,
		s.LightAccent1, s.LightAccent2,
		s.LightTextPrimary, s.LightTextSecondary,
		s.NavBorderColor,
	)
}

// ── Admin: Carousel ───────────────────────────────────────────────────────────

func handleGetCarousel(c *gin.Context) {
	var data any
	if err := readJSON("data/carousel.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleAdminGetCarousel(c *gin.Context) {
	var data any
	if err := readJSON("data/carousel.json", &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, data)
}

func handleAdminSaveCarousel(c *gin.Context) {
	var data any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}); return
	}
	if err := writeJSON("data/carousel.json", data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	c.JSON(200, gin.H{"message": "Carousel saved"})
}

func handleAdminUploadSlideImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, uploadMaxSize)
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image file required"}); return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(400, gin.H{"error": "only jpg, png, webp allowed"}); return
	}

	filename := fmt.Sprintf("slide-%d%s", time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join("static", filename))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	defer dst.Close()
	io.Copy(dst, file)

	c.JSON(200, gin.H{"url": "/static/" + filename})
}

// ── Admin: Photo upload ───────────────────────────────────────────────────────

func handleAdminUploadPhoto(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, uploadMaxSize)
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(400, gin.H{"error": "photo file required"}); return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(400, gin.H{"error": "only jpg, png, webp, gif allowed"}); return
	}

	filename := "sammy" + ext
	dst, err := os.Create(filepath.Join("static", filename))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()}); return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Update about.json photo field
	var about map[string]any
	if err := readJSON("data/about.json", &about); err == nil {
		about["photo"] = "/static/" + filename
		writeJSON("data/about.json", about)
	}

	c.JSON(200, gin.H{"message": "Photo uploaded", "url": "/static/" + filename})
}

// keep compiler happy — multipart is imported via FormFile
var _ multipart.File

// ── Slug helper ───────────────────────────────────────────────────────────────

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r == ' ' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// ── SPA fallback ──────────────────────────────────────────────────────────────

func spaHandler(c *gin.Context) {
	fsPath := "static" + c.Request.URL.Path
	if info, err := os.Stat(fsPath); err == nil && !info.IsDir() {
		c.File(fsPath)
		return
	}
	c.File("static/index.html")
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	_ = os.MkdirAll("data/posts", 0755)
	if _, err := os.Stat("data/contacts.json"); os.IsNotExist(err) {
		_ = os.WriteFile("data/contacts.json", []byte("[]"), 0644)
	}

	r := gin.Default()
	r.Use(corsMiddleware())

	// ── Public API
	pub := r.Group("/api")
	{
		pub.GET("/about", handleAbout)
		pub.GET("/skills", handleSkills)
		pub.GET("/projects", handleProjects)
		pub.GET("/blog", handleBlogList)
		pub.GET("/blog/:slug", handleBlogPost)
		pub.POST("/contact", handleContact)
		pub.GET("/carousel", handleGetCarousel)
	}

	// ── Admin auth (no session needed)
	r.POST("/api/admin/login", handleAdminLogin)
	r.POST("/api/admin/logout", handleAdminLogout)

	// ── Admin API (session required)
	adm := r.Group("/api/admin", authRequired())
	{
		adm.GET("/me", handleAdminMe)

		adm.GET("/about", handleAdminGetAbout)
		adm.PUT("/about", handleAdminSaveAbout)
		adm.POST("/about/photo", handleAdminUploadPhoto)

		adm.GET("/skills", handleAdminGetSkills)
		adm.PUT("/skills", handleAdminSaveSkills)

		adm.GET("/projects", handleAdminGetProjects)
		adm.PUT("/projects", handleAdminSaveProjects)

		adm.GET("/posts", handleAdminGetPosts)
		adm.GET("/posts/:slug", handleAdminGetPost)
		adm.POST("/posts", handleAdminCreatePost)
		adm.PUT("/posts/:slug", handleAdminUpdatePost)
		adm.DELETE("/posts/:slug", handleAdminDeletePost)

		adm.GET("/contacts", handleAdminGetContacts)
		adm.PUT("/contacts/:id/read", handleAdminMarkRead)
		adm.DELETE("/contacts/:id", handleAdminDeleteContact)

		adm.GET("/settings", handleAdminGetSettings)
		adm.PUT("/settings", handleAdminSaveSettings)

		adm.GET("/carousel", handleAdminGetCarousel)
		adm.PUT("/carousel", handleAdminSaveCarousel)
		adm.POST("/carousel/upload", handleAdminUploadSlideImage)
	}

	// Generate initial CSS override if settings exist but file doesn't
	if _, err := os.Stat("static/theme-override.css"); os.IsNotExist(err) {
		s := loadSettings()
		css := generateSettingsCSS(s)
		_ = os.WriteFile("static/theme-override.css", []byte(css), 0644)
	}

	// ── Static files
	r.Static("/static", "./static")
	r.StaticFile("/admin", "./admin/index.html")
	r.Static("/admin/assets", "./admin/assets")
	r.NoRoute(spaHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Portfolio  → http://localhost:%s\n", port)
	fmt.Printf("🔐 Admin CMS  → http://localhost:%s/admin\n", port)
	log.Fatal(r.Run(":" + port))
}
