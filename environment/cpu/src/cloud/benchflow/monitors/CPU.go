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
import "github.com/fsouza/go-dockerclient"

type Container struct {
	ID string
	statsChannel chan *docker.Stats
	flagsChannel chan bool
	last5 uint64
	last5D uint64
	last5Time time.Time
	last30 uint64
	last30D uint64
	last30Time time.Time
	last60 uint64
	last60D uint64
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
		container.last5 = 0
		container.last5D = 0
		var count5 int
		count5 = 0
		var count5Value uint64
		count5Value = 0
		
		container.last30 = 0
		container.last30D = 0
		var count30 int
		count30 = 0
		var count30Value uint64
		count30Value = 0
		
		container.last60 = 0
		container.last60D = 0
		var count60 int
		count60 = 0
		var count60Value uint64
		count60Value = 0
		
		var timestamp time.Time
		var prev uint64
		
		for true{
			stat := (<-container.statsChannel)
			count5Value += stat.CPUStats.CPUUsage.TotalUsage
			count5 += 1
			count30Value += stat.CPUStats.CPUUsage.TotalUsage
			count30 += 1
			count60Value += stat.CPUStats.CPUUsage.TotalUsage
			count60 += 1
			timestamp = stat.Read
			if count5 == 5 {
				prev = container.last5
				container.last5 = count5Value/5
				container.last5D = container.last5 - prev
				container.last5Time = timestamp
				count5 = 0
				}
			if count30 == 30 {
				prev = container.last30
				container.last30 = count30Value/30
				container.last30D = container.last30 - prev
				container.last30Time = timestamp
				count30 = 0
				}
			if count60 == 60 {
				prev = container.last60
				container.last60 = count60Value/60
				container.last60D = container.last60 - prev
				container.last60Time = timestamp
				count60 = 0
				}
			}
		}()
	}

func totalHandler(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprintf(w, each.ID + "\n")
				fmt.Fprintf(w, "Avg last 5 seconds: " + strconv.FormatUint(each.last5, 10) + "\n")
				fmt.Fprintf(w, "Avg last 30 seconds: " + strconv.FormatUint(each.last30, 10) + "\n")
				fmt.Fprintf(w, "Avg last 60 seconds: " + strconv.FormatUint(each.last60, 10) + "\n")
				fmt.Fprintf(w, "\n")
				}
			}
		}
}

func deltaHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("select")
	if query == "all" {
		for _, each := range containers {
			if each.ID != "" {
				fmt.Fprintf(w, each.ID + "\n")
				fmt.Fprintf(w, "Diff last 5 seconds: " + strconv.FormatUint(each.last5D, 10) + "\n")
				fmt.Fprintf(w, "Diff last 30 seconds: " + strconv.FormatUint(each.last30D, 10) + "\n")
				fmt.Fprintf(w, "Diff last 60 seconds: " + strconv.FormatUint(each.last60D, 10) + "\n")
				fmt.Fprintf(w, "\n")
				}
			}
		}
	if query != "all" {
		for _, each := range containers {
			if each.ID == query {
				fmt.Fprintf(w, each.ID + "\n")
				fmt.Fprintf(w, "Diff last 5 seconds: " + strconv.FormatUint(each.last5D, 10) + "\n")
				fmt.Fprintf(w, "Diff last 30 seconds: " + strconv.FormatUint(each.last30D, 10) + "\n")
				fmt.Fprintf(w, "Diff last 60 seconds: " + strconv.FormatUint(each.last60D, 10) + "\n")
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
    http.HandleFunc("/total", totalHandler)
    http.HandleFunc("/delta", deltaHandler)
    http.ListenAndServe(":8080", nil)
}