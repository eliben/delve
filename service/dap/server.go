package dap

import (
	"fmt"
	"net"

	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/debugger"
)

type DAPServer struct {
	config   *service.Config
	listener net.Listener
	stopChan chan struct{}
	debugger *debugger.Debugger
}

func NewServer(config *service.Config) *DAPServer {
	return &DAPServer{
		config:   config,
		listener: config.Listener,
		stopChan: make(chan struct{})}
}

func (s *DAPServer) Run() error {
	var err error
	// TODO(eliben): it's likely that we don't necessarily want to start the
	// debugger immediately here. We might need to wait for a launch request with
	// all the parameters.
	if s.debugger, err = debugger.New(&debugger.Config{
		AttachPid:            s.config.AttachPid,
		WorkingDir:           s.config.WorkingDir,
		CoreFile:             s.config.CoreFile,
		Backend:              s.config.Backend,
		Foreground:           s.config.Foreground,
		DebugInfoDirectories: s.config.DebugInfoDirectories,
		CheckGoVersion:       s.config.CheckGoVersion,
	},
		s.config.ProcessArgs); err != nil {
		return err
	}

	go func() {
		defer s.listener.Close()
		for {
			fmt.Println("$$ accepting")
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.stopChan:
					return
				default:
					panic(err)
				}
			}

			// TODO(eliben): to use canAccept here, refactor it to a more accessible
			// place (maybe at root of service? also export it)

			// TODO(eliben): actual serving goes here
			go func() { _ = conn }()

			if !s.config.AcceptMulti {
				break
			}
		}
	}()

	return nil
}

func (s *DAPServer) Stop() error {
	if s.config.AcceptMulti {
		close(s.stopChan)
		s.listener.Close()
	}
	kill := s.config.AttachPid == 0
	return s.debugger.Detach(kill)
}
