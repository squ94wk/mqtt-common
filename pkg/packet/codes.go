package packet

//ConnectReason is an alias for all defined connect reason codes a connack control packet can have.
type ConnectReason byte

//DisconnectReason is an alias for all defined disconnect reason codes a disconnect control packet can have.
type DisconnectReason byte

//Names for all defined connect reason codes a connack control packet can have.
const (
	ConnectSuccess                     ConnectReason = 0   // The Connection is accepted.
	ConnectUnspecifiedError                          = 128 // The Server does not wish to reveal the reason for the failure, or none of the other Reason Codes apply.
	ConnectMalformedPacket                           = 129 // Data within the CONNECT packet could not be correctly parsed.
	ConnectProtocolError                             = 130 // Data in the CONNECT packet does not conform to this specification.
	ConnectImplementationSpecificError               = 131 // The CONNECT is valid but is not accepted by this Server.
	ConnectUnsupportedProtocolVersion                = 132 // The Server does not support the version of the MQTT protocol requested by the Client.
	ConnectClientIdentifierNotValid                  = 133 // The Client Identifier is a valid string but is not allowed by the Server.
	ConnectBadUserNameOrPassword                     = 134 // The Server does not accept the User Name or Password specified by the Client
	ConnectNotAuthorized                             = 135 // The Client is not authorized to connect.
	ConnectServerUnavailable                         = 136 // The MQTT Server is not available.
	ConnectServerBusy                                = 137 // The Server is busy. Try again later.
	ConnectBanned                                    = 138 // This Client has been banned by administrative action. Contact the server administrator.
	ConnectBadAuthenticationMethod                   = 140 // The authentication method is not supported or does not match the authentication method currently in use.
	ConnectTopicNameInvalid                          = 144 // The Will Topic Name is not malformed, but is not accepted by this Server.
	ConnectPacketTooLarge                            = 149 // The CONNECT packet exceeded the maximum permissible size.
	ConnectQuotaExceeded                             = 151 // An implementation or administrative imposed limit has been exceeded.
	ConnectPayloadFormatInvalid                      = 153 // The Will Payload does not match the specified Payload Format Indicator.
	ConnectRetainNotSupported                        = 154 // The Server does not support retained messages, and Will Retain was set to 1.
	ConnectQoSNotSupported                           = 155 // The Server does not support the QoS set in Will QoS.
	ConnectUseAnotherServer                          = 156 // The Client should temporarily use another server.
	ConnectServerMoved                               = 157 // The Client should permanently use another server.
	ConnectConnectionRateExceeded                    = 159 // The connection rate limit has been exceeded.
)

//Names for all defined disconnect reason codes a disconnect control packet can have.
const (
	DisconnectNormalDisconnection                 DisconnectReason = 0   // Close the connection normally. Do not send the Will Message.
	DisconnectDisconnectWithWillMessage           DisconnectReason = 4   // The Client wishes to disconnect but requires that the Server also publishes its Will Message.
	DisconnectUnspecifiedError                    DisconnectReason = 128 // The Connection is closed but the sender either does not wish to reveal the reason, or none of the other Reason Codes apply.
	DisconnectMalformedPacket                     DisconnectReason = 129 // The received packet does not conform to this specification.
	DisconnectProtocolError                       DisconnectReason = 130 // An unexpected or out of order packet was received.
	DisconnectImplementationSpecificError         DisconnectReason = 131 // The packet received is valid but cannot be processed by this implementation.
	DisconnectNotAuthorized                       DisconnectReason = 135 // The request is not authorized.
	DisconnectServerBusy                          DisconnectReason = 137 // The Server is busy and cannot continue processing requests from this Client.
	DisconnectServerShuttingDown                  DisconnectReason = 139 // The Server is shutting down.
	DisconnectKeepAliveTimeout                    DisconnectReason = 141 // The Connection is closed because no packet has been received for 1.5 times the Keepalive time.
	DisconnectSessionTakenOver                    DisconnectReason = 142 // Another Connection using the same ClientID has connected causing this Connection to be closed.
	DisconnectTopicFilterInvalid                  DisconnectReason = 143 // The Topic Filter is correctly formed, but is not accepted by this Sever.
	DisconnectTopicNameInvalid                    DisconnectReason = 144 // The Topic Name is correctly formed, but is not accepted by this Client or Server.
	DisconnectReceiveMaximumExceeded              DisconnectReason = 147 // The Client or Server has received more than Receive Maximum publication for which it has not sent PUBACK or PUBCOMP.
	DisconnectTopicAliasInvalid                   DisconnectReason = 148 // The Client or Server has received a PUBLISH packet containing a Topic Alias which is greater than the Maximum Topic Alias it sent in the CONNECT or CONNACK packet.
	DisconnectPackettooLarge                      DisconnectReason = 149 // The packet size is greater than Maximum Packet Size for this Client or Server.
	DisconnectMessageRateTooHigh                  DisconnectReason = 150 // The received data rate is too high.
	DisconnectQuotaExceeded                       DisconnectReason = 151 // An implementation or administrative imposed limit has been exceeded.
	DisconnectAdministrativeAction                DisconnectReason = 152 // The Connection is closed due to an administrative action.
	DisconnectPayloadFormatInvalid                DisconnectReason = 153 // The payload format does not match the one specified by the Payload Format Indicator.
	DisconnectRetainNotSupported                  DisconnectReason = 154 // The Server has does not support retained messages.
	DisconnectQoSNotSupported                     DisconnectReason = 155 // The Client specified a QoS greater than the QoS specified in a Maximum QoS in the CONNACK.
	DisconnectUseAnotherServer                    DisconnectReason = 156 // The Client should temporarily change its Server.
	DisconnectServerMoved                         DisconnectReason = 157 // The Server is moved and the Client should permanently change its server location.
	DisconnectSharedSubscriptionsNotSupported     DisconnectReason = 158 // The Server does not support Shared Subscriptions.
	DisconnectConnectionRateExceeded              DisconnectReason = 159 // This connection is closed because the connection rate is too high.
	DisconnectMaximumConnectTime                  DisconnectReason = 160 // The maximum connection time authorized for this connection has been exceeded.
	DisconnectSubscriptionIdentifiersNotSupported DisconnectReason = 161 // The Server does not support Subscription Identifiers; the subscription is not accepted.
	DisconnectWildcardSubscriptionsNotSupported   DisconnectReason = 162 // The Server does not support Wildcard Subscriptions; the subscription is not accepted.
)
