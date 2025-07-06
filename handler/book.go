package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abaihaaqi/kutubu-api/model"
	"github.com/abaihaaqi/kutubu-api/repository"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{Repo: repo}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	userID := GetUserID(c)
	books, err := h.Repo.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	userID := GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.Repo.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	userID := GetUserID(c)

	// Ambil data dari form
	title := c.PostForm("title")
	author := c.PostForm("author")
	year, _ := strconv.Atoi(c.PostForm("year"))
	category := c.PostForm("category")

	// Inisialisasi model book
	book := model.Book{
		Title:    title,
		Author:   author,
		Year:     year,
		Category: category,
		UserID:   userID,
	}

	// Proses upload cover image
	file, err := c.FormFile("cover_image")
	if err == nil {
		// Simpan file dengan nama unik (misalnya pakai timestamp)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		path := "uploads/" + filename

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		// Simpan path ke field CoverImage
		book.CoverImage = path
	}

	// Simpan ke database
	if err := h.Repo.Create(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	userID := GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	// Ambil buku berdasarkan ID dan userID
	book, err := h.Repo.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Ambil data dari form
	title := c.PostForm("title")
	author := c.PostForm("author")
	year, _ := strconv.Atoi(c.PostForm("year"))
	category := c.PostForm("category")
	book.Title = title
	book.Author = author
	book.Year = year
	book.Category = category

	// Proses upload gambar baru (jika ada)
	file, err := c.FormFile("cover_image")
	if err == nil {
		// Hapus gambar lama jika ada
		if book.CoverImage != "" {
			if err := os.Remove(book.CoverImage); err != nil && !os.IsNotExist(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus gambar lama"})
				return
			}
		}

		// Simpan file baru
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		path := "uploads/" + filename

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar baru"})
			return
		}

		book.CoverImage = path
	}

	// Update ke database
	if err := h.Repo.Update(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	userID := GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	// Ambil data buku berdasarkan ID dan user
	book, err := h.Repo.GetByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Hapus gambar cover jika ada
	if book.CoverImage != "" {
		if err := os.Remove(book.CoverImage); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus gambar cover"})
			return
		}
	}

	// Hapus buku dari database
	if err := h.Repo.Delete(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku dan gambar berhasil dihapus"})
}

func (h *BookHandler) SearchBooks(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query pencarian tidak boleh kosong"})
		return
	}

	userID := GetUserID(c)
	books, err := h.Repo.Search(keyword, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
