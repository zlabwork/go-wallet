package btc

import (
    "github.com/zlabwork/go-zlibs"
)

type HandleConfigs struct {
    Host string // http://127.0.0.1:8332
    User string
    Pass string
}

type serviceHandle struct {
    cfg *HandleConfigs
    req *zlibs.HttpLib
}

func NewServiceHandle(c *HandleConfigs) *serviceHandle {
    return &serviceHandle{
        cfg: c,
        req: zlibs.NewHttpLib(),
    }
}
