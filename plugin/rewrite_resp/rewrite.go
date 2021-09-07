package rewrite_resp

import (
	"context"
	"fmt"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Result is the result of a rewrite
type Result int

const (
	// RewriteIgnored is returned when rewrite is not done on request.
	RewriteIgnored Result = iota
	// RewriteDone is returned when rewrite is done on request.
	RewriteDone
)

// These are defined processing mode.
const (
	// Processing should stop after completing this rule
	Stop = "stop"
	// Processing should continue to next rule
	Continue = "continue"
)

// Rewrite is a plugin to rewrite requests internally before being handled.
type Rewrite struct {
	Next  plugin.Handler
	Rules []Rule
}

// ResponseRule contains a rule to rewrite a response with.
type ResponseRule interface {
	RewriteResponse(rr dns.RR) Result
	NeedContinue() bool
}

// ResponseRules describes an ordered list of response rules to apply
// after a name rewrite
type ResponseRules = []ResponseRule

// ServeDNS implements the plugin.Handler interface.
func (rw Rewrite) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wr := NewResponseWriter(w, r)
	state := request.Request{W: w, Req: r}
	for _, rule := range rw.Rules {
		respRule := rule.Rewrite(ctx, state)
		if _, ok := dns.IsDomainName(state.Req.Question[0].Name); !ok {
			err := fmt.Errorf("invalid name after rewrite: %s", state.Req.Question[0].Name)
			state.Req.Question[0] = wr.originalQuestion
			return dns.RcodeServerFailure, err
		}
		wr.ResponseRules = append(wr.ResponseRules, respRule)
	}

	return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, wr, r)
}

// Name implements the Handler interface.
func (rw Rewrite) Name() string { return "rewrite_resp" }

// Rule describes a rewrite rule.
type Rule interface {
	// Rewrite rewrites the current request.
	Rewrite(ctx context.Context, state request.Request) ResponseRule
}

func newRule(args ...string) (Rule, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no rule type specified for rewrite")
	}

	arg0 := strings.ToLower(args[0])
	var ruleType string
	var startArg int
	mode := Stop
	switch arg0 {
	case Continue:
		mode = Continue
		ruleType = strings.ToLower(args[1])
		startArg = 2
	case Stop:
		ruleType = strings.ToLower(args[1])
		startArg = 2
	default:
		// for backward compatibility
		ruleType = arg0
		startArg = 1
	}

	switch ruleType {
	case "a":
		return newAnswerRule(mode, args[startArg:]...)
	default:
		return nil, fmt.Errorf("invalid rule type %q", args[0])
	}
}
