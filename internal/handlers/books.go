package handlers

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"bookshelf/internal/models"
	"bookshelf/internal/utils"
)

// GET /books/:title - Search books by title
func GetBooksByTitleHandler(c *fiber.Ctx) error {
	// Get the title parameter and decode any URL encoding
	titleParam, err := url.QueryUnescape(c.Params("title"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid title parameter",
		})
	}
	titleParam = strings.TrimSpace(titleParam)
	searchLower := strings.ToLower(titleParam)

	models.BooksStore.Sync.RLock()
	defer models.BooksStore.Sync.RUnlock()

	userKey := c.Locals("userKey").(string)

	var results []map[string]interface{}

	for _, b := range models.BooksStore.Data {
		// Check the book owner
		if b.User != userKey {
			continue
		}

		if strings.Contains(strings.ToLower(b.Title), searchLower) {
			// Return only the fields specified in the search response
			bookMap := map[string]interface{}{
				"isbn":      b.ISBN,
				"title":     b.Title,
				"cover":     b.Cover,
				"author":    b.Author,
				"published": b.PublishedYear,
			}
			results = append(results, bookMap)
		}
	}

	return c.JSON(models.Response{
		Data:    results,
		IsOk:    true,
		Message: "ok",
	})
}

// GET /books
func GetBooksHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	models.BooksStore.Sync.RLock()
	defer models.BooksStore.Sync.RUnlock()

	var result []models.Book
	for _, b := range models.BooksStore.Data {
		if b.User == userKey {
			result = append(result, b)
		}
	}
	return c.JSON(models.Response{
		Data:    result,
		IsOk:    true,
		Message: "ok",
	})
}

// POST /books
func AddBookHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	var payload struct {
		ISBN string `json:"isbn"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid request body",
		})
	}

	if payload.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "ISBN required",
		})
	}

	// Fetch detailed book info from OpenLibrary
	bookInfo, err := utils.FetchBookInfoFromOpenLibrary(payload.ISBN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Could not fetch book info",
		})
	}

	// Status for a newly created book is 0 ("new")
	status := 0
	// Create a new book entry
	models.BooksStore.Sync.Lock()
	newBook := models.Book{
		ID:            models.BooksStore.Next,
		ISBN:          payload.ISBN,
		Status:        status,
		User:          userKey,
		Title:         bookInfo.Title,
		Cover:         bookInfo.Cover,
		Author:        bookInfo.Author,
		PublishedYear: bookInfo.PublishedYear,
		Pages:         bookInfo.Pages,
	}
	models.BooksStore.Data = append(models.BooksStore.Data, newBook)
	models.BooksStore.Next++
	models.BooksStore.Sync.Unlock()

	// Construct response as per documentation
	responseData := map[string]interface{}{
		"book": map[string]interface{}{
			"id":        newBook.ID,
			"isbn":      newBook.ISBN,
			"title":     bookInfo.Title,
			"cover":     bookInfo.Cover,
			"author":    bookInfo.Author,
			"published": bookInfo.PublishedYear,
			"pages":     bookInfo.Pages,
		},
		"status": status,
	}

	return c.JSON(models.Response{
		Data:    responseData,
		IsOk:    true,
		Message: "ok",
	})
}

// PUT /books/:id
func UpdateBookHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	bookID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid ID",
		})
	}

	var payload struct {
		Status int `json:"status"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid request body",
		})
	}

	models.BooksStore.Sync.Lock()
	defer models.BooksStore.Sync.Unlock()

	found := false
	for i, b := range models.BooksStore.Data {
		if b.ID == bookID && b.User == userKey {
			models.BooksStore.Data[i].Status = payload.Status
			found = true
			return c.JSON(models.Response{
				Data:    models.BooksStore.Data[i],
				IsOk:    true,
				Message: "ok",
			})
		}
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Book not found",
		})
	}

	return nil
}

// DELETE /books/:id
func DeleteBookHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	bookID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid ID",
		})
	}

	models.BooksStore.Sync.Lock()
	defer models.BooksStore.Sync.Unlock()
	idx := -1
	for i, b := range models.BooksStore.Data {
		if b.ID == bookID && b.User == userKey {
			idx = i
			break
		}
	}

	if idx == -1 {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Book not found",
		})
	}

	// Delete the book
	models.BooksStore.Data = append(models.BooksStore.Data[:idx], models.BooksStore.Data[idx+1:]...)

	// return all books after deletion
	var result []models.Book
	for _, b := range models.BooksStore.Data {
		if b.User == userKey {
			result = append(result, b)
		}
	}

	return c.JSON(models.Response{
		Data:    result,
		IsOk:    true,
		Message: "ok",
	})
}

// GET /cleanup
func CleanupHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)
	models.BooksStore.Sync.Lock()
	defer models.BooksStore.Sync.Unlock()

	// remove all books of this user
	var filtered []models.Book
	for _, b := range models.BooksStore.Data {
		if b.User != userKey {
			filtered = append(filtered, b)
		}
	}
	models.BooksStore.Data = filtered

	return c.JSON(models.Response{
		Data:    nil,
		IsOk:    true,
		Message: "ok",
	})
}
