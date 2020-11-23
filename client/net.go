package pkg

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/register"
	etcd_kit "github.com/mhchlib/register/etcd-kit"
	"google.golang.org/grpc"
	"time"

	//etcdV "github.com/micro/go-micro/v2/registry/etcd_custom"
	"log"
)

func (m *Mconfig) initMconfigLink() {
	var reg register.Register
	reg = &etcd_kit.EtcdRegister{}
	reg.Init(func(options *register.Options) {
		options.NameSpace = m.opts.NameSpace
		options.Address = m.opts.RegistryUrl
	})
	request := &sdk.GetVRequest{
		AppId: m.opts.AppKey,
		Filters: &sdk.ConfigFilters{
			ConfigIds: m.opts.ConfigKeys,
			ExtraData: m.opts.ABFilters,
		},
	}

	//添加连接断开重试机制
	retryTime := m.opts.RetryTime
	once := true
	started := make(chan interface{})
	go func(m *Mconfig, started chan interface{}) {
		for {
			if !once {
				log.Println("[mconfig] ", "mconfig retry fail... it does not work now.... and will retry after ", retryTime)
			}
			<-time.After(retryTime)
			service, err := reg.GetService("mconfig-sdk")
			if err != nil {
				log.Println("[mconfig] ", err)
				continue
			}
			dial, err := grpc.Dial(service, grpc.WithInsecure())
			if err != nil {
				log.Println("[mconfig] ", err)
				continue
			}
			mConfigService := sdk.NewMConfigClient(dial)
			stream, err := mConfigService.GetVStream(context.Background())
			if err != nil {
				log.Println("[mconfig] ", err)
				continue
			}
			err = stream.SendMsg(request)
			if err != nil {
				log.Println("[mconfig] ", err)
				continue
			}
			for {
				recv, err := stream.Recv()
				if err != nil {
					log.Println("[mconfig] ", err)
					break
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
		}
	}(m, started)
	<-started
	close(started)
}
