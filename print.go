package cupsbackend

import (
    "os"
    "io"
    "bytes"
    "mime"
    "os/exec"
    "log"
    "seeddms.org/seeddms/apiclient"
)

func Print(cfg *Config, jobid string, user string, title string, options string) error {
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
    }
    contenttype := os.Getenv("CONTENT_TYPE")
    if extension, _ := mime.ExtensionsByType(contenttype); extension != nil {
        extraParams["filename"] += extension[0]
    }

    var r io.Reader
    // FINAL_CONTENT_TYPE contains the mime type output by the last filter
    // If it is postscript then run ghostscript otherwise just save it.
    // application/vnd.cups-postscript is the internal mime type of the output
    // of pstops. In a pdf centric workflow this backend should receive
    // application/vnd.cups-pdf and the ppd file should have cupsFilter
    // *cupsFilter: "application/vnd.cups-pdf 0 -"
    finalcontenttype := os.Getenv("FINAL_CONTENT_TYPE")
    if cfg.rungs && finalcontenttype == "application/vnd.cups-postscript" {
        log.Printf("Running ghostscript")
        cmd := exec.Command(
            "gs",
            "-q",
            "-dNOPAUSE",
            "-dQUIET",
            "-dBATCH",
            "-dSAFER",
            "-dAutoRotatePages=/PageByPage",
            "-sDEVICE=pdfwrite",
            "-dCompatibilityLevel=1.4",
            "-sPAPERSIZE=a4",
            "-dPDFSETTINGS=/printer",
            "-sOutputFile=-",
            "-c save pop",
            "-")
        cmd.Stdin = os.Stdin
        out, err := cmd.Output()
        if err != nil {
            log.Printf("Error running ghostscript command")
            return err
        }
        r = bytes.NewReader(out)
    } else {
        log.Printf("Not running ghostscript, reading document from stdin.")
        r = os.Stdin
    }
    res, err := c.Upload(r, extraParams, cfg.folderid)

    if err != nil {
        log.Printf("Error uploading file. Status: %d, ErrorMsg: %s\n", c.StatusCode, c.ErrorMsg)
        return err
    }
    log.Printf("Document (%d Bytes) saved with id=%d in folder id=%d", res.Data.Size, res.Data.Id, cfg.folderid)

    return nil
}


