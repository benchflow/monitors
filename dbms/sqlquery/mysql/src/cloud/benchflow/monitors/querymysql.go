package main
 
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    //"os"
    "net/http"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	//entryString := os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_USER_PASSWORD")+"@"+"tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB_NAME")
	entryString := "root:PASSWORD@tcp(dockerVM:3306)/stuff"
    db, err := sql.Open("mysql", entryString)
    if err != nil {
        panic(err)
    }
    if err == nil {
        fmt.Fprintf(w, "Connected to database \n")
    }
    query := r.FormValue("query")
    fmt.Println(query)
    value := r.FormValue("value")
    method := r.FormValue("method")
    fmt.Fprintf(w, "Performing query: "+query+"\n")
    rows , err := db.Query(query)
    if err != nil {
        panic(err)
    }
    if err == nil {
        fmt.Fprintf(w, "Queried database \n")
    }
    for rows.Next() {
        var name string 
        err = rows.Scan(&name)
        if err != nil {
            panic(err)
        }
        if name == value {
        	if(method == "equal") {
        		fmt.Fprintf(w, name+" matches "+value+" \n")
        		}
        	}
        if name != value {
        	if(method == "nequal") {
        		fmt.Fprintf(w, name+" doesn't match "+value+" \n")
        		}
        	}
    }
    db.Close()
}
 
func main() {
    http.HandleFunc("/status", queryHandler)
    http.ListenAndServe(":8080", nil)
}