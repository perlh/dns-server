package dns

import (
	"dns-server/util"
	"dns-server/web"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	MemoryHandler = "MEMORY"
)

// Handler handler
type Handler interface {
	//QueryQuestion 查询
	QueryQuestion(q dns.Question) (error, []dns.RR)
	//AddDnsParser 添加域名解析
	AddDnsParser(conf *DnsConfig) error
	//UpdateDnsParser 更新域名解析
	UpdateDnsParser(conf *DnsConfig) error
	//RemoveDnsParser 删除域名解析
	RemoveDnsParser(id uint64) error
	//QueryPageDnsParser 分页多条件查询dns解析
	QueryPageDnsParser(page, size int, conditions map[string]string) (*web.Result[[]*DnsConfig], error)

	//Login 登录
	Login(res http.ResponseWriter, req *http.Request) (string, error)

	//Logout 退出
	Logout(res http.ResponseWriter, req *http.Request)

	//CheckAuth 检测权限,校验通过,回调fun
	CheckAuth(res http.ResponseWriter, req *http.Request, fun func())
}

// DnsConfig 内存配置
type DnsConfig struct {
	Id       uint64  `json:"id" desc:"id"`
	Type     uint16  `json:"type" desc:"dns.TypeXXX"`
	Pattern  string  `json:"pattern" desc:"正则匹配规则，如type=dns.TypePtr时，这里匹配规则使用ip"`
	Values   []Value `json:"values" desc:"type=ptr,values=host names;type=ipv4 and ipv6,values=ips"`
	Ttl      uint32  `json:"ttl" desc:"unit seconds"`
	Priority int     `json:"priority" desc:"priority"`
}

type Value struct {
	Type uint16 `json:"type" desc:"返回的值类型,如ptr,会返回ipv4,ipv6类型的ip地址"`
	Val  string `json:"val" desc:"type=ptr时，此参数为域名"`
}

type DispatcherHandler struct {
	Handler Handler
}

var handler *DispatcherHandler

var remoteDnsServers []string

func init() {
	initEnv()
	handlerType := util.GetEnv("DNS_HANDLER_TYPE", MemoryHandler, "设置此环境变量,指定Dns服务器启动的类型,目前支持: MEMORY(内存方式)")
	handler = &DispatcherHandler{}
	switch handlerType {
	case MemoryHandler:
		(*handler).Handler = memoryHandler()
		break
	}
}

