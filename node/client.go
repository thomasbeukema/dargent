package node

type client struct {
    dest    string
    conn net.UDPConn
}
