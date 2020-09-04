package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	pb "github.com/usher2/u2byckbot/msg"
)

const (
	TBLOCK_UNKNOWN = iota
	TBLOCK_URL
	TBLOCK_HTTPS
	TBLOCK_DOMAIN
	TBLOCK_IP
)

const encodeCorc = "abcdefghjkmnpqrstvwxyz0123456789"

const PRINT_LIMIT = 5

const MAX_TIMESTAMP int64 = 1<<63 - 1

const (
	DTIME12 = 12 * 60 * 60
	DTIME3  = 3 * 60 * 60
)

const (
	_ = iota
	OFFSET_CONTENT
	OFFSET_IP4
	OFFSET_URL
	OFFSET_DOMAIN
)

type TPagination struct {
	Tag, Count int
}

type TReason struct {
	Id     int32
	Ip     []string
	Url    []string
	Domain []string
}

func printUpToDate(t int64) string {
	var r rune
	d := time.Now().Unix() - t
	switch {
	case d > DTIME12:
		r = 0x2b55
	case d > DTIME3:
		r = 0x000026a0
	default:
		r = 0x2705
	}
	return fmt.Sprintf("\n%c _Данные синхронизированы:_ %s\n", r, time.Unix(t, 0).In(time.FixedZone("UTC+3", 3*60*60)).Format(time.RFC3339))
}

func constructContentResult(a []*pb.Content, o TPagination) (res string, pages []TPagination) {
	var oldest int64 = MAX_TIMESTAMP
	if len(a) == 0 {
		return
	}
	for _, packet := range a {
		content := TContent{}
		err := json.Unmarshal(packet.Pack, &content)
		if err != nil {
			Error.Printf("Упс!!! %s\n", err)
			continue
		}
		if packet.RegistryUpdateTime < oldest {
			oldest = packet.RegistryUpdateTime
		}
		bt := ""
		switch packet.BlockType {
		case TBLOCK_URL:
			bt = "\U000026d4 (url) "
		case TBLOCK_HTTPS:
			bt = "\U0001f4db (https) "
		case TBLOCK_DOMAIN:
			bt = "\U0001f6ab (domain) "
		case TBLOCK_IP:
			bt = "\u274c (ip) "
		}
		res += fmt.Sprintf("%s /n\\_%d %s...\n", bt, content.Id, content.Description)
		res += fmt.Sprintf("внесено: %s\n", time.Unix(content.IncludeTime, 0).Format("2006-01-02"))
		res += "\n"
		cnt := 0
		for i, d := range content.Domain {
			if o.Tag == OFFSET_DOMAIN && i < o.Count {
				continue
			}
			if cnt >= PRINT_LIMIT {
				break
			}
			res += fmt.Sprintf("  domain: %s\n", Sanitize(d.Domain))
			cnt++
		}
		l := len(content.Domain)
		if l > PRINT_LIMIT {
			offset := 0
			if o.Tag == OFFSET_DOMAIN {
				switch {
				case l <= PRINT_LIMIT:
					offset = 0
				case o.Count > l-(l%PRINT_LIMIT):
					offset = l - (l % PRINT_LIMIT)
				default:
					offset = o.Count
				}
			}
			res += fmt.Sprintf("  \u2195 результаты с *%d* по *%d* из *%d*\n", offset+1, offset+cnt, l)
			pages = append(pages, TPagination{OFFSET_DOMAIN, l})
		}
		if cnt > 0 {
			res += "\n"
		}
		cnt = 0
		for i, u := range content.Url {
			if o.Tag == OFFSET_URL && i < o.Count {
				continue
			}
			if cnt >= PRINT_LIMIT {
				break
			}
			res += fmt.Sprintf("  url: %s\n", Sanitize(u.Url))
			cnt++
		}
		l = len(content.Url)
		if l > PRINT_LIMIT {
			offset := 0
			if o.Tag == OFFSET_URL {
				switch {
				case l <= PRINT_LIMIT:
					offset = 0
				case o.Count > l-(l%PRINT_LIMIT):
					offset = l - (l % PRINT_LIMIT)
				default:
					offset = o.Count
				}
			}
			res += fmt.Sprintf("  \u2195 результаты с *%d* по *%d* из *%d*\n", offset+1, offset+cnt, l)
			pages = append(pages, TPagination{OFFSET_URL, l})
		}
		if cnt > 0 {
			res += "\n"
		}
		cnt = 0
		for i, ip := range content.Ip4 {
			if o.Tag == OFFSET_IP4 && i < o.Count {
				continue
			}
			if cnt >= PRINT_LIMIT {
				break
			}
			res += fmt.Sprintf("  IP: %s\n", int2Ip4(ip.Ip4))
			cnt++
		}
		l = len(content.Ip4)
		if l > PRINT_LIMIT {
			offset := 0
			if o.Tag == OFFSET_IP4 {
				switch {
				case l <= PRINT_LIMIT:
					offset = 0
				case o.Count > l-(l%PRINT_LIMIT):
					offset = l - (l % PRINT_LIMIT)
				default:
					offset = o.Count
				}
			}
			res += fmt.Sprintf("  \u2195 результаты с *%d* по *%d* из *%d*\n", offset+1, offset+cnt, l)
			pages = append(pages, TPagination{OFFSET_IP4, l})
		}
		break
	}
	res += printUpToDate(oldest)
	return
}

