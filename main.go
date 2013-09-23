package main

import (
  "fmt"
  "os"
  "net/http"
  "log"
  "io/ioutil"
  "encoding/xml"
  "strconv"
  "time"
  "github.com/howeyc/gopass"
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
  switch command {
  case "login":
    loginCmd()
  case "log":
    logCmd()
  default:
    help()
  }
}

func help() {
  fmt.Println("Help & Usage for tickcli\n")
  fmt.Println("help    This message")
  fmt.Println("login   Login dawg")
}

func loginCmd() {
  var company, email, password string

  fmt.Print("Company: ")
  fmt.Scan(&company)

  fmt.Print("Email: ")
  fmt.Scan(&email)

  fmt.Print("Password: ")
  bytes := gopass.GetPasswd()
  password = string(bytes[:])

  resp, _ := http.Get("https://" + company + ".tickspot.com/api/clients?email=" + email + "&password=" + password)

  if resp.StatusCode != http.StatusOK {
    fmt.Println("Error logging in... maybe you messed up the password?")
    os.Exit(1);
  }

  config := Config{company, email, password}
  config.WriteToFile();

  fmt.Println("\nSuccess! saved your credentials to " + config.FilePath())
}

func logCmd() {
  config := Config{}
  config.ReadFromFile()

  // Todo here: accept the number of days to display with default of 10
  startDate := time.Now().Add((time.Duration(-24)*time.Hour) * 10).Format("2006-01-02") //This is 10 days ago

  requestUrl := "https://" + config.Company + ".tickspot.com/api/entries.json?email=" + config.Email + "&password=" + config.Password + "&updated_at=" + startDate

  response := getResponse(requestUrl)

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
