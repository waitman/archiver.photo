package main
import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "fmt"
        "time"
        "os"
)

func main(){

	readconfig()

	/* connect to database */
	db, e := sql.Open("mysql", config.Connstr)

        defer db.Close()
        tx, err := db.Begin()
        if (err!=nil){
                fmt.Print("Cannot Connect\n")
                fmt.Print(err)
                os.Exit(1)
        }
        st, err := tx.Prepare("SELECT * FROM dest")
        if (err!=nil){
                fmt.Print("Cannot Execute SQL\n")
                fmt.Print(e)
                os.Exit(1)
        }
        rows, err := st.Query()
        if (err!=nil){
                fmt.Print("Cannot Get Rows")
                fmt.Print(e)
                os.Exit(1)
        }
        i := 0
	fmt.Printf("%s\t%16s\t%24s\t%10s\t%s\n","idx","label","serial","size","tm")
	fmt.Printf("%s\t%16s\t%24s\t%10s\t%s\n","---","-----------","------------","--------","-----------------------------")
        for rows.Next() {
                i++
                var label string
                var size int
                var serial string
                var sequence int64
                var idx int
		
                err = rows.Scan(&idx,&label,&serial,&size,&sequence)
                if (err!=nil) {
                        fmt.Print("Cannot Scan Row")
                        fmt.Print(err)
                        os.Exit(1)
                }
                tm := time.Unix(sequence,0)


                fmt.Printf("%d\t%16s\t%24s\t%10d\t%s\n",idx,label,serial,size,tm)
        }
        db.Close()
}

