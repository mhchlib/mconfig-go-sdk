package mconfig

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"github.com/mhchlib/mconfig-go-sdk/adapter"
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/reg"
	"google.golang.org/grpc"
	"time"
)

func initAddressProvider(m *Mconfig) func(serviceName string) (*reg.ServiceVal, error) {
	log := m.opts.Logger
	if m.opts.EnableRegistry {
		regClient, err := register.InitRegister(func(options *reg.Options) {
			options.RegisterType = reg.RegistryType(RegisterType_Etcd)
			options.NameSpace = m.opts.NameSpace
			options.Address = m.opts.RegistryUrl
		})
		if err != nil {
			log.Fatal("register fail")
		}
		return func(serviceName string) (*reg.ServiceVal, error) {
			return regClient.GetService(serviceName)
		}
	}
	if m.opts.DirectLinkAddress == "" {
		log.Fatal("you should provider a direct link address or an register center address...")
	}
	return func(serviceName string) (*reg.ServiceVal, error) {
		return &reg.ServiceVal{
			Address: m.opts.DirectLinkAddress,
		}, nil
	}
}

func (m *Mconfig) initMconfigLink() {
	log := m.opts.Logger
	addressProvider := initAddressProvider(m)
	request := &server.WatchConfigStreamRequest{
		AppKey:     m.opts.AppKey,
		ConfigKeys: m.opts.ConfigKeys,
		Metadata:   m.opts.Metadata,
	}

	//添加连接断开重试机制
	retryTime := m.opts.RetryTime
	once := true
	enableRetry := false
	started := make(chan interface{})
	go func(m *Mconfig, started chan interface{}) {
		for {
			if enableRetry {
				log.Info("mconfig retry fail... it does not work now.... and will retry after ", retryTime)
				<-time.After(retryTime)
			}
			enableRetry = true
			service, err := addressProvider("mconfig-server")
			if err != nil {
				log.Info(err)
				continue
			}
			withTimeout, _ := context.WithTimeout(context.Background(), time.Second*3)
			dial, err := grpc.DialContext(withTimeout, service.Address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Info(err, " addr: ", service)
				continue
			}
			mConfigService := server.NewMConfigClient(dial)
			stream, err := mConfigService.WatchConfigStream(context.Background())
			if err != nil {
				log.Info(err)
				continue
			}
			err = stream.SendMsg(request)
			if err != nil {
				log.Info(err)
				continue
			}
			for {
				recv, err := stream.Recv()
				log.Info(recv)
				if err != nil {
					log.Info(err)
					break
				}
				configs := recv.Configs
				data := m.opts.ConfigsData
				data.Lock()
				for _, config := range configs {
					data.Data[config.ConfigKey] = config.Val
					adapter.ExecuteAdapterCallBack(config.ConfigKey, config.Val)
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
				log.Info(" try to refresh mconfig cache...")
			}
		}
	}(m, started)
	<-started
	close(started)
}
