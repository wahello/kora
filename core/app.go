package core

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/wecanooo/kora/core/context"
	"io/ioutil"
	"strings"
)

type Application struct {
	*echo.Echo
}

// NewApplication server instance 생성
func NewApplication(echo *echo.Echo) {
	application = &Application{Echo: echo}
}

// RunServer server 시작
func (app *Application) RunServer() {
	app.Logger.Fatal(app.Start(GetConfig().DefaultString("APP.ADDR", ":3000")))
}

// RoutePath 등록된 Route 정보
func (app *Application) RoutePath(name string, params ...interface{}) string {
	return app.Reverse(name, params...)
}

// PrintRoutes 등록되어 있는 Route 정보를 출력
func (app *Application) PrintRoutes(filename string) {
	routes := make([]*echo.Route, 0)
	for _, item := range app.Routes() {
		if strings.HasPrefix(item.Name, "github.com") || strings.HasSuffix(item.Name, "notFoundHandler") {
			continue
		}

		routes = append(routes, item)
	}

	routesStr, _ := json.MarshalIndent(struct {
		Count  int           `json:"count"`
		Routes []*echo.Route `json:"routes"`
	}{
		Count:  len(routes),
		Routes: routes,
	}, "", "  ")

	ioutil.WriteFile(filename, routesStr, 0644)
}

// RegisterRoutes Routes 등록
func (app *Application) RegisterRoutes(register func(*Application)) {
	app.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.AppContext{Context: c}
			return hf(cc)
		}
	})

	register(app)
}

// RegisterHandler 핸들러 등록
func (app *Application) RegisterHandler(fn context.AppRegisterFunc, path string, h context.AppHandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return fn(path, func(c echo.Context) error {
		cc, ok := c.(*context.AppContext)
		if !ok {
			cc = context.NewAppContext(c)
			return h(cc)
		}
		return h(cc)
	}, m...)
}
