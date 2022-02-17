package common

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

// ParamOrder ...
type ParamOrder struct {
	Field string
	By    string
}

// PaginateQuery ...
type PaginateQuery struct {
	Order  ParamOrder
	LastID int
	Limit  int
}

// PaginateLink ...
type PaginateLink struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

func PagingOperatorConditionSQL(pg PaginateQuery) string {
	operator := ">"
	if pg.Order.By == ASC {
		operator = ">"
	}
	if pg.Order.By == DESC && pg.LastID == 0 {
		operator = ">"
	}
	if pg.Order.By == DESC && pg.LastID != 0 {
		operator = "<"
	}

	return operator
}

// OrderByOffsetLimitSQL generate order by and add offset limit query
func OrderByLastIDLimitSQL(pg PaginateQuery) string {
	return fmt.Sprintf("order by %s %s limit %d;", pg.Order.Field, pg.Order.By, pg.Limit)
}

// Offset

// PaginateQueryOffset ...
type PaginateQueryOffset struct {
	Order  ParamOrder `json:"order"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
}

// OrderByOffsetLimitSQL generate order by and add offset limit query
func OrderByOffsetLimitSQL(pg PaginateQueryOffset) string {
	return fmt.Sprintf("order by %s %s limit %d offset %d;", pg.Order.Field, strings.ToUpper(pg.Order.By), pg.Limit, pg.Offset)
}

// OrderByOffsetLimitPaginateLink generate pagination link
func OrderByOffsetLimitPaginateLink(pg PaginateQueryOffset, prefix string, additional map[string]string) PaginateLink {
	var (
		prevOffset, nextOffset          int
		addParamStr, prevLink, nextLink string
	)
	// /v1/admin-users?order=id,asc|desc&offset=20&limit=50
	if additional != nil {
		addParams := url.Values{}
		for k, v := range additional {
			addParams.Add(k, v)
		}
		addParamStr = addParams.Encode()
		addParamStr = addParamStr + "&"

	}

	param := prefix + "%sorder=%s,%s&offset=%d&limit=%d"

	prevOffset = 0
	if pg.Offset > 0 {
		prevOffset = pg.Offset - pg.Limit
		if prevOffset < 0 {
			prevOffset = 0
		}

		prevLink = fmt.Sprintf(param, addParamStr, pg.Order.Field, pg.Order.By, prevOffset, pg.Limit)
	}

	// TODO: need to add a logic to handle latest link of the page
	nextOffset = pg.Offset + pg.Limit
	nextLink = fmt.Sprintf(param, addParamStr, pg.Order.Field, pg.Order.By, nextOffset, pg.Limit)

	return PaginateLink{
		Prev: prevLink,
		Next: nextLink,
	}
}

func GetParamOrder(c *fiber.Ctx) (ParamOrder, error) {
	param := c.Query("order")

	// sanitize param
	param = template.HTMLEscapeString(param)

	params := strings.Split(param, ",")
	po := ParamOrder{}
	if len(params) != 2 {
		return po, errors.New("wrong order parameters")
	}

	po.Field = params[0]
	po.By = params[1]

	return po, nil
}

// GetIntParam Parse the url param to get value as integer.
// for example, we need to get limit and offset param
func GetIntParam(c *fiber.Ctx, name string) (int, error) {
	param := c.Query(name)
	if len(param) == 0 {
		return 0, nil
	}

	// sanitize param
	param = template.HTMLEscapeString(param)

	return strconv.Atoi(param)
}

// GetStringParam Parse the url param to get value as string.
func GetStringParam(c *fiber.Ctx, name string) (string, error) {
	param := c.Query(name)

	if len(strings.TrimSpace(param)) == 0 {
		return "", nil
	}

	// sanitize param
	param = template.HTMLEscapeString(param)

	return param, nil
}

// GetBoolParam Parse the url param to get value as boolean
func GetBoolParam(c *fiber.Ctx, name string) (bool, error) {
	param := c.Query(name)
	if len(param) == 0 {
		return false, nil
	}

	// sanitize param
	param = template.HTMLEscapeString(param)

	result, err := strconv.ParseBool(param)
	if err != nil {
		return false, err
	}

	return result, nil
}

func GetPaginateQueryOffset(c *fiber.Ctx) (PaginateQueryOffset, error) {
	var pg PaginateQueryOffset

	// Define pagination
	order, err := GetParamOrder(c)
	if err != nil {
		return pg, err
	}
	pg.Order = order

	offset, err := GetIntParam(c, "offset")
	if err != nil {
		return pg, err
	}
	pg.Offset = offset

	limit, err := GetIntParam(c, "limit")
	if err != nil {
		return pg, err
	}
	pg.Limit = limit

	return pg, err
}
