package packet

/*
Package packet defines all mqtt control packets.
Each control packet has a method WriteTo(io.Writer) (int64, error).
To read control packets the package exports the ReadPacket() packet.Packet function.
To construct new control packets use the constructors and setter methods.
*/
