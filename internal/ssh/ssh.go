package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/git-proxy/go/internal/config"
	"github.com/git-proxy/go/internal/logger"
)

type Proxy struct {
	config   *config.Config
	logger   logger.Logger
	listener net.Listener
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(cfg *config.Config, log logger.Logger) *Proxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proxy{
		config: cfg,
		logger: log,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Proxy) Start() error {
	if !p.config.Target.SSH.Enabled {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", p.config.Server.Host, p.config.Target.SSH.Port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	p.listener = ln

	p.wg.Add(1)
	go p.handleConnections()

	p.logger.Info("SSH proxy listening on %s", addr)
	return nil
}

func (p *Proxy) handleConnections() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		p.listener.(*net.TCPListener).SetDeadline(time.Now().Add(10 * time.Second))

		conn, err := p.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
			if err != io.EOF {
				p.logger.Error("SSH accept error: %v", err)
			}
			return
		}

		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.handleConnection(conn)
		}()
	}
}

func (p *Proxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	p.logger.Debug("New SSH connection from %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		p.logger.Error("Failed to read SSH version: %v", err)
		return
	}

	if !strings.HasPrefix(line, "SSH-") {
		p.logger.Warn("Invalid SSH version string: %s", strings.TrimSpace(line))
		return
	}

	p.logger.Debug("SSH version: %s", strings.TrimSpace(line))

	targetAddr := fmt.Sprintf("%s:%d", p.config.Target.Host, p.config.Target.SSH.Port)
	targetConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		p.logger.Error("Failed to connect to target SSH: %v", err)
		return
	}
	defer targetConn.Close()

	targetConn.Write([]byte(line))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(targetConn, conn)
		targetConn.Close()
	}()

	go func() {
		defer wg.Done()
		io.Copy(conn, targetConn)
		conn.Close()
	}()

	wg.Wait()
	p.logger.Debug("SSH connection closed")
}

func (p *Proxy) Stop() {
	if p.listener != nil {
		p.listener.Close()
	}
	p.cancel()
	p.wg.Wait()
}
