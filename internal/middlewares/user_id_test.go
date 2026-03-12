package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"financialcontrol/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserIDMiddlewareSuccess(t *testing.T) {
	userID := uuid.New()

	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		ctxUserID, exists := c.Get(constants.UserID)
		if !exists {
			t.Error("UserID not found in context")
			c.Status(http.StatusInternalServerError)
			return
		}

		if ctxUserID.(uuid.UUID) != userID {
			t.Errorf("Expected UserID %v, got %v", userID, ctxUserID)
		}

		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: userID.String(),
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserIDMiddlewareNoCookie(t *testing.T) {
	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserIDMiddlewareInvalidUUID(t *testing.T) {
	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: "invalid-uuid",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserIDMiddlewareEmptyCookie(t *testing.T) {
	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: "",
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserIDMiddlewareNextCalled(t *testing.T) {
	userID := uuid.New()
	nextCalled := false

	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: userID.String(),
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("Expected next handler to be called")
	}
}

func TestUserIDMiddlewareNextNotCalledOnError(t *testing.T) {
	nextCalled := false

	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if nextCalled {
		t.Error("Expected next handler NOT to be called on error")
	}
}

func TestUserIDMiddlewareMultipleRequests(t *testing.T) {
	router := setupTestRouter()
	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// First request - valid
	userID1 := uuid.New()
	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req1.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: userID1.String(),
	})
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Request 1: Expected status %d, got %d", http.StatusOK, w1.Code)
	}

	// Second request - invalid
	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Request 2: Expected status %d, got %d", http.StatusUnauthorized, w2.Code)
	}

	// Third request - valid with different user
	userID3 := uuid.New()
	req3, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req3.AddCookie(&http.Cookie{
		Name:  constants.UserID,
		Value: userID3.String(),
	})
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Request 3: Expected status %d, got %d", http.StatusOK, w3.Code)
	}
}
