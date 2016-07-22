package dra

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"dra/dcca"
	"dra/gx"
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	peerAddr         string // remote ip:port to connect
	localAddr        string // local ip:port to listen on
	originHost       string
	originRealm      string
	destinationHost  string
	destinationRealm string
)

type TcpDump struct {
	process *os.Process
}

func (t *TcpDump) Start(name string) error {
	args := []string{"-i", "any", "-s0", "sctp", "port", "3868", "-w", name + ".pcap"}
	cmd := exec.Command("tcpdump", args...)
	// cmd.Dir = "/tmp" // work dir
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	t.process = cmd.Process
	return nil
}

func (t *TcpDump) Stop() {
	if t.process != nil {
		t.process.Signal(os.Interrupt)
	}
}
func TestBasicDcca(t *testing.T) {
	// dict.Default.LoadFile("credit_control.xml")  load by default
	// err := dict.Default.Load(bytes.NewReader([]byte(DccaXML)))
	// if err != nil {
	//	t.Fatal(err)
	// }
	var err error

	tcpdump := &TcpDump{}
	tcpdump.Start("basic_dcca")
	defer tcpdump.Stop()

	originStateId := uint32(time.Now().Unix())

	cmux, c := Client(peerAddr, originHost, originRealm, originStateId)
	// smux, s := server(":3868", "peer2.localdomain2", "localdomain2")
	smux, s := Client(peerAddr, destinationHost, destinationRealm, originStateId)
	defer c.Close()
	defer s.Close()

	time.Sleep(2 * time.Second)

	sendCCR := func(c diam.Conn) (n int64, err error) {
		// Build CCR
		// m := diam.NewRequest(272, 4, nil)
		m := diam.NewMessage(272, 0xc0, 4, 0, 0, nil)
		var ccr dcca.CCR
		// Add AVPs
		ccr.SessionId = originHost + ";25020007;1798;523100ae-602"
		ccr.OriginHost = originHost
		ccr.OriginRealm = originRealm
		ccr.DestinationRealm = destinationRealm
		ccr.AuthApplicationId = 4
		ccr.ServiceContextId = "32251@3gpp.org"
		ccr.CCRequestType = 4 // Initial = 1, Update =2, Term =3, event=4
		ccr.CCRequestNumber = 0
		ccr.UserName = "45678971@test.test.com"
		// ccr.ServiceIdentifier = 0
		ccr.OriginStateId = originStateId
		var et time.Time = time.Now()
		ccr.EventTimestamp = &et
		ccr.RequestedAction = 0 // DIRECT_DEBITING
		ccr.SubscriptionId = []dcca.SubscriptionId{
			{
				SubscriptionIdType: 0, //END_USER_E164
				SubscriptionIdData: "45678971",
			},
			{
				SubscriptionIdType: 1, //END_USER_IMSI
				SubscriptionIdData: "23611145678971",
			},
		}
		ccr.MultipleServicesIndicator = 1
		ccr.SPI = []dcca.ServiceParameterInfo{{1, "401"}, {2, "401"}}
		// ccr.RouteRecord = []string{"peer1.localdomain"}
		err = m.Marshal(&ccr)
		if err != nil {
			log.Print(err)
			return 0, err
		}
		return m.WriteTo(c)
	}
	handleCCR := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var ccr dcca.CCR
			if err := m.Unmarshal(&ccr); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode CCR")
			}

			rsp := m.Answer(diam.Success)
			var cca dcca.CCA
			cca.SessionId = ccr.SessionId
			cca.ResultCode = 2001
			cca.OriginHost = ""
			cca.OriginRealm = ""

			cca.CCRequestType = ccr.CCRequestType
			cca.CCRequestNumber = ccr.CCRequestNumber

			rsp.Marshal(&cca)
			rsp.WriteTo(c)
			// c.Close()
		}
	}
	handleCCA := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var cca dcca.CCA
			if err := m.Unmarshal(&cca); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode CCA")
			}
			if cca.ResultCode != 2001 {
				log.Printf("cca.ResultCode = %v\n", cca.ResultCode)
				errChan <- fmt.Errorf("Unexpected cca.ResultCode")
			} else {
				close(okChan) // test case pass
			}
		}
	}

	errChan := make(chan error, 1)
	okChan := make(chan struct{})
	cmux.Handle("CCA", handleCCA(errChan, okChan))
	smux.Handle("CCR", handleCCR(errChan, okChan))

	_, err = sendCCR(c)
	if err != nil {
		t.Fatal("Failed to send CCR", err)
	}
	select {
	case <-okChan: // test case pass
	case err := <-errChan:
		t.Fatal(err)
	case err := <-cmux.ErrorReports():
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out: no CCR or CCA received")
	}
}

// 7.2 DRA Definition
// the Sd interface for a certain IP-CAN session reach the same PCRF when
// multiple and separately addressable PCRFs have been deployed in a Diameter realm
func TestSdInterface(t *testing.T) {
}

//  The DRA also ensures that the NRR commands over Np interface reach the
// same PCRF for a certain IP-CAN session
func TestNpInterface(t *testing.T) {

}

