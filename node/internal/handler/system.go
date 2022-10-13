package handler

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
)

//获取node节点的信息
func WatchSystem(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, nodeUUID), clientv3.WithPrefix())
}
