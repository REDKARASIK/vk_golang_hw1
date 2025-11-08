package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type timeoutError string

func (e timeoutError) Error() string   { return string(e) }
func (e timeoutError) Timeout() bool   { return true }
func (e timeoutError) Temporary() bool { return false }

func TestFindUser_InvalidLimit_LessThanZero(t *testing.T) {
	client := SearchClient{
		AccessToken: accessToken,
		URL:         "http://example.com",
	}

	req := SearchRequest{
		Limit: -1,
	}

	if _, err := client.FindUsers(req); err == nil {
		t.Errorf("expected error for invalid limit, but got nil")
	}
}

func TestFindUser_InvalidLimit_MoreThanMax(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	client := SearchClient{
		AccessToken: accessToken,
		URL:         server.URL,
	}

	req := SearchRequest{
		Limit: 30,
	}

	data, err := client.FindUsers(req)
	if err != nil {
		t.Fatalf("expected error is nil")
	}
	if data == nil {
		t.Fatalf("expected not nil data")
	}
	if len(data.Users) != 25 {
		t.Fatalf("expected 25 users, but get %d", len(data.Users))
	}
}

func TestFindUser_InvalidOffset_LessThanZero(t *testing.T) {
	sc := SearchClient{
		AccessToken: accessToken,
		URL:         "http://example.com",
	}
	sr := SearchRequest{
		Offset: -1,
	}
	_, err := sc.FindUsers(sr)
	if err == nil {
		t.Errorf("expected error for invalid offset, got nil")
	}
}

func TestFindUsers_ClientTimeoutError(t *testing.T) {
	origClient := client
	defer func() { client = origClient }()

	client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, timeoutError("i/o timeout")
		}),
	}

	sc := SearchClient{URL: "http://example.com", AccessToken: accessToken}
	_, err := sc.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout for") {
		t.Errorf("expected timeout message, got: %v", err)
	}
}

func TestFindUsers_ClientUnknownError(t *testing.T) {
	origClient := client
	defer func() { client = origClient }()

	client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("some network failure")
		}),
	}

	sc := SearchClient{URL: "http://example.com", AccessToken: accessToken}
	_, err := sc.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Errorf("expected unknown error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown error") {
		t.Errorf("expected unknown error message, got: %v", err)
	}
}

func TestFindUsers_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_by=int&limit=1", nil)
	req.Header.Set("AccessToken", "t")

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected unauthorized")
	}
}

func TestFindUsers_BadRequest_OrderField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	sc := SearchClient{URL: ts.URL, AccessToken: accessToken}
	_, err := sc.FindUsers(SearchRequest{Limit: 1, OrderField: "Name"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "OrderFeld Name invalid") {
		t.Fatalf("expected error to contain 'OrderFeld Name invalid', got: %v", err)
	}
}

func TestFindUsers_BadRequest_OtherError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error":"some other error"}`))
	}))
	defer ts.Close()

	sc := SearchClient{URL: ts.URL}
	_, err := sc.FindUsers(SearchRequest{Limit: 1, OrderField: "Name"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown bad request error: some other error") {
		t.Fatalf("expected error to contain 'unknown bad request error: some other error', got: %v", err)
	}
}

func TestSearchServer_SuccessLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: accessToken,
	}

	req := SearchRequest{
		Limit: 3,
	}
	data, err := client.FindUsers(req)
	if err != nil {
		t.Fatalf("expected no error")
	}
	if len(data.Users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(data.Users))
	}
	for i, user := range data.Users {
		if user.ID < 0 {
			t.Fatalf("user %d has invalid ID: %d", i, user.ID)
		}
		if user.Name == "" {
			t.Fatalf("user %d has empty name", i)
		}
		if user.Age <= 0 {
			t.Fatalf("user %d has invalid age: %d", i, user.Age)
		}
	}
}

