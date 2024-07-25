package network

import (
	"errors"
	"fmt"
	"github.com/grandcat/zeroconf"
	"golang.org/x/net/context"
	"time"
)

// discoverService uses mDNS to find the FloppyDaemon service
func DiscoverService() (string, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return "", fmt.Errorf("Failed to initialize resolver: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	var serverAddr string

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Printf("\nFound service: %s\n", entry.ServiceRecord.Instance)
			if len(entry.AddrIPv4) > 0 {
				serverAddr = fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port)
				fmt.Printf("%s:%d \n", entry.AddrIPv4[0], entry.Port)
				break
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)
	if err != nil {
		return "", fmt.Errorf("Failed to browse: %v", err)
	}

	<-ctx.Done()
	if serverAddr == "" {
		return "", errors.New("service discovery failed or timed out")
	}

	return serverAddr, nil
}
