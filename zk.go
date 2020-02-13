package dist

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

/// @NOTICE 为简化代码，zookeeper 部署必须事先创建好对应的节点，
/// 而不通过应用去检测，自动创建，这是一种约定

/// 路径规则根据 unix 路径规则
type ZkCallback func(c *ZkClient, path string, ev zk.EventType)

// 一个到 zookeeper 的连接
type ZkClient struct {
	hosts       []string      // zookeeper 集群所有地址
	root        string        // 本应用根节点
	timeout     int           // 连接超时秒数
	conn        *zk.Conn      // 连接对象
	chQuitWatch chan struct{} // 退出 watcher
	cb          ZkCallback
	username    string
	password    string
}

// watcher 发送的通知
type Event struct {
	Err  error
	Path string
	Type zk.EventType
}

// 创建 zookeeper 连接
// timeout 秒
func NewZkClient(hosts []string, root string, user string, password string, timeout int, cb ZkCallback) (*ZkClient, error) {
	client := new(ZkClient)
	client.hosts = hosts
	client.root = root
	client.username = user
	client.password = password
	client.timeout = timeout
	client.chQuitWatch = make(chan struct{})
	client.cb = cb

	conn, _, err := zk.Connect(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	if len(client.username) > 0 {
		err = conn.AddAuth("digest", []byte(client.username+":"+client.password))
		if nil != err {
			return nil, err
		}
	}

	client.conn = conn

	return client, nil
}

func (self *ZkClient) fixName(name string) string {
	if strings.Index(name, "/") == 0 {
		return name
	}
	return self.root + "/" + name
}

// 节点是否存在
func (self *ZkClient) Exists(name string) bool {
	exists, _, err := self.conn.Exists(self.fixName(name))
	if err != nil {
		//@LOG
		return false
	}
	return exists
}

// 获得节点数据
func (self *ZkClient) GetNodeData(name string) ([]byte, error) {
	bytes, _, err := self.conn.Get(self.fixName(name))
	return bytes, err
}

func (self *ZkClient) SetNodeData(path string, data []byte) error {
	_, err := self.conn.Set(path, data, -1)
	return err
}

// 获取子节点列表
func (self *ZkClient) GetChildren(name string) ([]string, error) {
	children, _, err := self.conn.Children(self.fixName(name))
	return children, err
}

// 更新节点数据，不存在则创建
func (self *ZkClient) CreateEmpNode(name string, data []byte) error {
	path := self.fixName(name)
	_, err := self.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		//LOG
		return err
	}
	return nil
}

// 节点监控，目前通知不能详细到
func (self *ZkClient) WatchChildren(name string) error {
	go func() {
		for {
			_, _, chChild, err := self.conn.ChildrenW(self.fixName(name))
			if err != nil {
				//LOG
				//log.Warn("zk ChildrenW", zap.String("name", name), zap.Error(err))
				return
			}
			select {
			case <-self.chQuitWatch:
				//log.Info("zk ChildrenW quit", zap.String("name", name))
				return
			case event := <-chChild:
				if event.Err != nil {
					//log.Warn("zk event", zap.String("name", name), zap.Error(event.Err))
					return
				}
				self.cb(self, event.Path, event.Type)
			}
		}
	}()

	return nil
}

func (self *ZkClient) WatchNode(name string) error {
	go func() {
		for {
			_, _, ch, err := self.conn.GetW(self.fixName(name))
			if nil != err {
				//log.Warn("zk GetW", zap.String("name", name), zap.Error(err))
				return
			}
			select {
			case <-self.chQuitWatch:
				//log.Info("zk GetW quit", zap.String("name", name))
				return
			case event := <-ch:
				if event.Err != nil {
					//log.Warn("zk GetW", zap.String("name", name), zap.Error(err))
					return
				}
				self.cb(self, event.Path, event.Type)
			}
		}
	}()

	return nil
}

func (self *ZkClient) Close() {
	fmt.Println("zk closed")
	self.conn.Close()
	close(self.chQuitWatch)
}
