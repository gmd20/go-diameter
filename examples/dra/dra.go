package dra

import (
	// "encoding/xml"
	// "errors"
	"log"
	"math/rand"
	"net"
	// "strconv"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/diamtest"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
	// "github.com/fiorix/go-diameter/diam/sm/smpeer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleALL(c diam.Conn, m *diam.Message) {
	log.Printf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
}

// http://tools.ietf.org/html/rfc6733#section-8.16
// 8.16.  Origin-State-Id AVP
// If Origin-State-Id AVP is specified in a Diameter message, it's value must be
// the same in all messages, otherwise the Diameter entity will consider the
// entity is reboot and reconnect the connection
// If originStateId is 0,  diam library will omit this AVP in CER and DWR message,
// and the AVP should not exist in other messages.

func Client(addr string, host string, realm string, originStateId uint32) (*sm.StateMachine, diam.Conn) {
	return AppClient(addr, host, realm, originStateId, 4) // dcca application id
}
func AppClient(addr string, host string, realm string, originStateId uint32, applicationId uint32) (*sm.StateMachine, diam.Conn) {
	oid := originStateId
	// if oid == 0 {
	// oid = uint32(time.Now().Unix())
	// }
	cfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(host),
		OriginRealm:      datatype.DiameterIdentity(realm),
		VendorID:         0,
		ProductName:      "dra",
		OriginStateID:    datatype.Unsigned32(oid),
		FirmwareRevision: 1,
	}

	// Create the state machine (it's a diam.ServeMux) and client.
	mux := sm.New(cfg)

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: 10 * time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   10 * time.Second,
		// AcctApplicationID: []*diam.AVP{
		//     Advertise that we want support for both
		//     Accounting applications 4 and 999.
		//     diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)), // RFC 4006
		//     diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(helloApplication)),
		// },
		AuthApplicationID: []*diam.AVP{
			// 3GPP TS 29.212 version 7.7.0 Release 7
			// 5.1 Protocol support
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(applicationId)), // RFC 4006
			// diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(16777238)),      //  Gx Application-Id
		},
		VendorSpecificApplicationID: []*diam.AVP{
			// 3GPP TS 29.212 version 7.7.0 Release 7
			// 5.2 Initialization, maintenance and termination of connection and session
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)), // Vendor-Id = 3GPP
				},
			}),
		},
	}

	// Set message handlers.
	// done := make(chan struct{}, 1000)
	mux.HandleFunc("ALL", handleALL) // Catch all.

	// cli.DialTLS(addr, cert, key)
	c, err := cli.Dial(addr)
	if err != nil {
		log.Fatal("client dial error: ", err)
	}

	return mux, c
}

// http://tools.ietf.org/html/rfc6733#section-8.16
// 8.16.  Origin-State-Id AVP
// If Origin-State-Id AVP is specified in a Diameter message, it's value must be
// the same in all messages, otherwise the Diameter entity will consider the
// entity is reboot and reconnect the connection
// If originStateId is 0,  diam library will omit this AVP in CER and DWR message,
// and the AVP should not exist in other messages.
func Server(addr string, host string, realm string, originStateId uint32) (*sm.StateMachine, *diamtest.Server) {
	oid := originStateId
	// if oid == 0 {
	// oid = uint32(time.Now().Unix())
	// }
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(host),
		OriginRealm:      datatype.DiameterIdentity(realm),
		VendorID:         0,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(oid),
		FirmwareRevision: 1,
	}

	// Create the state machine (mux) and set its message handlers.
	mux := sm.New(settings)
	mux.HandleFunc("ALL", handleALL) // Catch all.

	// diam.ListenAndServeTLS(addr, cert, key, handler, nil)
	// err := diam.ListenAndServe(addr, mux, nil)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("server Listen error, addr=", addr, " err:", err)
	}
	// s := diamtest.NewUnstartedServer(mux, dict.Default)
	s := &diamtest.Server{
		Listener: l,
		Config: &diam.Server{
			Handler: mux,
			Dict:    dict.Default,
		},
	}

	s.Start()
	return mux, s
}
