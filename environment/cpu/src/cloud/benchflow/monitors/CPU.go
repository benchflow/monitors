package main
 
import (
    "fmt"
    "net/http"
    "log"
    "time"
    "strings"
    "strconv"
    "os"
    "sync"
)
import "github.com/fsouza/go-dockerclient"

type Container struct {
	ID string
	statsChannel chan *docker.Stats
	doneChannel chan bool
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

var containers []Container
var monitoring bool
var stopChannel chan bool
var waitGroup sync.WaitGroup

func attachToContainer(client docker.Client, container Container) {
	go func() {
		err := client.Stats(docker.StatsOptions{
			ID: container.ID,
   	 		Stats: container.statsChannel,
    		Stream: true,
    		Done: container.doneChannel,
    		Timeout: 0,
			})
		if err != nil {
    		log.Fatal(err)
    		}
		}()
	}

func monitorStats(container Container){
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
			select {
			case <- stopChannel:
				//container.doneChannel <- true;
				waitGroup.Done()
				return
			default:
				stat := <- container.statsChannel
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
			}
		}()
	}

func totalHandler(w http.ResponseWriter, r *http.Request) {
	if !monitoring {
		fmt.Fprintf(w, "Currently not Monitoring")
		return
	}
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
	if !monitoring {
		fmt.Fprintf(w, "Currently not Monitoring")
		return
	}
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

func startMonitoring(w http.ResponseWriter, r *http.Request) {
	if monitoring {
		fmt.Fprintf(w, "Already monitoring")
		return
	}
	client := createDockerClient()
	contEV := os.Getenv("CONTAINERS")
	conts := strings.Split(contEV, ":")
	containers = []Container{}
	stopChannel = make(chan bool)
	for _, each := range conts {
		statsChannel := make(chan *docker.Stats)
		doneChannel := make(chan bool)
		c := Container{ID: each, statsChannel: statsChannel, doneChannel: doneChannel}
		containers = append(containers, c)
		attachToContainer(client, c)
		monitorStats(c)
		waitGroup.Add(1)
	}
	monitoring = true
	fmt.Fprintf(w, "Started Monitoring")
}

func stopMonitoring(w http.ResponseWriter, r *http.Request) {
	if !monitoring {
		fmt.Fprintf(w, "Currently not Monitoring")
		return
	}
	close(stopChannel)
	waitGroup.Wait()
	monitoring = false
	fmt.Fprintf(w, "Stopped monitoring")
}

func createDockerClient() docker.Client {
	path := os.Getenv("DOCKER_CERT_PATH")
	endpoint := "tcp://"+os.Getenv("DOCKER_HOST")+":2376"
	endpoint = "tcp://192.168.99.100:2376"
    ca := fmt.Sprintf("%s/ca.pem", path)
    cert := fmt.Sprintf("%s/cert.pem", path)
    key := fmt.Sprintf("%s/key.pem", path)
    client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	//endpoint := "unix:///var/run/docker.sock"
    //client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
		}
	return *client
	}
 
func main() {
	monitoring = false
	
    http.HandleFunc("/total", totalHandler)
    http.HandleFunc("/delta", deltaHandler)
    http.HandleFunc("/start", startMonitoring)
	http.HandleFunc("/stop", stopMonitoring)
    http.ListenAndServe(":8080", nil)
}