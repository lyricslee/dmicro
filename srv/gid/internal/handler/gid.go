package handler

import (
	"context"
	"sync"

	"dmicro/common/log"
	"dmicro/pkg/idgen"
	gid "dmicro/srv/gid/api"
)

type GidHandler struct{}

var (
	gidHandler     *GidHandler
	onceGidHandler sync.Once
)

func GetGidHandler() *GidHandler {
	onceGidHandler.Do(func() {
		gidHandler = &GidHandler{}
	})
	return gidHandler
}

func (this *GidHandler) GetOne(ctx context.Context, req *gid.Request, rsp *gid.Response) error {
	log.Debug("Received Gid.Create request")
	rsp.Id = idgen.GetOne()
	return nil
}

func (this *GidHandler) GetMulti(ctx context.Context, req *gid.MultiRequest, rsp *gid.MultiResponse) error {
	log.Debug("Received Gid.Create request")
	rsp.Ids = idgen.GetMulti(int(req.Count))

	return nil
}