const smallXML = `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <row>
    <id>1</id>
    <first_name>Alice</first_name>
    <last_name>Anderson</last_name>
    <age>30</age>
    <about>Alpha</about>
    <gender>female</gender>
  </row>
  <row>
    <id>2</id>
    <first_name>Bob</first_name>
    <last_name>Baker</last_name>
    <age>25</age>
    <about>Beta</about>
    <gender>male</gender>
  </row>
  <row>
    <id>3</id>
    <first_name>Charlie</first_name>
    <last_name>Clark</last_name>
    <age>20</age>
    <about>Gamma</about>
    <gender>male</gender>
  </row>
</root>`

func writeTempDataset(t *testing.T, content string) (string, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "dataset-*.xml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		_ = f.Close()
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name(), func() { _ = os.Remove(f.Name()) }
}

func TestSearchServer_InvalidOrderField(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	searcherReq, _ := http.NewRequest("GET", ts.URL + "?order_field=Foo&order_by=0", nil)
	searcherReq.Header.Add("AccessToken", accessToken)
	client := http.Client{}

	res, err := client.Do(searcherReq)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if !strings.Contains(string(b), "OrderField invalid") {
		t.Fatalf("expected body to contain 'OrderField invalid', got %s", string(b))
	}
}

func TestSearchServer_LimitParseError(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	searcherReq, _ := http.NewRequest("GET", ts.URL + "?order_field=age&order_by=0&limit=H", nil)
	searcherReq.Header.Add("AccessToken", accessToken)
	client := http.Client{}

	res, err := client.Do(searcherReq)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if !strings.Contains(string(b), "cast limit to int") {
		t.Fatalf("expected cast limit to int, got %s", string(b))
	}
}

func TestSearchServer_OffsetParseError(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	searcherReq, _ := http.NewRequest("GET", ts.URL + "?order_field=age&order_by=0&limit=1&offset=int", nil)
	searcherReq.Header.Add("AccessToken", accessToken)
	client := http.Client{}

	res, err := client.Do(searcherReq)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if !strings.Contains(string(b), "cast offset to int") {
		t.Fatalf("expected cast offset to int, got %s", string(b))
	}
}

func TestSearchServer_QueryFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?query=Boyd&order_by=0&limit=10", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()


	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user with 'Boyd' in name, got %d", len(users))
	}
	if users[0].ID != 0 || users[0].Name != "Boyd Wolf" {
		t.Fatalf("expected Boyd Wolf (ID=0), got %s (ID=%d)", users[0].Name, users[0].ID)
	}
}

func TestSearchServer_OrderByNameAsc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_field=name&order_by=1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем сортировку по возрастанию имени из реального dataset
	for i := 0; i < len(users)-1; i++ {
		if users[i].Name > users[i+1].Name {
			t.Fatalf("expected ascending names, got %v > %v at position %d", users[i].Name, users[i+1].Name, i)
		}
	}
}

func TestSearchServer_OrderByIdAsc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_field=id&order_by=1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем конкретные ID из dataset: 0, 1, 2, 3, 4
	if users[0].ID != 0 || users[1].ID != 1 || users[2].ID != 2 || users[3].ID != 3 || users[4].ID != 4 {
		t.Fatalf("expected IDs 0,1,2,3,4, got %d,%d,%d,%d,%d", users[0].ID, users[1].ID, users[2].ID, users[3].ID, users[4].ID)
	}
}

func TestSearchServer_OffsetAndLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?offset=1&limit=1", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

type failingWriter struct {
	header http.Header
	first  bool
	buf    []byte
}

func (f *failingWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}
func (f *failingWriter) Write(p []byte) (int, error) {
	if !f.first {
		f.first = true
		return 0, errors.New("write failed")
	}
	f.buf = append(f.buf, p...)
	return len(p), nil
}
func (f *failingWriter) WriteHeader(status int) {}

func TestSearchServer_EncodeFailure(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	req := httptest.NewRequest("GET", "/?limit=3", nil)
	req.Header.Add("AccessToken", accessToken)
	w := &failingWriter{}

	SearchServer(w, req)

	if !strings.Contains(string(w.buf), "encoding json failed") {
		t.Fatalf("expected encoding failed message in response writer buffer, got: %s", string(w.buf))
	}
}

