package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain_NoMiddlewares(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	chained := Chain(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	chained(rr, req)

	if !called {
		t.Error("expected handler to be called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestChain_SingleMiddleware(t *testing.T) {
	order := []string{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}
	mw := func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw")
			f(w, r)
		}
	}

	chained := Chain(handler, mw)
	chained(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if len(order) != 2 || order[0] != "mw" || order[1] != "handler" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

func TestChain_MultipleMiddlewares_ExecutionOrder(t *testing.T) {
	// Chain(handler, mw1, mw2, mw3) should execute mw1 → mw2 → mw3 → handler
	order := []string{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}
	mw := func(name string) Middleware {
		return func(f http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				f(w, r)
			}
		}
	}

	chained := Chain(handler, mw("mw1"), mw("mw2"), mw("mw3"))
	chained(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	expected := []string{"mw1", "mw2", "mw3", "handler"}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("position %d: expected %q, got %q (full order: %v)", i, v, order[i], order)
		}
	}
}

func TestHttpMethod_AllowedMethod(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	mw := HttpMethod(http.MethodGet)(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw(rr, req)

	if !called {
		t.Error("expected handler to be called for allowed method")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHttpMethod_DisallowedMethod(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	mw := HttpMethod(http.MethodGet)(handler)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	mw(rr, req)

	if called {
		t.Error("expected handler NOT to be called for disallowed method")
	}
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestHttpMethod_PostAllowed(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	mw := HttpMethod(http.MethodPost)(handler)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	mw(rr, req)

	if !called {
		t.Error("expected handler to be called for POST")
	}
}

func TestJsonHeaderMiddleware_SetsContentType(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mw := JsonHeaderMiddleware()(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw(rr, req)

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestJsonHeaderMiddleware_HandlerStillExecutes(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	mw := JsonHeaderMiddleware()(handler)
	mw(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if !called {
		t.Error("expected handler to still be called after JsonHeaderMiddleware")
	}
}

func TestCorsHeaderMiddleware_SetsHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mw := CorsHeaderMiddleware()(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mw(rr, req)

	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin *, got %q", origin)
	}
}

func TestCorsHeaderMiddleware_HandlerStillExecutes(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	mw := CorsHeaderMiddleware()(handler)
	mw(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if !called {
		t.Error("expected handler to still be called after CorsHeaderMiddleware")
	}
}

func TestChain_CombinedMiddlewares(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	chained := Chain(handler,
		HttpMethod(http.MethodGet),
		JsonHeaderMiddleware(),
		CorsHeaderMiddleware(),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	chained(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin *, got %q", origin)
	}
}

func TestChain_CombinedMiddlewares_WrongMethod(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	chained := Chain(handler,
		HttpMethod(http.MethodGet),
		JsonHeaderMiddleware(),
		CorsHeaderMiddleware(),
	)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	chained(rr, req)

	if called {
		t.Error("handler should not be called when method is not allowed")
	}
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
