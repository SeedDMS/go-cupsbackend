# Cups backend for SeedDMS

This is a cups backend to add SeedDMS as a printer. Printing to
a printer based on this backend will upload the printed file into
a SeedDMS installation.

## How it workѕ

If you setup a printer in cups use the ppd file `seeddms-pdf.ppd`
shipped with this backend. It contains the line

    \*cupsFilter:    "application/vnd.cups-pdf 0 -"

which tells cups that application/vnd.cups-pdf can be send straight to
seeddms.  Without that line cups expects seeddms to be a postscript
printer and runs the filter pdftops before passing the output to the
backend.  In that case `FINAL_CONTENT_TYPE` will be set to
application/vnd.cups-postscript and the backend will convert it back
to pdf. That would be a lot of unnecessary converting.

Adding the above filter will tell cups that the output of the initial
filter pdftopdf can be processes by the backend and no more conversion
is done.  In that case `FINAL_CONTENT_TYPE` is set to
application/vnd.cups-pdf which keeps the backend from running
ghostscript.

You could even replace all the filtering done by cups by setting

    \*cupsFilter2:   "image/png image/png 0 -"

(please note the '2')

This will override cupsFilter and even turn off the initial `pdftopdf`
filtering done by cups. A file of mime-type `image/png` will be uploaded
right away.  The disadvantage of this approach is, that you need
cupsFilter2 for each mime-type you intend to support. Another
disadvantage is, that `pdftopdf` does some useful page management
(rotation and selecting page ranges) which isn't done anymore, unless
the printing application does it itself (like web browsers do). The
PPD file `seeddms-passthru.ppd` sets `cupsFilter2` for `image/png`,
`image/jpeg`, `application/pdf` and `application\postscript`.

## Configuration

The path and the name of the configuration file has been changed in 0.0.2 for
more flexibility and a more common paths. The configuration file in the
user's home directory must be located in `.config/seeddms-cups/printers.yaml`
The system wide configuration must be in `/etc/seeddms-cups/printers.yaml`.
Additionally, the new path `/etc/seeddms-cups/<username>/printers.yaml`
will be checked for a configuration file. The precedence of the files is

1. User's home directory
2. `/etc/seeddms-cups/<username>/printers.yaml`
3. `/etc/seeddms-cups/printers.yaml`

Before version 0.0.2 the
backend read the configuration file `.seeddms-cups.yaml` from
either `/etc/seeddms-cups` or the user's home directory.

The user is
the person issuing the print job.  If both, the cups server (having
the seeddms backend installed) and the client run on the same
computer, the user will be an existent user on the system and the
backend can easily access the user's home directory and
read the configuration file. But if the client runs on a different
computer, the user is likely to be somebody not available on the
server. Actually, if the client is for example a mobile phone, the user name
is often just be set to the name
of the phone. So, printing from your 'Redmi Note 8T' will set the user
to 'Redmi Note 8T'. In any of those cases the backend will not find the
user's home directory and therefore will read the system wide configuration.
Since version 0.0.2 the location of the configuration file can also be user
dependent in `/etc/seeddms-cups/<username>/printers.yaml`.

This configuration file may contain several sections. Each for a
configured printer in cups. If none of the section names match the
printer name, the `default` section will be used.

Each section must at least define the parameters:

  * `Url` (Url of restapi)
  * `ApiKey` or `User` and `Password`
  * `FolderId` (Id of folder to store printed documents)
  * `LogLevel` (set to `debug` for more verbose logging)

Example:

    default:
      Url: http://your-seeddms-host/restapi/index.php
      ApiKey: your-secret-key
      FolderId: 1
      LogLevel: debug
   
There is an example configuration `seeddms-cups.yaml`. Copy it into one
of the checked locations on your server and adjust it to your needs.

## Pitfalls

Systems having AppArmor installed may prevent calling the backend
because it request access on files not allowed for `/usr/sbin/cupsd`.
Check your system messages with `dmesg`. If it contains lines like

   [39772.258276] audit: type=1400 audit(1709830581.956:9349): apparmor="DENIED" operation="open" class="file" profile="/usr/sbin/cupsd" name="/proc/241961/cgroup" pid=241961 comm="cupsd" requested_mask="r" denied_mask="r" fsuid=0 ouid=0

If it does, then AppArmor blocks the execution of the backend script.
On Debian the above example will require to add a line

@{PROC}/*/cgroup r,

to your `/etc/apparmor.d/usr.sbin.cupsd`. There are already similar lines, but
none of them allows read access on `/proc/<number>/cgroup`

## Debuging

This cups backend logs many usefull information to the sys logger. If the
log level is set to `debug` it will also log the environment variables, which
contain most of the revelant data for the backend. On recent debian systems (>= 12)
just run 

     journalctl -f -u cups

to monitor the execution of the backend.
