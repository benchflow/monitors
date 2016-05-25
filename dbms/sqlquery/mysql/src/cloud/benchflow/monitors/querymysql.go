package main
 
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    //"fmt"
    "os"
    //"strconv"
    "net/http"
    "encoding/json"
)

type Response struct {
  Result bool
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	entryString := os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_USER_PASSWORD")+"@"+"tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB_NAME")
    db, err := sql.Open("mysql", entryString)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }
    query := r.FormValue("query")
    value := r.FormValue("value")
    method := r.FormValue("method")
    rows , err := db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	    return
    }
    rowI := 0
    response := Response{false}
    for rows.Next() {
    	rowI = rowI + 1
        var name string 
        err = rows.Scan(&name)
        if err != nil {
            panic(err)
        }
        if name == value {
        	if(method == "equal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" matches "+value+" \n")
        		response = Response{true}
        		}
        	}
        if name != value {
        	if(method == "nequal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" doesn't match "+value+" \n")
        		response = Response{true}
        		}
        	}
    }
    db.Close()
    
    js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}
 
func main() {
    http.HandleFunc("/data", queryHandler)
    http.ListenAndServe(":8080", nil)
}