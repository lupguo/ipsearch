package ipserach

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// MinResNumber, MaxResNumber 用于表示资源池中最小和最大可设置的资源数量
var (
	LimitMinResNumber = 5
	LimitMaxResNumber = 10
)

// Resource类型，特定类型的客户端资源，以及其存货了多久
type Resource struct {
	client *http.Client
	alive  time.Duration
}

// PoolCh 代表指定大小的资源池通道，可以从资源池获取资源，另外当资源被使用完后，需要回收到资源池通道
type ResCh chan *Resource

// Pool 代表资源池，其内主要是资源通道，availCh用于多个资源争用消费，recvCh用于多个资源争用回收(dispatch会来接管)
type Pool struct {
	availCh ResCh
	recvCh  ResCh
	resUsed int
	resIdle int
}

// NewPool 创建一个新的指定大小的资源池
func NewPool(size int) *Pool {
	if size < LimitMinResNumber || size > LimitMaxResNumber {
		log.Fatalln("the limit resource number for pool is error")
	}
	return &Pool{
		availCh: make(ResCh, size),
		recvCh:  make(ResCh),
		resUsed: 0,
		resIdle: LimitMaxResNumber,
	}
}

// AddResource 资源池初始化时候，添加资源到资源池
func (p *Pool) AddResource(r *Resource) {
	p.availCh <- r
}

// GetResource 从资源池中获取一个资源通道, waitTimeout设置等待资源获取的超时时间
func (p *Pool) GetResource(waitTimeout time.Duration) (r *Resource, err error) {
	tm := time.NewTimer(waitTimeout).C
	select {
	case r = <-p.availCh:
		return r, nil
	case <-tm:
		return nil, errors.New("get a resource timeout from the resource pool")
	}
}

// RecoverResource 当资源使用完成后，需要归还到资源池
func (p *Pool) RecoverResource(r *Resource)  {
	// check and refresh
	p.availCh <- r
}

