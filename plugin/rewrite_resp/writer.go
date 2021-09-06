package rewrite_resp

import "github.com/miekg/dns"

type ResponseWriter struct {
	dns.ResponseWriter
	originalQuestion dns.Question
	ResponseRules    ResponseRules
}

// NewResponseWriter returns a pointer to a new ResponseWriter.
func NewResponseWriter(w dns.ResponseWriter, r *dns.Msg) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter:   w,
		originalQuestion: r.Question[0],
	}
}

// WriteMsg records the status code and calls the underlying ResponseWriter's WriteMsg method.
func (r *ResponseWriter) WriteMsg(res1 *dns.Msg) error {
	// Deep copy 'res' as to not (e.g). rewrite a message that's also stored in the cache.
	res := res1.Copy()

	if len(r.ResponseRules) > 0 {
		for _, rr := range res.Answer {
			r.rewriteResourceRecord(res, rr)
		}
	}
	return r.ResponseWriter.WriteMsg(res)
}

// Write is a wrapper that records the size of the message that gets written.
func (r *ResponseWriter) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}

func (r *ResponseWriter) rewriteResourceRecord(res *dns.Msg, rr dns.RR) {
	for _, rule := range r.ResponseRules {
		rule.RewriteResponse(rr)
	}
}
