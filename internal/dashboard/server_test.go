package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grigata/chta-network-stats/internal/models"
)

func TestDashboardHandler(t *testing.T) {
	data := NewData("test", 42, time.Second, []models.NetworkBlock{{Height: 42, Hash: "abc", Type: "NORMAL", Pool: "TestPool", Difficulty: 2}})
	h, err := Handler(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/styles.css", "/app.js"} {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, w.Code)
		}
	}
	r := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	var got Data
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Height != 42 || len(got.Blocks) != 1 {
		t.Fatalf("unexpected data: %+v", got)
	}
}
