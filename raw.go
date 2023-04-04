package cupsbackend

import (
    "os"
    "log"
    "mime"
    "seeddms.org/seeddms/apiclient"
)

func Raw(cfg *Config, jobid string, user string, title string, options string, file string) error {
    c := apiclient.Connect(cfg.url, cfg.apikey)
    _, err := c.Login(cfg.user, cfg.pwd)
    if err != nil {
        log.Printf("Failed to login: %s\n", err)
        return err
    }
    extraParams := map[string]string{
		"name":        title,
        "comment":     "User: "+user+"\nJob-Id: "+jobid+"\nOptions: "+options,
        "filename":    "cups-"+user+"-"+jobid,
        "version_comment": "",
	}
    contenttype := os.Getenv("CONTENT_TYPE")
    if extension, _ := mime.ExtensionsByType(contenttype); extension != nil {
        extraParams["filename"] += extension[0]
    }

    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()

    res, err := c.Upload(f, extraParams, cfg.folderid)
    if err != nil {
        log.Printf("Error uploading file. Status: %d, ErrorMsg: %s", c.StatusCode, c.ErrorMsg)
        return err
    }
    log.Printf("Document (%d Bytes) saved with id=%d in folder id=%d", res.Data.Size, res.Data.Id, cfg.folderid)
    return nil
}