// startHttpServe 初始化http服务
func startHttpServe() {

	//检查是否已经登录
	http.HandleFunc("/auth/checkLogin", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json;charset=UTF-8")
		handler.Handler.CheckAuth(res, req, func() {
			_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultSuccessMsg[string]("success")))
		})
	})

	//登录
	http.HandleFunc("/auth/login", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json;charset=UTF-8")
		token, err := handler.Handler.Login(res, req)
		if err != nil {
			_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
			return
		}
		_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultSuccessData[string](token)))
	})

	//退出登录
	http.HandleFunc("/auth/logout", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json;charset=UTF-8")
		handler.Handler.Logout(res, req)
	})
	http.HandleFunc("/dns/search", func(res http.ResponseWriter, req *http.Request) {
		handler.Handler.CheckAuth(res, req, func() {
			res.Header().Set("Content-Type", "application/json;charset=UTF-8")
			params := req.URL.Query()
			typeCondition := util.Get(params, "type", "-1")
			pattern := util.Get(params, "pattern", "*")
			//asc desc
			ttlSort := util.Get(params, "ttlSort", "desc")
			ttl := util.Get(params, "ttl", "-1")
			ttlOperator := util.Get(params, "ttlOperator", "eq")
			//asc desc
			prioritySort := util.Get(params, "prioritySort", "desc")
			priority := util.Get(params, "priority", "-1")
			priorityOperator := util.Get(params, "priorityOperator", "eq")
			result, err := handler.Handler.QueryPageDnsParser(
				util.StringToInt(params.Get("page"), 1),
				util.StringToInt(params.Get("size"), 10),
				map[string]string{
					"type":             typeCondition,
					"pattern":          pattern,
					"ttlSort":          strings.ToLower(ttlSort),
					"prioritySort":     strings.ToLower(prioritySort),
					"ttl":              ttl,
					"ttlOperator":      strings.ToLower(ttlOperator),
					"priority":         priority,
					"priorityOperator": strings.ToLower(priorityOperator),
				})
			if err != nil {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
				return
			}
			_, _ = res.Write(web.UnSafeJsonSerializer(result))
		})
	})

	http.HandleFunc("/dns/del", func(res http.ResponseWriter, req *http.Request) {
		handler.Handler.CheckAuth(res, req, func() {
			res.Header().Set("Content-Type", "application/json;charset=UTF-8")
			requestMethod := strings.ToUpper(req.Method)
			query := req.URL.Query()
			if !strings.EqualFold(requestMethod, http.MethodDelete) {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string]("path: dns/del,supports requests method: DELETE ")))
				return
			}
			if !query.Has("id") {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string]("path: dns/del,url param: id cannot be null")))
				return
			}
			err := handler.Handler.RemoveDnsParser(uint64(util.StringToInt(util.Get(query, "id", "-1"), -1)))
			if !strings.EqualFold(requestMethod, http.MethodDelete) {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
				return
			}
			_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string]("删除成功")))
		})
	})

	//添加dns解析
	http.HandleFunc("/dns/parser", func(res http.ResponseWriter, req *http.Request) {
		handler.Handler.CheckAuth(res, req, func() {
			res.Header().Set("Content-Type", "application/json;charset=UTF-8")
			requestMethod := strings.ToUpper(req.Method)
			if !strings.EqualFold(requestMethod, http.MethodPost) && !strings.EqualFold(requestMethod, http.MethodPut) {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string]("path: dns/parser,supports requests method: POST,PUT ")))
				return
			}
			bytes, err := io.ReadAll(req.Body)
			if err != nil {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
				return
			}
			config := &DnsConfig{}
			err = json.Unmarshal(bytes, config)
			if err != nil {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
				return
			}
			if msg, ok := config.Check(); !ok {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](msg)))
				return
			}
			var msg string
			switch requestMethod {
			case http.MethodPost:
				err = handler.Handler.AddDnsParser(config)
				msg = "添加dns解析成功"
				break
			case http.MethodPut:
				err = handler.Handler.UpdateDnsParser(config)
				msg = "更新dns解析成功"
				break
			}
			if err != nil {
				_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
				return
			}
			_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultSuccessMsg[string](msg)))
		})
	})
	fileHandler := http.FileServer(http.Dir(util.GetEnv("HTTP_STATIC_DIR", "./static", "设置此环境变量,指定web控制台静态页面的地址,通过此参数,可以指向自定义的静态页面")))
	httpHandler := web.HttpHandler(func(res http.ResponseWriter, req *http.Request) {
		//res.Header().Set("Cache-Control", "max-age=86400")
		fileHandler.ServeHTTP(res, req)
	})
	http.Handle("/", httpHandler)
	go func() {
		err := http.ListenAndServe(util.GetEnv("HTTP_LISTENER_PORT", ":5000", "设置此环境变量,用于指定web服务的监听地址,默认使用: (:5000)"), nil)
		if err != nil {
			panic(err)
		}
	}()
}

func initEnv() {
	dnsServers := util.GetEnv("DNS_REMOTE_SERVERS", "114.114.114.114:53,8.8.8.8:53", "通过设置此环境变量,可以指定当未匹配到dns解析时,使用此参数决定使用那些外部dns服务器进行解析")
	remoteDnsServers = strings.Split(dnsServers, ",")
}
func Service() error {
	startHttpServe()
	dnsListenerPort := util.GetEnv("DNS_LISTENER_PORT", ":53", "设置此环境变量,用于指定dns服务器监听的udp地址,默认使用: (:53)")
	PrintHelp()
	return dns.ListenAndServe(dnsListenerPort, "udp", handler)
}

