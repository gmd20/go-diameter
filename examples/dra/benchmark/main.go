package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"dra"
	"dra/dcca"
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/sm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomSubscriptionId() string {
	t := rand.Intn(4)
	switch t {
	case 0:
		return "86138" + fmt.Sprintf("%09d", rand.Int63n(1000000000))
	case 1:
		return "86188" + fmt.Sprintf("%09d", rand.Int63n(1000000000))
	case 2:
		return "86159" + fmt.Sprintf("%09d", rand.Int63n(1000000000))
	case 3:
		return "86170" + fmt.Sprintf("%09d", rand.Int63n(1000000000))
	default:
		return "86189" + fmt.Sprintf("%09d", rand.Int63n(1000000000))
	}
}

type Session struct {
	mu              sync.Mutex
	numUpdates      uint32
	subscriptionId  string
	destinationHost string
	id              int   // session id
	sid             int64 // sesion id unix timestamp
	cid             int   // client id
	seq             uint32
	connection      diam.Conn
	ccrPending      bool
}

func NewSession(c diam.Conn, numUpdates int, id int, cid int) *Session {
	s := new(Session)
	s.id = id
	s.cid = cid
	s.connection = c
	s.numUpdates = uint32(numUpdates)

	return s
}

func (s *Session) SendRequest() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ccrPending == true {
		// log.Println("CCR pending, session_id=", s.id)
		// return
	}
	s.ccrPending = true

	seq := atomic.AddUint32(&s.seq, 1) - 1

	var ccr dcca.CCR
	m := diam.NewMessage(272, 0xc0, 4, 0, 0, nil)

	var et time.Time = time.Now()
	if seq == 0 { // Initial
		s.subscriptionId = RandomSubscriptionId()
		s.destinationHost = ""
		s.sid = et.UnixNano()
		ccr.OriginHost = fmt.Sprintf("peer%d.%s", s.cid, clientRealm)
		ccr.SessionId = fmt.Sprintf("%s;%d;%d", ccr.OriginHost, s.id, s.sid)
		ccr.OriginRealm = clientRealm
		ccr.DestinationRealm = serverRealm
		ccr.AuthApplicationId = 4
		ccr.ServiceContextId = "32251@3gpp.org"
		ccr.CCRequestType = 1 // Initial = 1, Update =2, Term =3, event=4
		ccr.CCRequestNumber = seq
		ccr.UserName = s.subscriptionId + "@test.test.com"
		ccr.CCSubSessionId = uint64(et.UnixNano()) // statistic time
		ccr.OriginStateId = originStateId
		ccr.EventTimestamp = &et
		ccr.RequestedAction = 0 // DIRECT_DEBITING
		ccr.SubscriptionId = []dcca.SubscriptionId{
			{
				SubscriptionIdType: 0, //END_USER_E164
				SubscriptionIdData: s.subscriptionId,
			},
			{
				SubscriptionIdType: 1, //END_USER_IMSI
				SubscriptionIdData: "2341" + s.subscriptionId,
			},
		}
		ccr.MultipleServicesIndicator = 1
		ccr.SPI = []dcca.ServiceParameterInfo{{1, "401"}, {2, "401"}}
	} else if seq < s.numUpdates { // Update
		ccr.OriginHost = fmt.Sprintf("peer%d.%s", s.cid, clientRealm)
		ccr.SessionId = fmt.Sprintf("%s;%d;%d", ccr.OriginHost, s.id, s.sid)
		ccr.OriginRealm = clientRealm
		ccr.DestinationRealm = serverRealm
		ccr.AuthApplicationId = 4
		ccr.ServiceContextId = "32251@3gpp.org"
		ccr.CCRequestType = 2 // Initial = 1, Update =2, Term =3, event=4
		ccr.CCRequestNumber = seq
		if len(s.destinationHost) > 0 {
			ccr.DestinationHost = s.destinationHost
		}
		ccr.UserName = s.subscriptionId + "@test.test.com"
		ccr.CCSubSessionId = uint64(et.UnixNano()) // statistic time
		ccr.OriginStateId = originStateId
		ccr.EventTimestamp = &et
		ccr.RequestedAction = 0 // DIRECT_DEBITING
		ccr.SubscriptionId = []dcca.SubscriptionId{
			{
				SubscriptionIdType: 0, //END_USER_E164
				SubscriptionIdData: s.subscriptionId,
			},
			{
				SubscriptionIdType: 1, //END_USER_IMSI
				SubscriptionIdData: "2341" + s.subscriptionId,
			},
		}
		ccr.MSCC = []dcca.MultipleServicesCreditControl{
			{
				RequestedServiceUnit: &dcca.RequestedServiceUnit{},
				UsedServiceUnit: []dcca.UsedServiceUnit{
					{
						CCTime:         12,
						CCTotalOctets:  540,
						CCInputOctets:  240,
						CCOutputOctets: 300,
					},
					{
						CCTime:                 1234,
						CCServiceSpecificUnits: 1,
					},
				},
				RatingGroup: 400,
			},
		}
		ccr.SPI = []dcca.ServiceParameterInfo{{1, "401"}, {2, "401"}}
	} else if seq == s.numUpdates { // Term
		ccr.OriginHost = fmt.Sprintf("peer%d.%s", s.cid, clientRealm)
		ccr.SessionId = fmt.Sprintf("%s;%d;%d", ccr.OriginHost, s.id, s.sid)
		ccr.OriginRealm = clientRealm
		ccr.DestinationRealm = serverRealm
		ccr.AuthApplicationId = 4
		ccr.ServiceContextId = "32251@3gpp.org"
		ccr.CCRequestType = 3 // Initial = 1, Update =2, Term =3, event=4
		ccr.CCRequestNumber = seq
		if len(s.destinationHost) > 0 {
			ccr.DestinationHost = s.destinationHost
		}
		ccr.UserName = s.subscriptionId + "@test.test.com"
		ccr.CCSubSessionId = uint64(et.UnixNano()) // statistic time
		ccr.OriginStateId = originStateId
		ccr.EventTimestamp = &et
		ccr.RequestedAction = 0 // DIRECT_DEBITING
		ccr.SubscriptionId = []dcca.SubscriptionId{
			{
				SubscriptionIdType: 0, //END_USER_E164
				SubscriptionIdData: s.subscriptionId,
			},
			{
				SubscriptionIdType: 1, //END_USER_IMSI
				SubscriptionIdData: "2341" + s.subscriptionId,
			},
		}
		ccr.TerminationCause = 4
		ccr.MSCC = []dcca.MultipleServicesCreditControl{
			{
				RequestedServiceUnit: &dcca.RequestedServiceUnit{},
				UsedServiceUnit: []dcca.UsedServiceUnit{
					{
						CCTime:                 1,
						CCTotalOctets:          2,
						CCInputOctets:          1,
						CCOutputOctets:         2,
						CCServiceSpecificUnits: 1,
					},
				},
				RatingGroup: 400,
			},
		}
		ccr.SPI = []dcca.ServiceParameterInfo{{1, "401"}, {2, "401"}}

		atomic.StoreUint32(&s.seq, 0)
	} else {
		atomic.StoreUint32(&s.seq, 0)
		return
	}

	// CCR + CCA's Marshal() cost about 30% CPU. To avoid the Marshal() ,
	// the message can be cached in each session, and FindAVP() can be used to
	// modify the specific AVP.
	err := m.Marshal(&ccr)
	if err != nil {
		log.Println("Marshal CCR Error: ", err)
	}

	m.WriteTo(s.connection)

	// stat.Mu.Lock()
	// stat.NumCCR++
	// stat.Mu.Unlock()
	atomic.AddUint64(&stat.NumCCR, 1)
}

