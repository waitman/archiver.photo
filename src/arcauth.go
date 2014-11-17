/*
arcauth.cgi 

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
	"fmt"
	"code.google.com/p/go.crypto/bcrypt"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "time"
        "os"
        "strconv"
        "net/http"
        "net/http/cgi"
)

func main() {

	readconfig()

	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		LogStatus := "N"

                header := w.Header()
                header.Set("Content-Type", "text/plain; charset=utf-8")
		
		password := []byte(r.FormValue("pass"))

	        /* connect to database */
	        db, e := sql.Open("mysql", config.Connstr)
	        defer db.Close()
	        if (e!=nil) {
	                fmt.Println("Connection Error")
	                fmt.Println(e)
			/* TODO: probably shouldn't bail like this on http/cgi */
	                os.Exit(1)
	        }
	        /* autocommit is on, transaction begin */
	        tx, err := db.Begin()
	        if (err!=nil){
	                fmt.Println("Cannot Connect")
	                fmt.Println(err)
	                os.Exit(1)
	        }

		sqlstr:="SELECT `passwd`,`idx` FROM accounts WHERE `user`=" + escape(r.FormValue("user"))

		/* check for user in database */
		rows, err := tx.Query(sqlstr)
		if (err!=nil) {
	                fmt.Println("SQL Error")
	                fmt.Println(err)
	                os.Exit(1)
	        }
	        var passwd string
		var idx int64
	        rows.Next()
	        err = rows.Scan(&passwd,&idx)
	        if (err!=nil) {
			http.Error(w,"401 Unauthorized",http.StatusUnauthorized)
       		} else {
		        rows.Close()
			hashPassword := []byte(passwd)
			err = bcrypt.CompareHashAndPassword(hashPassword,password)
			if err!=nil {
				http.Error(w,"401 Unauthorized",http.StatusUnauthorized)
			} else {
				LogStatus = "Y"
				/* Create Cookie */
				stamp := time.Now().Unix()
				cookie,err := bcrypt.GenerateFromPassword(hashPassword,10)
				sqlstr = "INSERT INTO cookies (`idx`,`account_idx`,`remote_addr`,`cookie`,`sequence`)" +
					" VALUES (NULL,"+string(strconv.FormatInt(idx,10))+","+escape(string(r.RemoteAddr))+
					","+escape(string(cookie))+","+string(strconv.FormatInt(stamp,10))+")";
				rows,err = tx.Query(sqlstr)
				if (err!=nil) {
					fmt.Println("Cannot Execute SQL: "+ sqlstr)
					fmt.Println(err)
					os.Exit(1)
				}
                		rows.Close()
				/* Auth Success */
				fmt.Fprintln(w,string(cookie)) 
			}
		}
		stamp := time.Now().Unix()
                sqlstr = "INSERT INTO audit (`idx`,`user`,`remote_addr`,`status`,`sequence`) " +
                        "VALUES (NULL,"+escape(r.FormValue("user"))+
			","+escape(string(r.RemoteAddr))+
			","+escape(LogStatus)+
                        ","+string(strconv.FormatInt(stamp,10))+")";
                rows,err = tx.Query(sqlstr)
                if (err!=nil) {
                        fmt.Println("SQL Error")
			fmt.Println(err)
			os.Exit(1)
                }
                rows.Close()


        })); err != nil {
                fmt.Println(err)
        }

}
//EOF
