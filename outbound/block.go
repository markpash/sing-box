package outbound

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

var _ adapter.Outbound = (*Block)(nil)

type Block struct {
	myOutboundAdapter
}

func NewBlock(logger log.ContextLogger, tag string) *Block {
	return &Block{
		myOutboundAdapter{
			protocol: C.TypeBlock,
			network:  []string{N.NetworkTCP, N.NetworkUDP},
			logger:   logger,
			slogger: slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
						if (a.Key == slog.TimeKey || a.Key == slog.LevelKey) && len(groups) == 0 {
							return slog.Attr{} // remove excess keys
						}
						return a
					},
				}),
			),
			tag: tag,
		},
	}
}

func (h *Block) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	h.logger.InfoContext(ctx, "blocked connection to ", destination)
	return nil, io.EOF
}

func (h *Block) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	h.logger.InfoContext(ctx, "blocked packet connection to ", destination)
	return nil, io.EOF
}

func (h *Block) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	conn.Close()
	h.logger.InfoContext(ctx, "blocked connection to ", metadata.Destination)
	h.slogger.Info("new connection",
		"inbound", metadata.Inbound,
		"outbound", h.Tag(),
		"user", metadata.User,
		"transport", metadata.Network,
		"protocol", metadata.Protocol,
		"source_ip", metadata.Source.Addr,
		"source_port", metadata.Source.Port,
		"destination_ip", metadata.Destination.Addr,
		"destination_hostname", metadata.Destination.Fqdn,
		"destination_port", metadata.Destination.Port,
	)
	return nil
}

func (h *Block) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	conn.Close()
	h.logger.InfoContext(ctx, "blocked packet connection to ", metadata.Destination)
	h.slogger.Info("new packet connection",
		"inbound", metadata.Inbound,
		"outbound", h.Tag(),
		"user", metadata.User,
		"transport", metadata.Network,
		"protocol", metadata.Protocol,
		"source_ip", metadata.Source.Addr,
		"source_port", metadata.Source.Port,
		"destination_ip", metadata.Destination.Addr,
		"destination_hostname", metadata.Destination.Fqdn,
		"destination_port", metadata.Destination.Port,
	)
	return nil
}
