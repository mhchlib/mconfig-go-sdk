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
	mConfigService := sdk.NewMConfigService(m.opts.NameSpace, mService.Client())
	request := &sdk.GetVRequest{
		AppId: m.opts.AppKey,
		Filters: &sdk.ConfigFilters{
			ConfigIds: m.opts.ConfigKeys,
			ExtraData: m.opts.ABFilters,
		},
	}

	//添加连接断开重试机制
	retryNum := m.opts.RetryNum
	once := true
	started := make(chan interface{})
	go func(m *Mconfig, started chan interface{}) {
		for retryNum >= 0 {
			if retryNum != Default_Retry_Num {
				log.Println("[mconfig] ", "mconfig retry link to mconfig server ... ")
			}
			stream, err := mConfigService.GetVStream(context.Background(), request)
			if err != nil {
				log.Println("[mconfig] ", err)
				retryNum = retryNum - 1
				continue
			}
			for {
				recv, err := stream.Recv()
				if err != nil {
					log.Println("[mconfig] ", err)
					retryNum = retryNum - 1
					break
				}
				retryNum = m.opts.RetryNum
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
		}
		log.Println("[mconfig] ", "mconfig retry fail... it does not work now....")
	}(m, started)
	<-started
	close(started)
}
