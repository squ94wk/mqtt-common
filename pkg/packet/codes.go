package packet

type ConnectReason byte

const (
	Success                     ConnectReason = 0   // The Connection is accepted.
	UnspecifiedError                          = 128 // The Server does not wish to reveal the reason for the failure, or none of the other Reason Codes apply.
	MalformedPacket                           = 129 // Data within the CONNECT packet could not be correctly parsed.
	ProtocolError                             = 130 // Data in the CONNECT packet does not conform to this specification.
	ImplementationSpecificError               = 131 // The CONNECT is valid but is not accepted by this Server.
	UnsupportedProtocolVersion                = 132 // The Server does not support the version of the MQTT protocol requested by the Client.
	ClientIdentifierNotValid                  = 133 // The Client Identifier is a valid string but is not allowed by the Server.
	BadUserNameOrPassword                     = 134 // The Server does not accept the User Name or Password specified by the Client
	NotAuthorized                             = 135 // The Client is not authorized to connect.
	ServerUnavailable                         = 136 // The MQTT Server is not available.
	ServerBusy                                = 137 // The Server is busy. Try again later.
	Banned                                    = 138 // This Client has been banned by administrative action. Contact the server administrator.
	BadAuthenticationMethod                   = 140 // The authentication method is not supported or does not match the authentication method currently in use.
	TopicNameInvalid                          = 144 // The Will Topic Name is not malformed, but is not accepted by this Server.
	PacketTooLarge                            = 149 // The CONNECT packet exceeded the maximum permissible size.
	QuotaExceeded                             = 151 // An implementation or administrative imposed limit has been exceeded.
	PayloadFormatInvalid                      = 153 // The Will Payload does not match the specified Payload Format Indicator.
	RetainNotSupported                        = 154 // The Server does not support retained messages, and Will Retain was set to 1.
	QoSNotSupported                           = 155 // The Server does not support the QoS set in Will QoS.
	UseAnotherServer                          = 156 // The Client should temporarily use another server.
	ServerMoved                               = 157 // The Client should permanently use another server.
	ConnectionRateExceeded                    = 159 // The connection rate limit has been exceeded.
)
