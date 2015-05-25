// heartbeat.go (c) 2010-2015 David Rook - all rights reserved - Released under
//    BSD 2-clause License

package main

// BUG(mdr): ServerIPstr should be gotten from command line or automatically,not hardcoded
// BUG(mdr): get ClientIP from flag/environment or system call?

import (
	// go 1.X stdlib pkgs
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	// nonlocal
	"github.com/hotei/mdr"
	"github.com/pmylund/sortutil"
)

const (
	MaxBeats = 3
)

type Client struct {
	name        string
	ipstring    string
	ipv4        net.IP
	hasReported bool
	beats       []time.Time
}

var (
	UDPPort        int
	ClientIPstr    string
	ServerIPstr    string
	ClientList     map[string]*Client
	alertThreshold int64         = 30
	alertDuration  time.Duration = time.Duration(alertThreshold) * time.Second
	// flags
	isServer       bool = false
	reportInterval int
	reportDuration time.Duration
	beatInterval   int
	beatDuration   time.Duration
	doVerbose      bool = false
)

func usage() {
	fmt.Printf("usage: heartbeat -server  -port=2345 -report=20: start server, report every 20 sec \n")
	fmt.Printf("       heartbeat -port=2345  -beat=15          : start client, beat every 15 sec \n")
}

func init() {
	flag.IntVar(&UDPPort, "port", 2667, "UDP port to use")
	flag.IntVar(&reportInterval, "report", 60*5, "seconds between reports")
	flag.IntVar(&beatInterval, "beat", 5, "seconds between beats")
	flag.BoolVar(&isServer, "server", false, "true for server mode")
	flag.BoolVar(&doVerbose, "v", false, "verbose mode")
	ClientList = make(map[string]*Client, 10)
	initClients()
}

// initClients contains a list of known (expected) clients here
//  non-expected clients can also appear but will be nameless
func initClients() {
	newClient("MARY", "10.1.2.213")
	newClient("JOHN", "10.1.2.124")
	newClient("FRED", "10.1.2.112")
	newClient("THAD", "10.1.2.115")
	newClient("WILL", "10.1.2.126")
}

func newClient(name, dotdec string) {
	var c Client
	c.ipstring = dotdec
	c.ipv4 = net.ParseIP(dotdec)
	c.name = name + " " + dotdec
	ClientList[dotdec] = &c
}

func (c *Client) dump() {
	fmt.Printf("%s\n", c.name)
}

func dumpAllClients() {
	for _, box := range ClientList {
		box.dump()
	}
}

func listener(c chan uint32) {
	udpaddr, err := net.ResolveUDPAddr("udp4", ServerIPstr)
	if err != nil {
		fmt.Printf("UDP parse error %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("UDP parse ok\n")
	data := make([]byte, 4)
	con, err := net.ListenUDP("udp4", udpaddr)
	if err != nil {
		fmt.Printf("UDP listener error %v for %v\n", err, udpaddr)
		os.Exit(1)
	}
	fmt.Printf("UDP listener started on port %d\n", UDPPort)

	for {
		n, addr, err := con.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("UDP listener connection error %v\n", err)
			os.Exit(1)
		}
		if false {
			fmt.Printf("%v %v\n", n, addr)
		}
		//fmt.Printf("pkt with %d bytes from %v = %v\n",n, addr, data)
		fromIP := mdr.ThirtyTwoNet(data[0:4])
		//fmt.Printf("rec pkt from %08x\n",fromIP)
		c <- fromIP
	}
}

func reporter() {
	reportCounter := 1
	bootTime := time.Now()
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		fmt.Printf("\nReport %d from server running since %s\n", reportCounter, bootTime.String())
		elapsed := time.Now().Sub(bootTime)
		fmt.Printf("Heartbeat server uptime is %s\n", mdr.HumanTime(elapsed))
		reportCounter++
		now := time.Now()
		var clients = []Client{}
		for _, box := range ClientList {
			clients = append(clients, *box)
		}
		sortutil.AscByField(clients, "name")
		for _, box := range clients {
			if box.hasReported == false {
				fmt.Printf("%15s has never reported itself as online \n", box.name)
				continue
			}
			if len(box.beats) < MaxBeats {
				continue
			}
			elapsed := now.Sub(box.beats[MaxBeats-1])
			fmt.Printf("%15s has a heartbeat",
				box.name)
			fmt.Printf(" : last one was %s ago\n",
				mdr.HumanTime(elapsed))
		}
	}
}

