# rewrite

## Name

*rewrite_resp* - performs internal response message rewriting.

## Description

Rewrites are invisible to the client. There are simple rewrites (fast) and complex rewrites
(slower), but they're powerful enough to accommodate most dynamic back-end applications.

## Syntax

A simplified/easy-to-digest syntax for *rewrite* is...
~~~
rewrite_resp [continue|stop] FIELD [TYPE] [(FROM TO)|TTL] [OPTIONS]
~~~

## Examples

###  A Record Rewrites

At times, the need to rewrite a TTL value could arise. For example, a DNS server
may not cache records with a TTL of zero (`0`). An administrator
may want to increase the TTL to ensure it is cached, e.g., by increasing it to 15 seconds.

In the below example, the TTL in the answers for `coredns.rocks` domain are
being set to `15`:

```
    rewrite_resp continue {
        a exact 1.117.245.97 10.100.0.15
    }
```

By the same token, an administrator may use this feature to prevent or limit caching by
setting the TTL value really low.


The syntax for the TTL rewrite rule is as follows. The meaning of
`exact|prefix|suffix|substring|regex` is the same as with the name rewrite rules.
An omitted type is defaulted to `exact`.

```
rewrite [continue|stop] ttl [exact|prefix|suffix|substring|regex] STRING SECONDS
```
