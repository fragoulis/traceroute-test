// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"time"
// )

// func main() {
// 	if len(os.Args) != 2 {
// 		fmt.Println("Usage: traceroute <hostname or IP address>")
// 		return
// 	}

// 	dest := os.Args[1]
// 	maxTTL := 1
// 	timeout := time.Second * 3

// 	fmt.Printf("Traceroute to %s (%s), %d hops max:\n", dest, net.ParseIP(dest), maxTTL)

// 	for ttl := 1; ttl <= maxTTL; ttl++ {
// 		recvAddr := ""
// 		startTime := time.Now()

// 		// Set up the UDP packet with a TTL of 'ttl'
// 		icmpConn, err := net.ListenPacket("icmp", "0.0.0.0")
// 		if err != nil {
// 			fmt.Printf("Error listening for ICMP packets: %s\n", err)
// 			return
// 		}
// 		icmpConn.IPv4PacketConn().SetTTL(ttl)

// 		// Set up the UDP packet to send to the destination
// 		udpConn, err := net.DialTimeout("udp", dest+":80", timeout)
// 		if err != nil {
// 			fmt.Printf("%d\t%s\t%s\n", ttl, "*", "Could not connect to destination")
// 			continue
// 		}

// 		// Send the UDP packet with the specified TTL
// 		_, err = udpConn.Write([]byte("Hello"))
// 		if err != nil {
// 			fmt.Printf("Error sending UDP packet: %s\n", err)
// 			return
// 		}

// 		// Wait for an ICMP Time Exceeded or Destination Unreachable message
// 		icmpBuf := make([]byte, 1500)
// 		icmpConn.SetReadDeadline(time.Now().Add(timeout))
// 		_, _, err = icmpConn.ReadFrom(icmpBuf)
// 		if err != nil {
// 			fmt.Printf("%d\t%s\t%s\n", ttl, "*", "Request timed out")
// 			continue
// 		}
// 		icmpType := icmpBuf[0]
// 		if icmpType == 11 {
// 			recvAddr = fmt.Sprintf("%s", icmpBuf[24:28])
// 		} else if icmpType == 3 {
// 			recvAddr = fmt.Sprintf("%s", icmpBuf[12:16])
// 		}

// 		// Calculate the round-trip time and print out the result
// 		elapsedTime := time.Since(startTime)
// 		fmt.Printf("%d\t%s\t%.3fms\n", ttl, recvAddr, elapsedTime.Seconds()*1000)

// 		// Close the connections
// 		udpConn.Close()
// 		icmpConn.Close()
// 	}
// }
