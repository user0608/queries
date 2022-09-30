package queries

import (
	"context"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type QrueryParmas struct {
	Preloads    []string                              `query:"preload"`
	preloadFuns map[string]func(tx *gorm.DB) *gorm.DB `query:"-"`
	Limit       int                                   `query:"limit"`
	Offset      int                                   `query:"offset"`
}
type contextkey string

type Preloader interface {
	/*
		Preload retorna los campos al los que se puede invocar un preload,
		estos campos tienen que estar en formato snakecase, el mismo nombre
		de campo json, ademas tiene que ser el mismo nombre de campo de modelo
		en formato camelcase, de tener nombres distintos se puede pasar el
		nombre del atributo de clase a la que hace referencia el campo en
		json separado por una coma: `"nombre_json,nombre_modelo"`
	*/
	Preload() []string
}

var query_param_key = contextkey("query-param-key")

func Wrapp(ctx context.Context, tx *gorm.DB) *gorm.DB {
	value := ctx.Value(query_param_key)
	if value == nil {
		return tx
	}
	qp, ok := value.(QrueryParmas)
	if !ok {
		return tx
	}
	if qp.Limit > 0 {
		tx = tx.Limit(qp.Limit)
	}
	tx = tx.Offset(qp.Offset)
	for _, field := range qp.Preloads {
		fn, ok := qp.preloadFuns[field]
		if ok {
			tx = tx.Preload(field, fn)
			continue
		}
		tx = tx.Preload(field)
	}
	return tx
}

func QueryParamMiddl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var qparams QrueryParmas
		(&echo.DefaultBinder{}).BindQueryParams(c, &qparams)
		ctx := context.WithValue(c.Request().Context(), query_param_key, qparams)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func PreloadFunc(ctx context.Context, querypreload string, fn func(tx *gorm.DB) *gorm.DB) context.Context {
	value := ctx.Value(query_param_key)
	if value == nil {
		return ctx
	}
	qp, ok := value.(QrueryParmas)
	if !ok {
		return ctx
	}
	if qp.preloadFuns == nil {
		qp.preloadFuns = make(map[string]func(tx *gorm.DB) *gorm.DB)
	}
	qp.preloadFuns[querypreload] = fn
	return context.WithValue(ctx, query_param_key, qp)
}

/*
Recibe un objeto que implemente la interfaz:

	type Preloader interface { Preload() []string }
*/
func Model(ctx context.Context, p Preloader) context.Context {
	value := ctx.Value(query_param_key)
	if value == nil {
		return ctx
	}
	qp, ok := value.(QrueryParmas)
	if !ok {
		return ctx
	}
	if len(qp.Preloads) == 0 {
		return ctx
	}
	var aux []string
	preloads := strings.Split(qp.Preloads[0], ",")
	for _, s := range preloads {
		r := standarize(p.Preload(), s)
		if r != "" {
			aux = append(aux, r)
		}
	}
	qp.Preloads = aux
	return context.WithValue(ctx, query_param_key, qp)
}

func standarize(items []string, s string) string {
	for _, item := range items {
		jsonname, modelname := fieldName(item)
		if jsonname == s {
			return modelname
		}
	}
	return ""
}

// fieldName out json part, in model part
func fieldName(s string) (jsonname string, modelname string) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ",")
	if len(parts) == 1 {
		return s, snakeCasetoCamelCase(s)
	}
	return parts[0], parts[1]
}

func snakeCasetoCamelCase(cadena string) string {
	cadena = strings.ReplaceAll(cadena, ".", "_._")
	rgx := regexp.MustCompile(`([A-Za-z0-9\.]+)`)
	var sb strings.Builder
	matchs := rgx.FindAllString(cadena, 100)
	for _, s := range matchs {
		if len(s) == 0 {
			continue
		}
		sb.WriteString(strings.ToUpper(s[0:1]))
		if len(s) == 1 {
			continue
		}
		sb.WriteString(s[1:])
	}
	return sb.String()
}
