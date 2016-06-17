package gx

import (
	"net"
	"time"
)

type SubscriptionId struct {
	SubscriptionIdType int32  `avp:"Subscription-Id-Type"` // enum
	SubscriptionIdData string `avp:"Subscription-Id-Data"`
}

type SupportedFeatures struct {
	VendorID      uint32 `avp:"Vendor-Id"`
	FeatureListID uint32 `avp:"Feature-List-ID"`
	FeatureList   uint32 `avp:"Feature-List"`
}

type UserEquipmentInfo struct {
	UserEquipmentInfoType  int32  `avp:"User-Equipment-Info-Type"`  // enum
	UserEquipmentInfoValue string `avp:"User-Equipment-Info-Value"` // OctetString
}

type QoSInformation struct {
	QoSClassIdentifier      *int32 `avp:"QoS-Class-Identifier"` // enum
	MaxRequestedBandwidthUL uint32 `avp:"Max-Requested-Bandwidth-UL,omitempty"`
	MaxRequestedBandwidthDL uint32 `avp:"Max-Requested-Bandwidth-DL,omitempty"`
	GuaranteedBitrateUL     uint32 `avp:"Guaranteed-Bitrate-UL,omitempty"`
	GuaranteedBitrateDL     uint32 `avp:"Guaranteed-Bitrate-DL,omitempty"`
	BearerIdentifier        string `avp:"Bearer-Identifier,omitempty"`
}

type TFTPacketFilterInformation struct {
	Precedence      uint32 `avp:"Precedence"`
	TFTFilter       string `avp:"TFT-Filter"` // IPFilterRule
	ToSTrafficClass string `avp:"ToS-Traffic-Class"`
}

type RedirectServer struct {
	RedirectAddressType   int32  `avp:"Redirect-Address-Type"` // enum
	RedirectServerAddress string `avp:"Redirect-Server-Address"`
}

type FinalUnitIndication struct {
	FinalUnitAction       int32           `avp:"Final-Unit-Action"`       // enum
	RestrictionFilterRule []string        `avp:"Restriction-Filter-Rule"` // IPFilterRule  or OctetString
	FilterId              []string        `avp:"Filter-Id"`
	RedirectServer        *RedirectServer `avp:"Redirect-Server"`
}

type ChargingRuleReport struct {
	ChargingRuleName     string               `avp:"Charging-Rule-Name,omitempty"`
	ChargingRuleBaseName string               `avp:"Charging-Rule-Base-Name,omitempty"`
	PCCRuleStatus        int32                `avp:"PCC-Rule-Status,omitempty"`   // enum
	RuleFailureCode      int32                `avp:"Rule-Failure-Code,omitempty"` // enum
	FinalUnitIndication  *FinalUnitIndication `avp:"Final-Unit-Indication"`
}

type AccessNetworkChargingIdentifierGx struct {
	AccessNetworkChargingIdentifierValue string `avp:"Access-Network-Charging-Identifier-Value"` // OctetString
	ChargingRuleBaseName                 string `avp:"Charging-Rule-Base-Name,omitempty"`
	ChargingRuleName                     string `avp:"Charging-Rule-Name,omitempty"`
}

type ProxyInfo struct {
	ProxyHost  string `avp:"Proxy-Host"`  // DiameterURI
	ProxyState string `avp:"Proxy-State"` //OctetString
}

