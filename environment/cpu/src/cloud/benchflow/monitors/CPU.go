package main
 
import (
    "fmt"
    "net/http"
    "log"
    "time"
    "strings"
    "strconv"
    "os"
)
import "fsouza/go-dockerclient"

type Container struct {
	ID string
	statsChannel chan *docker.Stats
	flagsChannel chan bool
	last5 uint64
	last5Time time.Time
	last30 uint64
	last30Time time.Time
	last60 uint64
	last60Time time.Time
	}

var containers [10]Container

func collectStats(client docker.Client, container Container) {
	go func() {
		err := client.Stats(docker.StatsOptions{
			ID: container.ID,
   	 		Stats: container.statsChannel,
    		Stream: true,
    		Done: container.flagsChannel,
    		Timeout: time.Duration(10),
			})
		if err != nil {
    		log.Fatal(err)
    		}
		}()
	}

func monitorStats(container *Container){
	go func() {
		var count5 uint64
		count5 = 0
		var count30 uint64
		count30 = 0
		var count60 uint64
		count60 = 0
		var timestamp time.Time
		for true{
			for i := 0; i<2; i++ {
				for j := 0; j<6; j++ {
					for k := 0; k<5; k++ {
						stat := (<-container.statsChannel)
						count5 += stat.CPUStats.CPUUsage.PercpuUsage[0]
						count30 += stat.CPUStats.CPUUsage.PercpuUsage[0]
						count60 += stat.CPUStats.CPUUsage.PercpuUsage[0]
						timestamp = stat.Read
						if k == 4 {
							container.last5 = count5/5
							container.last5Time = timestamp
							count5 = 0
							}
						time.Sleep(time.Second)
						}
					if j == 5 {
							container.last30 = count30/30
							container.last30Time = timestamp
							count30 = 0
							}
					}
				if i == 1 {
							container.last60 = count60/60
							container.last60Time = timestamp
							count60 = 0
							}
				}
			}
		}()
	}
 
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Greetings!")
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("select")
	if query == "all" {
		for _, each := range containers {
			if each.ID != "" {
				fmt.Fprintf(w, each.ID + "\n")
				fmt.Fprintf(w, "Avg last 5 seconds: " + strconv.FormatUint(each.last5, 10) + "\n")
				fmt.Fprintf(w, "Avg last 30 seconds: " + strconv.FormatUint(each.last30, 10) + "\n")
				fmt.Fprintf(w, "Avg last 60 seconds: " + strconv.FormatUint(each.last60, 10) + "\n")
				fmt.Fprintf(w, "\n")
				}
			}
		}
	if query != "all" {
		for _, each := range containers {
			if each.ID == query {
				fmt.Println(each.last5)
				fmt.Fprintf(w, each.ID + "\n")
				fmt.Fprintf(w, "Avg last 5 seconds: " + strconv.FormatUint(each.last5, 10) + "\n")
				fmt.Fprintf(w, "Avg last 30 seconds: " + strconv.FormatUint(each.last30, 10) + "\n")
				fmt.Fprintf(w, "Avg last 60 seconds: " + strconv.FormatUint(each.last60, 10) + "\n")
				fmt.Fprintf(w, "\n")
				}
			}
		}
}
 
func main() {
	path := os.Getenv("DOCKER_CERT_PATH")
	endpoint := os.Getenv("DOCKER_HOST")
    ca := fmt.Sprintf("%s/ca.pem", path)
    cert := fmt.Sprintf("%s/cert.pem", path)
    key := fmt.Sprintf("%s/key.pem", path)
    client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
    	log.Fatal(err)
	}
    contEV := os.Getenv("CONTAINERS")
    conts := strings.Split(contEV, ":")
    for i, each := range conts {
    	statsChannel := make(chan *docker.Stats)
		flagsChannel := make(chan bool)
		c := Container{ID: each, statsChannel: statsChannel, flagsChannel: flagsChannel}
		containers[i] = c
		collectStats(*client, containers[i])
		monitorStats(&containers[i])
    	}
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/status", queryHandler)
    http.ListenAndServe(":8080", nil)
}