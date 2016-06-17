package gx

const (
	GX_APPLICATION_ID = 16777238
)

var GxXML = `
<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="16777238" name="Gx">
    <vendor id="10415" name="3GPP" />

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.6.2 and 5.6.3 -->
    <command name="Credit-Control" short="CC" code="272" vendor-id="10415">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="Destination-Host" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Subscription-Id" required="false" max="1"/>
        <rule avp="Supported-Features" required="false" max="1"/>
        <rule avp="Network-Request-Support" required="false" max="1"/>
        <rule avp="Bearer-Identifier" required="false" max="1"/>
        <rule avp="Bearer-Operation" required="false" max="1"/>
        <rule avp="Framed-IP-Address" required="false" max="1"/>
        <rule avp="Framed-IPv6-Prefix" required="false" max="1"/>
        <rule avp="IP-CAN-Type" required="false" max="1"/>
        <rule avp="3GPP-RAT-Type" required="false" max="1"/>
        <rule avp="RAT-Type" required="false" max="1"/>
        <rule avp="Termination-Cause" required="false" max="1"/>
        <rule avp="User-Equipment-Info" required="false" max="1"/>
        <rule avp="QoS-Information" required="false" max="1"/>
        <rule avp="QoS-Negotiation" required="false" max="1"/>
        <rule avp="QoS-Upgrade" required="false" max="1"/>
        <rule avp="3GPP-SGSN-MCC-MNC" required="false" max="1"/>
        <rule avp="3GPP-SGSN-Address" required="false" max="1"/>
        <rule avp="3GPP-SGSN-IPv6-Address" required="false" max="1"/>
        <rule avp="RAI" required="false" max="1"/>
        <rule avp="3GPP-User-Location-Info" required="false" max="1"/>
        <rule avp="3GPP-MS-TimeZone" required="false" max="1"/>
        <rule avp="Called-Station-ID" required="false" max="1"/>
        <rule avp="Bearer-Usage" required="false" max="1"/>
        <rule avp="Online" required="false" max="1"/>
        <rule avp="Offline" required="false" max="1"/>
        <rule avp="TFT-Packet-Filter-Information" required="false" max="1"/>
        <rule avp="Charging-Rule-Report" required="false" max="1"/>
        <rule avp="Event-Trigger" required="false" max="1"/>
        <rule avp="Access-Network-Charging-Address" required="false" max="1"/>
        <rule avp="Access-Network-Charging-Identifier-Gx" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Result-Code" required="false" max="1"/>
        <rule avp="Experimental-Result" required="false" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="Supported-Features" required="false" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="Bearer-Control-Mode" required="false" max="1"/>
        <rule avp="Event-Trigger" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Charging-Rule-Remove" required="false" max="1"/>
        <rule avp="Charging-Rule-Install" required="false" max="1"/>
        <rule avp="Charging-Information" required="false" max="1"/>
        <rule avp="Online" required="false" max="1"/>
        <rule avp="Offline" required="false" max="1"/>
        <rule avp="QoS-Information" required="false" max="1"/>
        <rule avp="Revalidation-Time" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </answer>
    </command>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.6.4 and 5.6.5 -->
    <command name="Re-Auth" short="RA" code="258" vendor-id="10415">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Destination-Host" required="true" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Trigger" required="false" max="1"/>
        <rule avp="Charging-Rule-Remove" required="false" max="1"/>
        <rule avp="Charging-Rule-Install" required="false" max="1"/>
        <rule avp="QoS-Information" required="false" max="1"/>
        <rule avp="Revalidation-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Result-Code" required="false" max="1"/>
        <rule avp="Experimental-Result" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Charging-Rule-Report" required="false" max="1"/>
        <rule avp="Access-Network-Charging-Address" required="false" max="1"/>
        <rule avp="Access-Network-Charging-Identifier-Gx" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </answer>
    </command>

		<avp name="User-Equipment-Info" code="458" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.49-->
			<data type="Grouped">
				<rule avp="User-Equipment-Info-Type" required="true" max="1"/>
				<rule avp="User-Equipment-Info-Value" required="true" max="1"/>
			</data>
		</avp>

		<avp name="User-Equipment-Info-Type" code="459" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.50-->
			<data type="Enumerated">
				<item code="0" name="IMEISV"/>
				<item code="1" name="MAC"/>
				<item code="2" name="EUI64"/>
				<item code="3" name="MODIFIED_EUI64"/>
			</data>
		</avp>

		<avp name="User-Equipment-Info-Value" code="460" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.51-->
			<data type="OctetString"/>
		</avp>

    <!-- Ref = [3GPP - 29.214] 5.3.2 -->
    <avp name="Access-Network-Charging-Address" code="501" must="M" may-encrypt="yes" vendor-id="10415">
      <data type="Address" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.1 -->
    <avp name="Bearer-Usage" vendor-id="10415" code="1000" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="GENERAL" code="0" />
        <item name="IMS_SIGNALLING" code="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.23 -->
    <avp name="Bearer-Control-Mode" vendor-id="10415" code="1023" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="UE_ONLY" code="0" />
        <item name="RESERVED" code="1" />
        <item name="UE_NW" code="2" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.8 -->
    <avp name="Metering-Method" vendor-id="10415" code="1007" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="DURATION" code="0" />
        <item name="VOLUME" code="1" />
        <item name="DURATION_VOLUME" code="2" />
      </data>
    </avp>
    <!-- [3GPP TS 29.229 V7.7.0] , clause 6.3.20 -->
    <avp name="Primary-Event-Charging-Function-Name" vendor-id="10415" code="619" must="M" may-encrypt="yes">
      <data type="DiameterURI" />
    </avp>

    <!-- [3GPP TS 29.229 V7.7.0] , clause 6.3.21 -->
    <avp name="Secondary-Event-Charging-Function-Name" vendor-id="10415" code="620" must="M" may-encrypt="yes">
      <data type="DiameterURI" />
    </avp>

    <!-- [3GPP TS 29.229 V7.7.0] , clause 6.3.22 -->
    <avp name="Primary-Charging-Collection-Function-Name" vendor-id="10415" code="621" must="M" may-encrypt="yes">
      <data type="DiameterURI" />
    </avp>

    <!-- [3GPP TS 29.229 V7.7.0] , clause 6.3.23 -->
    <avp name="Secondary-Charging-Collection-Function-Name" vendor-id="10415" code="622" must="M" may-encrypt="yes">
      <data type="DiameterURI" />
    </avp>

    <!-- [3GPP TS 29.214 v7.3.0] , clause 5.3.4 -->
    <avp name="Access-Network-Charging-Identifier-Value" vendor-id="10415" code="503" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.061 V7.9.0] , clause 16.4.7 -->
    <avp name="3GPP-RAT-Type" vendor-id="10415" code="21" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>
    <avp name="RAT-Type" vendor-id="10415" code="1032" must="M" may-encrypt="yes">
      <data type="Enumerated">
				<item name="WLAN"		code="0"/>
				<item name="VIRTUAL"		code="1"/>
				<item name="UTRAN"		code="1000"/>
				<item name="GERAN"		code="1001"/>
				<item name="GAN"		code="1002"/>
				<item name="HSPA_EVOLUTION"	code="1003"/>
				<item name="EUTRAN"		code="1004"/>
				<item name="CDMA2000_1X"	code="2000"/>
				<item name="HRPD"		code="2001"/>
				<item name="UMB"		code="2002"/>
				<item name="EHRPD"		code="2003"/>
      </data>
    </avp>

    <!-- Ref = RFC 4006 chap 8.48 -->
    <avp name="Subscription-Id-Data" code="444" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>
    <!-- Ref = RFC 4006 chap 8.47 -->
    <avp name="Subscription-Id-Type" code="450" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="END_USER_E164" code="0" />
        <item name="END_USER_IMSI" code="1" />
        <item name="END_USER_SIP_URI" code="2" />
        <item name="END_USER_NAI" code="3" />
        <item name="END_USER_PRIVATE" code="4" />
      </data>
    </avp>

	<avp name="Supported-Features" code="628" must="M" vendor-id="10415" may-encrypt="no">
	  <data type="Grouped">
		  <rule name="Vendor-ID" required="true" min="1" max="1" />
		  <rule name="Feature-List-ID" required="true" min="1" max="1" />
		  <rule name="Feature-List" required="true" min="1" max="1" />
		</data>
	</avp>
	<avp name="Feature-List-ID" code="629" must="M" vendor-id="10415" may-encrypt="no">
		<data type="Unsigned32" />
	</avp>
	<avp name="Feature-List" code="630" must="M" vendor-id="10415" may-encrypt="no">
		<data type="Unsigned32" />
	</avp>
    <!-- [RFC 4006] , clause 8.2 -->
    <avp name="CC-Request-Number" code="415" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>

    <!-- [RFC 4006] , clause 8.3 -->
    <avp name="CC-Request-Type" code="416" must="M" may-encrypt="yes"> <data type="Enumerated">
        <item name="INITIAL_REQUEST" code="1" />
        <item name="UPDATE_REQUEST" code="2" />
        <item name="TERMINATION_REQUEST" code="3" />
        <item name="EVENT_REQUEST" code="4" />
      </data>
    </avp>

    <!-- [RFC 4006] , clause 8.38 -->
    <avp name="Redirect-Address-Type" code="433" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="IPV4_ADDRESS" code="0" />
        <item name="IPV6_ADDRESS" code="1" />
        <item name="URL" code="2" />
        <item name="SIP_URI" code="3" />
      </data>
    </avp>
<!-- [RFC 4006] , clause 8.39 -->
    <avp name="Redirect-Server-Address" code="435" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.19 -->
    <avp name="PCC-Rule-Status" vendor-id="10415" code="1019" must="M" may-encrypt="yes">
      <data type="Enumerated" />
      <item name="ACTIVE" code="0" />
      <item name="INACTIVE" code="1" />
      <item name="TEMPORARILY INACTIVE" code="2" />
    </avp>

    <!-- [RFC 4005] , clause 6.11.1 -->
    <avp name="Framed-IP-Address" code="8" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [RFC 4005] , clause 6.11.6 -->
    <avp name="Framed-IPv6-Prefix" code="97" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [RFC 4005] , clause 4.5 -->
    <avp name="Called-Station-Id" code="30" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>

    <!-- [RFC 4005] , clause 6.7 -->
    <avp name="Filter-Id" code="11" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.21 -->
    <avp name="Bearer-Operation" vendor-id="10415" code="1021" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="TERMINATION" code="0" />
        <item name="ESTABLISHMENT" code="1" />
        <item name="MODIFICATION" code="2" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.24 -->
    <avp name="Network-Request-Support" vendor-id="10415" code="1024" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="NETWORK_REQUEST_NOT_SUPPORTED" code="0" />
        <item name="NETWORK_REQUEST_SUPPORTED" code="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.5 -->
    <avp name="Charging-Rule-Base-Name" vendor-id="10415" code="1004" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.15 -->
    <avp name="ToS-Traffic-Class" vendor-id="10415" code="1014" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.27 -->
    <avp name="IP-CAN-Type" vendor-id="10415" code="1027" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="3GPP" code="0" />
        <item name="DOCSIS" code="1" />
        <item name="xDSL" code="2" />
        <item name="WiMAX" code="3" />
        <item name="3GPP2" code="4" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.28 -->
    <avp name="QoS-Negotiation" vendor-id="10415" code="1029" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="NO_QOS_NEGOTIATION" code="0" />
        <item name="QOS_NEGOTIATION_SUPPORTED" code="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.29 -->
    <avp name="QoS-Upgrade" vendor-id="10415" code="1030" must="M" may-encrypt="yes">
      <data type="Enumerated" />
      <item name="QOS_UPGRADE_NOT_SUPPORTED" code="0" />
      <item name="QOS_UPGRADE_SUPPORTED" code="1" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.30 -->
    <avp name="Rule-Failure-Code" vendor-id="10415" code="1031" must="M" may-encrypt="yes">
      <data type="Enumerated" />
      <item name="UNKNOWN_RULE_NAME" code="1" />
      <item name="RATING_GROUP_ERROR" code="2" />
      <item name="SERVICE_IDENTIFIER_ERROR" code="3" />
      <item name="GW/PCEF_MALFUNCTION" code="4" />
      <item name="RESOURCES_LIMITATION" code="5" />
      <item name="MAX_NR_BEARERS_REACHED" code="6" />
      <item name="UNKNOWN_BEARER_ID" code="7" />
      <item name="MISSING_BEARER_ID" code="8" />
      <item name="MISSING_FLOW_DESCRIPTION" code="9" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.34 -->
    <avp name="Session-Release-Cause" vendor-id="10415" code="1035" must="M" may-encrypt="yes">
      <data type="Enumerated" />
      <item name="UNSPECIFIED_REASONS" code="0" />
      <item name="UE_SUBSCRIPTION_REASON" code="1" />
      <item name="INSUFFICIENT_SERVER_RESOURCES" code="2" />
    </avp>

    <!-- [3GPP TS 29.061 V7.9.0] , clause 16.4.7. -->
    <avp name="3GPP-SGSN-MCC-MNC" vendor-id="10415" code="18" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>

    <!-- [3GPP TS 29.061 V7.9.0] , clause 16.4.7 -->
    <avp name="3GPP-User-Location-Info" vendor-id="10415" code="22" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.31 -->
    <avp name="Revalidation-Time" vendor-id="10415" code="1042" must="M" may-encrypt="yes">
      <data type="Time" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.32 -->
    <avp name="Rule-Activation-Time" code="1033" must="M" may-encrypt="yes">
      <data type="Time" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.33 -->
    <avp name="Rule-Deactivation-Time" vendor-id="10415" code="1034" must="M" may-encrypt="yes">
      <data type="Time" />
    </avp>

    <!-- [3GPP TS 29.061 V7.7.0] , clause 16.4.7.1 -->
    <avp name="3GPP-SGSN-Address" vendor-id="10415" code="6" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.061 V7.7.0] , clause 16.4.7.1 -->
    <avp name="3GPP-SGSN-IPv6-Address" vendor-id="10415" code="15" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.061 V7.7.0] , clause 17.7.12 -->
    <avp name="RAI" vendor-id="10415" code="909" must="M" may-encrypt="yes">
      <data type="UTF8String" />
    </avp>

    <!-- [3GPP TS 29.061 V7.7.0] , clause 16.4.7 -->
    <avp name="3GPP-MS-TimeZone" vendor-id="10415" code="23" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.20 -->
    <avp name="Bearer-Identifier" vendor-id="10415" code="1020" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.6 -->
    <avp name="Charging-Rule-Name" vendor-id="10415" code="1005" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- [RFC 4006] , clause 8.28 -->
    <avp name="Service-Identifier" code="439" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>

    <!-- [RFC 4006] , clause 8.29 -->
    <avp name="Rating-Group" code="432" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>
    <!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.11 -->
    <avp name="Flow-Status" vendor-id="10415" code="511" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="ENABLED-UPLINK" code="0" />
        <item name="ENABLED-DOWNLINK" code="1" />
        <item name="ENABLED" code="2" />
        <item name="DISABLED" code="3" />
        <item name="REMOVED" code="4" />
      </data>
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.17 -->
    <avp name="QoS-Class-Identifier" vendor-id="10415" code="1028" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="FINAL_SERVICE_INFORMATION" code="0" />
        <item name="PRELIMINARY_SERVICE_INFORMATION" code="1" />
      </data>
    </avp>
    <!-- Ref = [3GPP - 29.214] 5.3.15 -->
    <avp name="Max-Requested-Bandwidth-UL" code="516" vendor-id="10415" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.25 -->
    <avp name="Guaranteed-Bitrate-DL" vendor-id="10415" code="1025" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.26 -->
    <avp name="Guaranteed-Bitrate-UL" vendor-id="10415" code="1026" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.20 -->
    <avp name="Bearer-Identifier" vendor-id="10415" code="1020" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>
    <avp name="Max-Requested-Bandwidth-UL" code="516" must="M" vendor-id="10415" may-encrypt="yes">
      <data type="Unsigned32"/>
    </avp>
    <avp name="Max-Requested-Bandwidth-DL" code="515" must="M" vendor-id="10415" may-encrypt="yes">
      <data type="Unsigned32"/>
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.16 -->
    <avp name="QoS-Information" vendor-id="10415" code="1016" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="QoS-Class-Identifier" required="false"  min="1" max="1" />
        <rule name="Max-Requested-Bandwidth-UL" required="false" min="1" max="1" />
        <rule name="Max-Requested-Bandwidth-DL" required="false" min="1" max="1" />
        <rule name="Guaranteed-Bitrate-UL" required="false" min="1" max="1" />
        <rule name="Guaranteed-Bitrate-DL" required="false" min="1" max="1" />
        <rule name="Bearer-Identifier" required="false" min="1" max="1" />
      </data>
    </avp>
    <!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.8 -->
    <avp name="Flow-Description" vendor-id="10415" code="507" must="M" may-encrypt="yes">
      <data type="IPFilterRule" />
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.12 -->
    <avp name="Reporting-Level" vendor-id="10415" code="1011" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="SERVICE_IDENTIFIER_LEVEL" code="0" />
        <item name="RATING_GROUP_LEVEL" code="1" />
      </data>
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.9 -->
    <avp name="Offline" vendor-id="10415" code="1008" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="DISABLE_OFFLINE" code="0" />
        <item name="ENABLE_OFFLINE" code="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.10 -->
    <avp name="Online" vendor-id="10415" code="1009" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="DISABLE_ONLINE" code="0" />
        <item name="ENABLE_ONLINE" code="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.11 -->
    <avp name="Precedence" vendor-id="10415" code="1010" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>

    <!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.6 -->
    <avp name="AF-Charging-Identifier" vendor-id="10415" code="505" must="M" may-encrypt="yes">
      <data type="OctetString" />
    </avp>

    <!-- Ref = [3GPP - 29.214] 5.3.14 -->
    <avp name="Max-Requested-Bandwidth-DL" code="515" vendor-id="10415" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>
    <!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.17 -->
    <avp name="Media-Component-Number" vendor-id="10415" code="518" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>

    <!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.9 -->
    <avp name="Flow-Number" vendor-id="10415" code="509" must="M" may-encrypt="yes">
      <data type="Unsigned32" />
    </avp>
<!-- [3GPP TS 29.214 V7.3.0] , clause 5.3.10 -->
    <avp name="Flows" vendor-id="10415" code="510" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Media-Component-Number" required="true" min="1" max="1" />
        <rule name="Flow-Number" required="false" max="1" />
      </data>
    </avp>
    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.4 -->
    <avp name="Charging-Rule-Definition" vendor-id="10415" code="1003" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Charging-Rule-Name" required="true" min="1" max="1" />
        <rule name="Service-Identifier" required="true" max="1" />
        <rule name="Rating-Group" required="true" max="1" />
        <rule name="Flow-Description" required="false" max="1" />
        <rule name="Flow-Status" required="true" max="1" />
        <rule name="QoS-Information" required="true" max="1" />
        <rule name="Reporting-Level" required="true" max="1" />
        <rule name="Online" required="true" max="1" />
        <rule name="Offline" required="true" max="1" />
        <rule name="Metering-Method" required="true" max="1" />
        <rule name="Precedence" required="true" max="1" />
        <rule name="AF-Charging-Identifier" required="true" max="1" />
        <rule name="Flows" required="false" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.2 -->
    <avp name="Charging-Rule-Install" vendor-id="10415" code="1001" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Charging-Rule-Definition" required="false" max="1" />
        <rule name="Charging-Rule-Name" required="false" max="1" />
        <rule name="Charging-Rule-Base-Name" required="false" max="1" />
        <rule name="Bearer-Identifier" required="true" minumin="1" max="1" />
        <rule name="Rule-Activation-Time" required="false" max="1" />
        <rule name="Rule-Deactivation-Time" required="false" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.3 -->
    <avp name="Charging-Rule-Remove" vendor-id="10415" code="1002" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Charging-Rule-Name" required="false" max="1" />
        <rule name="Charging-Rule-Base-Name" required="false" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.7 -->
    <avp name="Event-Trigger" vendor-id="10415" code="1006" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="SGSN_CHANGE" code="0" />
        <item name="QOS_CHANGE" code="1" />
        <item name="RAT_CHANGE" code="2" />
        <item name="TFT_CHANGE" code="3" />
        <item name="PLMN_CHANGE" code="4" />
        <item name="LOSS_OF_BEARER" code="5" />
        <item name="RECOVERY_OF_BEARER" code="6" />
        <item name="IP-CAN_CHANGE" code="7" />
        <item name="PCEF_MALFUNCTION" code="8" />
        <item name="RESOURCES_LIMITATION" code="9" />
        <item name="MAX_NR_BEARERS_REACHED" code="10" />
        <item name="QOS_CHANGE_EXCEEDING_AUTHORIZATION" code="11" />
        <item name="RAI_CHANGE" code="12" />
        <item name="USER_LOCATION_CHANGE" code="13" />
        <item name="NO_EVENT_TRIGGERS" code="14" />
        <item name="OUT_OF_CREDIT" code="15" />
        <item name="REALLOCATION_OF_CREDIT" code="16" />
        <item name="REVALIDATION_TIMEOUT" code="17" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.13 -->
    <avp name="TFT-Filter" vendor-id="10415" code="1012" must="M" may-encrypt="yes">
      <data type="IPFilterRule" />
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.14 -->
    <avp name="TFT-Packet-Filter-Information" vendor-id="10415" code="1013" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Precedence" required="true" min="1" max="1" />
        <rule name="TFT-Filter" required="true" min="1" max="1" />
        <rule name="ToS-Traffic-Class" required="true" min="1" max="1" />
      </data>
    </avp>

<!-- [RFC 4006] , clause 8.35 -->
    <avp name="Final-Unit-Action" code="449" must="M" may-encrypt="yes">
      <data type="Enumerated">
        <item name="TERMINATE" code="0" />
        <item name="REDIRECT" code="1" />
        <item name="RESTRICT_ACCESS" code="2" />
      </data>
    </avp>

    <!-- [RFC 4006] , clause 8.36 -->
    <avp name="Restriction-Filter-Rule" code="438" must="M" may-encrypt="yes">
      <data type="IPFilterRule" />
    </avp>
    <!-- [RFC 4006] , clause 8.37 -->
    <avp name="Redirect-Server" code="434" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Redirect-Address-Type" required="true" min="1" max="1" />
        <rule name="Redirect-Server-Address" required="true" min="1" max="1" />
      </data>
    </avp>
    <!-- [RFC 4006] , clause 8.34 -->
    <avp name="Final-Unit-Indication" code="430" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Final-Unit-Action" required="true" min="1" max="1" />
        <rule name="Restriction-Filter-Rule" required="true" min="1" max="1" />
        <rule name="Filter-Id" required="true" min="1" max="1" />
        <rule name="Redirect-Server" required="true" min="1" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.18 -->
    <avp name="Charging-Rule-Report" vendor-id="10415" code="1018" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Charging-Rule-Name" required="false" max="1" />
        <rule name="Charging-Rule-Base-Name" required="false" max="1" />
        <rule name="PCC-Rule-Status" required="false" min="1" max="1" />
        <rule name="Rule-Failure-Code" required="false" max="1" />
        <rule name="Final-Unit-Indication" required="false" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.212 V7.7.0] , clause 5.3.22 -->
    <avp name="Access-Network-Charging-Identifier-Gx" vendor-id="10415" code="1022" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Access-Network-Charging-Identifier-Value" required="true" min="1" max="1" />
        <rule name="Charging-Rule-Base-Name" required="false" max="1" />
        <rule name="Charging-Rule-Name" required="false" max="1" />
      </data>
    </avp>

    <!-- [3GPP TS 29.229 V7.7.0] , clause 6.3.19 -->
    <avp name="Charging-Information" vendor-id="10415" code="618" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Primary-Event-Charging-Function-Name" required="true" min="1" max="1" />
        <rule name="Secondary-Event-Charging-Function-Name" required="true" min="1" max="1" />
        <rule name="Primary-Charging-Collection-Function-Name" required="true" min="1" max="1" />
        <rule name="Secondary-Charging-Collection-Function-Name" required="true" min="1" max="1" />
      </data>
    </avp>

    <!-- Ref = RFC 4006 chap 8.46 -->
    <avp name="Subscription-Id" code="443" must="M" may-encrypt="yes">
      <data type="Grouped">
        <rule name="Subscription-Id-Type" required="true"/>
        <rule name="Subscription-Id-Data" required="true"/>
      </data>
    </avp>

  </application>
</diameter>
`