// listens for heartbeats
func server() {
	lc := make(chan uint32)
	//	rc := make(chan uint32)
	go listener(lc)
	go reporter()
	for {
		heartbeat := <-lc
		heartIP := mdr.IPFromUint32(heartbeat)
		heartName := heartIP.String()
		Verbose.Printf("HB got %08x or %s\n", heartbeat, heartName)
		now := time.Now()
		if _, exists := ClientList[heartName]; !exists {
			newClient("Unknown Client", heartName)
		}
		ClientList[heartName].hasReported = true
		mybeats := ClientList[heartName].beats
		mybeats = append(mybeats, now)
		if len(mybeats) > MaxBeats {
			mybeats = mybeats[len(mybeats)-MaxBeats:]
		}
		ClientList[heartName].beats = mybeats
	}
}

// Given typical input IP string = 10.1.2.130:47162
//    this creates the heartbeat packet we send over UDP4
func beatFromIP(IP string) []byte {
	beat := make([]byte, 4)
	pcs := strings.Split(IP, ":")      // IP on left, port on right
	part := strings.Split(pcs[0], ".") // four parts of dotted decimal addr

	i, _ := strconv.Atoi(part[0])
	beat[0] = byte(i)
	i, _ = strconv.Atoi(part[1])
	beat[1] = byte(i)
	i, _ = strconv.Atoi(part[2])
	beat[2] = byte(i)
	i, _ = strconv.Atoi(part[3])
	beat[3] = byte(i)
	return beat
}

// sends heartbeats to server
//    will send heartbeats even if server is down
func client() {
	udpaddr, err := net.ResolveUDPAddr("udp4", ServerIPstr)
	if err != nil {
		fmt.Printf("UDP parse error %v\n", err)
		os.Exit(1)
	}
	con, err := net.DialUDP("udp4", nil, udpaddr)
	if err != nil {
		fmt.Printf("Client cant connect due to UDP error %v\n", err)
		os.Exit(1)
	}
	laddr := con.LocalAddr()
	// Verbose.Printf("laddr = %s %s\n",laddr.Network(), laddr.String())
	// return should be something like: udp 10.1.2.130:46511
	beat := beatFromIP(laddr.String())
	Verbose.Printf("beat = %v\n", beat)
	for {
		Verbose.Printf("Sending heartbeat to %s\n", ServerIPstr)
		n, err := con.Write(beat)
		if err != nil {
			fmt.Printf("UDP write error %v\n", err)
		}
		if n != len(beat) {
			// should this be Verbose.Printf()?
			fmt.Printf("Can't write enough (?Server Down?)\n")
		}
		time.Sleep(beatDuration)
	}
}

func main() {
	fmt.Printf("<start HeartBeat.go>\n")
	flag.Parse()
	if doVerbose {
		Verbose = true
	}
	beatDuration = time.Duration(int64(beatInterval) * int64(time.Second))
	reportDuration = time.Duration(int64(reportInterval) * int64(time.Second))
	ClientIPstr = "127.0.0.1:" + strconv.Itoa(UDPPort)
	ServerIPstr += ("10.1.2.213:" + strconv.Itoa(UDPPort))

	fmt.Printf("The alert threshold is %d seconds\n", alertThreshold)
	fmt.Printf("Next are the defaults that can be overridden\n")
	flag.PrintDefaults()
	if doVerbose {
		dumpAllClients()
		fmt.Printf("May make sense to the client in background with nohup\n")
		fmt.Printf("Server is usually run in a window that's left visible\n")
	}

	if isServer {
		fmt.Printf("Starting server at port %d with reporting interval of %d seconds\n",
			UDPPort, reportInterval)
		server()
	} else {
		fmt.Printf("Starting client using port %d with beat interval of %d seconds\n",
			UDPPort, beatInterval)
		client()
	}
	fmt.Printf("<End HeartBeat.go>\n")
}
