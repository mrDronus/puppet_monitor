package impl

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
)

type Service struct {
	waitGroup *sync.WaitGroup
	listener  *net.TCPListener
	env       *EnvironmentCollection
}

func (s *Service) SetEnvCollection(envs *EnvironmentCollection) {
	s.env = envs
}

func (s *Service) Stop() {
	s.waitGroup.Wait()
	s.listener.Close()
}

func (s Service) NewService() *Service {
	var srv = &Service{
		waitGroup: &sync.WaitGroup{},
	}
	return srv
}

func (s *Service) HandleListener(listener *net.TCPListener) {
	s.listener = listener
	for {
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		//log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)
		go s.HandleConnection(conn)
	}
}

const (
	timeout = 5 * time.Second
)

func (s *Service) HandleConnection(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	conn.SetDeadline(time.Now().Add(timeout))
	data, _ := ioutil.ReadAll(conn)

	res := s.env.ProcessReport(data)
	if !res {
		w := bufio.NewWriter(conn)
		w.WriteString("HTTP/1.1 404 Not Found\r\n\r\n")
		w.Flush()
		/*_, err := fmt.Fprint(conn, )
		if err != nil {
			log.Println(err)
		}*/
		//conn.Write([]byte("HTTP/1.0 404 Not Found"))
		//conn.CloseWrite()
	}
}
