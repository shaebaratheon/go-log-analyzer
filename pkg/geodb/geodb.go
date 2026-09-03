package geodb

import (
	"fmt"
	"net"
	"sync"
)

type Location struct {
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type trieNode struct {
	children [2]*trieNode
	location *Location
}

// IPPrefixTrie indexes IP subnets for fast O(32) longest-prefix-match lookup.
type IPPrefixTrie struct {
	mu   sync.RWMutex
	root *trieNode
}

func NewIPPrefixTrie() *IPPrefixTrie {
	return &IPPrefixTrie{root: &trieNode{}}
}

func (t *IPPrefixTrie) InsertCIDR(cidr string, loc *Location) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR %s: %w", cidr, err)
	}

	ipv4 := ipNet.IP.To4()
	if ipv4 == nil {
		return fmt.Errorf("only IPv4 currently supported: %s", cidr)
	}

	maskOnes, _ := ipNet.Mask.Size()

	t.mu.Lock()
	defer t.mu.Unlock()

	curr := t.root
	for i := 0; i < maskOnes; i++ {
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		bit := (ipv4[byteIdx] >> bitIdx) & 1

		if curr.children[bit] == nil {
			curr.children[bit] = &trieNode{}
		}
		curr = curr.children[bit]
	}
	curr.location = loc
	return nil
}

func (t *IPPrefixTrie) Lookup(ipStr string) *Location {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	curr := t.root
	var lastMatched *Location

	for i := 0; i < 32; i++ {
		if curr == nil {
			break
		}
		if curr.location != nil {
			lastMatched = curr.location
		}
		byteIdx := i / 8
		bitIdx := 7 - (i % 8)
		bit := (ipv4[byteIdx] >> bitIdx) & 1
		curr = curr.children[bit]
	}

	if curr != nil && curr.location != nil {
		lastMatched = curr.location
	}

	return lastMatched
}
