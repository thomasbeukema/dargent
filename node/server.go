package node

type server struct {
	addr	*net.UDPAddr
	conn	net.UDPConn
}

func newServer(host string) {
	//server, _ := net.ListenUDP("udp", host)
}
