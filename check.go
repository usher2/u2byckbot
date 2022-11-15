package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	pb "github.com/usher2/u2byckbot/msg"
)

func Ping(c pb.CheckClient) string {
	Info.Printf("Looking for Ping\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.Ping(ctx, &pb.PingRequest{Ping: "ping"})
	if err != nil {
		Debug.Printf("%v.Ping(_) = _, %v\n", c, err)
		return "\U00002620 Что-то пошло не так! Повторите попытку позже\n"
	}
	if r.Error != "" {
		Debug.Printf("ERROR: %s\n", r.Error)
		return fmt.Sprintf("\u23f3 Повторите попытку позже: %s\n", r.Error)
	}
	return fmt.Sprintf("\U0001f919 *%s*%s", r.Pong, printUpToDate(r.RegistryUpdateTime))
}

func searchID(c pb.CheckClient, id uint64) (int64, []*pb.Content, error) {
	Info.Printf("Looking for content: #%d\n", id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SearchID(ctx, &pb.IDRequest{Query: id})
	if err != nil {
		Debug.Printf("%v.SearchContent(_) = _, %v\n", c, err)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\U00002620 Что-то пошло не так! Повторите попытку позже\n")
	}
	if r.Error != "" {
		Debug.Printf("ERROR: %s\n", r.Error)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\u23f3 Повторите попытку позже: %s\n", r.Error)
	}
	return r.RegistryUpdateTime, r.Results[:], nil
}

func searchIP4(c pb.CheckClient, ip string) (int64, []*pb.Content, error) {
	Info.Printf("Looking for %s\n", ip)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SearchIP4(ctx, &pb.IP4Request{Query: parseIp4(ip)})
	if err != nil {
		Debug.Printf("%v.SearchIP4(_) = _, %v\n", c, err)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\U00002620 Что-то пошло не так! Повторите попытку позже\n")
	}
	if r.Error != "" {
		Debug.Printf("ERROR: %s\n", r.Error)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\u23f3 Try again later: %s\n", r.Error)
	}
	return r.RegistryUpdateTime, r.Results[:], nil
}

func searchURL(c pb.CheckClient, u string) (int64, []*pb.Content, error) {
	_url := NormalizeUrl(u)
	if _url != u {
		fmt.Printf("Input was %s\n", u)
	}
	Info.Printf("Looking for %s\n", _url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SearchURL(ctx, &pb.URLRequest{Query: _url})
	if err != nil {
		Debug.Printf("%v.SearchURL(_) = _, %v\n", c, err)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\U00002620 Что-то пошло не так! Повторите попытку позже\n")
	}
	if r.Error != "" {
		Debug.Printf("ERROR: %s\n", r.Error)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\u23f3 Повторите попытку позже: %s\n", r.Error)
	}
	return r.RegistryUpdateTime, r.Results[:], nil
}

func searchDomain(c pb.CheckClient, s string) (int64, []*pb.Content, error) {
	domain := NormalizeDomain(s)
	Info.Printf("Looking for %s\n", domain)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	r, err := c.SearchDomain(ctx, &pb.DomainRequest{Query: domain})
	if err != nil {
		Debug.Printf("%v.SearchDomain(_) = _, %v\n", c, err)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\U00002620 Что-то пошло не так! Повторите попытку позже\n")
	}
	if r.Error != "" {
		Debug.Printf("ERROR: %s\n", r.Error)
		return MAX_TIMESTAMP, nil, fmt.Errorf("\u23f3 Повторите попытку позже: %s\n", r.Error)
	}
	return r.RegistryUpdateTime, r.Results[:], nil
}

func refSearch(c pb.CheckClient, s string) (int64, []*pb.Content, []string, error) {
	var (
		err    error
		oldest int64 = MAX_TIMESTAMP
		utime  int64
		a, a2  []*pb.Content
		ips4   []string
	)
	domain := NormalizeDomain(s)
	ips4 = getIP4(domain)
	for _, ip := range ips4 {
		utime, a2, err = searchIP4(c, ip)
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			a = append(a, a2...)
		} else {
			break
		}
	}
	if err != nil {
		return oldest, nil, ips4, err
	}
	return oldest, a, ips4, nil
}

func numberSearch(c pb.CheckClient, s string, o TPagination) (res string, pages []TPagination) {
	var (
		a      []*pb.Content
		oldest int64 = MAX_TIMESTAMP
		utime  int64
		_res   string
	)
	if len(s) == 0 {
		res = "\U0001f914 Что имелось ввиду?..\n"
		return
	}
	n, err := Base32ToUint64(s)
	switch {
	case err == nil && n != 0:
		utime, a, err = searchID(c, n)
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			if len(a) == 0 {
				res = fmt.Sprintf("\U0001f914 %s *не найден*\n", s)
				res += printUpToDate(oldest)
			}
		}
		if err != nil {
			res = err.Error() + "\n"
		} else {
			_res, pages = constructContentResult(a, o)
			res += _res
		}
	case err != nil:
		res = fmt.Sprintf("\U0001f914 Что имелось ввиду?.. /n\\_%s: %s\n", s, err.Error())
	default:
		res = fmt.Sprintf("\U0001f914 Что имелось ввиду?.. /n\\_%s\n", s)
	}
	return
}

func mainSearch(c pb.CheckClient, s string, o TPagination) (res string, pages []TPagination) {
	var (
		err    error
		a, a2  []*pb.Content
		oldest int64 = MAX_TIMESTAMP
		utime  int64
		_res   string
	)
	if len(s) == 0 {
		res = "\U0001f914 Что имелось ввиду?..\n"
		return
	}
	domain := NormalizeDomain(s)
	if len(s) > 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			s = s[1 : len(s)-2]
			domain = s
		}
	}
	_u, _ur := url.Parse(s)
	if _ur == nil && _u.IsAbs() &&
		(_u.Scheme == "http" || _u.Scheme == "https") &&
		(_u.Port() == "80" || _u.Port() == "443" || _u.Port() == "") &&
		(_u.RequestURI() == "" || _u.RequestURI() == "/") {
		s = _u.Hostname()
		domain = NormalizeDomain(s)
		_ur = fmt.Errorf("fake")
	}
	ip := net.ParseIP(s)
	switch {
	case ip != nil:
		if ip.To4() != nil {
			utime, a, err = searchIP4(c, s)
			if err == nil {
				if utime < oldest {
					oldest = utime
				}
				utime, a2, err = searchDomain(c, s)
				if err == nil {
					if utime < oldest {
						oldest = utime
					}
					if len(a2) > 0 {
						a = append(a, a2...)
					}
				}
			}
		}
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			if len(a) > 0 {
				res = fmt.Sprintf("\U0001f525 %s *заблокирован*\n\n", Sanitize(s))
			} else {
				res = fmt.Sprintf("\u2705 %s *не заблокирован*\n", Sanitize(s))
				res += printUpToDate(oldest)
			}
		}
		if err != nil {
			res = err.Error() + "\n"
		} else {
			_res, pages = constructResult(a, o)
			res += _res
		}
	case isDomainName(domain):
		utime, a, err = searchDomain(c, s)
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			if strings.HasPrefix(s, "www.") {
				utime, a2, err = searchDomain(c, s[4:])
			} else {
				utime, a2, err = searchDomain(c, "www."+s)
			}

		}
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			if len(a2) > 0 {
				a = append(a, a2...)
			}
			if len(a) > 0 {
				res = fmt.Sprintf("\U0001f525 %s *заблокирован*\n\n", Sanitize(s))
			} else {
				res = fmt.Sprintf("\u2705 %s *не заблокирован*\n", Sanitize(s))
				var ips4 []string
				utime, a, ips4, err = refSearch(c, s)
				if err == nil && len(a) > 0 {
					if utime < oldest {
						oldest = utime
					}
					res += "\n\U0001f525 но может быть ограничен по IP-адресу:\n"
					for _, ip := range ips4 {
						res += fmt.Sprintf("    %s\n", ip)
					}
				} else {
					res += printUpToDate(oldest)
				}
			}
		}
		if err != nil {
			res = err.Error() + "\n"
		} else {
			_res, pages = constructResult(a, o)
			res += _res
		}
		fallthrough
	case _ur == nil:
		if len(a) != 0 {
			break
		}
		if _u.Scheme != "https" && _u.Scheme != "http" {
			utime, a, err = searchURL(c, s)
			if err == nil {
				if utime < oldest {
					oldest = utime
				}
				_u.Scheme = "http"
				utime, a2, err = searchURL(c, _u.String())
				if err == nil {
					if len(a2) > 0 {
						a = append(a, a2...)
					}
					if utime < oldest {
						oldest = utime
					}
					_u.Scheme = "https"
					utime, a2, err = searchURL(c, _u.String())
					if len(a2) > 0 {
						a = append(a, a2...)
					}
					if utime < oldest {
						oldest = utime
					}
				}
			}
		} else {
			_u.Scheme = "https"
			utime, a, err = searchURL(c, _u.String())
			if err == nil {
				if utime < oldest {
					oldest = utime
				}
				_u.Scheme = "http"
				utime, a2, err = searchURL(c, _u.String())
				if len(a2) > 0 {
					a = append(a, a2...)
				}
				if utime < oldest {
					oldest = utime
				}
			}
		}
		if err == nil {
			if len(a) > 0 {
				res = fmt.Sprintf("\U0001f525 URL %s *заблокирован*\n\n", Sanitize(s))
			} else {
				res = fmt.Sprintf("\u2705 URL %s *не заблокирован*\n", Sanitize(s))
				res += printUpToDate(oldest)
			}
		}
		if err != nil {
			res = err.Error() + "\n"
		} else {
			_res, pages = constructResult(a, o)
			res += _res
		}
	default:
		utime, a, err = searchURL(c, s)
		if err == nil {
			if utime < oldest {
				oldest = utime
			}
			utime, a2, err = searchDomain(c, s)
			if err == nil {
				if utime < oldest {
					oldest = utime
				}
				if len(a2) > 0 {
					a = append(a, a2...)
				}
			}
		}
		if err != nil {
			res = err.Error() + "\n"
		} else {
			if len(a) > 0 {
				res = fmt.Sprintf("\U0001f525 URL %s *заблокирован*\n\n", Sanitize(s))
				_res, pages = constructResult(a, o)
				res += _res
			} else {
				res = fmt.Sprintf("\U0001f914 Что имелось ввиду?.. %s\n", s)
				res += printUpToDate(oldest)
			}
		}
	}
	return
}
