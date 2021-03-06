
package main

import (
	// "os"
	// "encoding/json"
	// "strconv"
    // "io/ioutil"
    // "encoding/hex"
	// "fmt"
	"flag"
	"strings"

    "os/exec"

	riffle "github.com/exis-io/core/riffle"
)

var backend_location string
var can_interface string
var batch_size string

// fmt.Printf("port = %d", port)
func parse(str string) []string{
	//TODO: This may be able to be optimized also need to handle errors
	//Split output from candump
	pstring := strings.Split(str, " ")
	if len(pstring) < 3{
		return nil
	}
	pdata := strings.Split(pstring[2], "#")
	if len(pdata)<2{
		return nil
	}

	replacer := strings.NewReplacer("(", "", ")", "")

	//[timestamp,sid,type,data]
	ts := replacer.Replace(pstring[0])
	return []string{ts,pdata[0],pdata[1][0:2],pdata[1][2:len(pdata[1])]}
}

func parse_batch(str string) [][]string{
	message_list := strings.Split(str, "\n")
	batch := [][]string{}
	for i, _ := range message_list {
		parsed_message := parse(message_list[i])
		if parsed_message != nil {
    		batch = append(batch,parsed_message)
    	}
	}
	return batch
}

func listen_can(app riffle.Domain){
		//Heartbeat message ids
		// is_hb := map[string]bool{
  //   		"01": true,
  //   		"02": true,
  //   		"03": true,
  //   		"04": true,
  //   		"19": true,
		// }
		for {
			// Make the call to get the data we need
			//if out, err := exec.Command("python3","listen.py","can0").Output(); err != nil {
			if out, err := exec.Command("candump","-L","-n",batch_size,can_interface).Output(); err != nil {
				riffle.Error("Error %v", err)
			} else {
				str := string(out)
				send_data := parse_batch(str)
				riffle.Info("Sending %s", send_data)
				app.Publish("can", send_data)
				//TODO: add hb validator logic
				// if (is_hb[send_data[2]]){
				// 	app.Publish("hb", send_data)
				// }
			}
		//time.Sleep(5 * time.Second)
		}
}

func cmd(command string) {
	riffle.Info("Command received: %s", command)
	if _, err := exec.Command("cansend",can_interface,command).Output(); err != nil {
		riffle.Error("Error %v", err)
	} else {
		riffle.Info("Sent command to CAN: %s", command)
	}
}

// func hb(heartbeat []string) {
// 	riffle.Info("Heartbeat from nuc received: %s", heartbeat)

// }

func listen(app riffle.Domain){
        app.Listen()
}

func main() {
	// set flags for testing
	//riffle.SetFabricLocal()
	flag.StringVar(&backend_location, "l", "ws://192.168.1.99:9000", "Location of Exis backend \nDefaults to ws://192.168.1.99:9000")
	flag.StringVar(&can_interface, "i", "can1", "CAN network interface \nCan be can0, can1, can2 or can3 \nDefaults to can1 ")
	flag.StringVar(&batch_size, "b", "10", "Size of can message batches to be sent \nDefaults to 10 ")
	flag.Parse()


	// file, e := ioutil.ReadFile("./parser.json")
	// if e != nil {
	//     fmt.Printf("File error: %v\n", e)
	//     os.Exit(1)
	// }
	//fmt.Printf("%s\n", string(file))

	//m := new(Dispatch)
	//var m interface{}
	// json.Unmarshal(file, &Parser)


	riffle.SetLogLevelInfo()
	riffle.SetFabric(backend_location)
	// Domain objects
	core := riffle.NewDomain("xs")
	app := core.Subdomain("node")
	// heartbeat := app.Subdomain("hb")
	// sender := app.Subdomain("can")
	// receiver := app.Subdomain("cmd")
	// Connect
	app.Join()
	// heartbeat.Join()
	// sender.Join()

	// if e := receiver.Subscribe("cmd", cmd); e != nil {
	if e := app.Subscribe("cmd", cmd); e != nil {
		riffle.Info("Unable to subscribe to cmd endpoint", e.Error())
	} else {
		riffle.Info("Subscribed to cmd endpoint!")
	}

	// if e := heartbeat.Subscribe("hb", hb); e != nil {
	// if e := app.Subscribe("cmd", cmd); e != nil {
	// 	riffle.Info("Unable to subscribe to heartbeat endpoint ", e.Error())
	// } else {
	// 	riffle.Info("Subscribed to heartbeat!")
	// }

	go listen_can(app)

	// Handle until the connection closes
	listen(app)
	//If connection closed reopen
    defer func(){listen(app)}()
}