func PrintHelp() {
	fmt.Println(`
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
APIS:
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
以下除了登录以外，所有请求可以通过url?token=xxx或者header: Authorization: Token方式携带token

检查是否已经登录
任意方式请求: /auth/check[?token=xxx]
Header:
[Authorization: Token]
Resp: 
{
    "data": "",
    "code": 200,
    "state": "ok",
    "msg": "success",
    "attr": null
}

登录(返回token,每次请求在Header中添加: Authorization: Token):
Token有效期 默认7天
POST: /auth/login
Header: 
Content-Type: application/w-xxx-form-urlencoded
Body:
username=admin&password=123456
Resp:
{
    "data": "mNMncSFBoJh3GHRjp/SEVg==",
    "code": 200,
    "state": "ok",
    "msg": "success",
    "attr": null
}

退出登录(内存模式不做实际退出,前端做本地退出即可):
GET: /auth/logout

查询dns解析配置:
GET: /dns/search?page=xxx&size=xxx&type=xxx&ttlSort=asc|desc&ttlOperator=eq|lt|lte|gt|gte|between&ttl=xxx[,xxx]&prioritySort=asc|desc&priority=xxx[,xxx]&priorityOperator=eq|lt|lte|gt|gte|between

删除dns解析配置,根据id
DELETE: /dns/del?id=xxx

添加dns解析配置
POST: /dns/parser
{
    "type": 1,
    "pattern": "www.lhstack.com",
    "values": [
        {
            "type": 1,
            "val": "192.168.101.1"
        }
    ],
    "ttl": 450,
    "priority": 130
}

根据id更新dns解析配置
PUT: /dns/parser
{
	"id": 1,
    "type": 12,
    "pattern": "192.168.101.1",
    "values": [
        {
            "type": 1,
            "val": "www.lhstack.com"
        }
    ],
    "ttl": 450,
    "priority": 130
}
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
APIS:
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
COMMAND:
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

nslookup -q=ptr 192.168.101.5 127.0.0.1(dns) 通过ip解析ip的域名
对应配置
{
    "id": 3,
    "type": 12,
    "pattern": "192.168\\.\\d{1,3}\\.\\d{1,3}",
    "values": [
        {
            "type": 1,
            "val": "www.lhstack.com."
        }
    ],
    "ttl": 450,
    "priority": 130
}

nslookup www.lhstack.com 127.0.0.1(dns) 通过域名获取ip
对应配置,可以使用正则匹配,如www\.lh.*\.com
{
    "type": 1,
    "pattern": "www.lhstack.com.",
    "values": [
        {
            "type": 1,
            "val": "192.168.101.1"
        }
    ],
    "ttl": 450,
    "priority": 130
}
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
COMMAND:
--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------`)
}

func (c *DnsConfig) Check() (string, bool) {
	if c.Type == dns.TypeNone {
		return "dns类型不能为空", false
	}
	if len(strings.TrimSpace(c.Pattern)) == 0 {
		return "dns解析表达式不能为空", false
	}
	if len(c.Values) == 0 {
		return "dns解析结果不能为空", false
	}
	return "", true
}

func (h *DispatcherHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	rr := make([]dns.RR, 0)
	flag := false
	for i := range r.Question {
		question := r.Question[i]
		if question.Qtype == dns.TypePTR {
			ip := util.ResolveIp(question.Name)
			if strings.EqualFold(ip, "127.0.0.1") || strings.EqualFold(ip, "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.1") {
				flag = true
				continue
			}
		}
		err, r := h.Handler.QueryQuestion(question)
		if err != nil {
			continue
		} else {
			flag = true
			rr = append(rr, r...)
		}
	}
	if flag {
		r.MsgHdr.Rcode = dns.RcodeSuccess
		r.MsgHdr.Response = true
		r.MsgHdr.Opcode = dns.OpcodeQuery
		r.MsgHdr.RecursionAvailable = true
		r.Answer = rr
		util.PrintError(w.WriteMsg(r))
	} else {
		exchange(w, r)
	}
}

// 请求转发,将dns请求交由第三方dns处理
func exchange(w dns.ResponseWriter, r *dns.Msg) {
	for i := 0; i < len(remoteDnsServers); i++ {
		for j := range r.Question {
			question := r.Question[j]
			fmt.Printf("using third party dns server query {%v};question name {%v};question query type {%v}\r\n", remoteDnsServers[i], question.Name, question.Qtype)
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		msg, err := dns.ExchangeContext(ctx, r, remoteDnsServers[i])
		cancelFunc()
		if err != nil {
			util.PrintError(err)
			continue
		}
		if msg.Rcode == dns.RcodeNameError {
			continue
		}
		err = w.WriteMsg(msg)
		util.PrintError(err)
		return
	}
}
