package dcca

import (
	"time"
)

type SubscriptionId struct {
	SubscriptionIdType int32  `avp:"Subscription-Id-Type"` // enum
	SubscriptionIdData string `avp:"Subscription-Id-Data"`
}

type UnitValue struct {
	ValueDigits int64 `avp:"Value-Digits"`
	Exponent    int32 `avp:"Exponent"`
}
type CCMoney struct {
	UnitValue    `avp:"Unit-Value"`
	CurrencyCode uint32 `avp:"Currency-Code"`
}

type RequestedServiceUnit struct {
	CCTime                 uint32   `avp:"CC-Time,omitempty"`
	CCMoney                *CCMoney `avp:"CC-Money,omitempty"`
	CCTotalOctets          uint64   `avp:"CC-Total-Octets,omitempty"`
	CCInputOctets          uint64   `avp:"CC-Input-Octets,omitempty"`
	CCOutputOctets         uint64   `avp:"CC-Output-Octets,omitempty"`
	CCServiceSpecificUnits uint64   `avp:"CC-Service-Specific-Units,omitempty"`
}

type UsedServiceUnit struct {
	TariffChangeUsage      int32    `avp:"Tariff-Change-Usage,omitempty"` //enum
	CCTime                 uint32   `avp:"CC-Time,omitempty"`
	CCMoney                *CCMoney `avp:"CC-Money,omitempty"`
	CCTotalOctets          uint64   `avp:"CC-Total-Octets,omitempty"`
	CCInputOctets          uint64   `avp:"CC-Input-Octets,omitempty"`
	CCOutputOctets         uint64   `avp:"CC-Output-Octets,omitempty"`
	CCServiceSpecificUnits uint64   `avp:"CC-Service-Specific-Units,omitempty"`
}

type GrantedServiceUnit struct {
	TariffTimeChange       *time.Time `avp:"Tariff-Time-Change"`
	CCTime                 uint32     `avp:"CC-Time"`
	CCMoney                *CCMoney   `avp:"CC-Money"`
	CCTotalOctets          uint64     `avp:"CC-Total-Octets"`
	CCInputOctets          uint64     `avp:"CC-Input-Octets"`
	CCOutputOctets         uint64     `avp:"CC-Output-Octets"`
	CCServiceSpecificUnits uint64     `avp:"CC-Service-Specific-Units"`
}

type GSUPoolReference struct {
	GSUPoolIdentifier uint32 `avp:"G-S-U-Pool-Identifier"`
	CCUnitType        int32  `avp:"CC-Unit-Type"` // enum
	UnitValue         `avp:"Unit-Value"`
}
type RedirectServer struct {
	RedirectAddressType   int32  `avp:"Redirect-Address-Type"` // enum
	RedirectServerAddress string `avp:"Redirect-Server-Address"`
}
type FinalUnitIndication struct {
	FinalUnitAction       int32           `avp:"Final-Unit-Action"`       // enum
	RestrictionFilterRule string          `avp:"Restriction-Filter-Rule"` // IPFilterRule  or OctetString
	FilterId              string          `avp:"Filter-Id"`
	RedirectServer        *RedirectServer `avp:"Redirect-Server"`
}

type MultipleServicesCreditControl struct {
	GrantedServiceUnit   *GrantedServiceUnit   `avp:"Granted-Service-Unit,omitempty"`
	RequestedServiceUnit *RequestedServiceUnit `avp:"Requested-Service-Unit,omitempty"`
	UsedServiceUnit      []UsedServiceUnit     `avp:"Used-Service-Unit,omitempty"`
	TariffChangeUsage    int32                 `avp:"Tariff-Change-Usage,omitempty"` // enum
	ServiceIdentifier    uint32                `avp:"Service-Identifier,omitempty"`
	RatingGroup          uint32                `avp:"Rating-Group,omitempty"`
	GSUPoolReference     *GSUPoolReference     `avp:"G-S-U-Pool-Reference,omitempty"`
	ValidityTime         uint32                `avp:"Validity-Time,omitempty"`
	ResultCode           uint32                `avp:"Result-Code,omitempty"`
	FinalUnitIndication  *FinalUnitIndication  `avp:"Final-Unit-Indication,omitempty"`
}

type ServiceParameterInfo struct {
	ServiceParameterType  uint32 `avp:"Service-Parameter-Type"`
	ServiceParameterValue string `avp:"Service-Parameter-Value"` // OctetString
}

type UserEquipmentInfo struct {
	UserEquipmentInfoType  int32  `avp:"User-Equipment-Info-Type"`  // enum
	UserEquipmentInfoValue string `avp:"User-Equipment-Info-Value"` // OctetString
}
type ProxyInfo struct {
	ProxyHost  string `avp:"Proxy-Host"`  // DiameterURI
	ProxyState string `avp:"Proxy-State"` //OctetString
}

