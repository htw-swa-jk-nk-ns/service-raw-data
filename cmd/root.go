package cmd

import (
	"github.com/htw-swa-jk-nk-ns/service-raw-data/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	cobra.OnInitialize(func() {
		rep := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(rep)

		viper.SetEnvPrefix("SERVICE_RAW_DATA")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()
	})

	//api
	rootCMD.PersistentFlags().String("api-format", "json", "json format ('json' or 'xml')")
	rootCMD.PersistentFlags().Int("api-port", 8889, "api port")

	//db
	rootCMD.PersistentFlags().String("db-type", "redis", "Database type to communicate with")

	//redis
	rootCMD.PersistentFlags().String("redis-addr", "service-redis:6379", "Database address if using the redis driver")
	rootCMD.PersistentFlags().String("redis-pass", "", "Database password if using the redis driver")
	rootCMD.PersistentFlags().Int("redis-db", 0, "Database to use if using the redis driver")

	//mysql
	rootCMD.PersistentFlags().String("mysql-dataSourceName", "tcp:service-mysql:3306", "mysql dataSourceName")

	//api
	err := viper.BindPFlag("api.format", rootCMD.PersistentFlags().Lookup("api-format"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag api format")
		return
	}

	err = viper.BindPFlag("api.port", rootCMD.PersistentFlags().Lookup("api-port"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag api port")
		return
	}

	//db
	err = viper.BindPFlag("db.type", rootCMD.PersistentFlags().Lookup("db-type"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag db-drivername")
		return
	}

	//redis
	err = viper.BindPFlag("redis.addr", rootCMD.PersistentFlags().Lookup("redis-addr"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag redis-addr")
		return
	}

	err = viper.BindPFlag("redis.password", rootCMD.PersistentFlags().Lookup("redis-pass"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag redis-password")
		return
	}

	err = viper.BindPFlag("redis.db", rootCMD.PersistentFlags().Lookup("redis-db"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag redis-db")
		return
	}

	//mysql
	err = viper.BindPFlag("mysql.dataSourceName", rootCMD.PersistentFlags().Lookup("mysql-dataSourceName"))
	if err != nil {
		log.Error().
			AnErr("Error", err).
			Msg("Can't bind flag mysql-addr")
		return
	}
}

var rootCMD = &cobra.Command{
	Use:   "service-raw-data",
	Short: "This tool serves as a gateway to communicate with a database for saving and reading out votes.",
	Long: "It starts an API which mainly offers two functionalities, inserting a vote and reading out all existing votes.\n" +
		"Parameters to reach a database must be provided via CLI flags.\n",
	DisableFlagsInUseLine: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !(viper.GetString("api.format") == "json" || viper.GetString("format") == "xml") {
			return errors.New("invalid api format set")
		}
		if viper.GetString("api.username") != "" && viper.GetString("api.password") == "" {
			return errors.New("username but no password for api authorization set")
		}
		if viper.GetString("api.username") == "" && viper.GetString("api.password") != "" {
			return errors.New("password but no username for api authorization set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.StartAPI()
	},
}

// Execute is the entrypoint for the CLI interface.
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}
