package utils

import (
	"testing"
	"os"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func TestNormalizeDate(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		ok       bool
	}{
		{"2025-08-10T22:10:00Z", "2025-08-10", true},
		{"2025-08-10 22:10:00", "2025-08-10", true},
		{"2025-08-10", "2025-08-10", true},
		{"2025-08-10 22:10:00.000", "2025-08-10", true},
		{"2025-08-10T22:10:00.000", "2025-08-10", true},
		{"not-a-date", "", false},
	}
	for _, c := range cases {
		out, err := NormalizeDate(c.in)
		if c.ok {
			assert.NoError(t, err, c.in)
			assert.Equal(t, c.expected, out, c.in)
		} else {
			assert.Error(t, err, c.in)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	assert.Equal(t, "foo", SanitizeString(" foo "))
	assert.Equal(t, "bar", SanitizeString("bar"))
}

func TestSanitizeInt(t *testing.T) {
	assert.Equal(t, 42, SanitizeInt(42))
}

func TestSanitizeFloat(t *testing.T) {
	assert.Equal(t, 3.14, SanitizeFloat(3.14))
}

func TestRoundFloat(t *testing.T) {
	assert.Equal(t, 3.14, RoundFloat(3.14159, 2))
	assert.Equal(t, 3.1, RoundFloat(3.14, 1))
	assert.Equal(t, 3.0, RoundFloat(3.01, 0))
}

func TestGetenv(t *testing.T) {
	os.Setenv("FOO_BAR", "baz")
	defer os.Unsetenv("FOO_BAR")
	assert.Equal(t, "baz", Getenv("FOO_BAR"))
}

func TestParseQueryInt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	req := httptest.NewRequest("GET", "/?num=42", nil)
	c.Request = req
	assert.Equal(t, 42, ParseQueryInt(c, "num", 99))
	req = httptest.NewRequest("GET", "/", nil)
	c.Request = req
	assert.Equal(t, 99, ParseQueryInt(c, "missing", 99))
	req = httptest.NewRequest("GET", "/?bad=notanint", nil)
	c.Request = req
	assert.Equal(t, 99, ParseQueryInt(c, "bad", 99))
}
