package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/micro/go-micro/metadata"

	"dmicro/common/errors"
	"dmicro/common/log"
	"dmicro/common/typ"
)

func WriteError(w http.ResponseWriter, err error) {
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", "application/json")
	}

	e := errors.Parse(err.Error())
	response := map[string]interface{}{
		"t": time.Now().UnixNano(),
	}
	response["errno"] = e.Errno
	response["errmsg"] = e.Errmsg

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	if e.Errno == -1 {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(499)
	}

	w.Write(b)
}

func GetHeaderFromContext(ctx context.Context) (*typ.Header, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Error("metadata.FromContext error")
		return nil, fmt.Errorf("GetHeaderFromContext: metadata.FromContext error")
	}

	header := typ.Header{}
	header.Token = md["Token"]
	header.Uid, _ = strconv.ParseInt(md["Uid"], 10, 64)
	header.AppId, _ = strconv.Atoi(md["App-Id"])
	//header.AppVersion = md["App-Version"]
	//header.OsType = md["Os-Type"]
	//header.OsVersion = md["Os-Version"]
	//header.Resolution = md["Resolution"]
	//header.Model = md["Model"]
	//header.Channel = md["Channel"]
	//header.Net = md["Net"]
	//header.DeviceId = md["Device-Id"]

	return &header, nil
}
