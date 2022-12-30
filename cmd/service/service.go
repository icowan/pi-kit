/**
 * @Time : 2022/12/29 9:34 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package service

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kit/log"
	"github.com/icowan/pi-kit/src/api"
	"github.com/icowan/pi-kit/src/logging"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type PiConnect string

const (
	PiConnectLocal  PiConnect = "Local"
	PiConnectRemote PiConnect = "Remote"
)

func (c PiConnect) String() string {
	return strings.ToLower(string(c))
}

const (
	DefaultCmdCameraJpegBin = "/usr/bin/libcamera-jpeg" // LibCameraJpeg 命令
	DefaultCameraJpegOutput = "/tmp/"                   // 照片存储路径

	DefaultPiConnect  = "Local"     // Local: 本地; Remote: 远程
	DefaultPiHost     = "127.0.0.1" // 树莓派地址
	DefaultPiSSHPort  = 22          // 树莓派SSH端口
	DefaultPiUser     = "pi"        // 树莓派SSH用户
	DefaultPiPassword = ""          // 树莓派SSH密码

	DefaultServerLogLevel = "all"
	DefaultServerLogPath  = ""
	DefaultServerLogName  = "pi-kit.log"

	EnvCmdCameraJpegBin = "ENV_CMD_CAMERA_JPEG_BIN"
	EnvCameraJpegOutput = "ENV_CAMERA_JPEG_OUTPUT"

	EnvPiConnect  = "ENV_PI_CONNECT"
	EnvPiHost     = "ENV_PI_HOST"
	EnvPiSSHPort  = "ENV_PI_SSH_PORT"
	EnvPiUser     = "ENV_PI_USER"
	EnvPiPassword = "ENV_PI_PASSWORD"

	EnvServerLogLevel = "ENV_SERVER_LOG_LEVEL"
	EnvServerLogPath  = "ENV_SERVER_LOG_PATH"
	EnvServerLogName  = "ENV_SERVER_LOG_NAME"
)

var (
	version = ""
	logger  log.Logger
	gormDB  *gorm.DB
	db      *sql.DB
	err     error

	apiSvc api.Service

	goOS                            = runtime.GOOS
	goArch                          = runtime.GOARCH
	goVersion                       = runtime.Version()
	compiler                        = runtime.Compiler
	buildDate, gitCommit, gitBranch string

	serverLogPath, serverLogName, serverLogLevel string

	cmdCameraJpegBin, cameraJpegOutput    string
	piConnect, piHost, piUser, piPassword string
	piSSHPort                             int

	rootCmd = &cobra.Command{
		Use:               "pi-kit",
		Short:             "",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# pi-kit工具集
可用的配置类型：
[start, camera]
有关本系统的相关概述，请参阅 https://github.com/icowan/pi-kit

本工具支持环境变量、cmd入参两种方案传入配置
cmd args > env > default
如果cmd args取不到值会读取环境变量，如果环境变量没有值则取默认值
`,
		Version: version,
	}
)

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
GitCommit: ` + gitCommit + `
GitBranch: ` + gitBranch + `
BuildDate: ` + buildDate + `
Compiler: ` + compiler + `
GoVersion: ` + goVersion + `
Platform: ` + goOS + "/" + goArch + `
`)

	//startCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的http端口")
	//startCmd.PersistentFlags().BoolVar(&webEmbed, "web.embed", true, "是否使用embed.FS")
	//startCmd.PersistentFlags().BoolVar(&enableCORS, "cors.enable", DefaultEnableCORS, "是否开启跨域访问")
	//startCmd.PersistentFlags().StringVar(&corsAllowOrigins, "cors.allow.origins", DefaultCORSAllowOrigins, "允许跨域访问的域名")
	//startCmd.PersistentFlags().StringVar(&corsAllowMethods, "cors.allow.methods", DefaultCORSAllowMethods, "允许跨域访问的方法")
	//startCmd.PersistentFlags().StringVar(&corsAllowHeaders, "cors.allow.headers", DefaultCORSAllowHeaders, "允许跨域访问的头部")
	//startCmd.PersistentFlags().StringVar(&corsExposeHeaders, "cors.expose.headers", DefaultCORSExposeHeaders, "允许跨域访问的头部")
	//startCmd.PersistentFlags().BoolVar(&corsAllowCredentials, "cors.allow.credentials", DefaultCORSAllowCredentials, "是否允许跨域访问的凭证")
	//startCmd.PersistentFlags().BoolVar(&tracerEnable, "tracer.enable", DefaultJaegerEnable, "是否启用Tracer")
	//startCmd.PersistentFlags().StringVar(&tracerDrive, "tracer.drive", DefaultJaegerDrive, "Tracer驱动")
	//startCmd.PersistentFlags().StringVar(&tracerName, "tracer.name", DefaultJaegerName, "Tracer名称")
	//startCmd.PersistentFlags().StringVar(&tracerJaegerHost, "tracer.jaeger.host", DefaultJaegerHost, "Tracer Jaeger Host")
	//startCmd.PersistentFlags().Float64Var(&tracerJaegerParam, "tracer.jaeger.param", DefaultJaegerParam, "Tracer Jaeger Param")
	//startCmd.PersistentFlags().StringVar(&tracerJaegerType, "tracer.jaeger.type", DefaultJaegerType, "采样器的类型 const: 固定采样, probabilistic: 随机取样, ratelimiting: 速度限制取样, remote: 基于Jaeger代理的取样")
	//startCmd.PersistentFlags().BoolVar(&tracerJaegerLogSpans, "tracer.jaeger.log.spans", DefaultJaegerLogSpans, "Tracer Jaeger Log Spans")
	//startCmd.PersistentFlags().StringVar(&githubClientId, "github.client.id", DefaultGithubClientId, "默认连接的GitHub地址")
	//startCmd.PersistentFlags().StringVar(&githubSecretKey, "github.secret.key", DefaultGithubSecretKey, "默认连接的GitHub地址")

	rootCmd.PersistentFlags().StringVar(&piConnect, "pi.connect", DefaultPiConnect, "树莓派连接方式")
	rootCmd.PersistentFlags().StringVar(&piHost, "pi.host", DefaultPiHost, "树莓派连接Host")
	rootCmd.PersistentFlags().IntVar(&piSSHPort, "pi.ssh.port", DefaultPiSSHPort, "树莓派SSH端口")
	rootCmd.PersistentFlags().StringVar(&piUser, "pi.user", DefaultPiUser, "树莓派登陆用户")
	rootCmd.PersistentFlags().StringVar(&piPassword, "pi.password", DefaultPiPassword, "树莓派登陆密码")
	//rootCmd.PersistentFlags().StringVar(&dbDrive, "db.drive", DefaultDbDrive, "数据库驱动")
	//rootCmd.PersistentFlags().StringVar(&mysqlHost, "db.mysql.host", DefaultMysqlHost, "mysql数据库地址: mysql")
	//rootCmd.PersistentFlags().IntVar(&mysqlPort, "db.mysql.port", DefaultMysqlPort, "mysql数据库端口")
	//rootCmd.PersistentFlags().StringVar(&mysqlUser, "db.mysql.user", DefaultMysqlUser, "mysql数据库用户")
	//rootCmd.PersistentFlags().StringVar(&mysqlPassword, "db.mysql.password", DefaultMysqlPassword, "mysql数据库密码")
	//rootCmd.PersistentFlags().StringVar(&mysqlDatabase, "db.mysql.database", DefaultMysqlDatabase, "mysql数据库")
	//rootCmd.PersistentFlags().StringVar(&redisHosts, "redis.hosts", DefaultRedisHosts, "连接Redis地址")
	//rootCmd.PersistentFlags().IntVar(&redisDb, "redis.db", DefaultRedisDb, "连接Redis DB")
	//rootCmd.PersistentFlags().StringVar(&redisAuth, "redis.auth", DefaultRedisPassword, "连接Redis密码")
	//rootCmd.PersistentFlags().StringVar(&redisPrefix, "redis.prefix", DefaultRedisPrefix, "Redis写入Cache的前缀")
	//rootCmd.PersistentFlags().StringVar(&serverName, "server.name", DefaultServerName, "本系统服务名称")
	//rootCmd.PersistentFlags().StringVar(&serverKey, "server.key", DefaultServerKey, "本系统服务密钥")
	rootCmd.PersistentFlags().StringVar(&serverLogLevel, "server.log.level", DefaultServerLogLevel, "本系统日志级别")
	rootCmd.PersistentFlags().StringVar(&serverLogPath, "server.log.path", DefaultServerLogPath, "本系统日志路径")
	rootCmd.PersistentFlags().StringVar(&serverLogName, "server.log.name", DefaultServerLogName, "本系统日志名称")
	//rootCmd.PersistentFlags().StringVar(&serverDomain, "server.domain", DefaultServerDomain, "本系统域名")
	//rootCmd.PersistentFlags().StringVar(&serverDomainSuffix, "server.domain.suffix", DefaultServerDomainSuffix, "生成域名后缀")
	//rootCmd.PersistentFlags().StringVar(&serverUploadPath, "server.upload.path", DefaultServerUploadPath, "本系统上传文件路径")
	//rootCmd.PersistentFlags().StringVar(&serverHubAddr, "server.hub.addr", DefaultServerHubAddr, "生成镜像仓库的域名")
	//rootCmd.PersistentFlags().StringVar(&serverDefaultCluster, "server.default.cluster", DefaultServerDefaultCluster, "新注册用户默认集群")
	//rootCmd.PersistentFlags().StringVar(&serverDefaultRole, "server.default.role", DefaultServerDefaultRole, "新注册用户默认角色")
	//rootCmd.PersistentFlags().StringVar(&serverDefaultNamespace, "server.default.namespace", DefaultServerDefaultNamespace, "新注册用户默认命名空间")
	//rootCmd.PersistentFlags().Int64Var(&serverSessionTimeout, "server.session.timeout", DefaultServerSessionTimeout, "本系统session超时时间")
	//rootCmd.PersistentFlags().Int64Var(&serverTerminalSessionTimeout, "server.terminal.session.timeout", DefaultServerTerminalSessionTimeout, "本系统终端session超时时间")
	//rootCmd.PersistentFlags().BoolVar(&serverDebug, "server.debug", DefaultServerDebug, "是否开启Debug模式")
	//rootCmd.PersistentFlags().StringVar(&serverHttpProxy, "server.http.proxy", DefaultServerHttpProxy, "请求外部服务的Http代理地址")
	//rootCmd.PersistentFlags().BoolVar(&serverSelfQueue, "server.self.queue", DefaultServerSelfQueue, "是否使用http服务启动自动监听队列")
	//rootCmd.PersistentFlags().IntVar(&corsMaxAge, "cors.max.age", DefaultCORSMaxAge, "允许跨域访问的最大时间")
	//rootCmd.PersistentFlags().IntVar(&maxCacheTTL, "cache.max.ttl", DefaultMaxCacheTTL, "缓存最大存活时间")

	cameraJpegCmd.PersistentFlags().StringVarP(&cameraJpegOutput, "camera.jpeg.output", "o", DefaultCameraJpegOutput, "照片存储")
	cameraJpegCmd.PersistentFlags().StringVar(&cmdCameraJpegBin, "cmd.camera.jpeg.bin", DefaultCmdCameraJpegBin, "拍照命令路径")

	cameraCmd.AddCommand(cameraJpegCmd)

	addFlags(rootCmd)
	rootCmd.AddCommand(cameraCmd)
}

func Run() {
	serverLogLevel = envString(EnvServerLogLevel, DefaultServerLogLevel)
	serverLogPath = envString(EnvServerLogPath, DefaultServerLogPath)
	serverLogName = envString(EnvServerLogName, DefaultServerLogName)

	cameraJpegOutput = envString(EnvCameraJpegOutput, DefaultCameraJpegOutput)
	cmdCameraJpegBin = envString(EnvCmdCameraJpegBin, DefaultCmdCameraJpegBin)

	piConnect = envString(EnvPiConnect, DefaultPiConnect)
	piHost = envString(EnvPiHost, DefaultPiHost)
	piSSHPort, _ = strconv.Atoi(envString(EnvPiSSHPort, "22"))
	piUser = envString(EnvPiUser, DefaultPiUser)
	piPassword = envString(EnvPiPassword, DefaultPiPassword)

	if err = rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute", err.Error())
		os.Exit(-1)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func addFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		fmt.Println(gf.Name, gf.Value)
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}

func prepare(ctx context.Context) error {
	logger = logging.SetLogging(logger, serverLogPath, serverLogName, serverLogLevel)
	apiSvc = api.New()
	return nil
}
