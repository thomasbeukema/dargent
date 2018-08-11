package node

type client struct {
    dest    string
    conn net.UDPConn
}

func (c *client) Send(msg string) {

}