// The DRA (Diameter Routing Agent) is a functional element that ensures that
// all Diameter sessions established over the Gx, S9, Gxx, Rx interfaces and
// for unsolicited application reporting
func TestGxInterface(t *testing.T) {
	err := dict.Default.Load(bytes.NewReader([]byte(gx.GxXML)))
	if err != nil {
		t.Fatal(err)
	}

	tcpdump := &TcpDump{}
	tcpdump.Start("GxInterface")
	defer tcpdump.Stop()

	originStateId := uint32(time.Now().Unix())

	cmux, c := AppClient("192.168.4.231:3868", "peer-gx-client.localdomain.net", "localdomain.net",
		originStateId, gx.GX_APPLICATION_ID)
	// smux, s := server(":3868", "peer2.localdomain2", "localdomain2")
	smux, s := AppClient("192.168.4.231:3868", "peer-gx-server.localdomain2.net", "localdomain2.net",
		originStateId, gx.GX_APPLICATION_ID)
	defer c.Close()
	defer s.Close()
	time.Sleep(2 * time.Second) // wait for freeDiameter peer state

	var requestNumber uint32 = 0

	newInt32 := func(v int32) *int32 {
		i := new(int32)
		*i = v
		return i
	}
	sendCCR := func(c diam.Conn, requestType int32) (n int64, err error) {
		// Build CCR
		// m := diam.NewRequest(272, 4, nil)
		m := diam.NewMessage(272, 0xc0, gx.GX_APPLICATION_ID, 0, 0, nil)
		var ccr gx.CCR
		// Add AVPs
		ccr.SessionId = "peer-gx-client.localdomain.net;25020007;1798;192.168.4.85"
		ccr.AuthApplicationId = gx.GX_APPLICATION_ID
		ccr.OriginHost = "peer-gx-client.localdomain.net"
		ccr.OriginRealm = "localdomain.net"
		ccr.DestinationRealm = "localdomain2.net"
		ccr.CCRequestType = requestType // Initial = 1, Update =2, Term =3, event=4
		ccr.CCRequestNumber = requestNumber
		requestNumber++
		ccr.OriginStateId = originStateId
		ccr.SubscriptionId = []gx.SubscriptionId{
			{
				SubscriptionIdType: 0, //END_USER_E164
				SubscriptionIdData: "13800138000",
			},
			// {
			//	SubscriptionIdType: 1, //END_USER_IMSI
			//	SubscriptionIdData: "2123413800138000",
			// },
		}
		ccr.FramedIPAddress = "\xc0\xa8\x05\x55"
		ccr.UserEquipmentInfo = &gx.UserEquipmentInfo{
			UserEquipmentInfoType:  0,
			UserEquipmentInfoValue: "\x33\x35\x32\x36\x34\x38\x30\x35\x37\x38\x36\x39\x35\x38\x30\x31",
		}
		ccr.CalledStationID = "test.net"
		ancAddr := net.ParseIP("192.168.4.85")
		ccr.AccessNetworkChargingAddress = &ancAddr
		ccr.AccessNetworkChargingIdentifierGx = &gx.AccessNetworkChargingIdentifierGx{
			AccessNetworkChargingIdentifierValue: "\x73\x00\x03\x40",
		}

		if requestType == 3 {
			ccr.TerminationCause = newInt32(1)
		} else {
			ccr.SupportedFeatures = []gx.SupportedFeatures{
				{
					VendorID:      10415,
					FeatureListID: 1,
					FeatureList:   11,
				},
			}
			var notSupport int32 = 0
			ccr.NetworkRequestSupport = &notSupport
			ccr.IPCANType = newInt32(0)
			ccr.RATType = newInt32(1000)
			ccr.QoSInformation = &gx.QoSInformation{
				MaxRequestedBandwidthUL: 32000,
				MaxRequestedBandwidthDL: 32000,
			}
			ccr.QoSNegotiation = newInt32(1)
			ccr.QoSUpgrade = newInt32(1)
			ccr.TGPPSGSNAddress = "\xc0\xa8\x04\xd9"
			ccr.TGPPUserLocationInfo = "\x00\x64\xf6\x79\x00\x01\xea\x6c"
			ccr.BearerUsage = newInt32(0)
		}

		// ccr.RouteRecord = []string{"peer1.localdomain.net"}
		err = m.Marshal(&ccr)
		if err != nil {
			log.Print(err)
			return 0, err
		}
		return m.WriteTo(c)
	}
	handleCCR := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var ccr gx.CCR
			if err := m.Unmarshal(&ccr); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode CCR")
			}

			rsp := m.Answer(diam.Success)
			var cca gx.CCA
			cca.SessionId = ccr.SessionId
			cca.AuthApplicationId = gx.GX_APPLICATION_ID
			cca.OriginHost = "peer-gx-server.localdomain2.net"
			cca.OriginRealm = "localdomain2.net"
			cca.ResultCode = 2001
			cca.CCRequestType = ccr.CCRequestType
			cca.CCRequestNumber = ccr.CCRequestNumber
			cca.OriginStateId = originStateId

			if cca.CCRequestType != 3 {
				cca.BearerControlMode = newInt32(1)
				cca.EventTrigger = []int32{1, 2, 0, 33}
				cca.ChargingRuleInstall = &gx.ChargingRuleInstall{
					ChargingRuleName: "100",
				}
				cca.QoSInformation = &gx.QoSInformation{
					MaxRequestedBandwidthUL: 153600000,
					MaxRequestedBandwidthDL: 153600000,
				}
				cca.SupportedFeatures = []gx.SupportedFeatures{
					{
						VendorID:      10415,
						FeatureListID: 1,
						FeatureList:   11,
					},
				}
			}

			rsp.Marshal(&cca)
			rsp.WriteTo(c)
			// c.Close()
		}
	}
	handleCCA := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var cca gx.CCA
			if err := m.Unmarshal(&cca); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode CCA")
			}

			if cca.ResultCode != 2001 {
				log.Printf("cca.ResultCode = %v\n", cca.ResultCode)
				errChan <- fmt.Errorf("Unexpected cca.ResultCode")
			} else if cca.CCRequestType == 3 {
				close(okChan) // test case pass
			}
		}
	}

	errChan := make(chan error, 1)
	okChan := make(chan struct{})
	cmux.Handle("CCA", handleCCA(errChan, okChan))
	smux.Handle("CCR", handleCCR(errChan, okChan))

	_, err = sendCCR(c, 1) // CCR-Innitial
	if err != nil {
		t.Fatal("Failed to send CCR", err)
	}
	time.Sleep(2 * time.Second) // wait for freeDiameter peer state
	_, err = sendCCR(c, 3)      // CCR-Terminiation
	if err != nil {
		t.Fatal("Failed to send CCR", err)
	}
	select {
	case <-okChan: // test case pass
	case err := <-errChan:
		t.Fatal(err)
	case err := <-cmux.ErrorReports():
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out: no CCR or CCA received")
	}
}

