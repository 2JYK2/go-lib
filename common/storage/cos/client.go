package cos

import (
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	cosClients = make(map[string]*cos.Client)
	mu         sync.Mutex
)

// GetClient 获取或创建 cos Client（按 bucket 缓存）
func (c *COSTokenManager) InitClient() error {
	// TODO
	return nil
}
