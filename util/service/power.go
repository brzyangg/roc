package rocserv

import (
    "net"
    "net/http"


    "github.com/julienschmidt/httprouter"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/shawnfeng/sutil/snetutil"
	"github.com/shawnfeng/sutil/slog"
)


func powerHttp(addr string, router *httprouter.Router) (string, error) {
	fun := "powerHttp"

	if len(addr) == 0 {
		addr = ":"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return "", err
	}

	netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if err != nil {
		return "", err
	}

	laddr, err := snetutil.GetServAddr(netListen.Addr())
	if err != nil {
		netListen.Close()
		return "", err
	}

	go func() {
		err := http.Serve(netListen, router)
		if err != nil {
			slog.Panicf("%s --> laddr[%s]", fun, laddr)
		}
	}()

	return laddr, nil
}



func powerThrift(addr string, processor thrift.TProcessor) (string, error) {
	fun := "powerThrift"

	if len(addr) == 0 {
		addr = ":"
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return "", err
	}

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	// Listen后就可以拿到端口了
	//err = server.Listen()
	err = serverTransport.Listen()
	if err != nil {
		return "", err
	}

	laddr, err := snetutil.GetServAddr(serverTransport.Addr())
	if err != nil {
		return "", err
	}

	go func() {
		err := server.Serve()
		if err != nil {
			slog.Panicf("%s --> laddr[%s]", fun, laddr)
		}
	}()


	return laddr, nil

}