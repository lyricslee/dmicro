package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/micro/go-micro/v2/metadata"

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

func GetMetaDataFromContext(ctx context.Context) (*typ.MetaData, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Error("metadata.FromContext error")
		return nil, fmt.Errorf("GetMetaDataFromContext: metadata.FromContext error")
	}

	log.Debug(md)
	header := typ.MetaData{}
	header.Token = md["Token"]
	header.Appid, _ = strconv.Atoi(md["App-Id"])
	header.Uid, _ = strconv.ParseInt(md["Uid"], 10, 64)
	header.Plat, _ = strconv.Atoi(md["Plat"])
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
