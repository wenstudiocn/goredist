package dist

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)
/// 进程通用设置解析
/// 规定: 配置使用 protos/conf.proto 文件，配置文件出现该字段，则自动启用该功能

// 一个程序的基本信息
type AppInfo struct {
	Name string				// 程序名，和文件名无关，编译时指定
	InstName string		// 进程实例名称
	BootAt time.Time	// 启动时间
	Ver *Version				// 版本号
	ConfFile string			// 配置文件名
	Log *SLogger			// 日志对象
	Rdb *redis.Client		// redis 客户端
	Ev *Eventbus			// 事件总线
	Mq *Mq						// 消息队列
	Zk *ZkClient				// zookeeper
	Debug bool
	Mode string
	Offline int32 				// 是否下线 0 不下线
	ChQuit chan struct{} 	// 进程退出

	enableLogger bool
	logLevel int32
	logPath string
	logStdout bool

	enableRdb bool
	enableEv bool
	enableMq bool
	enableZk bool
}

type AppInfoOptions func(info *AppInfo)

func AppEnableLogger(enable bool, logLevel int32, logPath string, logStdout bool) AppInfoOptions {
	return func(ai *AppInfo) {
		ai.enableLogger = enable
		ai.logLevel = logLevel
		ai.logPath = logPath
		ai.logStdout = logStdout
	}
}

func NewAppInfo(programeName string, major, minor, revision int, id uint64, confFile string, mode string) *AppInfo {
	ai := &AppInfo{}
	ai.Name = programeName
	ai.Ver = NewVersion(major, minor, revision)
	ai.InstName = ai.Name + "-" + strconv.FormatUint(id, 10)
	ai.BootAt = time.Now()
	ai.ConfFile = confFile
	ai.Mode = mode
	ai.Debug = true
	ai.ChQuit = make(chan struct{})

	if strings.ToLower(ai.Mode) == "release" {
		ai.Debug = false
	}
	return ai
}

func (self *AppInfo)EnableLog(console bool, level int32, logPath string) error {
	// level
	lv, ok := LogLevelMap[int(level)]
	if !ok {
		return ErrInvalidParams
	}

	// path
	f := self.InstName + ".log"
	if len(logPath) > 0{
		if _, err := os.Stat(logPath); err != nil {
			return ErrInvalidParams
		}
		f = path.Join(logPath, f)
	}
	// logger
	self.Log = NewSLogger(console, f, lv)
	SetDefaultLogger(self.Log)

	return nil
}

// add a node to zookeeper
func (self *AppInfo) RegisterSelf(parentPath string) error {
	node := parentPath + "/" + self.InstName
	fh, err := os.Open(self.ConfFile)
	if err != nil {
		return err
	}
	data ,err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}
	return self.Zk.CreateEmpNode(node, data)
}