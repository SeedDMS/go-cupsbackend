This is a cups backend to add SeedDMS as a printer. Printing to
a printer based on this backend will upload the printed file into
a SeedDMS installation.

How it workѕ
============

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

Configuration
==============

This backend reads a configuration file `.seeddms-cups.yaml` from
either `/etc/seeddms-cups` or the user's home directory. The user is
the person issuing the print job.  If both, the cups server (having
the seeddms backend installed) and the client run on the same
computer, the backend can easily access the user's home directory and
read the configuration file. But if the client runs on a different
computer, the user is likely to be somebody not available on the
server. If the client is a mobile phone, the user may just be the name
of the phone. In any of those case the backend will not find the
user's home directory and therefore will read the configuration from
`/etc/seeddms-cup/.seeddms-cups.yaml`

This configuration file may contain several sections. Each for a
configured printer in cups. If none of the section names match the
printer name, the default section will be used.

Each section must at least define the parameters:

  * Url (Url of restapi)
  * ApiKey or User and Password
  * FolderId (Id of folder to store printed documents)

There is an example configuration `seeddms-cups.yaml`.

