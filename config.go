package cupsbackend

import (
    "fmt"
    "strings"
    "os"
    "os/user"
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    url string
    user string
    pwd string
    apikey string
    folderid int
    rungs bool
    LogLevel string
}

func NewConfig(username string, options string) (*Config, error) {

    viper.SetConfigName(".seeddms-cups.yaml")
    viper.SetConfigType("yaml")
    if user, err := user.Lookup(username); err == nil {
        viper.AddConfigPath(user.HomeDir)
    }
    viper.AddConfigPath("/etc/seeddms-cups/")
    viper.AddConfigPath(".")
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            return nil, fmt.Errorf("could not find configuration file: %w\n", err)
        } else {
            return nil, fmt.Errorf("fatal error in config file: %w\n", err)
        }
    }
    log.Printf("Using configuration file \"%s\"", viper.ConfigFileUsed())

    printer := os.Getenv("PRINTER")
    if printer == "" {
        log.Print("Environment varialbe PRINTER not set, using defaults")
        printer = "default"
    }
    cfgSection := viper.Sub(printer)
    if cfgSection == nil {
        if printer != "default" {
            log.Printf("Printer \"%s\" not set in configuration, trying default", printer)
            printer = "default"
            cfgSection = viper.Sub("default")
        }
        if cfgSection == nil {
            return nil, fmt.Errorf("Configuration for printer \"%s\" not found\n", printer)
        }
    }

    cfg := Config{}
    cfg.url = cfgSection.GetString("Url")
    cfg.user = cfgSection.GetString("User")
    cfg.pwd = cfgSection.GetString("Password")
    cfg.apikey = cfgSection.GetString("ApiKey")
    cfg.folderid = cfgSection.GetInt("FolderId")
    if !cfgSection.IsSet("RunGS") || cfgSection.GetBool("RunGS") {
        cfg.rungs = true
    } else {
        cfg.rungs = false
    }
    if !cfgSection.IsSet("LogLevel") {
        cfg.LogLevel = "info"
    } else {
        cfg.LogLevel = cfgSection.GetString("LogLevel")
    }

    s := strings.Split(options, " ")
    for _, e := range s {
        o := strings.SplitN(e, "=", 2)
        switch o[0] {
        case "LogLevel":
            cfg.LogLevel = o[1]
            log.Printf("Overriding LogLevel with value from printer options (%s)", o[1])
        }
    }
    return &cfg, nil
}