type CCR struct {
	SessionId                         string                             `avp:"Session-Id"`
	AuthApplicationId                 uint32                             `avp:"Auth-Application-Id"`
	OriginHost                        string                             `avp:"Origin-Host"`
	OriginRealm                       string                             `avp:"Origin-Realm"`
	DestinationRealm                  string                             `avp:"Destination-Realm"`
	CCRequestType                     int32                              `avp:"CC-Request-Type"` // enum
	CCRequestNumber                   uint32                             `avp:"CC-Request-Number"`
	DestinationHost                   string                             `avp:"Destination-Host,omitempty"`
	OriginStateId                     uint32                             `avp:"Origin-State-Id,omitempty"`
	SubscriptionId                    []SubscriptionId                   `avp:"Subscription-Id,omitempty"`
	SupportedFeatures                 []SupportedFeatures                `avp:"Supported-Features,omitempty"`
	NetworkRequestSupport             *int32                             `avp:"Network-Request-Support"` // enum
	BearerIdentifier                  string                             `avp:"Bearer-Identifier,omitempty"`
	BearerOperation                   *int32                             `avp:"Bearer-Operation"` //enum
	FramedIPAddress                   string                             `avp:"Framed-IP-Address,omitempty"`
	FramedIPv6Prefix                  string                             `avp:"Framed-IPv6-Prefix,omitempty"`
	IPCANType                         *int32                             `avp:"IP-CAN-Type"`
	TGPPRATType                       string                             `avp:"3GPP-RAT-Type,omitempty"`
	RATType                           *int32                             `avp:"RAT-Type,omitempty"` // enum
	TerminationCause                  *int32                             `avp:"Termination-Cause"`  // enum
	UserEquipmentInfo                 *UserEquipmentInfo                 `avp:"User-Equipment-Info,omitempty"`
	QoSInformation                    *QoSInformation                    `avp:"QoS-Information,omitempty"`
	QoSNegotiation                    *int32                             `avp:"QoS-Negotiation"` // enum
	QoSUpgrade                        *int32                             `avp:"QoS-Upgrade"`     // enum
	TGPPSGSNMCCMNC                    string                             `avp:"3GPP-SGSN-MCC-MNC,omitempty"`
	TGPPSGSNAddress                   string                             `avp:"3GPP-SGSN-Address,omitempty"`
	TGPPSGSNIPv6Address               string                             `avp:"3GPP-SGSN-IPv6-Address,omitempty"`
	RAI                               string                             `avp:"RAI,omitempty"`
	TGPPUserLocationInfo              string                             `avp:"3GPP-User-Location-Info,omitempty"`
	TGPPMSTimeZone                    string                             `avp:"3GPP-MS-TimeZone,omitempty"`
	CalledStationID                   string                             `avp:"Called-Station-Id,omitempty"`
	BearerUsage                       *int32                             `avp:"Bearer-Usage"` // enum
	Online                            *int32                             `avp:"Online"`       // enum
	Offline                           *int32                             `avp:"Offline"`      // enum
	TFTPacketFilterInformation        *TFTPacketFilterInformation        `avp:"TFT-Packet-Filter-Information"`
	ChargingRuleReport                *ChargingRuleReport                `avp:"Charging-Rule-Report"`
	EventTrigger                      *int32                             `avp:"Event-Trigger"`                   // enum
	AccessNetworkChargingAddress      *net.IP                            `avp:"Access-Network-Charging-Address"` // address
	AccessNetworkChargingIdentifierGx *AccessNetworkChargingIdentifierGx `avp:"Access-Network-Charging-Identifier-Gx"`
	ProxyInfo                         []ProxyInfo                        `avp:"Proxy-Info,omitempty"`
	RouteRecord                       []string                           `avp:"Route-Record,omitempty"` // DiameterIdentity
}

type ChargingRuleRemove struct {
	ChargingRuleName     string `avp:"Charging-Rule-Name,omitempty"`
	ChargingRuleBaseName string `avp:"Charging-Rule-Base-Name.omitempty"`
}

type Flows struct {
	MediaComponentNumber uint32  `avp:"Media-Component-Number"`
	FlowNumber           *uint32 `avp:"Flow-Number,omitempty"`
}

type ChargingRuleDefinition struct {
	ChargingRuleName     string         `avp:"Charging-Rule-Name"`
	ServiceIdentifier    uint32         `avp:"Service-Identifier"`
	RatingGroup          uint32         `avp:"Rating-Group"`
	FlowDescription      string         `avp:"Flow-Description,omitempty"` // IPFilterRule
	FlowStatus           *int32         `avp:"Flow-Status,omitempty"`      // enum
	QoSInformation       QoSInformation `avp:"QoS-Information"`
	ReportingLevel       int32          `avp:"Reporting-Level"` // enum
	Online               int32          `avp:"Online"`          // enum
	Offline              int32          `avp:"Offline"`         // enum
	MeteringMethod       int32          `avp:"Metering-Method"` // enum
	Precedence           uint32         `avp:"Precedence-"`
	AFChargingIdentifier string         `avp:"AF-Charging-Identifier"`
	Flows                *Flows         `avp:"Flows"`
}

type ChargingRuleInstall struct {
	ChargingRuleDefinition *ChargingRuleDefinition `avp:"Charging-Rule-Definition"`
	ChargingRuleName       string                  `avp:"Charging-Rule-Name,omitempty"`
	ChargingRuleBaseName   string                  `avp:"Charging-Rule-Base-Name,omitempty"`
	BearerIdentifier       string                  `avp:"Bearer-Identifier,omitempty"`
	RuleActivationTime     *time.Time              `avp:"Rule-Activation-Time"`
	RuleDeactivationTime   *time.Time              `avp:"Rule-Deactivation-Time"`
}

