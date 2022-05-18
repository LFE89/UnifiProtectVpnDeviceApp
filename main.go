// Copyright (C) 2022 Lars D. Feicho
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
// GNU General Public License for more details.
// You should have received a copy of the GNU General Public License
// along with this program.If not, see<http://www.gnu.org/licenses/>.

package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

var discoveryPayload []byte

func main() {
	payloadHex := flag.String("payload", "", "Discovery payload (hex value, e.g. 0100009102000...)")
	flag.Parse()

	if *payloadHex == "" {
		fmt.Println("Discovery payload cannot be empty.")
		os.Exit(1)
	}

	payload, err := hex.DecodeString(*payloadHex)
	if err != nil {
		fmt.Println("Ooops. Something went wrong.")
		fmt.Println(err)
		os.Exit(1)
	}

	copy(discoveryPayload, payload)
	startUdpListening()
}

func startUdpListening() {
	connection, err := net.ListenPacket("udp4", "0.0.0.0:1338")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer connection.Close()

	for {

		// wait for 6 byte udp packet
		// byte 0-3 (4 byte): destination ip address
		// byte 4-5 (2 byte): destination port

		buf := make([]byte, 6)
		_, _, err := connection.ReadFrom(buf)
		if err != nil {
			continue
		}

		// Improve failure handling here...

		destinationPort := binary.LittleEndian.Uint16([]byte{buf[5], buf[4]})
		destinationIpAddress := net.IP{buf[0], buf[1], buf[2], buf[3]}

		fmt.Println(destinationIpAddress.String())
		fmt.Println(destinationPort)

		go sendUnifiDiscoverResponse(destinationIpAddress.String() + ":" + strconv.Itoa(int(destinationPort)))
	}
}

func sendUnifiDiscoverResponse(address string) {

	// Just for "reasons", the UDP packet is being sent 5 times

	for i := 1; i < 6; i++ {
		conn, _ := net.Dial("udp4", address)
		conn.Write(discoveryPayload)

	}
}
