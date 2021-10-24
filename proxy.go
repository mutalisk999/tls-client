package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"net"
	"sort"
)

func handleTcpProxyConn(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	var connToTarget *tls.Conn = nil
	var targetCopy LBTargetCopy

	conn := a.(*net.TCPConn)
	targetsCopy := LBTargetsMgrP.DumpTargetsCopy()
	sort.Sort(LBTargetCopys(targetsCopy))

	cert, err := tls.LoadX509KeyPair(LBConfig.Tls.TlsCert, LBConfig.Tls.TlsKey)
	if err != nil {
		Error.Fatalf("LoadX509KeyPair fail: %s", err)
	}

	for _, t := range targetsCopy {
		if !t.Active {
			continue
		}
		if t.ConnCount >= t.MaxConnCount {
			continue
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		connToTarget, err = tls.Dial("tcp", t.EndPointConn, &config)
		if err != nil {
			Warn.Printf("TLS Dial fail: %s", err)
			continue
		}
		Info.Printf("TLS Dial Connected To: %s Success!", conn.RemoteAddr().String())

		state := connToTarget.ConnectionState()
		for _, v := range state.PeerCertificates {
			pKIXPublicKey, _ := x509.MarshalPKIXPublicKey(v.PublicKey)
			Info.Printf("pKIXPublicKey: %s", hex.EncodeToString(pKIXPublicKey))
		}
		Info.Println("TLS: Handshake Complete: ", state.HandshakeComplete)
		Info.Println("TLS: Mutual: ", state.NegotiatedProtocolIsMutual)

		break
	}

	if connToTarget == nil {
		Warn.Printf("Can not connect to any target endpoint, Close node connection")
		err := conn.Close()
		if err == nil {
			LBNodeP.ConsumeNewConn()
		}

	} else {
		targetId := calcTargetId(targetCopy.EndPointConn)

		var nodeConn NodeConnection
		nodeConn.Initialise(conn, LBNodeP.timeout)

		var targetConn TargetConnection
		targetConn.Initialise(connToTarget, targetCopy.Timeout, targetId)

		LBConnectionPairMgrP.AddConnectionPair(&nodeConn, &targetConn)

		LBGoroutineManagerP.GoroutineCreateP2("tcp_node_data", handleNodeData, &nodeConn, &targetConn)
		LBGoroutineManagerP.GoroutineCreateP2("tcp_target_data", handleTargetData, &nodeConn, &targetConn)
	}
}

func startTcpProxy(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	cfg := a.(*Config)
	addr, err := net.ResolveTCPAddr("tcp", cfg.Node.ListenEndPoint)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	defer server.Close()
	Info.Printf("Node listening on %s", cfg.Node.ListenEndPoint)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}
		_ = conn.SetKeepAlive(true)

		// TODO
		//ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		// banned and need to ban

		LBNodeP.ProductNewConn()
		LBGoroutineManagerP.GoroutineCreateP1("tcp_proxy_conn", handleTcpProxyConn, conn)
	}
}
