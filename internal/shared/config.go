package shared

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Verbose            bool
	Domain             string
	SOAFqdn            string
	HostExpirationDays int
	ListenFrontend     string
	ListenBackend      string
	RedisHost          string
	RecordTTL          int
}

func (c *Config) Initialize() {
	flag.StringVar(&c.Domain, "domain", "",
		"The subdomain which should be handled by DDNS")

	flag.StringVar(&c.SOAFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	flag.StringVar(&c.ListenBackend, "listen-backend", ":8053",
		"Which socket should the backend web service use to bind itself")

	flag.StringVar(&c.ListenFrontend, "listen-frontend", ":8080",
		"Which socket should the frontend web service use to bind itself")

	flag.StringVar(&c.RedisHost, "redis", ":6379",
		"The Redis socket that should be used")

	flag.IntVar(&c.HostExpirationDays, "expiration-days", 90,
		"The number of days after a host is released when it is not updated")

	flag.IntVar(&c.RecordTTL, "dns-ttl", 300,
		"The TTL for a Record, which is returned by the ddns service")

	flag.BoolVar(&c.Verbose, "verbose", false,
		"Be more verbose")

	flag.Parse()
	parseEnvConfig(c)
}

func (c *Config) Validate() {

	if c.Domain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(c.Domain, ".") {
		// get the domain in the right format
		c.Domain = "." + c.Domain
	}

	if c.SOAFqdn == "" {
		log.Fatal("You have to supply the server FQDN via --soa_fqdn=FQDN")
	}
}

// parseEnvConfig overwrites configuration with values from the environment.
func parseEnvConfig(cfg *Config) {
	domain, got := os.LookupEnv("DDNS_DOMAIN")
	if got {
		cfg.Domain = domain
	}
	soaDomain, got := os.LookupEnv("DDNS_SOA_DOMAIN")
	if got {
		cfg.SOAFqdn = soaDomain
	}
	redisHost, got := os.LookupEnv("DDNS_REDIS_HOST")
	if got {
		cfg.RedisHost = redisHost
	}
	expDays, got := os.LookupEnv("DDNS_EXPIRATION_DAYS")
	if got {
		days, err := strconv.Atoi(expDays)
		if err != nil {
			panic(fmt.Errorf("Unexpected value for 'DDNS_EXPIRATION_DAYS' '%s': %w", expDays, err))
		}
		cfg.HostExpirationDays = days
	}
	ttl, got := os.LookupEnv("DDNS_TTL")
	if got {
		ttlInt, err := strconv.Atoi(ttl)
		if err != nil {
			panic(fmt.Errorf("Unexpected value for 'DDNS_TTL' '%s': %w", ttl, err))
		}
		cfg.RecordTTL = ttlInt
	}

}