func (s *Session) HandleCCA(c diam.Conn, m *diam.Message) {

	var cca dcca.CCA
	if err := m.Unmarshal(&cca); err == nil {
		s.mu.Lock()
		s.ccrPending = false
		s.destinationHost = cca.OriginHost
		s.mu.Unlock()
		t1 := cca.CCSubSessionId
		if t1 != 0 {
			t2 := uint64(time.Now().UnixNano())
			duration := (t2 - t1) / (1000 * 1000)

			atomic.AddUint64(&stat.NumCCA, 1)
			atomic.AddUint64(&stat.TotalTime, duration)
			if duration > stat.MaxTime {
				atomic.StoreUint64(&stat.MaxTime, duration)
			}
			if duration < stat.MinTime {
				atomic.StoreUint64(&stat.MinTime, duration)
			}
		} else {
			log.Printf("CCA error, CCSubSessionId doesn't exist. session_id=%v, cca.ResultCode=%d",
				cca.SessionId, cca.ResultCode)
		}
	} else {
		s.mu.Lock()
		s.ccrPending = false
		s.mu.Unlock()
		log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
	}
	SendRequests()
}

func handleCCR(server_id int) diam.HandlerFunc {
	sid := server_id
	return func(c diam.Conn, m *diam.Message) {
		var ccr dcca.CCR
		if err := m.Unmarshal(&ccr); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
		}

		rsp := m.Answer(diam.Success)
		var cca dcca.CCA
		cca.SessionId = ccr.SessionId
		cca.ResultCode = 2001
		cca.OriginHost = fmt.Sprintf("peer%d.%s", sid, serverRealm)
		cca.OriginRealm = serverRealm
		cca.CCRequestType = ccr.CCRequestType
		cca.CCRequestNumber = ccr.CCRequestNumber
		cca.CCSubSessionId = ccr.CCSubSessionId
		cca.MSCC = []dcca.MultipleServicesCreditControl{
			{
				RequestedServiceUnit: &dcca.RequestedServiceUnit{},
				ResultCode:           2001,
				RatingGroup:          400,
				GrantedServiceUnit: &dcca.GrantedServiceUnit{
					CCTotalOctets: 22345,
				},
			},
		}
		err := rsp.Marshal(&cca)
		if err != nil {
			log.Println("Marshal CCA Error: ", err)
		}
		rsp.WriteTo(c)
	}
}

