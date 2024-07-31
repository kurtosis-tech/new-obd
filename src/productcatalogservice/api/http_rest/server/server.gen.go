// Package productcatalogservice_server_rest_server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.1-0.20240604070534-2f0ff757704b DO NOT EDIT.
package productcatalogservice_server_rest_server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	. "github.com/kurtosis-tech/new-obd/src/productcatalogservice/api/http_rest/types"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List products
	// (GET /products)
	GetProducts(ctx echo.Context) error
	// Get product by id
	// (GET /products/{id})
	GetProductsId(ctx echo.Context, id Id) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetProducts converts echo context to params.
func (w *ServerInterfaceWrapper) GetProducts(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProducts(ctx)
	return err
}

// GetProductsId converts echo context to params.
func (w *ServerInterfaceWrapper) GetProductsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id Id

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProductsId(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/products", wrapper.GetProducts)
	router.GET(baseURL+"/products/:id", wrapper.GetProductsId)

}

type NotOkJSONResponse ResponseInfo

type GetProductsRequestObject struct {
}

type GetProductsResponseObject interface {
	VisitGetProductsResponse(w http.ResponseWriter) error
}

type GetProducts200JSONResponse []Product

func (response GetProducts200JSONResponse) VisitGetProductsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetProductsdefaultJSONResponse struct {
	Body       ResponseInfo
	StatusCode int
}

func (response GetProductsdefaultJSONResponse) VisitGetProductsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetProductsIdRequestObject struct {
	Id Id `json:"id"`
}

type GetProductsIdResponseObject interface {
	VisitGetProductsIdResponse(w http.ResponseWriter) error
}

type GetProductsId200JSONResponse Product

func (response GetProductsId200JSONResponse) VisitGetProductsIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetProductsIddefaultJSONResponse struct {
	Body       ResponseInfo
	StatusCode int
}

func (response GetProductsIddefaultJSONResponse) VisitGetProductsIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// List products
	// (GET /products)
	GetProducts(ctx context.Context, request GetProductsRequestObject) (GetProductsResponseObject, error)
	// Get product by id
	// (GET /products/{id})
	GetProductsId(ctx context.Context, request GetProductsIdRequestObject) (GetProductsIdResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetProducts operation middleware
func (sh *strictHandler) GetProducts(ctx echo.Context) error {
	var request GetProductsRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetProducts(ctx.Request().Context(), request.(GetProductsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetProducts")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetProductsResponseObject); ok {
		return validResponse.VisitGetProductsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetProductsId operation middleware
func (sh *strictHandler) GetProductsId(ctx echo.Context, id Id) error {
	var request GetProductsIdRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetProductsId(ctx.Request().Context(), request.(GetProductsIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetProductsId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetProductsIdResponseObject); ok {
		return validResponse.VisitGetProductsIdResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xV0U4bOxD9FWvulXhZsrlw1Yd9QypFUduAAlUfEELGniSmWdsdj1GjaP+9sneTELKh",
	"VJX6EmW9xzNnzpmZXYFytXcWLQeoVuAlyRoZKT8ZnX41BkXGs3EWKvDkdFQsjIYCTD6RPIcCrKwRKsjn",
	"hN+jIdRQMUUsIKg51jIF46VPqMBk7Ayapkng4J0NmFOOHV9+S3+Us4yW01/p/cIomQiUjyGxWD2L+C/h",
	"FCr4p9xWUrZvQznpQo/s1LXJdov5YvGHR8WoBRI5ggTpLqfYn53FZdaFnEdi05JUkQitWt4rp7GnqiSG",
	"dRk6dVRLTrpYPj2BYg01lnGGlLDRGt7Dvvu/B9tsjtzDIypOt69aP3pISsaZo+7JMNahl2p3IInkEl4q",
	"1INvm6Kn4rpfCm8URzrwjozC+xj0r4xsjegVYMfjfRU6hzbaxsNG1BiCnOErKr2t224Stu3s9RjctgG2",
	"OYqW2d0rBd10KdHGOkU4n0wuJ1DAaPzhEgr4ejYZj8YXz0I8HyrTqbHb75Pz65tpXIizq5EIHpWZdnMl",
	"po4Ez1F07SSUZLlwMxGQnozCQhg+CiIG1IKdkJHd8QwtkmQUamHQsrh+//EoCGl1voR0HIxGkcssgA0v",
	"EscD8aGAJ6TQshwO/hsMkxLOo5XeQAWng+HgFIq8a7KvZbeH8sMMc/8n33M1Iw0VXCBfrTEvlszJcPhb",
	"K2YzPK+5vx7EvZHaXzvXUSkMITmx5tVO3lTGBR9KtCmhbJdk3lWxriUtoYJPJrDYiJLebSQqV0Y3b9Fp",
	"pLPE20/AbT+TLaQ0Gpq7P5T3Tar+HRUvcCOieFimj1xO3DZ0K8guiQP9nCYMCoi0gArmzD5U5dqPDtoh",
	"S2jump8BAAD//0FxNqeHBwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
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
	res := make(map[string]func() ([]byte, error))
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
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
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