func TestSearchServer_NoLimitToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_by=-1&offset=1", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected bad request")
	}
}

func TestSearchServer_InvalidOrderBy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_by=int&limit=1", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected bad request")
	}
}

func TestSearchServer_InvalidHTTPMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"?order_by=int&limit=1", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected not found")
	}
}

func TestSearchServer_OrderByNameDesc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_by=-1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем сортировку по убыванию имени из реального dataset
	for i := 0; i < len(users)-1; i++ {
		if users[i].Name < users[i+1].Name {
			t.Fatalf("expected descending names, got %v < %v at position %d", users[i].Name, users[i+1].Name, i)
		}
	}
}

func TestSearchServer_OrderByIdDesc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	// Сортировка по ID по убыванию из реального dataset

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_field=id&order_by=-1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем убывание ID: должны быть самые большие ID в начале
	for i := 0; i < len(users)-1; i++ {
		if users[i].ID < users[i+1].ID {
			t.Fatalf("expected descending IDs, got %d < %d at position %d", users[i].ID, users[i+1].ID, i)
		}
	}
}

func TestSearchServer_OrderByAgeAsc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_field=age&order_by=1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем сортировку по возрасту из реального dataset
	for i := 0; i < len(users)-1; i++ {
		if users[i].Age > users[i+1].Age {
			t.Fatalf("expected ascending ages, got %d > %d at position %d", users[i].Age, users[i+1].Age, i)
		}
	}
}

func TestSearchServer_OrderByAgeDesc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_field=age&order_by=-1&limit=5", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
	// Проверяем сортировку по убыванию возраста из реального dataset
	for i := 0; i < len(users)-1; i++ {
		if users[i].Age < users[i+1].Age {
			t.Fatalf("expected descending ages, got %d < %d at position %d", users[i].Age, users[i+1].Age, i)
		}
	}
}

func TestSearchServer_OrderByAsIs(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?o?order_by=0&limit=3", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
	if users[0].ID != 1 || users[1].ID != 2 || users[2].ID != 3 {
		t.Fatalf("expected AsIs order (1,2,3), got %v, %v, %v", users[0].ID, users[1].ID, users[2].ID)
	}
}

func TestFindUsers_InternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sc := SearchClient{URL: ts.URL}
	_, err := sc.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "SearchServer fatal error") {
		t.Fatalf("expected error to contain 'SearchServer fatal error', got: %v", err)
	}
}

func TestFindUsers_BadJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json`))
	}))
	defer ts.Close()

	sc := SearchClient{URL: ts.URL}
	_, err := sc.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "cant unpack result json") {
		t.Fatalf("expected error to contain 'cant unpack result json', got: %v", err)
	}
}

func TestSearchServer_OffsetGreaterThanLen(t *testing.T) {
	orig := datasetPath
	tmp, cleanup := writeTempDataset(t, smallXML)
	datasetPath = tmp
	defer func() { datasetPath = orig; cleanup() }()

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?offset=100&limit=25", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users (offset > len), got %d", len(users))
	}
}

func TestSearchServer_ComplexQueryFromSpec(t *testing.T) {
	// Используем реальный dataset.xml без изменения пути
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"?order_by=-1&order_field=age&limit=1&offset=0&query=on", nil)
	req.Header.Set("AccessToken", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil { t.Fatalf("get failed: %v", err) }
	defer res.Body.Close()

	var users []User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if !strings.Contains(users[0].Name, "on") && !strings.Contains(users[0].About, "on") {
		t.Fatalf("expected user with 'on' in Name or About, got Name=%s, About=%s", users[0].Name, users[0].About)
	}
}

func TestFindUsers_BadRequestUnpackError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`not valid json`))
	}))
	defer ts.Close()

	sc := SearchClient{URL: ts.URL}
	_, err := sc.FindUsers(SearchRequest{Limit: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "cant unpack error json") {
		t.Fatalf("expected error to contain 'cant unpack error json', got: %v", err)
	}
}