type Statistics struct {
	Mu             sync.Mutex
	NumCCR         uint64
	NumCCA         uint64
	TotalTime      uint64
	MaxTime        uint64
	MinTime        uint64
	LastReportTime int64
}

func (s *Statistics) Report() (num int, rnum int, avg int, max int, min int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	num = int(s.NumCCA)
	if s.NumCCA != 0 {
		avg = int(s.TotalTime / s.NumCCA)
	}
	max = int(s.MaxTime)
	min = int(s.MinTime)

	duration := s.LastReportTime
	s.LastReportTime = time.Now().UnixNano()
	duration = (s.LastReportTime - duration) / (1000 * 1000 * 1000) // seconds
	if duration < 1 {
		duration = 1
	}
	rnum = int(s.NumCCA / uint64(duration))

	log.Printf("Statistics: duration:%d, target req/s=%d, real req/s=%d, ccr=%d, cca=%d, avg=%dms, max=%dms, min=%dms",
		duration, reqPerSec, rnum, s.NumCCR, s.NumCCA, avg, max, min)

	s.NumCCR = 0
	s.NumCCA = 0
	s.TotalTime = 0
	s.MaxTime = 0
	s.MinTime = 1000 * 1000 * 1000 * 1000

	if rnum+500 < int(reqPerSec) {
		log.Println("The real req/s it too small, it can not reach the target req/s, may be you need to increase the number of sessions and re-run the test.")
		// os.Exit(0)
	}

	return
}

var (
	originStateId uint32 = uint32(time.Now().Unix())
	peerAddr      string
	clientRealm   string
	serverRealm   string
	firstPeerId   int = 1 // avoid using the same peer id in all test
	Sessions      []*Session
	curSession    uint32
	reqPerSec     int32
	reqBoost      int32 = 1
	stat          Statistics
)

func SendRequest() {
	i := atomic.AddUint32(&curSession, 1)
	i = i % uint32(len(Sessions))
	Sessions[i].SendRequest()
	// log.Println("SendRequest session=", i)
}

func SendRequests() {
	perSec := atomic.LoadInt32(&reqPerSec)
	boost := atomic.LoadInt32(&reqBoost)
	n := atomic.LoadUint64(&stat.NumCCR)
	t := (time.Now().UnixNano() - stat.LastReportTime) / (1000 * 1000) // ms
	if t == 0 {
		t = 1
	}
	if (float64(n) / float64(t)) < float64(perSec)/1000.0 {
		boost += 4
	} else {
		boost -= 4
		atomic.StoreInt32(&reqBoost, boost)
		return
	}
	if boost <= 0 {
		boost = 1
	}
	atomic.StoreInt32(&reqBoost, boost)
	for i := 0; i < int(boost); i++ {
		SendRequest()
	}
}

