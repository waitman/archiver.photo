
Archiver.Photo

------------------------------------------
        Edit Configuration File
        /etc/archiver.photo.json
------------------------------------------

This software was developed to solve the problem of inexpensively 
indexing a photo archive of several years of work. The strategy is to 
use an array of disconnected drives as storage medium for the 
archive of RAW images, then create the index and thumbnails on 
a cloud based platform. 

Cost analysis of 32TB array

a) disconnected.  10x3TB
(using device similar to Sabrent USB 3.0 SATA III docking station.)
cheapest - $1,000
good quality/performance $1,500


b) connected. 8x4TB NAS
cheapest - $2,350
good quality/performance $5,500
+electricity
+administration
+maintenance



1. Prerequisites

go version 1.3 or later

mysql client 5.5 or later

go add-on packages:

code.google.com/p/go.crypto/bcrypt
github.com/go-sql-driver/mysql
github.com/dchest/uniuri

The add-on packages may be installed using the following 
commands:

# go get code.google.com/p/go.crypto/bcrypt
# go get github.com/go-sql-driver/mysql
# go get github.com/dchest/uniuri

You will also need access to a MySQL Server.

This software was built on a Debian VM from the Google 
Compute System. The database server is Google Cloud SQL
The images are stored on Google Cloud Storage.

You will need an HTTP server to run the CGI programs.
This system was developed using the lighttpd web server, 
however it should also work on Apache or NGINX.

There are configuration examples for lighttpd in the 
'share' directory of the software distribution.

The client machine which processes the image archive
will need dcraw, ImageMagick and exif installed.


2. Build and Install

'make' will build the software. The binary executables will 
be in the ./build directory.

'make install' will install the software to the CGI-BIN and
PREFIX/bin directories. It will also copy the configuration 
file from ./etc into /etc/. The install routine will create
a backup of /etc/archiver.photo.json if it already exists.
The backup file is named archiver.photo.json~ and located 
in the /etc directory.

'make clean' will remove the binary executables from the
./build directory

'make deinstall' will erase the binary executables from 
PREFIX/bin, and CGI-BIN directories. It will not erase 
the configuration file.


3. Environment Variables

The following list are the environment variables used by 
the build script. The defaults indicated will be used if
the variable is not set.

GO ?=           go
GO_FLAGS ?=     -ldflags "-s -w"
PREFIX ?=       /usr
CGI-BIN ?=      $(PREFIX)/lib/cgi-bin
ETC ?=          /etc
INSTALL ?=      install

For example, if you want the executables in /usr/local/bin
then set the PREFIX as follows:

export PREFIX=/usr/local


4. Usage

Note: The CGI programs are intented to be accessed over a 
network using an HTTP client.

arctoken - use this to get and cache your oauth token
arcuser [username] [password] - create user
arcdest [label] [serial] [size] - create a destination device id
arcdestlist - list destinations

arcauth.cgi [username] [password] - attempt login, returns cookie
arcposter.cgi [] - post data to archive

