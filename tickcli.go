package main

import ("fmt"
        "os"
        "net/http"
        "log"
        "io/ioutil"
        "encoding/xml"
        "strconv"
        "time"
)

type TickLog struct {
  Entries []Entry `xml:"entry"`
}

type Entry struct {
  Date string `xml:"date"`
  Note string `xml:"notes"`
  Project string `xml:"project_name"`
  Task string `xml:"task_name"`
  Hours float64 `xml:"hours"`
}

func main() {
  args := os.Args

  if len(args) > 1 {
    executeCommand(args[1], args[2:])
  } else {
    executeCommand("help", args[:1])
  }
}

func executeCommand(command string, args []string) {
  fmt.Println("Processing command: " + command)

  switch command {
  default:
    help()
  case "log":
    log_cmd()
  }
}

func help() {
  fmt.Println("Help & Usage for tickcli")
  fmt.Println("")
  fmt.Println("help    This message")
}

func log_cmd() {
  // Todo here: accept the number of days to display with default of 10
  start_date := time.Now().Add((time.Duration(-24)*time.Hour) * 10).Format("2006-01-02") //This is 10 days ago
  // Todo: Move these to a configuration file instead
  company := ""
  email := ""
  pass := ""

  request_url := "https://" + company + ".tickspot.com/api/entries.json?email=" + email + "&password=" + pass + "&updated_at=" + start_date

  response := getResponse(request_url)

  var tickLog TickLog
  err := xml.Unmarshal(response, &tickLog)
  checkError(err)

  for _,entry := range tickLog.Entries {
    fmt.Println(entry.Date + " - " + entry.Project + " - " + entry.Task + " - " + strconv.FormatFloat(entry.Hours, 'f', 1, 64) + " hours")
  }

}

func getResponse(url string) []byte {
  resp, err := http.Get(url)
  checkError(err)

  body, err := ioutil.ReadAll(resp.Body)
  checkError(err)
  return body
}

func checkError(err error) {
  if err != nil {
    fmt.Print(err)
    log.Fatal(err)
  }
}