type CCR struct {
	SessionId                 string                          `avp:"Session-Id"`
	OriginHost                string                          `avp:"Origin-Host"`
	OriginRealm               string                          `avp:"Origin-Realm"`
	DestinationRealm          string                          `avp:"Destination-Realm"`
	AuthApplicationId         uint32                          `avp:"Auth-Application-Id"`
	ServiceContextId          string                          `avp:"Service-Context-Id"`
	CCRequestType             int32                           `avp:"CC-Request-Type"` //enum
	CCRequestNumber           uint32                          `avp:"CC-Request-Number"`
	DestinationHost           string                          `avp:"Destination-Host,omitempty"`
	UserName                  string                          `avp:"User-Name"`
	CCSubSessionId            uint64                          `avp:"CC-Sub-Session-Id,omitempty"`
	AcctMultiSessionId        string                          `avp:"Acct-Multi-Session-Id,omitempty"`
	OriginStateId             uint32                          `avp:"Origin-State-Id,omitempty"`
	EventTimestamp            *time.Time                      `avp:"Event-Timestamp"`
	SubscriptionId            []SubscriptionId                `avp:"Subscription-Id"`
	ServiceIdentifier         uint32                          `avp:"Service-Identifier,omitempty"`
	TerminationCause          int32                           `avp:"Termination-Cause,omitempty"` // enum
	RequestedServiceUnit      *RequestedServiceUnit           `avp:"Requested-Service-Unit,omitempty"`
	RequestedAction           int32                           `avp:"Requested-Action,omitempty"` // enum
	UsedServiceUnit           []UsedServiceUnit               `avp:"Used-Service-Unit,omitempty"`
	MultipleServicesIndicator int32                           `avp:"Multiple-Services-Indicator,omitempty"` //enum
	MSCC                      []MultipleServicesCreditControl `avp:"Multiple-Services-Credit-Control,omitempty"`
	SPI                       []ServiceParameterInfo          `avp:"Service-Parameter-Info,omitempty"`
	CcCorrelationId           string                          `avp:"CC-Correlation-Id,omitempty"` // OctetString
	UserEquipmentInfo         *UserEquipmentInfo              `avp:"User-Equipment-Info,omitempty"`
	ProxyInfo                 []ProxyInfo                     `avp:"Proxy-Info,omitempty"`
	RouteRecord               []string                        `avp:"Route-Record,omitempty"` // DiameterIdentity
}

type CostInformation struct {
	UnitValue    `avp:"Unit-Value"`
	CurrencyCode uint32 `avp:"Currency-Code"`
	CostUnit     string `avp:"Cost-Unit"`
}

type CCA struct {
	SessionId                     string                          `avp:"Session-Id"`
	ResultCode                    uint32                          `avp:"Result-Code"`
	OriginHost                    string                          `avp:"Origin-Host"`
	OriginRealm                   string                          `avp:"Origin-Realm"`
	CCRequestType                 int32                           `avp:"CC-Request-Type"`
	CCRequestNumber               uint32                          `avp:"CC-Request-Number"`
	UserName                      string                          `avp:"User-Name,omitempty"`
	CCSessionFailover             int32                           `avp:"CC-Session-Failover,omitempty"` // enum
	CCSubSessionId                uint64                          `avp:"CC-Sub-Session-Id,omitempty"`
	AcctMultiSessionId            string                          `avp:"Acct-Multi-Session-Id,omitempty"`
	OriginStateId                 uint32                          `avp:"Origin-State-Id,omitempty"`
	EventTimestamp                *time.Time                      `avp:"Event-Timestamp,omitempty"`
	GrantedServiceUnit            *GrantedServiceUnit             `avp:"Granted-Service-Unit,omitempty"`
	MSCC                          []MultipleServicesCreditControl `avp:"Multiple-Services-Credit-Control,omitempty"`
	CostInformation               *CostInformation                `avp:"Cost-Information,omitempty"`
	FinalUnitIndication           *FinalUnitIndication            `avp:"Final-Unit-Indication,omitempty"`
	CheckBalanceResult            int32                           `avp:"Check-Balance-Result,omitempty"`             // enum
	CreditControlFailureHandling  int32                           `avp:"Credit-Control-Failure-Handling,omitempty"`  // enum
	DirectDebitingFailureHandling int32                           `avp:"Direct-Debiting-Failure-Handling,omitempty"` // enum
	ValidityTime                  uint32                          `avp:"Validity-Time,omitempty"`
	RedirectHost                  []string                        `avp:"Redirect-Host,omitempty"`       // DiameterURI
	RedirectHostUsage             int32                           `avp:"Redirect-Host-Usage,omitempty"` // enum
	RedirectMaxCacheTime          uint32                          `avp:"Redirect-Max-Cache-Time,omitempty"`
	ProxyInfo                     []ProxyInfo                     `avp:"Proxy-Info,omitempty"`
	RouteRecord                   []string                        `avp:"Route-Record,omitempty"` // DiameterIdentity
	FailedAVP                     []byte                          `avp:"Failed-AVP,omitempty"`   // diamtype.Grouped
}
