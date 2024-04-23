package api

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

var _ Context = (*MockContext)(nil)

type MockContext struct {
	mock.Mock
}

// Attachment implements echo.Context.
func (m *MockContext) Attachment(file string, name string) error {
	panic("unimplemented")
}

// Blob implements echo.Context.
func (m *MockContext) Blob(code int, contentType string, b []byte) error {
	panic("unimplemented")
}

// Cookie implements echo.Context.
func (m *MockContext) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

// Cookies implements echo.Context.
func (m *MockContext) Cookies() []*http.Cookie {
	panic("unimplemented")
}

// Echo implements echo.Context.
func (m *MockContext) Echo() *echo.Echo {
	panic("unimplemented")
}

// Error implements echo.Context.
func (m *MockContext) Error(err error) {
	panic("unimplemented")
}

// File implements echo.Context.
func (m *MockContext) File(file string) error {
	panic("unimplemented")
}

// FormFile implements echo.Context.
func (m *MockContext) FormFile(name string) (*multipart.FileHeader, error) {
	panic("unimplemented")
}

// FormParams implements echo.Context.
func (m *MockContext) FormParams() (url.Values, error) {
	panic("unimplemented")
}

// FormValue implements echo.Context.
func (m *MockContext) FormValue(name string) string {
	panic("unimplemented")
}

// Get implements echo.Context.
func (m *MockContext) Get(key string) interface{} {
	panic("unimplemented")
}

// HTML implements echo.Context.
func (m *MockContext) HTML(code int, html string) error {
	panic("unimplemented")
}

// HTMLBlob implements echo.Context.
func (m *MockContext) HTMLBlob(code int, b []byte) error {
	panic("unimplemented")
}

// Handler implements echo.Context.
func (m *MockContext) Handler() echo.HandlerFunc {
	panic("unimplemented")
}

// Inline implements echo.Context.
func (m *MockContext) Inline(file string, name string) error {
	panic("unimplemented")
}

// IsTLS implements echo.Context.
func (m *MockContext) IsTLS() bool {
	panic("unimplemented")
}

// IsWebSocket implements echo.Context.
func (m *MockContext) IsWebSocket() bool {
	panic("unimplemented")
}

// JSONBlob implements echo.Context.
func (m *MockContext) JSONBlob(code int, b []byte) error {
	panic("unimplemented")
}

// JSONP implements echo.Context.
func (m *MockContext) JSONP(code int, callback string, i interface{}) error {
	panic("unimplemented")
}

// JSONPBlob implements echo.Context.
func (m *MockContext) JSONPBlob(code int, callback string, b []byte) error {
	panic("unimplemented")
}

// JSONPretty implements echo.Context.
func (m *MockContext) JSONPretty(code int, i interface{}, indent string) error {
	panic("unimplemented")
}

// Logger implements echo.Context.
func (m *MockContext) Logger() echo.Logger {
	panic("unimplemented")
}

// MultipartForm implements echo.Context.
func (m *MockContext) MultipartForm() (*multipart.Form, error) {
	panic("unimplemented")
}

// NoContent implements echo.Context.
func (m *MockContext) NoContent(code int) error {
	panic("unimplemented")
}

// Param implements echo.Context.
func (m *MockContext) Param(name string) string {
	panic("unimplemented")
}

// ParamNames implements echo.Context.
func (m *MockContext) ParamNames() []string {
	panic("unimplemented")
}

// ParamValues implements echo.Context.
func (m *MockContext) ParamValues() []string {
	panic("unimplemented")
}

// Path implements echo.Context.
func (m *MockContext) Path() string {
	panic("unimplemented")
}

// QueryParam implements echo.Context.
func (m *MockContext) QueryParam(name string) string {
	panic("unimplemented")
}

// QueryParams implements echo.Context.
func (m *MockContext) QueryParams() url.Values {
	panic("unimplemented")
}

// QueryString implements echo.Context.
func (m *MockContext) QueryString() string {
	panic("unimplemented")
}

// RealIP implements echo.Context.
func (m *MockContext) RealIP() string {
	panic("unimplemented")
}

// Redirect implements echo.Context.
func (m *MockContext) Redirect(code int, url string) error {
	panic("unimplemented")
}

// Render implements echo.Context.
func (m *MockContext) Render(code int, name string, data interface{}) error {
	panic("unimplemented")
}

// Request implements echo.Context.
func (m *MockContext) Request() *http.Request {
	panic("unimplemented")
}

// Reset implements echo.Context.
func (m *MockContext) Reset(r *http.Request, w http.ResponseWriter) {
	panic("unimplemented")
}

// Response implements echo.Context.
func (m *MockContext) Response() *echo.Response {
	panic("unimplemented")
}

// Scheme implements echo.Context.
func (m *MockContext) Scheme() string {
	panic("unimplemented")
}

// Set implements echo.Context.
func (m *MockContext) Set(key string, val interface{}) {
	panic("unimplemented")
}

// SetCookie implements echo.Context.
func (m *MockContext) SetCookie(cookie *http.Cookie) {
	panic("unimplemented")
}

// SetHandler implements echo.Context.
func (m *MockContext) SetHandler(h echo.HandlerFunc) {
	panic("unimplemented")
}

// SetLogger implements echo.Context.
func (m *MockContext) SetLogger(l echo.Logger) {
	panic("unimplemented")
}

// SetParamNames implements echo.Context.
func (m *MockContext) SetParamNames(names ...string) {
	panic("unimplemented")
}

// SetParamValues implements echo.Context.
func (m *MockContext) SetParamValues(values ...string) {
	panic("unimplemented")
}

// SetPath implements echo.Context.
func (m *MockContext) SetPath(p string) {
	panic("unimplemented")
}

// SetRequest implements echo.Context.
func (m *MockContext) SetRequest(r *http.Request) {
	panic("unimplemented")
}

// SetResponse implements echo.Context.
func (m *MockContext) SetResponse(r *echo.Response) {
	panic("unimplemented")
}

// Stream implements echo.Context.
func (m *MockContext) Stream(code int, contentType string, r io.Reader) error {
	panic("unimplemented")
}

// String implements echo.Context.
// Subtle: this method shadows the method (Mock).String of MockContext.Mock.
func (m *MockContext) String(code int, s string) error {
	panic("unimplemented")
}

func (m *MockContext) Validate(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

// XML implements echo.Context.
func (m *MockContext) XML(code int, i interface{}) error {
	panic("unimplemented")
}

// XMLBlob implements echo.Context.
func (m *MockContext) XMLBlob(code int, b []byte) error {
	panic("unimplemented")
}

// XMLPretty implements echo.Context.
func (m *MockContext) XMLPretty(code int, i interface{}, indent string) error {
	panic("unimplemented")
}

func (m *MockContext) Bind(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

func (m *MockContext) JSON(status int, i interface{}) error {
	args := m.Called(status, i)
	return args.Error(0)
}
