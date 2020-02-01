// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: srv/gid/api/gid.proto

package go_micro_srv_gid

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Gid service

type GidService interface {
	GetOne(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	GetMulti(ctx context.Context, in *MultiRequest, opts ...client.CallOption) (*MultiResponse, error)
}

type gidService struct {
	c    client.Client
	name string
}

func NewGidService(name string, c client.Client) GidService {
	return &gidService{
		c:    c,
		name: name,
	}
}

func (c *gidService) GetOne(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Gid.GetOne", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gidService) GetMulti(ctx context.Context, in *MultiRequest, opts ...client.CallOption) (*MultiResponse, error) {
	req := c.c.NewRequest(c.name, "Gid.GetMulti", in)
	out := new(MultiResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Gid service

type GidHandler interface {
	GetOne(context.Context, *Request, *Response) error
	GetMulti(context.Context, *MultiRequest, *MultiResponse) error
}

func RegisterGidHandler(s server.Server, hdlr GidHandler, opts ...server.HandlerOption) error {
	type gid interface {
		GetOne(ctx context.Context, in *Request, out *Response) error
		GetMulti(ctx context.Context, in *MultiRequest, out *MultiResponse) error
	}
	type Gid struct {
		gid
	}
	h := &gidHandler{hdlr}
	return s.Handle(s.NewHandler(&Gid{h}, opts...))
}

type gidHandler struct {
	GidHandler
}

func (h *gidHandler) GetOne(ctx context.Context, in *Request, out *Response) error {
	return h.GidHandler.GetOne(ctx, in, out)
}

func (h *gidHandler) GetMulti(ctx context.Context, in *MultiRequest, out *MultiResponse) error {
	return h.GidHandler.GetMulti(ctx, in, out)
}