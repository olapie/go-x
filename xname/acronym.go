package xname

import (
	"strings"
	"sync"
)

var acronymsMapRWMutex sync.RWMutex
var acronymsU2L = map[string]string{
	"TCP":   "tcp",
	"HTTP":  "http",
	"UDP":   "udp",
	"ID":    "id",
	"SSL":   "ssl",
	"TLS":   "tls",
	"CPU":   "cpu",
	"DOB":   "dob",
	"TTL":   "ttl",
	"SSO":   "sso",
	"HTTPS": "https",
	"IP":    "ip",
	"XSS":   "xss",
	"OS":    "os",
	"SIP":   "sip",
	"XML":   "xml",
	"JSON":  "json",
	"HTML":  "html",
	"XHTML": "xhtml",
	"XSL":   "xsl",
	"XSLT":  "xslt",
	"YAML":  "yaml",
	"TOML":  "toml",
	"WLAN":  "wlan",
	"WIFI":  "wifi",
	"VM":    "vm",
	"JVM":   "jvm",
	"UI":    "ui",
	"URI":   "uri",
	"URL":   "url",
	"SLA":   "sla",
	"SCP":   "scp",
	"SMTP":  "smtp",
	"SOA":   "soa",
	"OA":    "oa",
	"SVG":   "svg",
	"PNG":   "png",
	"JPG":   "jpg",
	"JPEG":  "jpeg",
	"PDF":   "pdf",
	"IO":    "io",
}
var acronymsL2U = map[string]string{
	"tcp":   "TCP",
	"http":  "HTTP",
	"udp":   "UDP",
	"id":    "ID",
	"ssl":   "SSL",
	"tls":   "TLS",
	"cpu":   "CPU",
	"dob":   "DOB",
	"ttl":   "TTL",
	"sso":   "SSO",
	"https": "HTTPS",
	"ip":    "IP",
	"xss":   "XSS",
	"os":    "OS",
	"sip":   "SIP",
	"xml":   "XML",
	"json":  "JSON",
	"html":  "HTML",
	"xhtml": "XHTML",
	"xsl":   "XSL",
	"xslt":  "XSLT",
	"yaml":  "YAML",
	"toml":  "TOML",
	"wlan":  "WLAN",
	"wifi":  "WIFI",
	"vm":    "VM",
	"jvm":   "JVM",
	"ui":    "UI",
	"uri":   "URI",
	"url":   "URL",
	"sla":   "SLA",
	"scp":   "SCP",
	"smtp":  "SMTP",
	"soa":   "SOA",
	"oa":    "OA",
	"svg":   "SVG",
	"png":   "PNG",
	"jpg":   "JPG",
	"jpeg":  "JPEG",
	"pdf":   "PDF",
	"io":    "IO",
}

// AddAcronym adds acronyms for Pascal and Camel cases
func AddAcronym(acronyms ...string) {
	acronymsMapRWMutex.Lock()
	for _, acronym := range acronyms {
		lc := strings.ToLower(acronym)
		acronymsL2U[lc] = acronym
		acronymsU2L[acronym] = lc
	}
	acronymsMapRWMutex.Unlock()
}

// RemoveAcronym removes acronyms for Pascal and Camel cases
func RemoveAcronym(acronyms ...string) {
	acronymsMapRWMutex.Lock()
	for _, acronym := range acronyms {
		delete(acronymsL2U, acronymsU2L[acronym])
		delete(acronymsU2L, acronym)
	}
	acronymsMapRWMutex.Unlock()
}
