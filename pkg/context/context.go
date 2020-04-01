package context

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"dmicro/pkg/util/convert"
)

func New(ctx *gin.Context) *DmContext {
	return &DmContext{gctx: ctx}
}

type DmContext struct {
	gctx *gin.Context
}

func (c *DmContext) Context() *gin.Context {
	return c.gctx
}

func (c *DmContext) Param(key string) string {
	return c.gctx.Param(key)
}

func (c *DmContext) Query(key string) string {
	return c.gctx.Query(key)
}

// 请求参数 Json 转换成 obj 对象
func (c *DmContext) ParseJSON(obj interface{}) error {
	if err := c.gctx.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func (c *DmContext) Response(obj interface{}) {
	//if obj == nil {
	//	obj = gin.H{}
	//}
	//
	//m := make(map[string]interface{})
	//m["data"] = obj
	//m["t"] = time.Now().UnixNano()
	//c.response(http.StatusOK, m)
	c.response(http.StatusOK, obj)
}

func (c *DmContext) ResponseError(err error) {
	c.response(499, err)
}

func (c *DmContext) response(status int, obj interface{}) {
	c.gctx.JSON(status, obj)
	c.gctx.Abort()
}

var (
	DefaultPageSize int64 = 10
	MaxPageSize     int64 = 100
)

func (c *DmContext) GetUidStr() string {
	return c.gctx.GetHeader("Uid")
}

func (c *DmContext) GetUid() int64 {
	uidstr := c.gctx.GetHeader("Uid")
	if uidstr == "" {
		return 0
	}
	uid, _ := convert.ConvertInt(uidstr)
	return uid
}

func (c *DmContext) GetToken() string {
	return c.gctx.GetHeader("Token")
}

func (c *DmContext) GetPage() int64 {
	if v := c.Query("page"); len(v) > 0 {
		if i, _ := convert.ConvertInt(v); i > 0 {
			return i
		}
	}

	return 1
}

func (c *DmContext) GetPageSize() int64 {
	if v := c.Query("per_page"); len(v) > 0 {
		if i, _ := convert.ConvertInt(v); i > 0 {
			if i > MaxPageSize {
				i = MaxPageSize
			}
			return i
		}
	}

	return DefaultPageSize
}

type list struct {
	List       interface{} `json:"list,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

type pagination struct {
	Total   int64 `json:"total,omitempty"`
	Page    int64 `json:"page,omitempty"`
	PerPage int64 `json:"per_page,omitempty"`
}

func (c *DmContext) ResponseList(obj interface{}) {
	c.Response(list{List: obj})
}

func (c *DmContext) ResponsePage(total int64, obj interface{}) {
	c.Response(list{
		List: obj,
		Pagination: &pagination{
			Total:   total,
			Page:    c.GetPage(),
			PerPage: c.GetPageSize(),
		},
	})
}
