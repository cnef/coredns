package rewrite_resp

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

const (
	// ExactMatch matches only on exact match of the name in the question section of a request
	ExactMatch = "exact"
	// PrefixMatch matches when the name begins with the matching string
	PrefixMatch = "prefix"
	// RegexMatch matches when the name in the question section of a request matches a regular expression
	RegexMatch = "regex"
)

type answerRuleBase struct {
	nextAction string
	targetIP   net.IP
	match      func(string) bool
}

func (r *answerRuleBase) RewriteResponse(rr dns.RR) Result {
	switch rr.Header().Rrtype {
	case dns.TypeA:
		if r.match(rr.(*dns.A).A.String()) {
			rr.(*dns.A).A = r.targetIP
			clog.Debugf("matched rule, forward to: %s", r.targetIP.String())
			return RewriteDone
		}
	}
	return RewriteIgnored
}

func (r *answerRuleBase) NeedContinue() bool {
	return r.nextAction == Continue
}

func newAnswerRuleBase(nextAction string, ip net.IP) answerRuleBase {
	return answerRuleBase{
		nextAction: nextAction,
		targetIP:   ip,
	}
}

type exactAnswerRule struct {
	answerRuleBase
	From string
}

type prefixAnswerRule struct {
	answerRuleBase
	Prefix string
}

type regexAnswerRule struct {
	answerRuleBase
	Pattern *regexp.Regexp
}

// Rewrite rewrites the current request based upon exact match of the name
// in the question section of the request.
func (rule *exactAnswerRule) Rewrite(ctx context.Context, state request.Request) ResponseRule {
	rule.match = func(ip string) bool {
		clog.Debugf("exact rule: %s, resp: %s", rule.From, ip)
		return rule.From == ip
	}
	return rule
}

// Rewrite rewrites the current request when the name begins with the matching string.
func (rule *prefixAnswerRule) Rewrite(ctx context.Context, state request.Request) ResponseRule {
	rule.match = func(ip string) bool {
		clog.Debugf("prefix rule: %s, resp: %s", rule.Prefix, ip)
		return strings.HasPrefix(ip, rule.Prefix)
	}
	return rule
}

// Rewrite rewrites the current request when the name in the question
// section of the request matches a regular expression.
func (rule *regexAnswerRule) Rewrite(ctx context.Context, state request.Request) ResponseRule {
	rule.match = func(ip string) bool {
		clog.Debugf("regex rule: %s, resp: %s", rule.Pattern, ip)
		return len(rule.Pattern.FindStringSubmatch(ip)) != 0
	}
	return rule
}

// newAnswerRule creates a name matching rule based on exact, partial, or regex match
func newAnswerRule(nextAction string, args ...string) (Rule, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few (%d) arguments for a answer rule", len(args))
	}
	var s string
	if len(args) == 2 {
		s = args[1]
	}
	if len(args) == 3 {
		s = args[2]
	}
	ip, valid := isValidIP(s)
	if !valid {
		return nil, fmt.Errorf("invalid IP address '%s' for a record rule", s)
	}
	if len(args) == 3 {
		switch strings.ToLower(args[0]) {
		case ExactMatch:
			return &exactAnswerRule{
				newAnswerRuleBase(nextAction, ip),
				args[1],
			}, nil
		case PrefixMatch:
			return &prefixAnswerRule{
				newAnswerRuleBase(nextAction, ip),
				args[1],
			}, nil
		case RegexMatch:
			regexPattern, err := regexp.Compile(args[1])
			if err != nil {
				return nil, fmt.Errorf("invalid regex pattern in the a record rule: %s", args[1])
			}
			return &regexAnswerRule{
				newAnswerRuleBase(nextAction, ip),
				regexPattern,
			}, nil
		default:
			return nil, fmt.Errorf("a record rule supports only exact, prefix, and regex ip matching")
		}
	}
	if len(args) > 3 {
		return nil, fmt.Errorf("many few arguments for a record rule")
	}
	return &exactAnswerRule{
		newAnswerRuleBase(nextAction, ip),
		plugin.Name(args[0]).Normalize(),
	}, nil
}

// isValidIP returns true if v is valid ip value.
func isValidIP(v string) (net.IP, bool) {
	ip := net.ParseIP(v)
	if ip.Equal(net.IPv4zero) {
		return nil, false
	}
	return ip, true
}
