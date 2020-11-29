package client

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/mregister"
	"google.golang.org/grpc"
	"log"
	"time"
)

func (m *Mconfig) initMconfigLink() {
	reg, err := register.InitRegister(string(RegisterType_Etcd), func(options *mregister.Options) {
		options.NameSpace = m.opts.NameSpace
		options.Address = m.opts.RegistryUrl
	})
	if err != nil {
		log.Fatal("[mconfig] ", "register fail")
	}
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
	enableRetry := false
	started := make(chan interface{})
	go func(m *Mconfig, started chan interface{}) {
		for {
			if enableRetry {
				<-time.After(retryTime)
				log.Println("[mconfig] ", "mconfig retry fail... it does not work now.... and will retry after ", retryTime)
			}
			enableRetry = true
			service, err := reg.GetService("mconfig-sdk")
			if err != nil {
				log.Println("[mconfig] ", err)
				continue
			}
			withTimeout, _ := context.WithTimeout(context.Background(), time.Second*3)
			dial, err := grpc.DialContext(withTimeout, service, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Println("[mconfig] ", err, " addr: ", service)
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
				log.Println("refresh mconfig cache...")
			}
		}
	}(m, started)
	<-started
	close(started)
}
