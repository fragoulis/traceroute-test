package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	// switch runtime.GOOS {
	// case "darwin", "ios":
	// case "linux":
	// 	log.Println("you may need to adjust the net.ipv4.ping_group_range kernel state")
	// default:
	// 	log.Println("not supported on", runtime.GOOS)
	// 	return
	// }

	replyBytes := make([]byte, 1500)
	for hop := 1; hop <= 30; hop++ { // up to 64 hops

		// Parse the destination IP address
		dstAddr, err := net.ResolveIPAddr("ip", "8.8.8.8")
		if err != nil {
			fmt.Printf("ResolveIPAddr failed: %v\n", err)
			os.Exit(1)
		}

		// Create a connection to the destination
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			fmt.Printf("ListenPacket failed: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		// Set the TTL on the socket
		// conn.IPv4PacketConn().SetTTL(hop)
		// p := ipv4.NewPacketConn(conn)
		p := conn.IPv4PacketConn()

		if err := p.SetTTL(hop); err != nil {
			fmt.Printf("SetTTL failed: %v\n", err)
			os.Exit(1)
		}

		if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
			fmt.Printf("SetControlMessage failed: %v\n", err)
			os.Exit(1)
		}

		// Create the ICMP echo message
		echo := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  hop,
				Data: []byte("hello router"),
			},
		}

		// echo.Body.(*icmp.Echo).Seq = i

		// Marshal the ICMP echo message
		echoBytes, err := echo.Marshal(nil)
		if err != nil {
			fmt.Printf("Echo Marshal failed: %v\n", err)
			os.Exit(1)
		}

		// Send the ICMP echo message to the destination
		start := time.Now()
		_, err = conn.WriteTo(echoBytes, dstAddr)
		if err != nil {
			fmt.Printf("WriteTo failed: %v\n", err)
			os.Exit(1)
		}

		// Wait for the ICMP echo reply from the destination
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("SetReadDeadline failed: %v\n", err)
			os.Exit(1)
		}

		// n, _, err := conn.ReadFrom(replyBytes)
		n, cm, peer, err := p.ReadFrom(replyBytes)
		if err != nil {
			fmt.Printf("ReadFrom failed: %v\n", err)
			os.Exit(1)
		}
		rtt := time.Since(start)

		// Parse the ICMP echo reply
		reply, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), replyBytes[:n])
		if err != nil {
			fmt.Printf("ParseMessage failed: %v\n", err)
			os.Exit(1)
		}

		_, err = p.TTL()
		if err != nil {
			fmt.Printf("TTL() failed: %v\n", err)
			os.Exit(1)
		}

		// In the real world you need to determine whether the
		// received message is yours using ControlMessage.Src,
		// ControlMessage.Dst, icmp.Echo.ID and icmp.Echo.Seq.
		switch reply.Type {
		case ipv4.ICMPTypeTimeExceeded:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d (e)\t%v %+v %v\n\t%+v\n", hop, peer, names, rtt, cm)
		case ipv4.ICMPTypeEchoReply:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d (r)\t%v %+v %v\n\t%+v\n", hop, peer, names, rtt, cm)
			return
		default:
			log.Printf("unknown ICMP message: %+v\n", reply)
		}
	}

	// Print the results
	// switch reply.Type {
	// case ipv4.ICMPTypeEchoReply:
	// 	fmt.Printf("Received ICMP echo reply from %s with TTL %d in %v:\n", dstAddr.String(), ttl, rtt)
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeDestinationUnreachable:
	// 	fmt.Printf("%s\n", "Destination Unreachable")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeRedirect:
	// 	fmt.Printf("%s\n", "Redirect")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeEcho:
	// 	fmt.Printf("%s\n", "Echo")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeRouterAdvertisement:
	// 	fmt.Printf("%s\n", "Router Advertisement")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeRouterSolicitation:
	// 	fmt.Printf("%s\n", "Router Solicitation")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeTimeExceeded:
	// 	fmt.Printf("%s\n", "Time Exceeded")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.TimeExceeded).Data))
	// case ipv4.ICMPTypeParameterProblem:
	// 	fmt.Printf("%s\n", "Parameter Problem")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeTimestamp:
	// 	fmt.Printf("%s\n", "Timestamp")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeTimestampReply:
	// 	fmt.Printf("%s\n", "Timestamp Reply")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypePhoturis:
	// 	fmt.Printf("%s\n", "Photuris")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeExtendedEchoRequest:
	// 	fmt.Printf("%s\n", "Extended Echo Request")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// case ipv4.ICMPTypeExtendedEchoReply:
	// 	fmt.Printf("%s\n", "Extended Echo Reply")
	// 	fmt.Printf("%s\n", hex.Dump(reply.Body.(*icmp.Echo).Data))
	// default:
	// 	fmt.Printf("Received unexpected ICMP message type %v\n", reply.Type)
	// }
}
