# rewrite

## Name

*rewrite_resp* - 匹配返回 A 记录，重写 A 记录

## Description

根据规则匹配返回的 A 记录并重写，使用在灰度场景，有大量域名解析到某个地址，但是不能通过域名进行规则匹配时，
可以通过该插件对解析结果进行判断，域名解析结果符合规则的域名重写为灰度网关的IP

## Syntax

A simplified/easy-to-digest syntax for *rewrite_resp* is...
~~~
rewrite_resp [continue|stop] a [exact|prefix|regex] [RULE] [TARGET]
~~~

## Examples

###  A Record Rewrites

例如：

1. 解析结果为 1.117.245.97 的改写为 10.100.0.15
2. 解析结果为 121.40.49.30 的改写为 192.168.2.123
3. 解析结果前缀为 12.40.149. 的改写为 192.168.2.124

```
    rewrite_resp stop a exact 1.117.245.97 10.100.0.15
    rewrite_resp stop a exact 121.40.49.30 192.168.2.123
    rewrite_resp stop a prefix 12.40.149. 192.168.2.124
```

`exact|prefix|regex` 分别表示：精确匹配，前缀匹配和正则匹配