func TestGxRAR(t *testing.T) {
	err := dict.Default.Load(bytes.NewReader([]byte(gx.GxXML)))
	if err != nil {
		t.Fatal(err)
	}

	tcpdump := &TcpDump{}
	tcpdump.Start("GxRAR")
	defer tcpdump.Stop()

	originStateId := uint32(time.Now().Unix())

	cmux, c := AppClient("192.168.4.231:3868", "peer-gx-client-rar.localdomain.net", "localdomain.net",
		originStateId, gx.GX_APPLICATION_ID)
	// smux, s := server(":3868", "peer2.localdomain2", "localdomain2")
	smux, s := AppClient("192.168.4.231:3868", "peer-gx-server-rar.localdomain2.net", "localdomain2.net",
		originStateId, gx.GX_APPLICATION_ID)
	defer c.Close()
	defer s.Close()
	time.Sleep(2 * time.Second) // wait for freeDiameter peer state

	sendRAR := func(c diam.Conn, requestType int32) (n int64, err error) {
		// Build CCR
		// m := diam.NewRequest(272, 4, nil)
		m := diam.NewMessage(258, 0xc0, gx.GX_APPLICATION_ID, 0, 0, nil)
		var rar gx.RAR
		// Add AVPs
		rar.SessionId = "peer-gx-server-rar.localdomain.net;25020007;1798;192.168.4.231"
		rar.AuthApplicationId = gx.GX_APPLICATION_ID
		rar.OriginHost = "peer-gx-server-rar.localdomain2.net"
		rar.OriginRealm = "localdomain2.net"
		rar.DestinationRealm = "localdomain.net"
		rar.DestinationHost = "peer-gx-client-rar.localdomain.net"
		rar.OriginStateId = originStateId

		rar.EventTrigger = []int32{1, 2}
		if requestType == 0 {
			rar.ChargingRuleInstall = &gx.ChargingRuleInstall{
				ChargingRuleName: "100",
			}
		} else {
			rar.ChargingRuleRemove = &gx.ChargingRuleRemove{
				ChargingRuleName: "100",
			}
		}
		rar.QoSInformation = &gx.QoSInformation{
			MaxRequestedBandwidthUL: 32000,
			MaxRequestedBandwidthDL: 32000,
		}
		rar.RevalidationTime = &time.Time{}
		*rar.RevalidationTime = time.Now().Add(30 * time.Second)

		// rar.RouteRecord = []string{"peer1.localdomain.net"}
		err = m.Marshal(&rar)
		if err != nil {
			log.Print(err)
			return 0, err
		}
		return m.WriteTo(c)
	}
	handleRAR := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var rar gx.RAR
			if err := m.Unmarshal(&rar); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode RAR")
			}

			rsp := m.Answer(diam.Success)
			var raa gx.RAA
			raa.SessionId = rar.SessionId
			raa.AuthApplicationId = gx.GX_APPLICATION_ID
			raa.OriginHost = "peer-gx-client-rar.localdomain2.net"
			raa.OriginRealm = "localdomain.net"
			if rar.ChargingRuleInstall != nil {
				raa.ResultCode = 2001
			} else {
				raa.ResultCode = 2002
				raa.FailedAVP = []byte{0x00, 0x00, 0x03, 0xf8, 0xc0, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x28, 0xaf, 0x00, 0x00, 0x02, 0x04, 0xc0, 0x00, 0x00, 0x10, 0x00, 0x00, 0x28, 0xaf, 0x00, 0x00, 0x7d, 0x00, 0x00, 0x00, 0x02, 0x03, 0xc0, 0x00, 0x00, 0x10, 0x00, 0x00, 0x28, 0xaf, 0x00, 0x00, 0x7d, 0x00}
			}
			raa.OriginStateId = originStateId
			ancAddr := net.ParseIP("192.168.4.85")
			raa.AccessNetworkChargingAddress = &ancAddr
			raa.AccessNetworkChargingIdentifierGx = &gx.AccessNetworkChargingIdentifierGx{
				AccessNetworkChargingIdentifierValue: "\x73\x00\x03\x40",
			}

			rsp.Marshal(&raa)
			rsp.WriteTo(c)
			// c.Close()
		}
	}
	handleRAA := func(errChan chan error, okChan chan struct{}) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var raa gx.RAA
			if err := m.Unmarshal(&raa); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				errChan <- fmt.Errorf("Failed to decode RAA")
			}

			if raa.ResultCode == 2001 {
			} else if raa.ResultCode == 2002 {
				close(okChan) // test case pass
			} else {
				log.Printf("raa.ResultCode = %v\n", raa.ResultCode)
				errChan <- fmt.Errorf("Unexpected raa.ResultCode")
			}
		}
	}

	errChan := make(chan error, 1)
	okChan := make(chan struct{})
	cmux.Handle("RAR", handleRAR(errChan, okChan))
	smux.Handle("RAA", handleRAA(errChan, okChan))

	_, err = sendRAR(s, 0)
	if err != nil {
		t.Fatal("Failed to send RAR", err)
	}
	time.Sleep(1 * time.Second)
	_, err = sendRAR(s, 1)
	if err != nil {
		t.Fatal("Failed to send RAR", err)
	}
	select {
	case <-okChan: // test case pass
	case err := <-errChan:
		t.Fatal(err)
	case err := <-cmux.ErrorReports():
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out: no RAR or RAR received")
	}
}

