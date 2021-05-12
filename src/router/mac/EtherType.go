package mac

import "github.com/mdlayher/ethernet"

type EtherType ethernet.EtherType

const (
	IPv4EtherType = EtherType(0x0800)

	// Make use of an unassigned EtherType to differentiate between different types of traffic
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	JoinQueryEtherType = EtherType(0x0901)
	JoinReplyEtherType = EtherType(0x0902)

	SARPEtherType = EtherType(0x0809)

	DZRequestEtherType = EtherType(0x080A)

	DZResponseEtherType = EtherType(0x080B)

	ZoneFloodEtherType = EtherType(0x7033)

	ZIDDataEtherType = EtherType(0x7031)

	InterzoneLSREtherType = EtherType(0x7032)
)
