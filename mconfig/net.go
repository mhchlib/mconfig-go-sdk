package mconfig

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	etcdV "github.com/micro/go-micro/v2/registry/etcd"
	"log"
)

func (m *Mconfig) initMconfigLink() {
	var reg registry.Registry
	if m.opts.RegistryType == RegisterType_Etcd {
		reg = etcdV.NewRegistry(func(options *registry.Options) {
			options.Addrs = []string{m.opts.RegistryUrl} //地址
		})
	}
	mService := micro.NewService(
		micro.Registry(reg),
	)
	mService.Init()
	mConfigService := sdk.NewMConfigService("", mService.Client())
	request := &sdk.GetVRequest{
		AppId: m.opts.AppKey,
		Filters: &sdk.ConfigFilters{
			ConfigIds: m.opts.ConfigKeys,
			ExtraData: m.opts.ABFilters,
		},
	}
	stream, err := mConfigService.GetVStream(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	started := make(chan interface{})
	go func(stream sdk.MConfig_GetVStreamService, m *Mconfig, started chan interface{}) {
		once := true
		for {
			recv, err := stream.Recv()
			if err != nil {
				log.Println(err)
			}
			configs := recv.Configs
			data := m.opts.ConfigsData
			data.Lock()
			for _, config := range configs {
				data.Data[config.Id] = config.Config
			}
			data.Unlock()
			if once {
				started <- &struct{}{}
				once = false
			}
			//refer the cache
			//todo: 这里可以优化为主动去更新cache中内容，但是这个需要加大cache大小，带上类型
			// support soon....
			m.opts.Cache.Lock()
			m.opts.Cache.Cache = map[string]*FieldInterface{}
			m.opts.Cache.Unlock()
		}
	}(stream, m, started)
	<-started
	close(started)
}
