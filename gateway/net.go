package main

// func listenTCP() {
// 	addr := fmt.Sprintf(":%d", config.SrvCfg.TCPPort)
// 	l, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		log.Printf("tcp listen failed: %s\n", err.Error())
// 		os.Exit(-1)
// 	}

// 	pname := utils.ProcessName()
// 	xlog.Infof("%s listen on %s", pname, l.Addr().String())

// 	go func() {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			xlog.Errorf("accept connection failed: %s", err.Error())
// 			continue
// 		}

// 		go handleConn(conn)
// 	}()
// }

// func handleConn(c net.Conn) {
// 	conn := netframe.NewTCPConn(c)
// 	conn.SetReadDeadline(time.Now().Add(time.Second))
// 	buff, err := conn.ReadHead()
// 	if err != nil {
// 		xlog.Debugf("read packet head failed: %s", err.Error())
// 		conn.Close()
// 		return
// 	}

// 	// TODO:
// 	// decode packet head, and get packet head length
// 	var bl int

// 	conn.SetReadDeadline(time.Now().Add(time.Second * 3))
// 	buff, err = conn.ReadBody(bl)
// 	if err != nil {
// 		conn.Close()
// 		return
// 	}
// }
