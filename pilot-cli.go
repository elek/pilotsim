package main

import (
   "fmt"
   "os/exec"
   "encoding/json"
   "os/user"
   "strings"
   "encoding/base64"
)

type Value struct {
   Name string
   Value string
   Command string
   Query string
}
type Response struct {
     Version int8
     Title string
     Type string
     Status string
     Message string
     Text string
     Values []Value
}


func run(command string) {
    usr, e1 := user.Current()
    if e1 != nil {
        panic( e1 )
    }
    fmt.Println("Executing",usr.HomeDir,"/",command)
    commandArray := strings.Split(usr.HomeDir + "/" + command," ")
    var cmd *exec.Cmd
    if len(commandArray) == 1 {
       cmd = exec.Command(commandArray[0])
    } else {
       cmd = exec.Command(commandArray[0], commandArray[1:]...)
    }
    out, err := cmd.Output()
    if err != nil {
       panic(err)
    }
    var lst Response
    err2 := json.Unmarshal(out, &lst)
    fmt.Println(string(out))
    if err2 != nil {
    	panic(err2)
    }

    handle(&lst)
}

func print(context *Response) {
    fmt.Println("****", context.Title, "****")
    for i,v := range context.Values {
       if len(v.Command) > 0 {
          fmt.Printf("[%d]", i)
       }
       fmt.Print(v.Name)
       if len(v.Value) > 0 {
          fmt.Println(": ",v.Value)
       }
       fmt.Println()
    }
}

func handle(context *Response) {

   switch context.Type {
      case "text":
         handleText(context)
      case "status":
         handleStatus(context)
      case "commands":
         handleCommand(context)
   }
}
func handleText(context *Response) {
   fmt.Println(context.Title)
   fmt.Println()
   data, err := base64.StdEncoding.DecodeString(context.Text)
   if err != nil {
      fmt.Println("error:", err)
   }
   fmt.Printf("%q\n", data)
}
func handleStatus(context *Response) {
   fmt.Println("STATUS: ", context.Status)
   fmt.Println(context.Title, ": ", context.Message)
}
func handleCommand(context *Response) {
   print(context)

   var i int
   var e error
   _, e = fmt.Scanf("%d", &i)
   for e != nil {
      _, e = fmt.Scanf("%d", &i)
      if e != nil {
      	 fmt.Println(e)
      }
   }

   v := context.Values[i]
   if len(v.Query) > 0 {
      fmt.Println(v.Query)
      var param string
      _, e = fmt.Scanf("%s", &param)
      if e != nil {
      	 fmt.Println(e)
      }

      run(v.Command + " " + param)
   } else {
      run(v.Command)
   }

}

func main() {
   run(".pilotssh/index")

}