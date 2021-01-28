package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tianhongw/tinyid/config"
	"github.com/tianhongw/tinyid/conn"
	"github.com/tianhongw/tinyid/handler"
	handlerOption "github.com/tianhongw/tinyid/handler/option"
	tlog "github.com/tianhongw/tinyid/pkg/log"
	"github.com/tianhongw/tinyid/repository"
	"github.com/tianhongw/tinyid/server"
	"github.com/tianhongw/tinyid/service"
	"github.com/tianhongw/tinyid/version"
)

const (
	defaultCfgFile = "$HOME/.tinyid.toml"
	defaultCfgType = "toml"
)

var (
	cfgFile string
	cfgType string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "x",
		Version: fmt.Sprintf(
			`%s
Git branch: %s
Git commit: %s
Git summary: %s
Commit time: %s
Build time: %s`,
			version.Version,
			version.GitBranch,
			version.GitCommit,
			version.GitSummary,
			version.GitCommitTime,
			version.BuildTime,
		),
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default is %s)", defaultCfgFile))
	flags.StringVarP(&cfgType, "type", "t", "", fmt.Sprintf("config file type (default is %s)", defaultCfgType))

	addProfilingFlags(flags)

	return cmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("get home dir failed, ", err)
		}

		cfgFile = strings.Replace(defaultCfgFile, "$HOME", home, 1)
	}

	if cfgType == "" {
		cfgType = defaultCfgType
	}

	if cfg, err := config.Init(cfgFile, cfgType); err != nil {
		log.Fatal("init config file failed, ", err)
	} else {
		cfgFile = cfg
	}
}

func serve() {
	cfg := config.GetConfig()

	logger, err := tlog.NewLogger(
		cfg.Log.Type,
		tlog.WithFormat(cfg.Log.Format),
		tlog.WithLevel(cfg.Log.Level),
		tlog.WithOutputs(cfg.Log.Outputs),
		tlog.WithErrorOutputs(cfg.Log.ErrorOutputs),
	)
	if err != nil {
		log.Fatal("init logger failed, ", err)
	}

	db, err := conn.NewConn(
		conn.WithAddressAndAuth(cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.Username,
			cfg.DB.Password,
			cfg.DB.Name),
		conn.WithMaxIdleConns(cfg.DB.MaxIdleConns),
		conn.WithMaxOpenConns(cfg.DB.MaxOpenConns))
	if err != nil {
		log.Fatal("init database failed, ", err)
	}

	repo, err := repository.NewRepository(db, repository.WithLogger(logger))
	if err != nil {
		log.Fatal("init repositories failed, ", err)
	}

	services, err := service.NewService(
		repo,
		service.WithConfig(cfg),
		service.WithLogger(logger),
	)
	if err != nil {
		log.Fatal("init service failed, ", err)
	}

	handlers, err := handler.NewHandler(services, handlerOption.WithLogger(logger))
	if err != nil {
		log.Fatal("init handlers failed, ", err)
	}

	srv, err := server.NewServer(handlers, cfg.Port)
	if err != nil {
		log.Fatal("init server failed, ", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal("start server failed, ", err)
	}
}