func constructResult(a []*pb.Content, o TPagination) (res string, pages []TPagination) {
	var mass string
	var oldest int64 = MAX_TIMESTAMP
	var ra []TReason
	if len(a) == 0 {
		return
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Id < a[j].Id
	})
	ra = make([]TReason, 1)
	ra[0].Id = a[0].Id
	if a[0].Ip4 != 0 {
		ra[0].Ip = append(ra[0].Ip, int2Ip4(a[0].Ip4))
	}
	if a[0].Domain != "" {
		ra[0].Domain = append(ra[0].Domain, PrintedDomain(a[0].Domain))
	}
	if a[0].Url != "" {
		ra[0].Url = append(ra[0].Url, PrintedDomain(a[0].Url))
	}
	for i := 0; i < len(a)-1; i++ {
		if a[i].Id == a[i+1].Id {
			if a[i+1].Ip4 != 0 {
				ra[i].Ip = append(ra[i].Ip, int2Ip4(a[i+1].Ip4))
			}
			if a[i+1].Domain != "" {
				ra[i].Domain = append(ra[i].Domain, PrintedDomain(a[i+1].Domain))
			}
			if a[i+1].Url != "" {
				ra[i].Url = append(ra[i].Url, a[i+1].Url)
			}
			a = append(a[:i], a[i+1:]...)
			i--
		} else {
			ra = append(ra, TReason{})
			ra[i+1].Id = a[i+1].Id
			if a[i+1].Ip4 != 0 {
				ra[i+1].Ip = append(ra[i+1].Ip, int2Ip4(a[i+1].Ip4))
			}
			if a[i+1].Domain != "" {
				ra[i+1].Domain = append(ra[i+1].Domain, PrintedDomain(a[i+1].Domain))
			}
			if a[i+1].Url != "" {
				ra[i+1].Url = append(ra[i+1].Url, a[i+1].Url)
			}

		}
	}
	sort.Slice(a, func(j, i int) bool {
		switch {
		case a[i].BlockType == TBLOCK_URL && a[j].BlockType != TBLOCK_URL:
			return true
		case a[i].BlockType == TBLOCK_HTTPS &&
			(a[j].BlockType != TBLOCK_URL &&
				a[j].BlockType != TBLOCK_HTTPS):
			return true
		case a[i].BlockType == TBLOCK_DOMAIN &&
			(a[j].BlockType != TBLOCK_URL &&
				a[j].BlockType != TBLOCK_HTTPS &&
				a[j].BlockType != TBLOCK_DOMAIN):
			return true
		default:
			return false
		}
	})
	offset := 0
	if o.Tag == OFFSET_CONTENT {
		switch {
		case len(a) <= PRINT_LIMIT:
			offset = 0
		case o.Count > len(a)-(len(a)%PRINT_LIMIT):
			offset = len(a) - (len(a) % PRINT_LIMIT)
		default:
			offset = o.Count
		}
	}
	var cnt, cbu, cbh, cbd, cbi int
	for i, packet := range a {
		if o.Tag == OFFSET_CONTENT && i < offset {
			continue
		}
		content := TContent{}
		err := json.Unmarshal(packet.Pack, &content)
		if err != nil {
			Error.Printf("Упс!!! %s\n", err)
			continue
		}
		if packet.RegistryUpdateTime < oldest {
			oldest = packet.RegistryUpdateTime
		}
		var req TReason
		for _, req = range ra {
			if req.Id == packet.Id {
				break
			}
		}
		if cnt < PRINT_LIMIT {
			bt := ""
			switch packet.BlockType {
			case TBLOCK_URL:
				bt = "\U000026d4 "
				cbu++
			case TBLOCK_HTTPS:
				bt = "\U0001f4db "
				cbh++
			case TBLOCK_DOMAIN:
				bt = "\U0001f6ab "
				cbd++
			case TBLOCK_IP:
				bt = "\u274c "
				cbi++
			}
			res += fmt.Sprintf("%s /n\\_%d %s...\n", bt, content.Id, string([]rune(content.Description)[:127]))
			if len(req.Ip) != 0 {
				for _, ip := range req.Ip {
					res += fmt.Sprintf("    _как ip_ %s\n", ip)
				}
			}
			if len(req.Domain) != 0 {
				for _, domain := range req.Domain {
					res += fmt.Sprintf("    _как domain_ %s\n", Sanitize(PrintedDomain(domain)))
				}
			}
			if len(req.Url) != 0 {
				for _, u := range req.Url {
					res += fmt.Sprintf("    _как url_ %s\n", Sanitize(PrintedDomain(u)))
				}
			}
			res += "\n"
		}
		cnt++
	}
	if mass != "" {
		res = mass + res
	}
	if len(a) > PRINT_LIMIT {
		pages = append(pages, TPagination{OFFSET_CONTENT, len(a)})
		//rest := cnt - PRINT_LIMIT
		res += fmt.Sprintf("\u2195 результаты с *%d* по *%d* из *%d*\n", offset+1, offset+PRINT_LIMIT, len(a))
		/*if cbu > 0 && cbu < rest {
			res += fmt.Sprintf(" url=%d", cbu)
		} else if cbu > 0 {
			res += fmt.Sprintf(" url=%d", rest)
		}
		if cbh > 0 && cbu+cbh < rest {
			res += fmt.Sprintf(" https=%d", cbu)
		} else if cbh > 0 {
			res += fmt.Sprintf(" https=%d", rest-cbu)
		}
		if cbd > 0 && cbd+cbu+cbh < rest {
			res += fmt.Sprintf(" domain=%d", cbu)
		} else if cbd > 0 {
			res += fmt.Sprintf(" domain=%d", rest-cbh-cbu)
		}
		if cbm > 0 && cbm+cbd+cbu+cbh < rest {
			res += fmt.Sprintf(" wildcard=%d", cbu)
		} else if cbm > 0 {
			res += fmt.Sprintf(" wildcard=%d", rest-cbd-cbh-cbu)
		}
		if cbi > 0 && cbm+cbd+cbu+cbh < rest {
			res += fmt.Sprintf(" ip=%d", rest-cbm-cbd-cbh-cbu)
		}*/
		res += "\n"
	}
	var abt []string
	if cbu > 0 {
		abt = append(abt, "url: \U000026d4")
	}
	if cbh > 0 {
		abt = append(abt, "https: \U0001f4db")
	}
	if cbd > 0 {
		abt = append(abt, "domain: \U0001f6ab")
	}
	if cbi > 0 {
		abt = append(abt, "ip: \u274c")
	}
	res += "*типы блокировки:* " + strings.Join(abt, " | ")
	res += printUpToDate(oldest)
	return
}
