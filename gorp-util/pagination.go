package gorpUtil

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/gorp.v2"

	"github.com/juju/errors"
)

// Direction represents the sort direction
// swagger:strfmt Direction
// enum:ASC,DESC
type Direction string

// sort directions
const (
	Asc  Direction = "ASC"
	Desc Direction = "DESC"
)

const (
	DefaultPage = 1
	DefaultSize = 20
)

// NewPage creates *Page from query parameters(e.g. page=2&size=50&sort=rank,DESC&sort=name).
func NewPage(params url.Values) *Page {
	page := &Page{}
	if params == nil {
		return page
	}

	n, err := strconv.ParseInt(params.Get("page"), 10, 64)
	if err == nil {
		page.Page = int(n)
	}
	if page.Page == 0 {
		page.Page = DefaultPage
	}
	n, err = strconv.ParseInt(params.Get("size"), 10, 64)
	if err == nil {
		page.Size = int(n)
	}
	if page.Size == 0 {
		page.Size = DefaultSize
	}
	for _, s := range params["sort"] {
		page.Sort = append(page.Sort, NewOrder(s))
	}
	page.Keyword = params.Get("keyword")

	return page
}

// Page contains pagination information
type Page struct {
	// page index, start from 1
	// example:1
	// default:1
	Page int `db:"page" json:"page"`
	// size of page
	// example:10
	// default:10
	Size int `db:"size" json:"size"`
	// result sorts
	Sort []Sort `db:"sort" json:"sort"`
	// keyword to quary
	// example:arch
	Keyword string `db:"keyword" json:"keyword"`
}

func (p Page) ToParams() url.Values {
	params := make(url.Values)
	params.Add("page", strconv.Itoa(p.Page))
	params.Add("size", strconv.Itoa(p.Size))
	for _, s := range p.Sort {
		params.Add("sort", s.String())
	}
	if len(p.Keyword) > 0 {
		params.Add("keyword", p.Keyword)
	}
	return params
}

// String returns readable string represents the page
func (p Page) String() string {
	return p.ToParams().Encode()
}

// Validate implements Validatable interface
func (p Page) Validate() error {
	if p.Page < 1 {
		return errors.Errorf("page must greater than 0, got %v", p.Page)
	}
	if p.Size < 1 {
		return errors.Errorf("size must greater than 0, got %v", p.Size)
	}

	if p.Sort != nil && len(p.Sort) > 0 {
		for _, o := range p.Sort {
			if o.Name == "" {
				return errors.Errorf("sort field is blank")
			}
			if o.Direction != Asc && o.Direction != Desc {
				return errors.Errorf("invalid sort direction: %v", o.Direction)
			}
		}
	}

	return nil
}

func NewOrder(s string) Sort {
	tmp := strings.Split(s, ",")
	dir := Asc
	if len(tmp) > 1 {
		dir = Direction(strings.ToUpper(tmp[1]))
	}
	return Sort{
		Name:      tmp[0],
		Direction: dir,
	}
}

// Sort determines how a filed sorted
// swagger:strfmt
type Sort struct {
	// field name to sort
	// example:appid
	// required:true
	Name string `db:"name" json:"name"`
	// sort direction
	// example:DESC
	// default:ASC
	// enum:ASC,DESC
	Direction Direction `db:"direction" json:"direction"`
}

func (o Sort) String() string {
	return fmt.Sprintf("%s,%s", o.Name, o.Direction)
}

// PageResponse is the response to a page request
// swagger:model
type PageResponse struct {
	*Page
	// total item
	//   example:100
	Total int `db:"total" json:"total"`
	// total page
	//   example:10
	TotalPage int `db:"total_page" json:"total_page"`
	// has previous page
	//   example:true
	HasPrevPage bool `db:"has_prev_page" json:"has_prev_page"`
	// has next page
	//   example:true
	HasNextPage bool `db:"has_next_page" json:"has_next_page"`
	// payload
	Data interface{} `db:"data" json:"data"`
}

// NewPageResponse creates a PageResponse
func NewPageResponse(page *Page, total int, data interface{}) *PageResponse {
	totalPage := total / page.Size
	if total%page.Size > 0 {
		totalPage++
	}
	return &PageResponse{
		Page:        page,
		Total:       total,
		TotalPage:   totalPage,
		HasPrevPage: page.Page > 1,
		HasNextPage: page.Page < totalPage,
		Data:        data,
	}
}

func LoadPage(tx gorp.SqlExecutor, q *Query, page *Page, holder interface{}, countQueries ...*Query) (int, error) {
	cq := q
	if len(countQueries) > 0 {
		if v := countQueries[0]; v != nil {
			cq = v
		}
	}
	n, err := cq.Count(tx)
	if err != nil {
		return -1, errors.Trace(err)
	}
	if page != nil && page.Validate() == nil {
		for _, s := range page.Sort {
			q.OrderByString(s.Name, string(s.Direction))
		}
		q.Pagination(page.Page, page.Size)
	}
	if _, err = q.FetchAll(tx, holder); err != nil {
		return -1, errors.Trace(err)
	}
	return int(n), nil
}