type ChargingInformation struct {
	PrimaryEventChargingFunctionName        string `avp:"Primary-Event-Charging-Function-Name"`        // DiameterURI
	SecondaryEventChargingFunctionName      string `avp:"Secondary-Event-Charging-Function-Name"`      // DiameterURI
	PrimaryChargingCollectionFunctionName   string `avp:"Primary-Charging-Collection-Function-Name"`   // DiameterURI
	SecondaryChargingCollectionFunctionName string `avp:"Secondary-Charging-Collection-Function-Name"` // DiameterURI
}
type CCA struct {
	SessionId           string               `avp:"Session-Id"`
	AuthApplicationId   uint32               `avp:"Auth-Application-Id"`
	OriginHost          string               `avp:"Origin-Host"`
	OriginRealm         string               `avp:"Origin-Realm"`
	ResultCode          uint32               `avp:"Result-Code,omitempty"`
	ExperimentalResult  uint32               `avp:"Experimental-Result,omitempty"`
	CCRequestType       int32                `avp:"CC-Request-Type"` // enum
	CCRequestNumber     uint32               `avp:"CC-Request-Number"`
	SupportedFeatures   []SupportedFeatures  `avp:"Supported-Features,omitempty"`
	BearerControlMode   *int32               `avp:"Bearer-Control-Mode"` // enum
	EventTrigger        []int32              `avp:"Event-Trigger"`       // enum
	OriginStateId       uint32               `avp:"Origin-State-Id,omitempty"`
	ChargingRuleRemove  *ChargingRuleRemove  `avp:"Charging-Rule-Remove"`
	ChargingRuleInstall *ChargingRuleInstall `avp:"Charging-Rule-Install"`
	ChargingInformation *ChargingInformation `avp:"Charging-Information"`
	Online              *int32               `avp:"Online"`  // enum
	Offline             *int32               `avp:"Offline"` // enum
	QoSInformation      *QoSInformation      `avp:"QoS-Information,omitempty"`
	RevalidationTime    *time.Time           `avp:"Revalidation-Time,omitempty"`
	ErrorMessage        string               `avp:"Error-Message,omitempty"`
	ErrorReportingHost  string               `avp:"Error-Reporting-Host,omitempty"` // DiameterIdentity
	FailedAVP           []byte               `avp:"Failed-AVP,omitempty"`           // diamtype.Grouped
	ProxyInfo           []ProxyInfo          `avp:"Proxy-Info,omitempty"`
	RouteRecord         []string             `avp:"Route-Record,omitempty"` // DiameterIdentity	SessionId                     string                          `avp:"Session-Id"`
}

type RAR struct {
	SessionId           string               `avp:"Session-Id"`
	AuthApplicationId   uint32               `avp:"Auth-Application-Id"`
	OriginHost          string               `avp:"Origin-Host"`
	OriginRealm         string               `avp:"Origin-Realm"`
	DestinationRealm    string               `avp:"Destination-Realm"`
	DestinationHost     string               `avp:"Destination-Host"`
	OriginStateId       uint32               `avp:"Origin-State-Id,omitempty"`
	EventTrigger        *int32               `avp:"Event-Trigger"` // enum
	ChargingRuleRemove  *ChargingRuleRemove  `avp:"Charging-Rule-Remove"`
	ChargingRuleInstall *ChargingRuleInstall `avp:"Charging-Rule-Install"`
	QoSInformation      *QoSInformation      `avp:"QoS-Information,omitempty"`
	RevalidationTime    *time.Time           `avp:"Revalidation-Time,omitempty"`
	ProxyInfo           []ProxyInfo          `avp:"Proxy-Info,omitempty"`
	RouteRecord         []string             `avp:"Route-Record,omitempty"` // DiameterIdentity	SessionId                     string                          `avp:"Session-Id"`
}

type RAA struct {
	SessionId                         string                             `avp:"Session-Id"`
	AuthApplicationId                 uint32                             `avp:"Auth-Application-Id"`
	OriginHost                        string                             `avp:"Origin-Host"`
	OriginRealm                       string                             `avp:"Origin-Realm"`
	ResultCode                        uint32                             `avp:"Result-Code,omitempty"`
	ExperimentalResult                uint32                             `avp:"Experimental-Result,omitempty"`
	OriginStateId                     uint32                             `avp:"Origin-State-Id,omitempty"`
	ChargingRuleReport                *ChargingRuleReport                `avp:"Charging-Rule-Report"`
	AccessNetworkChargingAddress      *net.IP                            `avp:"Access-Network-Charging-Address"` // address
	AccessNetworkChargingIdentifierGx *AccessNetworkChargingIdentifierGx `avp:"Access-Network-Charging-Identifier-Gx"`
	ErrorMessage                      string                             `avp:"Error-Message,omitempty"`
	ErrorReportingHost                string                             `avp:"Error-Reporting-Host,omitempty"` // DiameterIdentity
	FailedAVP                         []byte                             `avp:"Failed-AVP,omitempty"`           // diamtype.Grouped
	ProxyInfo                         []ProxyInfo                        `avp:"Proxy-Info,omitempty"`
	RouteRecord                       []string                           `avp:"Route-Record,omitempty"` // DiameterIdentity	SessionId                     string                          `avp:"Session-Id"`
}
