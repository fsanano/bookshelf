package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func MD5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// BookInfo holds extracted details from OpenLibrary
type BookInfo struct {
	Title         string
	Author        string
	Cover         string
	PublishedYear int
	Pages         int
}

// FetchBookInfoFromOpenLibrary fetches and extracts required book info from OpenLibrary
func FetchBookInfoFromOpenLibrary(isbn string) (*BookInfo, error) {
	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", isbn)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-OK response from OpenLibrary: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}

	key := "ISBN:" + isbn
	val, ok := data[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("ISBN not found in response")
	}

	info := &BookInfo{}

	// Title
	if title, ok := val["title"].(string); ok {
		info.Title = title
	} else {
		info.Title = ""
	}

	// Author
	if authors, ok := val["authors"].([]interface{}); ok && len(authors) > 0 {
		if firstAuthor, ok := authors[0].(map[string]interface{}); ok {
			if authorName, ok := firstAuthor["name"].(string); ok {
				info.Author = authorName
			}
		}
	}

	// Cover
	if cover, ok := val["cover"].(map[string]interface{}); ok {
		// Prefer large cover if available
		if coverURL, ok := cover["large"].(string); ok {
			info.Cover = coverURL
		} else if coverURL, ok := cover["medium"].(string); ok {
			info.Cover = coverURL
		} else if coverURL, ok := cover["small"].(string); ok {
			info.Cover = coverURL
		}
	}

	// Published Year
	// If publish_date is something like "2012" or "May 1, 2012"
	// We'll try to extract the year:
	if pubDate, ok := val["publish_date"].(string); ok && pubDate != "" {
		// Extract year using a regex
		re := regexp.MustCompile(`\b(\d{4})\b`)
		yearMatch := re.FindString(pubDate)
		if yearMatch != "" {
			if yearInt, err := strconv.Atoi(yearMatch); err == nil {
				info.PublishedYear = yearInt
			}
		}
	}

	// Pages
	// Check number_of_pages first
	if numPages, ok := val["number_of_pages"].(float64); ok {
		info.Pages = int(numPages)
	} else if pagination, ok := val["pagination"].(string); ok && pagination != "" {
		// pagination might be something like "221 pages"
		pagesStr := strings.TrimSpace(strings.Replace(pagination, "pages", "", -1))
		pagesStr = strings.TrimSpace(pagesStr)
		if p, err := strconv.Atoi(pagesStr); err == nil {
			info.Pages = p
		}
	}

	return info, nil
}
