package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type TcpConn struct {
	conn net.Conn
}

func (conn *TcpConn) Read(b []byte) (int, error) {
	return conn.conn.Read(b)
}

func (conn *TcpConn) bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func (conn *TcpConn) ReadMsg() ([]byte, error) {
	head := make([]byte, 2)
	_, err := io.ReadFull(conn, head)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint16(head)
	data := make([]byte, size)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	return conn.bytesCombine(head, data), nil
}

func (conn *TcpConn) WriteMsg(arg []byte) error {
	_, err := conn.conn.Write(arg)
	return err
}

func (conn *TcpConn) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

func (conn *TcpConn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

func (conn *TcpConn) SetReadDeadline(t time.Time) error {
	return conn.conn.SetReadDeadline(t)
}

func (conn *TcpConn) SetWriteDeadline(t time.Time) error {
	return conn.conn.SetWriteDeadline(t)
}

func (conn *TcpConn) Close() {
	// default drop data
	conn.conn.(*net.TCPConn).SetLinger(0)
	conn.conn.Close()
}
