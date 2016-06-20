package dra

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"dra/dcca"
	"dra/gx"
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/dict"
)

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
			var ccr gx.CCR
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
		ccr.SessionId = "peeri-gx-client.localdomain.net;25020007;1798;192.168.4.85"
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
