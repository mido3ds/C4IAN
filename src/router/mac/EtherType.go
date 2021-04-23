package mac

import "github.com/mdlayher/ethernet"

type EtherType ethernet.EtherType

const (
	IPv4EtherType = EtherType(0x0800)

	// Make use of an unassigned EtherType to differentiate between odmrp traffic and other traffic
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	JoinQueryEtherType = EtherType(0x0901)
	JoinReplyEtherType = EtherType(0x0902)

	SARPReqEtherType = EtherType(0x0809)
	SARPResEtherType = EtherType(0x080A)

	ZIDEtherType = EtherType(0x7031)
)
