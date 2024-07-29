package telnet

import (
	"io"
)

//go:generate stringer -type Operator

type Operator byte

const (
	Se Operator = 0xF0 + iota
	Nop
	DataMark
	Break
	Interrupt
	Abort
	AreYouThere
	EraseChar
	EraseLine
	GoAhead
	SubNegotiation
	Will
	Wont
	Do
	Dont
	Iac
)

const (
	BinaryTransmission Operator = iota
	Echo
	Reconnection
	SuppressGoAhead
	ApproxMessageSizeNegotiation
	Status
	TimingMark
	RemoteControlledTransAndEcho
	OutputLineWidth
	OutputPageSize
	OutputCarriageReturnDisposition
	OutputHorizontalTabStops
	OutputHorizontalTabDisposition
	OutputFormFeedDisposition
	OutputVerticalTabStops
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
	EndOfRecord
	TACACSUserIdentification
	OutputMarking
	TerminalLocationNumber
	Telnet3270Regime
	X3PAD
	NegotiateAboutWindowSize
	TerminalSpeed
	RemoteFlowControl
	LineMode
	XDisplayLocation
	ExtendedOptionsList
)

func Bytes(cmds ...Operator) []byte {
	bytes := make([]byte, 0, len(cmds))
	for _, cmd := range cmds {
		bytes = append(bytes, byte(cmd))
	}
	return bytes
}

func Write(w io.Writer, cmds ...Operator) (int, error) {
	return w.Write(Bytes(cmds...))
}

func WriteAndClear(w io.Writer, cmds ...Operator) (int, error) {
	return w.Write(append(Bytes(cmds...), []byte("\r\x1B[K")...))
}
