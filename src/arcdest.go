/*
	arcdest - Create Archiver.Photo destination
	usage: arcdest [label] [serial] [size]

	label: the label you are giving the disk.
		You might want to set this with glabel / GEOM system
		(but not required)

	serial: the drive serial number. not essential but potentiall useful.
		You can get the serial number of the drive using the included
		'arcserial' command (FreeBSD systems). 

	size: block size of the drive. Used to calculate an 'idea' of the 
		size of your archive.

	'label' is the only essential parameter (however all parameters
		are required.) 'label' is used when posting images to the 
		archive, so that the index will reference the correct 
		storage medium / destination. the label should match 
		what you have written physically on the outside of the
		actual hard disk, so that you can easily identify 
		and locate the device for retrieval purposes.


Copyright (c) 2014, Waitman Gobble <ns@waitman.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer. 
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

package main
import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "fmt"
        "time"
        "os"
	"flag"
	"strconv"
)

func usage() {
	/* generate usage */
	fmt.Fprintf(os.Stderr, "arcdest [label] [serial] [size]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main(){

	readconfig()

	/* get command line params */
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if (len(args)<3) {
		fmt.Println("Parameters are missing.")
		usage()
		os.Exit(1)
	}
		
	/* connect to database */
        db, e := sql.Open("mysql", config.Connstr)
        defer db.Close()
	if (e!=nil) {
		fmt.Println("Connection Error")
		fmt.Println(e)
		os.Exit(1)
	}
	/* autocommit is on, transaction begin */
        tx, err := db.Begin()
        if (err!=nil){
                fmt.Println("Cannot Connect")
                fmt.Println(err)
                os.Exit(1)
        }

	/* check for user already in database */
	sqlstr := string("SELECT COUNT(idx) as `count` FROM dest WHERE `label`=" + string(escape(args[0])))
        rows, err := tx.Query(sqlstr)
	if (err!=nil) {
		fmt.Println("SQL Error: " + sqlstr)
		fmt.Println(err)
		os.Exit(1)
	}
	var count int
	rows.Next()
	err = rows.Scan(&count)
	if (err!=nil) {
		fmt.Println("Cannot Scan Row")
		fmt.Println(err)
		os.Exit(1)
	}
	rows.Close()

        if (count<1){
		/* label is not found in database */
		/* add this device to dest db */

		stamp := time.Now().Unix()
		sqlstr := "INSERT INTO dest (`idx`,`label`,`serial`,`size`,`sequence`) " + 
			"VALUES (NULL,"+string(escape(args[0]))+
			","+string(escape(args[1]))+
			","+string(escape(args[2]))+
			","+string(strconv.FormatInt(stamp,10))+")";
		rows,err = tx.Query(sqlstr)
		if (err!=nil) {
			fmt.Println("Cannot Execute SQL: "+ sqlstr)
			fmt.Println(err)
			os.Exit(1)
		}
		rows.Close()
		fmt.Println("Dest label: "+args[0]+" Created.")
        } else {
		/* user already exists */
		rows.Close()
		fmt.Println("This device already exists in the database.")
		os.Exit(1)
	}
        db.Close()
}
//EOF
