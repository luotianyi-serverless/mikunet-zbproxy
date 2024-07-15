//go:build !linux && !windows && !darwin && !freebsd

package network

func NewDialerControlFromOptions(*OutboundSocketOptions) ControlFunc {
	return nil
}

func NewListenerControlFromOptions(*InboundSocketOptions) ControlFunc {
	return nil
}
