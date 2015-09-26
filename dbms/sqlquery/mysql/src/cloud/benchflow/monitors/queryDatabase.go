package main
 
import (
    _ "go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "net/http"
)
 
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Greetings!")
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("mysql", "root:PASSWORD@tcp(db:3306)/stuff")
    if err != nil {
        panic(err)
    }
    if err == nil {
        fmt.Fprintf(w, "Connected to database \n")
    }
    query := r.FormValue("query")
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
        	if(method == "nequal") {
        		fmt.Fprintf(w, name+" doesn't match "+value+" \n")
        		}
        	}
        if name != value {
        	if(method == "equal") {
        		fmt.Fprintf(w, name+" doesn't match "+value+" \n")
        		}
        	if(method == "nequal") {
        		fmt.Fprintf(w, name+" matches "+value+" \n")
        		}
        	}
    }
    db.Close()
}
 
func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/status", queryHandler)
    http.ListenAndServe(":8080", nil)
}