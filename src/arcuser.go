/*
	arcuser - Create Archiver.Photo user
	usage: arcuser [username] [password]

	username, password not checked for validity,
	ie username: _ and password: _ are allowed.


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
	"code.google.com/p/go.crypto/bcrypt"
)

func usage() {
	/* generate usage */
	fmt.Fprintf(os.Stderr, "usage: arcuser [username] [password]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main(){

	readconfig()

	/* get command line params */
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if (len(args)<2) {
		fmt.Println("Username and/or password are missing.")
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
	sqlstr := string("SELECT COUNT(idx) as `count` FROM accounts WHERE `user`=" + string(escape(args[0])))
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
		/* user is not found in database */
		/* create password hash and insert user into db */
		password := []byte(args[1])
		hashedPassword,err := bcrypt.GenerateFromPassword(password,10)
        	if err!= nil {
			panic(err)
		}
		stamp := time.Now().Unix()
		sqlstr := "INSERT INTO accounts (`idx`,`user`,`passwd`,`sequence`) " + 
			"VALUES (NULL,"+string(escape(args[0]))+","+string(escape(string(hashedPassword)))+
			","+string(strconv.FormatInt(stamp,10))+")";
		rows,err = tx.Query(sqlstr)
		if (err!=nil) {
			fmt.Println("Cannot Execute SQL: "+ sqlstr)
			fmt.Println(err)
			os.Exit(1)
		}
		rows.Close()
		fmt.Println("User "+args[0]+" Created.")
        } else {
		/* user already exists */
		rows.Close()
		fmt.Println("Username already exists in the database.")
		os.Exit(1)
	}
        db.Close()
}
//EOF
