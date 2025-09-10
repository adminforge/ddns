package backend

import (
	"log"
	"strings"

	"github.com/adminforge/ddns/internal/shared"
	"github.com/gin-gonic/gin"
)

type Backend struct {
	config *shared.Config
	lookup *HostLookup
}

func NewBackend(config *shared.Config, lookup *HostLookup) *Backend {
	return &Backend{
		config: config,
		lookup: lookup,
	}
}

func (b *Backend) Run() error {
	r := gin.New()
	r.Use(gin.Recovery())

	if b.config.Verbose {
		r.Use(gin.Logger())
	}

	// Lookup Method, which is called by PowerDNS when resolving a Domain
	r.GET("/dnsapi/lookup/:qname/:qtype", func(c *gin.Context) {
		// Get DNS Zone Domain,
		host := strings.TrimPrefix(b.config.Domain, ".")

		// PowerDNS sends requests with "." in the end. Strip i toff
		request := &Request{
			QName: strings.TrimRight(c.Param("qname"), "."),
			QType: c.Param("qtype"),
		}

		// Lookup Domain in database
		response, err := b.lookup.Lookup(request)
		if err == nil {
			// When no error return record
			c.JSON(200, gin.H{
				"result": []*Response{response},
			})
		} else if request.QName == host {
			/*
				If the Requested Domain is the DNS Zone, then return SOA and NS record.
				This is needed to have a valid Zone configuration for the Zone-Cache in PowerDNS.
				If Requested Zone/Domain is not listed in Zone-Cache and feature is enabled,
				PowerDNS rejects the request.
			*/
			request.QType = "SOA"
			responseSOA, _ := b.lookup.Lookup(request)
			request.QType = "NS"
			responseNS, _ := b.lookup.Lookup(request)
			c.JSON(200, gin.H{
				"result": []*Response{responseSOA, responseNS},
			})
		} else {
			// Else case, means Domains is not found nor the Zone. Return "false"
			if b.config.Verbose {
				log.Printf("Error during lookup: %v", err)
			}

			c.JSON(200, gin.H{
				"result": false,
			})
		}
	})

	r.GET("/dnsapi/getDomainMetadata/:name/:kind", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": []string{"0"},
		})
	})

	r.GET("/dnsapi/getAllDomainMetadata/:name", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": gin.H{"PRESIGNED": []string{"0"}},
		})
	})

	// return the Zone-Domain, which is the ddns Domain where all the dynamic domains are served as subdomains from
	r.GET("/dnsapi/getAllDomains", func(c *gin.Context) {
		// get DNS Zone Domain and check if the Zone ends with an ".", if not add it on the end
		host := strings.TrimPrefix(b.config.Domain, ".")
		if !strings.HasSuffix(host, ".") {
			host += "."
		}

		c.JSON(200, gin.H{
			"result": []map[string]string{
				{"id": "1", "zone": host, "type": "NATIVE"},
			},
		})
	})

	// If PowerDNS requests DomainInfo for the Zone, return Zone Configuration, else "false"
	r.GET("/dnsapi/getDomainInfo/:name", func(c *gin.Context) {
		host := strings.TrimPrefix(b.config.Domain, ".")
		if !strings.HasSuffix(host, ".") {
			host += "."
		}

		if c.Param("name") == host {
			c.JSON(200, gin.H{
				"result": []map[string]string{
					{"id": "1", "zone": host, "type": "NATIVE"},
				},
			})
		} else {
			c.JSON(200, gin.H{
				"result": false,
			})
		}
	})

	return r.Run(b.config.ListenBackend)
}