// Route-Record and Proxy-Info
func TestRouteRecord(t *testing.T) {
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

func newInt32(v int32) *int32 {
	i := new(int32)
	*i = v
	return i
}

type GxSession struct {
	OriginStateId    uint32
	ClientConnection diam.Conn
	ClientId         int
	SeqNo            uint32
	Sid              int
	SessionId        string
	SubscriptionId   string
	DestinationHost  string
	ServerId         int
	WrongServerId    int // the CCR is sent to a wrong server
	WrongClientId    int // the CCA is sent to a wrong client
	CCAError         int
}

func NewGxSession(c diam.Conn, cid int, sid int, originStateId uint32) *GxSession {
	s := new(GxSession)
	s.OriginStateId = originStateId
	s.ClientConnection = c
	s.ClientId = cid
	s.Sid = sid
	s.Init()
	return s
}

func (s *GxSession) Init() {
	s.SeqNo = 0
	s.SubscriptionId = RandomSubscriptionId()
	s.DestinationHost = ""
	s.ServerId = -1
	s.WrongServerId = -1
	s.WrongClientId = -1
}

func (s *GxSession) SendRequest() {
	s.CCAError++

	m := diam.NewMessage(272, 0xc0, gx.GX_APPLICATION_ID, 0, 0, nil)
	var ccr gx.CCR
	ccr.AuthApplicationId = gx.GX_APPLICATION_ID
	ccr.OriginHost = fmt.Sprintf("peer%d-gx-client.%s", s.ClientId, originRealm)
	ccr.OriginRealm = originRealm
	ccr.DestinationRealm = destinationRealm
	if s.SeqNo%2 == 0 {
		s.SessionId = fmt.Sprintf("%v;%d;%d", ccr.OriginHost, s.SeqNo, s.Sid)
		s.DestinationHost = ""
		ccr.CCRequestType = 1 // Initial = 1, Update =2, Term =3, event=4
	} else {
		ccr.CCRequestType = 3 // Initial = 1, Update =2, Term =3, event=4
		ccr.DestinationHost = s.DestinationHost
	}
	ccr.SessionId = s.SessionId
	ccr.CCRequestNumber = s.SeqNo % 2
	s.SeqNo++
	ccr.OriginStateId = s.OriginStateId
	ccr.SubscriptionId = []gx.SubscriptionId{
		{
			SubscriptionIdType: 0, //END_USER_E164
			SubscriptionIdData: s.SubscriptionId,
		},
		{
			SubscriptionIdType: 1, //END_USER_IMSI
			SubscriptionIdData: "21234" + s.SubscriptionId,
		},
	}
	ccr.FramedIPAddress = "\xc0\xa8\x05\x55"
	ccr.UserEquipmentInfo = &gx.UserEquipmentInfo{
		UserEquipmentInfoType:  0,
		UserEquipmentInfoValue: "\x33\x35\x32\x36\x34\x38\x30\x35\x37\x38\x36\x39\x35\x38\x30\x31",
	}
	ccr.CalledStationID = "test.net"
	ancAddr := net.ParseIP("192.168.4.85")
	ccr.AccessNetworkChargingAddress = &ancAddr
	ccr.AccessNetworkChargingIdentifierGx = &gx.AccessNetworkChargingIdentifierGx{
		AccessNetworkChargingIdentifierValue: "\x73\x00\x03\x40",
	}

	if ccr.CCRequestType == 3 {
		ccr.TerminationCause = newInt32(1)
	} else {
		ccr.SupportedFeatures = []gx.SupportedFeatures{
			{
				VendorID:      10415,
				FeatureListID: 1,
				FeatureList:   11,
			},
		}
		var notSupport int32 = 0
		ccr.NetworkRequestSupport = &notSupport
		ccr.IPCANType = newInt32(0)
		ccr.RATType = newInt32(1000)
		ccr.QoSInformation = &gx.QoSInformation{
			MaxRequestedBandwidthUL: 32000,
			MaxRequestedBandwidthDL: 32000,
		}
		ccr.QoSNegotiation = newInt32(1)
		ccr.QoSUpgrade = newInt32(1)
		ccr.TGPPSGSNAddress = "\xc0\xa8\x04\xd9"
		ccr.TGPPUserLocationInfo = "\x00\x64\xf6\x79\x00\x01\xea\x6c"
		ccr.BearerUsage = newInt32(0)
	}

	// ccr.RouteRecord = []string{"peer1.localdomain.net"}
	err := m.Marshal(&ccr)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = m.WriteTo(s.ClientConnection)
	if err != nil {
		log.Print(err)
	}
}

func (s *GxSession) HandleCCA(clientId int, serverId int, m *diam.Message) {
	s.CCAError--
	var cca gx.CCA
	if err := m.Unmarshal(&cca); err != nil {
		s.CCAError++
		log.Printf("Failed to parse CCA message from server %v: %s\n%s", serverId, err, m)
		return
	}

	if cca.ResultCode != 2001 {
		s.CCAError++
		log.Printf("cca.ResultCode = %v\n", cca.ResultCode)
		return
	} else if cca.CCRequestType == 3 {
	}

	s.DestinationHost = cca.OriginHost

	if clientId != s.ClientId {
		s.CCAError++
		s.WrongClientId = clientId
		log.Println("CCA is not sent back to the right client.",
			"ClientId=", s.ClientId,
			"SeqNo=", s.SeqNo,
			"Sid=", s.Sid,
			"SessionId=", s.SessionId,
			"SubscriptionId=", s.SubscriptionId,
			"DestinationHost=", s.DestinationHost,
			"ServerId=", s.ServerId,
			"WrongServerId=", s.WrongServerId,
			"WrongClientId=", s.WrongClientId)
	}

	if s.ServerId >= 0 && serverId != s.ServerId {
		s.CCAError++
		s.WrongServerId = serverId
		log.Println("CCR is sent to differetnt server.",
			"ClientId=", s.ClientId,
			"SeqNo=", s.SeqNo,
			"Sid=", s.Sid,
			"SessionId=", s.SessionId,
			"SubscriptionId=", s.SubscriptionId,
			"DestinationHost=", s.DestinationHost,
			"ServerId=", s.ServerId,
			"WrongServerId=", s.WrongServerId,
			"WrongClientId=", s.WrongClientId)
	}
	s.ServerId = serverId
}

type GySession struct {
	OriginStateId    uint32
	ClientConnection diam.Conn
	ClientId         int
	SeqNo            uint32
	Sid              int
	SessionId        string
	SubscriptionId   string
	DestinationHost  string
	ServerId         int
	WrongServerId    int // the CCR is sent to a wrong server
	WrongClientId    int // the CCA is sent to a wrong client
	CCAError         int
}

func NewGySession(c diam.Conn, cid int, sid int, originStateId uint32) *GySession {
	s := new(GySession)
	s.OriginStateId = originStateId
	s.ClientConnection = c
	s.ClientId = cid
	s.Sid = sid
	s.SeqNo = 0
	s.SubscriptionId = RandomSubscriptionId()
	s.ServerId = -1
	s.WrongServerId = -1
	s.WrongClientId = -1
	return s
}

func (s *GySession) SendRequest() {
	s.CCAError++

	m := diam.NewMessage(272, 0xc0, dcca.DCCA_APPLICATION_ID, 0, 0, nil)
	var ccr dcca.CCR
	ccr.OriginHost = fmt.Sprintf("peer%d-gy-client.%s", s.ClientId, originRealm)
	ccr.OriginRealm = originRealm
	ccr.DestinationRealm = destinationRealm
	ccr.AuthApplicationId = 4
	s.SessionId = fmt.Sprintf("%v;%d;%d", ccr.OriginHost, s.SeqNo, s.Sid)
	ccr.SessionId = s.SessionId
	ccr.ServiceContextId = "32251@3gpp.org"
	ccr.CCRequestType = 4 // Initial = 1, Update =2, Term =3, event=4
	ccr.CCRequestNumber = 0
	s.SeqNo++
	ccr.UserName = s.SubscriptionId + "@test.com"
	// ccr.ServiceIdentifier = 0
	ccr.OriginStateId = s.OriginStateId
	var et time.Time = time.Now()
	ccr.EventTimestamp = &et
	ccr.RequestedAction = 0 // DIRECT_DEBITING
	ccr.SubscriptionId = []dcca.SubscriptionId{
		{
			SubscriptionIdType: 0, //END_USER_E164
			SubscriptionIdData: s.SubscriptionId,
		},
		{
			SubscriptionIdType: 1, //END_USER_IMSI
			SubscriptionIdData: "236111" + s.SubscriptionId,
		},
	}
	ccr.MultipleServicesIndicator = 1
	ccr.SPI = []dcca.ServiceParameterInfo{{1, "401"}, {2, "401"}}
	// ccr.RouteRecord = []string{"peer1.localdomain"}
	err := m.Marshal(&ccr)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = m.WriteTo(s.ClientConnection)
	if err != nil {
		log.Print(err)
	}
}

func (s *GySession) HandleCCA(clientId int, serverId int, m *diam.Message) {
	s.CCAError--
	var cca dcca.CCA
	if err := m.Unmarshal(&cca); err != nil {
		s.CCAError++
		log.Printf("Failed to parse CCA message from server %v: %s\n%s", serverId, err, m)
		return
	}

	if cca.ResultCode != 2001 {
		s.CCAError++
		log.Printf("cca.ResultCode = %v\n", cca.ResultCode)
		return
	} else if cca.CCRequestType == 3 {
	}

	s.DestinationHost = cca.OriginHost

	if clientId != s.ClientId {
		s.CCAError++
		s.WrongClientId = clientId
		log.Println("CCA is not sent back to the right client.",
			"ClientId=", s.ClientId,
			"SeqNo=", s.SeqNo,
			"Sid=", s.Sid,
			"SessionId=", s.SessionId,
			"SubscriptionId=", s.SubscriptionId,
			"DestinationHost=", s.DestinationHost,
			"ServerId=", s.ServerId,
			"WrongServerId=", s.WrongServerId,
			"WrongClientId=", s.WrongClientId)
	}

	if s.ServerId >= 0 && serverId != s.ServerId {
		s.CCAError++
		s.WrongServerId = serverId
		log.Println("CCR is not sent to the same server.",
			"ClientId=", s.ClientId,
			"SeqNo=", s.SeqNo,
			"Sid=", s.Sid,
			"SessionId=", s.SessionId,
			"SubscriptionId=", s.SubscriptionId,
			"DestinationHost=", s.DestinationHost,
			"ServerId=", s.ServerId,
			"WrongServerId=", s.WrongServerId,
			"WrongClientId=", s.WrongClientId)
	}
	s.ServerId = serverId
}

func sum(numbers []int) (total float64) {
	for _, x := range numbers {
		total += float64(x)
	}
	return total
}

func mean(numbers []int) (mean float64) {
	mean = sum(numbers) / float64(len(numbers))
	return mean
}

// standard deviation
func stdDev(numbers []int, mean float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += math.Pow(float64(number)-mean, 2)
	}
	variance := total / float64(len(numbers)-1)
	return math.Sqrt(variance)
}

