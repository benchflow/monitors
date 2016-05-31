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

// Struct for the Json respone for this monitor
type Response struct {
  Result bool
  Query_response string
}

// Handler for when the monitor is called
func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Connecting to database
	entryString := os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_USER_PASSWORD")+"@"+"tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB_NAME")
    db, err := sql.Open("mysql", entryString)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }
    // Taking the query to perform, the value to compare, and the method to apply (equal, nequal)
    query := r.FormValue("query")
    value := r.FormValue("value")
    method := r.FormValue("method")
    rows , err := db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	    return
    }
    // Checks over the rows if the comparison holds
    result := false
    queryResponse := ""
    for rows.Next() {
        var name string 
        err = rows.Scan(&name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
	    	return
        }
        queryResponse += name
        queryResponse += " "
        if name == value {
        	if(method == "equal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" matches "+value+" \n")
        		result = true
        		} else
        	if(method == "nequal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" doesn't match "+value+" \n")
        		result = false
        		}
        	} else
        if name != value {
        	if(method == "equal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" matches "+value+" \n")
        		result = false
        		} else
        	if(method == "nequal") {
        		//fmt.Fprintf(w, "Row "+strconv.Itoa(rowI)+" doesn't match "+value+" \n")
        		result = true
        		}
        	}
        }
    response := Response{result, queryResponse}
    // Closes db and sends response to client
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