// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// User defines model for User.
type User struct {
	Email     *openapi_types.Email `json:"email,omitempty"`
	FirstName *string              `json:"firstName,omitempty"`
	Id        *int64               `json:"id,omitempty"`
	LastName  *string              `json:"lastName,omitempty"`
	Phone     *string              `json:"phone,omitempty"`
	Username  *string              `json:"username,omitempty"`
}

// CreateUserJSONBody defines parameters for CreateUser.
type CreateUserJSONBody = User

// UpdateUserJSONBody defines parameters for UpdateUser.
type UpdateUserJSONBody = User

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserJSONBody

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = UpdateUserJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create user
	// (POST /user)
	CreateUser(w http.ResponseWriter, r *http.Request)

	// (DELETE /user/{userId})
	DeleteUser(w http.ResponseWriter, r *http.Request, userId int64)

	// (GET /user/{userId})
	FindUserById(w http.ResponseWriter, r *http.Request, userId int64)

	// (PUT /user/{userId})
	UpdateUser(w http.ResponseWriter, r *http.Request, userId int64)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// CreateUser operation middleware
func (siw *ServerInterfaceWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateUser(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// DeleteUser operation middleware
func (siw *ServerInterfaceWrapper) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int64

	err = runtime.BindStyledParameter("simple", false, "userId", chi.URLParam(r, "userId"), &userId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteUser(w, r, userId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// FindUserById operation middleware
func (siw *ServerInterfaceWrapper) FindUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int64

	err = runtime.BindStyledParameter("simple", false, "userId", chi.URLParam(r, "userId"), &userId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindUserById(w, r, userId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// UpdateUser operation middleware
func (siw *ServerInterfaceWrapper) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userId" -------------
	var userId int64

	err = runtime.BindStyledParameter("simple", false, "userId", chi.URLParam(r, "userId"), &userId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userId", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateUser(w, r, userId)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/user", wrapper.CreateUser)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/user/{userId}", wrapper.DeleteUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{userId}", wrapper.FindUserById)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/user/{userId}", wrapper.UpdateUser)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWSY/bNhT+K8Rrj6qlzJK2OiWpU9RFNzTNqfWBFp8lTiSS5SOdMQb+7wUfZY8XNQmK",
	"HiYwYEs0+ZZvIfkAjR2cNWgCQf0A1HQ4SH587b316cF569AHjTzcWIXpd239IAPUoE24voICwtZhfsUW",
	"PewKGJBItjx7/JOC16aF3a4Aj39H7VFB/WeO+Th/eQhmV3fYhBTrLeFEMThI3Z9Uk0eK84QFrLWn8Isc",
	"psopQKvznp7fTPbUyw9EcZ01p+DkkYlyIqE3Y5xB3v+Epg0d1Fe3z4sJsM7gSEParG3mwwTZhCM0mMSg",
	"jdHv7OZFmwZnjR1SVoXUeO2CtgZq+KPTJDQJ0oPrUTS9RhPEy98Wf5lUsg59SpqQF2/Qb3STOtmgp7z8",
	"2ayaVSmqdWik01DDNQ8V4GTomKAy7nmzxDVOVNBII6zpt2KFQlmDYrUVoUPR27ZFJbQRKcgMOJGXaelC",
	"QQ3feZQBWRhZTkjhlVXbPSpoOKN0rtcNLyvvKKV9ALyXqeUseX78al8oxWGQfgs1vM6TUs+yj3iM8Aop",
	"3NnOvFAWGdsTfcGPtksIPmoF5jbFGfXx9bOquqqq6+qmOhYCpIjK4u0332aKsxlT3i89rqGGL8pHt5aj",
	"VUsGgOefYpvhUQyeGJVzbLvgI7IPyVlDGQuFaxn7CaIoNg0SrWMvDiSwMh/hyvk4XVKPbClZm1+XaSZL",
	"oXxI3wu1yyl6DHiZLI+TkIK0afscU6wkoRLWsDYWc0ExEYvqQhdzXj7qwkkvBwzoUzXniRZzYdf7ivHe",
	"9by3rWVPmAyW/CtDBwWMBOXaL1A8ZuoTthAKWzZWth3slmckXFU3l5gwBBkYlZ18YOojWv80FeXdfkJG",
	"0eC9wyYpCcc5F+wW0OKEaH7HEL1JPJ4SeOB1MS+EXjOhuT+LJIwNopMbFJIlJ4I9TLig+nttVCL61ZZp",
	"+SzJrv43Dv9tJ2Bs91mfoHhcnBDPW6f224l4r0Mn+Bz6kO/zis/M90/p1Iq9xo8cWzdX/PkPx9NuWvgT",
	"So3M49Pb5bhn9JtpTY24802JdRV9wrwLwdVlOVI3U3aQ2pTS6XLzrIRdcR7mzXvZtuh/iKt0ERMvY7Di",
	"Z9u8SzfB45hUl+VG+5ALmlFe1sVVIra0IdLtLZ+4VI73tOWhofOcv+5tREKubAx7ZxypH3bL3T8BAAD/",
	"/wSi7H8sDAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}