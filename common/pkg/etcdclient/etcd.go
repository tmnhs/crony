package etcdclient

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/logger"
	"strings"
	"time"
)

var _defalutEtcd *Client

type Client struct {
	*clientv3.Client
	reqTimeout time.Duration
}

func Init(e models.Etcd) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: time.Duration(e.DialTimeout) * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return nil, err
	}
	_defalutEtcd = &Client{
		Client:     cli,
		reqTimeout: time.Duration(e.ReqTimeout) * time.Second,
	}
	return _defalutEtcd, nil
}

func GetEtcdClient() *Client {
	if _defalutEtcd == nil {
		logger.Errorf("mysql database is not initialized")
		return nil
	}
	return _defalutEtcd
}

func Put(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext()
	defer cancel()
	return _defalutEtcd.Put(ctx, key, val, opts...)
}

func PutWithModRev(key, val string, rev int64) (*clientv3.PutResponse, error) {
	if rev == 0 {
		return Put(key, val)
	}

	ctx, cancel := NewEtcdTimeoutContext()
	tresp, err := _defalutEtcd.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", rev)).
		Then(clientv3.OpPut(key, val)).
		Commit()
	cancel()
	if err != nil {
		return nil, err
	}

	if !tresp.Succeeded {
		return nil, ErrValueMayChanged
	}

	resp := clientv3.PutResponse(*tresp.Responses[0].GetResponsePut())
	return &resp, nil
}

func Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext()
	defer cancel()
	return _defalutEtcd.Get(ctx, key, opts...)
}

func Delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext()
	defer cancel()
	return _defalutEtcd.Delete(ctx, key, opts...)
}

func Watch(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return _defalutEtcd.Watch(context.Background(), key, opts...)
}

func Grant(ttl int64) (*clientv3.LeaseGrantResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext()
	defer cancel()
	return _defalutEtcd.Grant(ctx, ttl)
}

func Revoke(id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _defalutEtcd.reqTimeout)
	defer cancel()
	return _defalutEtcd.Revoke(ctx, id)
}

func KeepAliveOnce(id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext()
	defer cancel()
	return _defalutEtcd.KeepAliveOnce(ctx, id)
}

func GetLock(key string, id clientv3.LeaseID) (bool, error) {
	key = fmt.Sprintf(KeyEtcdLock, key)
	ctx, cancel := NewEtcdTimeoutContext()
	resp, err := _defalutEtcd.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(id))).
		Commit()
	cancel()

	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

func DelLock(key string) error {
	_, err := Delete(fmt.Sprintf(KeyEtcdLock, key))
	return err
}

func IsValidAsKeyPath(s string) bool {
	return strings.IndexAny(s, "/\\") == -1
}

// etcdTimeoutContext return better error info
type etcdTimeoutContext struct {
	context.Context
	etcdEndpoints []string
}

func (c *etcdTimeoutContext) Err() error {
	err := c.Context.Err()
	if err == context.DeadlineExceeded {
		err = fmt.Errorf("%s: etcd(%v) maybe lost",
			err, c.etcdEndpoints)
	}
	return err
}

// NewEtcdTimeoutContext return a new etcdTimeoutContext
func NewEtcdTimeoutContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), _defalutEtcd.reqTimeout)
	etcdCtx := &etcdTimeoutContext{}
	etcdCtx.Context = ctx
	etcdCtx.etcdEndpoints = config.GetConfigModels().Etcd.Endpoints
	return etcdCtx, cancel
}
