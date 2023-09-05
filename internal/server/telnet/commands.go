package telnet

const (
	Se = 0xF0 + iota
	Nop
	DataMark
	Break
	Interrupt
	Abort
	AreYouThere
	EraseChar
	EraseLine
	GoAhead
	Subnegotiation
	Will
	Wont
	Do
	Dont
	Iac
)

const (
	BinaryTransmission = iota
	Echo
	Reconnection
	SuppressGoAhead
	ApproxMessageSizeNegotiation
	Status
	TimingMark
	RemoteControlledTransandEcho
	OutputLineWidth
	OutputPageSize
	OutputCarriageReturnDisposition
	OutputHorizontalTabStops
	OutputHorizontalTabDisposition
	OutputFormfeedDisposition
	OutputVerticalTabstops
	OutputVerticalTabDisposition
	OutputLinefeedDisposition
	ExtendedASCII
	Logout
	ByteMacro
	DataEntryTerminal
	SUPDUP
	SUPDUPOutput
	SendLocation
	TerminalType
	EndofRecord
	TACACSUserIdentification
	OutputMarking
	TerminalLocationNumber
	Telnet3270Regime
	X3PAD
	NegotiateAboutWindowSize
	TerminalSpeed
	RemoteFlowControl
	Linemode
	XDisplayLocation
	ExtendedOptionsList
)