// 3GPP TS 29.213 versio n 13.4.0 Release 13
// 7.3.2 DRA Information Storage
// The DRA shall maintain PCRF routing information per IP-CAN session or per
// UE-NAI, depending on the operator"s configuration.  The DRA shall select
// the same PCRF for all the Diameter sessions established for the same UE
// in case 2a.
// The DRA has information about the user identity (UE NAI), the UE Ipv4 address
// and/or Ipv6 prefix, the APN (if available), the PCEF identity (if available)
// and the selected PCRF identity for a certain IP-CAN Session.
func TestPCRFSelection(t *testing.T) {
	err := dict.Default.Load(bytes.NewReader([]byte(gx.GxXML)))
	if err != nil {
		t.Fatal(err)
	}
	tcpdump := &TcpDump{}
	tcpdump.Start("PCRFSelection")
	defer tcpdump.Stop()

	var counterMutex sync.Mutex
	var gxCounter [16]int
	var gyCounter [16]int

	resetCounter := func() {
		for i := 0; i < 16; i++ {
			gxCounter[i] = 0
			gyCounter[i] = 0
		}
	}

	originStateId := uint32(time.Now().Unix())
	gxClients := make([]diam.Conn, 16)
	gxServers := make([]diam.Conn, 16)
	gxClientSm := make([]*sm.StateMachine, 16)
	gxServerSm := make([]*sm.StateMachine, 16)
	const numSessions int = 4000
	gxSessions := make([]*GxSession, numSessions)

	gyClients := make([]diam.Conn, 16)
	gyServers := make([]diam.Conn, 16)
	gyClientSm := make([]*sm.StateMachine, 16)
	gyServerSm := make([]*sm.StateMachine, 16)
	gySessions := make([]*GySession, numSessions)

	handleGxCCR := func(serverId int) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			counterMutex.Lock()
			gxCounter[serverId]++
			counterMutex.Unlock()

			var ccr gx.CCR
			if err := m.Unmarshal(&ccr); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				return
			}

			rsp := m.Answer(diam.Success)
			var cca gx.CCA
			cca.SessionId = ccr.SessionId
			cca.AuthApplicationId = gx.GX_APPLICATION_ID
			cca.OriginHost = fmt.Sprintf("peer%d-gx-server.%s", serverId, destinationRealm)
			cca.OriginRealm = destinationRealm
			cca.ResultCode = 2001
			cca.CCRequestType = ccr.CCRequestType
			cca.CCRequestNumber = ccr.CCRequestNumber
			cca.OriginStateId = originStateId

			if cca.CCRequestType != 3 {
				cca.BearerControlMode = newInt32(1)
				cca.EventTrigger = []int32{1, 2, 0, 33}
				cca.ChargingRuleInstall = &gx.ChargingRuleInstall{
					ChargingRuleName: "100",
				}
				cca.QoSInformation = &gx.QoSInformation{
					MaxRequestedBandwidthUL: 153600000,
					MaxRequestedBandwidthDL: 153600000,
				}
				cca.SupportedFeatures = []gx.SupportedFeatures{
					{
						VendorID:      10415,
						FeatureListID: 1,
						FeatureList:   11,
					},
				}
			}

			rsp.Marshal(&cca)
			rsp.WriteTo(c)
			// c.Close()
		}
	}
	handleGxCCA := func(clientId int) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var sid int
			var clientId2 int
			var seqNo int
			var serverId int

			if m.AVP[0].Code == 263 {
				sessionIdAvp := m.AVP[0].Data.(datatype.UTF8String)
				n, err := fmt.Sscanf(string(sessionIdAvp), "peer%d-gx-client."+originRealm+";%d;%d", &clientId2, &seqNo, &sid)
				if err != nil || n != 3 || sid < 0 || sid >= numSessions || gxSessions[sid] == nil {
					log.Println("Failed to extract sessionId from Gx CCA. ", m)
					return
				}
			} else {
				log.Println("Failed to extract SessionId from Gx CCA. ", m)
				return
			}

			if m.AVP[2].Code == 264 {
				originHostAvp := m.AVP[2].Data.(datatype.DiameterIdentity)
				n, err := fmt.Sscanf(string(originHostAvp), "peer%d-gx-server."+destinationRealm, &serverId)
				if err != nil || n != 1 || serverId < 0 || gxServers[serverId] == nil {
					log.Println("Failed to extract serverId from Gx CCA. ", m)
					return
				}
			} else {
				log.Println("Failed to extract serverId from Gx CCA. ", m)
				return
			}

			gxSessions[sid].HandleCCA(clientId, serverId, m)
		}
	}

	handleGyCCR := func(serverId int) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			counterMutex.Lock()
			gyCounter[serverId]++
			counterMutex.Unlock()

			var ccr dcca.CCR
			if err := m.Unmarshal(&ccr); err != nil {
				log.Printf("Failed to parse message from %s: %s\n%s", c.RemoteAddr(), err, m)
				return
			}

			rsp := m.Answer(diam.Success)
			var cca dcca.CCA
			cca.SessionId = ccr.SessionId
			cca.ResultCode = 2001
			cca.OriginHost = fmt.Sprintf("peer%d-gy-server.%s", serverId, destinationRealm)
			cca.OriginRealm = destinationRealm

			cca.CCRequestType = ccr.CCRequestType
			cca.CCRequestNumber = ccr.CCRequestNumber

			rsp.Marshal(&cca)
			rsp.WriteTo(c)
			// c.Close()
		}
	}
	handleGyCCA := func(clientId int) diam.HandlerFunc {
		return func(c diam.Conn, m *diam.Message) {
			var sid int
			var clientId2 int
			var seqNo int
			var serverId int

			if m.AVP[0].Code == 263 {
				sessionIdAvp := m.AVP[0].Data.(datatype.UTF8String)
				n, err := fmt.Sscanf(string(sessionIdAvp), "peer%d-gy-client."+originRealm+";%d;%d", &clientId2, &seqNo, &sid)
				if err != nil || n != 3 || sid < 0 || gySessions[sid] == nil {
					log.Println("Failed to extract sessionId from Gy CCA. ", m)
					return
				}
			} else {
				log.Println("Failed to extract SessionId from Gy CCA. ", m)
				return
			}

			if m.AVP[2].Code == 264 {
				originHostAvp := m.AVP[2].Data.(datatype.DiameterIdentity)
				n, err := fmt.Sscanf(string(originHostAvp), "peer%d-gy-server."+destinationRealm, &serverId)
				if err != nil || n != 1 || serverId < 0 || gyServers[serverId] == nil {
					log.Println("Failed to extract serverId from Gy CCA. ", m)
					return
				}
			} else {
				log.Println("Failed to extract serverId from Gy CCA. ", m)
				return
			}
			gySessions[sid].HandleCCA(clientId, serverId, m)
		}
	}

	for i := 0; i < 8; i++ {
		gxClientHost := fmt.Sprintf("peer%d-gx-client.%s", i, originRealm)
		gxServerHost := fmt.Sprintf("peer%d-gx-server.%s", i, destinationRealm)
		gxCmux, gxC := AppClient(peerAddr, gxClientHost, originRealm, originStateId, gx.GX_APPLICATION_ID)
		gxSmux, gxS := AppClient(peerAddr, gxServerHost, destinationRealm, originStateId, gx.GX_APPLICATION_ID)
		defer gxC.Close()
		defer gxS.Close()
		gxClientSm[i] = gxCmux
		gxServerSm[i] = gxSmux
		gxClients[i] = gxC
		gxServers[i] = gxS
		gxCmux.HandleFunc("CCA", handleGxCCA(i))
		gxSmux.Handle("CCR", handleGxCCR(i))

		gyClientHost := fmt.Sprintf("peer%d-gy-client.%s", i, originRealm)
		gyServerHost := fmt.Sprintf("peer%d-gy-server.%s", i, destinationRealm)
		gyCmux, gyC := AppClient(peerAddr, gyClientHost, originRealm, originStateId, dcca.DCCA_APPLICATION_ID)
		gySmux, gyS := AppClient(peerAddr, gyServerHost, destinationRealm, originStateId, dcca.DCCA_APPLICATION_ID)
		defer gyC.Close()
		defer gyS.Close()
		gyClientSm[i] = gyCmux
		gyServerSm[i] = gySmux
		gyClients[i] = gyC
		gyServers[i] = gyS
		gyCmux.HandleFunc("CCA", handleGyCCA(i))
		gySmux.Handle("CCR", handleGyCCR(i))
	}
	okChan := make(chan struct{})
	defer close(okChan)
	go func() {
		for {
			for i := 0; i < 8; {
				select {
				case <-okChan: // test case pass
					return
				case err := <-gyClientSm[i].ErrorReports():
					log.Println("client Error: ", err)
				case err := <-gyServerSm[i].ErrorReports():
					log.Println("server Error: ", err)
				default:
					i++
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	for i := 0; i < numSessions; i++ {
		s := NewGxSession(gxClients[i%8], i%8, i, originStateId)
		gxSessions[i] = s
	}

	for i := 0; i < numSessions; i++ {
		s := NewGySession(gyClients[i%8], i%8, i, originStateId)
		gySessions[i] = s
	}

	time.Sleep(10 * time.Second) // wait for freeDiameter peer state

	// 1. Gx load balance
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < 16; i++ {
		if gyCounter[i] > 0 {
			t.Fatal("Gx CCR were sent to Gy server ", i)
		}
	}
	if true {
		sum := sum(gxCounter[:8])
		mean := sum / 8
		if sum == 0 || mean == 0 {
			t.Fatal("Gx CCR counter is zero.")
		}
		rsd := stdDev(gxCounter[:8], mean) / mean
		if rsd > 0.2 {
			t.Fatal("Gx CCR doesn't spread equally among the 8 servers. RSD=", rsd, ", counter[8] =", gxCounter[:8])
		}
		t.Log("Gx CCR spread equally among the 8 servers. RSD=", rsd, ", counter[8] =", gxCounter[:8])
	}
	for i := 0; i < numSessions; i++ {
		if gxSessions[i].CCAError > 0 {
			t.Fatal("Gx CCR is sent to differetnt server.",
				"ClientId=", gxSessions[i].ClientId,
				"SeqNo=", gxSessions[i].SeqNo,
				"Sid=", gxSessions[i].Sid,
				"SessionId=", gxSessions[i].SessionId,
				"SubscriptionId=", gxSessions[i].SubscriptionId,
				"DestinationHost=", gxSessions[i].DestinationHost,
				"ServerId=", gxSessions[i].ServerId,
				"WrongServerId=", gxSessions[i].WrongServerId,
				"WrongClientId=", gxSessions[i].WrongClientId)
		}
	}
	resetCounter()

	// 2. Gy
	for i := 0; i < numSessions; i++ {
		gySessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gySessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < 16; i++ {
		if gxCounter[i] > 0 {
			t.Fatal("Gy CCR were sent to Gx server ", i)
		}
	}

	if true {
		sum := sum(gyCounter[:8])
		mean := sum / 8
		if sum == 0 || mean == 0 {
			t.Fatal("Gy CCR counter is zero.")
		}
		rsd := stdDev(gyCounter[:8], mean) / mean
		if rsd > 0.2 {
			t.Fatal("Gy CCR doesn't spread equally among the 8 servers. RSD=", rsd, ", counter[8] =", gyCounter[:8])
		}
		t.Log("Gy CCR spread equally among the 8 servers. RSD=", rsd, ", counter[8] =", gyCounter[:8])
	}
	for i := 0; i < numSessions; i++ {
		if gySessions[i].CCAError > 0 {
			t.Fatal("Gy CCR is sent to differetnt server.",
				"ClientId=", gySessions[i].ClientId,
				"SeqNo=", gySessions[i].SeqNo,
				"Sid=", gySessions[i].Sid,
				"SessionId=", gySessions[i].SessionId,
				"SubscriptionId=", gySessions[i].SubscriptionId,
				"DestinationHost=", gySessions[i].DestinationHost,
				"ServerId=", gySessions[i].ServerId,
				"WrongServerId=", gySessions[i].WrongServerId,
				"WrongClientId=", gySessions[i].WrongClientId)
		}
	}
	resetCounter()

	// 3. The same User should reach the same PCRF
	user1 := RandomSubscriptionId()
	user2 := RandomSubscriptionId()
	sessionA := gxSessions[0]
	sessionB := gxSessions[1]
	sessionC := gxSessions[2]
	sessionD := gxSessions[3]
	sessionA.Init() // should reset serverId before SubscriptionId has changed
	sessionB.Init()
	sessionC.Init()
	sessionD.Init()
	sessionA.SubscriptionId = user1
	sessionB.SubscriptionId = user1
	sessionC.SubscriptionId = user2
	sessionD.SubscriptionId = user2

	for i := 0; i < 400; i++ {
		sessionA.SendRequest()
		sessionB.SendRequest()
		sessionC.SendRequest()
		sessionD.SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < 400; i++ {
		sessionA.SendRequest()
		sessionB.SendRequest()
		sessionC.SendRequest()
		sessionD.SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < 400; i++ {
		sessionA.SendRequest()
		sessionB.SendRequest()
		sessionC.SendRequest()
		sessionD.SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < 400; i++ {
		sessionA.SendRequest()
		sessionB.SendRequest()
		sessionC.SendRequest()
		sessionD.SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	if sessionA.ServerId != sessionB.ServerId {
		t.Fatal("Gx CCR of the same user didn't reach the same PCRF server.",
			"sessionA.SessionId=", sessionA.SessionId,
			"sessionA.ServerId=", sessionA.ServerId,
			"sessionB.SessionId=", sessionB.SessionId,
			"sessionB.ServerId=", sessionB.ServerId)
	}
	if sessionC.ServerId != sessionD.ServerId {
		t.Fatal("Gx CCR of the same user didn't reach the same PCRF server.",
			"sessionC.SessionId=", sessionC.SessionId,
			"sessionC.ServerId=", sessionC.ServerId,
			"sessionD.SessionId=", sessionD.SessionId,
			"sessionD.ServerId=", sessionD.ServerId)
	}
	sessionA.Init() // This reset serverId and SubscriptionId
	sessionB.Init()
	sessionC.Init()
	sessionD.Init()
	resetCounter()

	// 4.  Add a server
	for i := 0; i < numSessions; i++ {
		gxSessions[i].Init()
	}
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	resetCounter()
	for i := 8; i < 9; i++ {
		gxClientHost := fmt.Sprintf("peer%d-gx-client.%s", i, originRealm)
		gxServerHost := fmt.Sprintf("peer%d-gx-server.%s", i, destinationRealm)
		gxCmux, gxC := AppClient(peerAddr, gxClientHost, originRealm, originStateId, gx.GX_APPLICATION_ID)
		gxSmux, gxS := AppClient(peerAddr, gxServerHost, destinationRealm, originStateId, gx.GX_APPLICATION_ID)
		defer gxC.Close()
		defer gxS.Close()
		gxClientSm[i] = gxCmux
		gxServerSm[i] = gxSmux
		gxClients[i] = gxC
		gxServers[i] = gxS
		gxCmux.HandleFunc("CCA", handleGxCCA(i))
		gxSmux.Handle("CCR", handleGxCCR(i))
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < 16; i++ {
		if gyCounter[i] > 0 {
			t.Fatal("Gx CCR were sent to Gy server ", i)
		}
	}
	if true {
		sum := sum(gxCounter[:9])
		mean := sum / 9
		if sum == 0 || mean == 0 {
			t.Fatal("Gx CCR counter is zero.")
		}
		rsd := stdDev(gxCounter[:9], mean) / mean
		if rsd > 0.2 {
			t.Fatal("Gx CCR doesn't spread equally among the 9 servers. RSD=", rsd, ", counter[9] =", gxCounter[:9])
		}
		t.Log("Gx CCR spread equally among the 9 servers. RSD=", rsd, ", counter[9] =", gxCounter[:9])
	}
	numSessionAffected := 0
	for i := 0; i < numSessions; i++ {
		if gxSessions[i].CCAError > 0 {
			numSessionAffected++
			gxSessions[i].CCAError = 0
		}
	}
	if float64(numSessionAffected)/float64(numSessions) > 0.2 {
		t.Fatal("Aftter add 1 server, more than 20% sessions were affected. numSessionAffected/totalSessions = ", numSessionAffected, "/", numSessions)
	}
	t.Log("Aftter add 1 server, numSessionAffected/totalSessions = ", numSessionAffected, "/", numSessions)
	resetCounter()

	// 5. remove 1 server
	for i := 0; i < numSessions; i++ {
		gxSessions[i].Init()
	}
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	resetCounter()
	gxServers[0].Close() // close server 1
	time.Sleep(60 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Second)
	for i := 0; i < numSessions; i++ {
		gxSessions[i].SendRequest()
		if i%128 == 1 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
	for i := 0; i < 16; i++ {
		if gyCounter[i] > 0 {
			t.Fatal("Gx CCR were sent to Gy server ", i)
		}
	}
	if true {
		sum := sum(gxCounter[1:9])
		mean := sum / 8
		if sum == 0 || mean == 0 {
			t.Fatal("Gx CCR counter is zero.")
		}
		rsd := stdDev(gxCounter[1:9], mean) / mean
		if rsd > 0.2 {
			t.Fatal("Gx CCR doesn't spread equally among the 9 servers. counter[8] =", gxCounter[1:9])
		}
		t.Log("Gx CCR spread equally among the 8 servers. RSD=", rsd, ", counter[8] =", gxCounter[1:9])
	}
	numSessionAffected = 0
	for i := 0; i < numSessions; i++ {
		if gxSessions[i].CCAError > 0 {
			numSessionAffected++
			gxSessions[i].CCAError = 0
		}
	}
	if float64(numSessionAffected)/float64(numSessions) > 0.2 {
		t.Fatal("Aftter remove 1 server, more than 20% sessions were affected. numSessionAffected/totalSessionss = ", numSessionAffected, "/", numSessions)
	}
	t.Log("Aftter remove 1 server, numSessionAffected/totalSessions = ", numSessionAffected, "/", numSessions)
	resetCounter()
}

// The PCRF routing information stored
// 7.3.3 Capabilities Exchange
// In addition to the capabilities exchange procedures defined in IETF RFC 3588 [14],
// the Redirect DRA and Proxy DRA shall advertise the specific applications it
// supports (e.g., Gx, Gxx, Rx, Np,S9 and for unsolicited application reporting, Sd)
// by including the value of the application identifier in the Auth-Application-Id
// AVP and the value of the 3GPP (10415) in the Vendor-Id AVP of the
// Vendor-Specific-Application-Id AVP contained in the Capabilities-ExchangeRequest
// and Capabilities-Exchange-Answer commands.
func TestDRACEA(t *testing.T) {

}

// 7.4.1.2 Modification of Diameter Sessions
func TestSessionModification(t *testing.T) {

}

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	flag.StringVar(&peerAddr, "peer_addr", "192.168.4.231:3868", "remote peer address in the form of ip:port to connect")
	flag.StringVar(&localAddr, "local_addr", ":3868", "local address in the form of ip:port to listen on")
	flag.StringVar(&originHost, "origin_host", "peer1.localdomain.net", "diameter iorigin_host")
	flag.StringVar(&originRealm, "origin_realm", "localdomain.net", "diameter origin_realm")
	flag.StringVar(&destinationHost, "destination_host", "peer2.localdomain2.net", "diameter destination_Host")
	flag.StringVar(&destinationRealm, "destination_realm", "localdomain2.net", "diameter destination_realm")
	// certFile := flag.String("cert_file", "", "tls certificate file (optional)")
	// keyFile := flag.String("key_file", "", "tls key file (optional)")

	flag.Parse()
	if len(peerAddr) == 0 {
		flag.Usage()
		os.Exit(m.Run())
	}

	m.Run()
}
