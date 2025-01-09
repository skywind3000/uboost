// =====================================================================
//
// UdpInternal.go -
//
// Last Modified: 2025/01/09 16:20:11
//
// =====================================================================
package forward

type ForwardSide int

const (
	ForwardSideServer ForwardSide = 0
	ForwardSideClient ForwardSide = 1
)

var _ForwardSideNames = map[ForwardSide]string{
	ForwardSideServer: "server",
	ForwardSideClient: "client",
}

func (f ForwardSide) String() string {
	return _ForwardSideNames[f]
}