func benchmark(connections, sessions, updates int) {
	f, err := os.Create(fmt.Sprintf("connections-%d sessions-%d updates-%d.txt", connections, sessions, updates))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	statFile := bufio.NewWriter(f)
	defer statFile.Flush()

	clients := make([]diam.Conn, connections)
	servers := make([]diam.Conn, connections)
	sm := make([]*sm.StateMachine, connections*2)
	Sessions = make([]*Session, sessions)

	handleCCA := func(clientId int) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var sid int
			var clientId2 int
			var sid_unix_timestamp int64

			if m.AVP[0].Code == 263 {
				sessionIdAvp := m.AVP[0].Data.(datatype.UTF8String)
				n, err := fmt.Sscanf(string(sessionIdAvp), "peer%d."+clientRealm+";%d;%d", &clientId2, &sid, &sid_unix_timestamp)
				if err != nil || n != 3 || sid < 0 || sid >= sessions || Sessions[sid] == nil {
					log.Println("Failed to extract sessionId from dcca CCA. ", m)
					return
				}
			} else {
				log.Println("Failed to extract SessionId from dcca CCA. ", m)
				return
			}

			Sessions[sid].HandleCCA(c, m)
		}
	}

	for i := 0; i < connections; i++ {
		clientHost := fmt.Sprintf("peer%d.%s", firstPeerId+i, clientRealm)
		serverHost := fmt.Sprintf("peer%d.%s", firstPeerId+i, serverRealm)
		cmux, c := dra.Client(peerAddr, clientHost, clientRealm, originStateId)
		smux, s := dra.Client(peerAddr, serverHost, serverRealm, originStateId)
		defer c.Close()
		defer s.Close()
		sm[i] = cmux
		sm[connections+i] = smux

		clients[i] = c
		servers[i] = s
		cmux.HandleFunc("CCA", handleCCA(firstPeerId+i))
		smux.Handle("CCR", handleCCR(firstPeerId+i))
	}

	for i := 0; i < sessions; i++ {
		s := NewSession(clients[i%connections], updates, i, firstPeerId+i%connections)
		Sessions[i] = s
	}

	firstPeerId += connections

	time.Sleep(10 * time.Second) // wait for CEA ready?
	stat.LastReportTime = time.Now().UnixNano()
	atomic.StoreUint32(&curSession, 0)
	atomic.StoreInt32(&reqBoost, 1)
	if sessions > 1600 {
		atomic.StoreInt32(&reqPerSec, int32(sessions/20)) // avoid session timeout
	} else if sessions > 100 {
		atomic.StoreInt32(&reqPerSec, 100)
	} else {
		atomic.StoreInt32(&reqPerSec, 10)
	}

	lastRate := 0
	lastAvgTime := 0
	var lastStep int32 = 10
	iskneePoint := false

	tick := 0
	SendRequests()
	for {
		select {
		case <-time.After(200 * time.Millisecond):
			SendRequests()
			tick++
			if tick%100 != 0 {
				break
			}

			_, rnum, a, _, _ := stat.Report()

			if !iskneePoint && ((rnum < lastRate && a > lastAvgTime) || a-lastAvgTime > 100) {
				iskneePoint = true
				// retest on knee point using small steps
				time.Sleep(8 * time.Second)
				stat.Report()
				reqPerSec -= lastStep
				break
			} else {
				statFile.WriteString(fmt.Sprintf("%d %d\n", rnum, a))
				statFile.Flush()
			}

			if a > 500 {
				return
			}
			if iskneePoint {
				lastStep = 50
			} else if rnum < 4000 {
				lastStep = 800
			} else if a-lastAvgTime < 20 {
				lastStep = 400
			} else {
				lastStep = 200
			}
			reqPerSec += lastStep
			lastRate = rnum
			lastAvgTime = a
		}

		for i := 0; i < connections*2; {
			select {
			case err := <-sm[i].ErrorReports():
				log.Println("StateMachine Error: ", err)
			default:
				i++
			}
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	flag.StringVar(&peerAddr, "peer_addr", "192.168.4.231:3868", "address in the form of ip:port to connect")
	flag.StringVar(&clientRealm, "client_realm", "localdomain.net", "diameter client realm")
	flag.StringVar(&serverRealm, "server_realm", "localdomain2.net", "diameter server realm")

	connections := flag.Int("connections", 0, "number of clients and servers.")
	sessions := flag.Int("sessions", 0, "number of sessions per connection.")
	updates := flag.Int("updates", 0, "number of Update CCR.")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *connections == 0 {
		for _, c := range []int{1, 8, 16, 32} {
			for _, s := range []int{32, 400, 3200, 16000} {
				for _, u := range []int{100} {
					log.Printf("%d-%d-%d\n", c, s, u)
					benchmark(c, s, u)
					log.Printf("%d-%d-%d finished\n\n\n\n", c, s, u)
					time.Sleep(60 * time.Second) // wait for freeDiameter peer state ready?
				}
			}
		}
	} else {
		log.Printf("%d-%d-%d", *connections, *sessions, *updates)
		benchmark(*connections, *sessions, *updates)
	}
}
