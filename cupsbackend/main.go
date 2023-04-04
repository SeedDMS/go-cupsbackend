package main

import (
    "fmt"
    "os"
    "log"
    "log/syslog"
    "strings"
    "seeddms.org/seeddms/cupsbackend"
)

const CupsBackendOk = 0              /* Job completed successfully */
const CupsBackendFailed = 1          /* Job failed, use error-policy */
const CupsBackendAuthRequired = 2    /* Job failed, authentication required */
const CupsBackendHold = 3            /* Job failed, hold job */
const CupsBackendStop = 4            /* Job failed, stop queue */
const CupsBackendCancel = 5          /* Job failed, cancel job */
const CupsBackendRetry = 6           /* Job failed, retry this job later */
const CupsBackendRetryCurrent = 7    /* Job failed, retry this job immediately */

func main() {
    syslogger, err := syslog.New(syslog.LOG_INFO, "seeddms-backend")
    if err != nil {
        log.Fatalln(err)
    }

    log.SetOutput(syslogger)

    args := os.Args
    if len(args) == 1 {
        fmt.Println("direct", "seeddms:/", "\"Unknown\" \"SeedDMS\"")
        os.Exit(0)
    }
    // There are either 6 or 7 arguments
    // If the backend is called as a raw printer the optional last argument
    // is the unchanged file to be printed, otherwise the postscript file
    // need to be read from stdin
    if len(args) < 6 || len(args) > 7 {
        fmt.Println("Usage:", args[0], "job-id user title copies options [file]")
        os.Exit(1)
    }

    log.Printf("Starting seeddms backend jobid=%s user=%s title=%s", args[1], args[2], args[3])

    // Pass the user name, because the config file is searched in the user's home
    // Also pass the options because it used for setting some things
    cfg, err := cupsbackend.NewConfig(args[2], args[5])
    if err != nil {
        log.Printf("Error reading configuration %s", err)
        os.Exit(1)
    }

    if cfg.LogLevel == "debug" {
        s := strings.Split(args[5], " ")
        for _, e := range s {
            log.Print(e)
        }
        for _, e := range os.Environ() {
            log.Print(e)
        }
    }

    if  len(args) == 7 {
        log.Printf("Running in raw mode")
        err = cupsbackend.Raw(cfg, args[1], args[2], args[3], args[5], args[6])
    } else {
        log.Printf("Running in print mode")
        err = cupsbackend.Print(cfg, args[1], args[2], args[3], args[5])
    }
    if(err == nil) {
        os.Exit(0)
    } else {
        log.Printf("Error executing backend %s", err)
        os.Exit(1)
    }

}

