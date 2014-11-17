package main
import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"io"
	"github.com/dchest/uniuri"
	"database/sql"
        _ "github.com/go-sql-driver/mysql"
)

func main() {

	readconfig()

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


	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "text/plain; charset=utf-8")

		cookie := r.FormValue("cookie")
		label := r.FormValue("label")
		path := r.FormValue("path")

	        /* check for cookie. no cookie, no nookie. */
		sqlstr := string("SELECT `account_idx` FROM cookies WHERE `cookie`=" + escape(cookie))
	        rows, err := tx.Query(sqlstr)
	        if (err!=nil) {
	                fmt.Println("SQL Error: " + sqlstr)
	                fmt.Println(err)
	                os.Exit(1)
	        }
	        var account_idx int
		rows.Next()
	        err = rows.Scan(&account_idx)
		if (err!=nil) {
			http.Error(w,"401 Unauthorized",http.StatusUnauthorized)
	        } else {
		        rows.Close()

			fname := uniuri.New()

		        out, err := os.Create("/uploads/"+fname)
			if err != nil {
				fmt.Fprintln(w,err)
			}
			defer out.Close()

			file,fh,err := r.FormFile("file")
			if err!=nil {
				fmt.Fprintln(w,err)
			}
			defer file.Close()
			_, err = io.Copy(out, file)
		        if err != nil {
		        	fmt.Fprintln(w,err)
		        }
			/* record in database */
			

		}
	})); err != nil {
		fmt.Println(err)
	}

}
