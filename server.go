package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Alert AlertManager API acceptable JSON Data
type Alert struct {
	Labels      map[string]interface{} `json:"labels"`      // Map of Labels for each alert
	Annotations map[string]interface{} `json:"annotations"` // Map of Annotations for each alert
	StartsAt    time.Time              `json:"startsAt"`    // Starting time of an alert
	EndsAt      time.Time              `json:"endsAt"`      // Ending time of an alert
}

// AMRecord represents AlertManager response record
type AMRecord struct {
	Status string
	Data   []Alert
}

// Rule contains AM rule to process the alert
type Rule struct {
	Name         string `json:"name"`           // name of the alert to match
	Namespace    string `json:"namespace"`      // pod namespace, e.g. dbs
	Action       string `json:"action"`         // action value, e.g. restart
	Pod          string `json:"pod"`            // name of pod attribute, e.g. apod
	Env          string `json:"env"`            // k8s environment
	PodNameMatch string `json:"pod_name_match"` // matching string for pod name
}

// Match matches given alert name with alert record labels
func (r *Rule) Match(alert Alert, verbose int) string {
	if name, ok := alert.Labels["alertname"]; ok {
		if r.Name == name {
			// check if environment between alert and rule are matching
			envMatch := false
			if env, ok := alert.Labels["env"]; ok {
				if env == r.Env {
					envMatch = true
				}
			}
			if env, ok := alert.Annotations["env"]; ok {
				if env == r.Env {
					envMatch = true
				}
			}
			if r.Env == "" { // match any environment
				envMatch = true
			}

			// check if alert is still valid
			//             diff := alert.EndsAt.Sub(alert.StartsAt)
			//             now := time.Now()
			//             validAlert := false
			//             if diff.Hours() < 1 && alert.EndsAt.After(now) && alert.StartsAt.Before(now) {
			//                 validAlert = true
			//             }

			// check if pod name exists in annotations or labels
			if envMatch {
				if verbose > 0 {
					data, err := json.Marshal(alert)
					if err == nil {
						log.Printf("found alert %s for rule %+v", string(data), r)
					} else {
						log.Printf("found alert %+v for rule %+v, error %v", alert, r, err)
					}
				}
				var podName string

				if pod, ok := alert.Annotations[r.Pod]; ok {
					podName = pod.(string)
				}
				if pod, ok := alert.Labels[r.Pod]; ok {
					podName = pod.(string)
				}
				// check if pod name should be matched too
				if r.PodNameMatch != "" {
					if strings.Contains(podName, r.PodNameMatch) {
						return podName
					}
				} else {
					return podName
				}
			}
		}
	}
	return ""
}

// server
func server(configFile string) {
	err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// log time, filename, and line number
	if Config.Verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// main loop
	count := 0
	for {
		if count > 0 {
			time.Sleep(time.Duration(Config.Interval) * time.Second)
		}
		count += 1
		alerts, err := getAlerts(Config.AlertManager)
		if err != nil {
			log.Println("unable to get alerts", err)
			continue
		}
		if Config.Verbose > 0 {
			log.Printf("fetched %d alerts", len(alerts))
		}
		for _, alert := range alerts {
			for _, rule := range Config.Rules {
				pod := rule.Match(alert, Config.Verbose)
				if pod != "" {
					process(alert, pod, rule.Namespace, rule.Action, Config.Verbose)
				}
			}
		}
	}
}

// helper function to get AM alerts data
func getAlerts(rurl string) ([]Alert, error) {
	var records []Alert
	aurl := fmt.Sprintf("%s/api/v1/alerts?active=true&silenced=false&inhibited=false&unprocessed=false", rurl)
	var headers [][]string
	headers = append(headers, []string{"Accept-Encoding", "identify"})
	headers = append(headers, []string{"Accept", "application/json"})
	resp := HttpCall("GET", aurl, headers, nil)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read JSON Data from AlertManager GET API, error: %v\n", err)
		return records, err
	}

	// parse obtained data
	var rec AMRecord
	err = json.Unmarshal(data, &rec)
	return rec.Data, err

}

// helper function to check if pod exists in given namespace
func podExist(pod, namespace string) bool {
	args := []string{"get", "pod", pod, "-n", namespace}
	cmd := exec.Command("kubectl", args...)
	log.Println("execute", cmd)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("unable to execute %v %v, error %v", cmd, args, err)
		return false
	}
	if strings.Contains(string(out), "not found") {
		return false
	}
	return true
}

// process given alert
func process(alert Alert, pod, namespace, action string, verbose int) {
	if verbose > 0 {
		data, err := json.Marshal(alert)
		if err != nil {
			log.Println("unable to marshal alert data, error", err)
			return
		}
		log.Println(string(data))
		log.Printf("perform action '%s' on pod '%s' within namespace '%s'\n", action, pod, namespace)
	}
	if action == "restart" {
		// make sure that pod exists
		if podExist(pod, namespace) {
			args := []string{"delete", "pod", pod, "-n", namespace}
			cmd := exec.Command("kubectl", args...)
			log.Println("execute", cmd)
			out, err := cmd.Output()
			if err != nil {
				log.Printf("unable to execute %v %v, error %v", cmd, args, err)
				return
			}
			log.Println(out)
		} else {
			log.Printf("pod %s does not exists in namespace %s", pod, namespace)
		}
	} else {
		cmd := fmt.Sprintf("kubectl delete pod %s -n %s", pod, namespace)
		log.Println(cmd)
	}
}
