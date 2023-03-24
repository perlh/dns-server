package dns

import (
	"dns-server/util"
	"dns-server/web"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	_ "github.com/sirupsen/logrus"
	"github.com/wumansgy/goEncrypt/aes"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Memory 内存配置
type Memory struct {
	memoryConfigs []*DnsConfig
}

var memoryFile *os.File

var lock = &sync.Mutex{}

var authUsername string

var authPassword string

var tokenKey string

var tokenExpireSeconds int64

// Login 登录
func (m *Memory) Login(res http.ResponseWriter, req *http.Request) (string, error) {
	method := strings.ToUpper(req.Method)
	if !strings.EqualFold(method, http.MethodPost) {
		return "", errors.New("login: only the HTTP POST method is supported")
	}
	contentType := req.Header.Get("Content-Type")
	if strings.Index(contentType, "multipart/form-data") >= 0 {
		err := req.ParseMultipartForm(128)
		if err != nil {
			return "", err
		}
	} else if strings.Index(contentType, "x-www-form-urlencoded") >= 0 {
		err := req.ParseForm()
		if err != nil {
			return "", err
		}
	}

	form := req.Form
	if !form.Has("username") || !form.Has("password") {
		return "", errors.New("用户名密码不能为空")
	}
	if strings.EqualFold(authUsername, form.Get("username")) && strings.EqualFold(authPassword, form.Get("password")) {
		encodeStr := fmt.Sprintf("%v:%v:%v", authUsername, authPassword, time.Now().Unix())
		//这里需要配置key
		return aes.AesEcbEncryptHex([]byte(encodeStr), []byte(tokenKey)[:16])
	} else {
		return "", errors.New("用户名密码输入有误")
	}
}

// Logout 退出
func (m *Memory) Logout(res http.ResponseWriter, req *http.Request) {
	_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultSuccessMsg[string]("退出登录成功")))
}

// CheckAuth 检测权限,校验通过,回调fun
func (m *Memory) CheckAuth(res http.ResponseWriter, req *http.Request, fun func()) {
	token := strings.TrimSpace(req.Header.Get("Authorization"))
	queryToken := strings.TrimSpace(req.URL.Query().Get("token"))
	if len(token) == 0 && len(queryToken) == 0 {
		_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string]("请携带token进行访问,如在Header中携带Authorization: token或者在url上面携带/xxx?token=xxx")))
		return
	}
	if len(token) == 0 {
		token = queryToken
	}
	var tokenData string
	tokenVal, err := aes.AesEcbDecryptByHex(token, []byte(tokenKey))
	if err != nil {
		_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
		return
	}
	tokenData = string(tokenVal)
	dataArray := strings.Split(tokenData, ":")
	createByStr := dataArray[2]
	createBy, err := util.StringToInt64(createByStr)
	if err != nil {
		_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureMsg[string](err.Error())))
		return
	}
	timeSpace := time.Now().Unix() - createBy
	//如果时间间隔超过token过期时间了，则代表token已经过期了
	if timeSpace > tokenExpireSeconds {
		_, _ = res.Write(web.UnSafeJsonSerializer(web.ResultFailureCodeMsg[string](403, "token已经失效,请重新登录")))
		return
	}
	fun()
}

func (m *Memory) QueryPageDnsParser(page, size int, conditions map[string]string) (*web.Result[[]*DnsConfig], error) {
	var configs []*DnsConfig
	typeCondition := conditions["type"]
	if typeCondition != "-1" {
		if strings.Index(typeCondition, ",") > 0 {
			types := strings.Split(typeCondition, ",")
			configs = util.ArrayFilter(m.memoryConfigs, func(item *DnsConfig) bool {
				return util.ArrayContains(types, strconv.Itoa(int(item.Type)))
			})
		} else {
			typeInt := uint16(util.StringToInt(typeCondition, -1))
			configs = util.ArrayFilter(m.memoryConfigs, func(item *DnsConfig) bool {
				return item.Type == typeInt
			})
		}
	}
	pattern := strings.ReplaceAll(strings.ReplaceAll(conditions["pattern"], ".", "\\."), "*", ".*")
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	//处理type条件
	if typeCondition == "-1" {
		configs = util.ArrayFilter(m.memoryConfigs, func(item *DnsConfig) bool {
			return regex.MatchString(item.Pattern)
		})
	} else {
		configs = util.ArrayFilter(configs, func(item *DnsConfig) bool {
			return regex.MatchString(item.Pattern)
		})
	}

	//处理ttl条件
	ttlOperator := conditions["ttlOperator"]
	var ttl int
	var ttlb int
	if strings.EqualFold(ttlOperator, "between") {
		ttlArray := strings.Split(conditions["ttl"], ",")
		if len(ttlArray) != 2 {
			return nil, errors.New("TtlOperator: between, [ttl] supports only two parameters")
		}
		ttl = util.StringToInt(ttlArray[0], -1)
		ttlb = util.StringToInt(ttlArray[1], -1)
		if ttl == -1 {
			return nil, errors.New("TtlOperator: between, ttl.0 cannot be equals -1")
		}
		if ttlb == -1 {
			return nil, errors.New("TtlOperator: between, ttl.1 cannot be equals -1")
		}
	} else {
		ttl = util.StringToInt(conditions["ttl"], -1)
	}

	if ttl != -1 {
		configs = util.ArrayFilter(configs, func(item *DnsConfig) bool {
			switch ttlOperator {
			case "eq":
				return item.Ttl == uint32(ttl)
			case "lte": //小于等于
				return item.Ttl <= uint32(ttl)
			case "lt": //小于
				return item.Ttl < uint32(ttl)
			case "gt": //大于
				return item.Ttl > uint32(ttl)
			case "gte": //大于等于
				return item.Ttl >= uint32(ttl)
			case "between": //介于
				return item.Ttl >= uint32(ttl) && item.Ttl <= uint32(ttlb)
			}
			return item.Ttl == uint32(ttl)
		})
	}

	//处理ttl条件
	priorityOperator := conditions["priorityOperator"]
	var priority int
	var priorityb int
	if strings.EqualFold(priorityOperator, "between") {
		priorityArray := strings.Split(conditions["priority"], ",")
		if len(priorityArray) != 2 {
			return nil, errors.New("priorityOperator: between, [priority] supports only two parameters")
		}
		priority = util.StringToInt(priorityArray[0], -1)
		priorityb = util.StringToInt(priorityArray[1], -1)
		if priority == -1 {
			return nil, errors.New("priorityOperator: between, priority.0 cannot be equals -1")
		}
		if priorityb == -1 {
			return nil, errors.New("priorityOperator: between, priority.1 cannot be equals -1")
		}
	} else {
		priority = util.StringToInt(conditions["priority"], -1)
	}

	if priority != -1 {
		configs = util.ArrayFilter(configs, func(item *DnsConfig) bool {
			switch priorityOperator {
			case "eq":
				return item.Priority == priority
			case "lte": //小于等于
				return item.Priority <= priority
			case "lt": //小于
				return item.Priority < priority
			case "gt": //大于
				return item.Priority > priority
			case "gte": //大于等于
				return item.Priority >= priority
			case "between": //介于
				return item.Priority >= priority && item.Priority <= priorityb
			}
			return item.Priority == priority
		})
	}
	ttlSort := conditions["ttlSort"]
	if ttlSort == "desc" {
		util.ArrayInsertSort(configs, func(a, b *DnsConfig) bool {
			return a.Ttl > b.Ttl
		})
	} else {
		util.ArrayInsertSort(configs, func(a, b *DnsConfig) bool {
			return b.Ttl > a.Ttl
		})
	}

	prioritySort := conditions["prioritySort"]
	if prioritySort == "desc" {
		util.ArrayInsertSort(configs, func(a, b *DnsConfig) bool {
			return a.Priority > b.Priority
		})
	} else {
		util.ArrayInsertSort(configs, func(a, b *DnsConfig) bool {
			return b.Priority > a.Priority
		})
	}

	startIndex := (page - 1) * size
	dataLen := len(configs)
	if startIndex > dataLen {
		result := web.ResultSuccessData[[]*DnsConfig]([]*DnsConfig{})
		result.Attr = make(map[string]interface{}, 1)
		result.Attr["total"] = 0
		return result, nil
	}
	endIndex := 0
	if dataLen <= startIndex+size {
		endIndex = dataLen
	} else {
		endIndex = startIndex + size
	}
	result := web.ResultSuccessData[[]*DnsConfig](configs[startIndex:endIndex])
	result.Attr = make(map[string]interface{}, 1)
	result.Attr["total"] = dataLen
	return result, nil
}

func (m *Memory) UpdateDnsParser(conf *DnsConfig) error {
	lock.Lock()
	defer lock.Unlock()
	if conf.Id <= 0 {
		return errors.New("id cannot be null")
	}
	for i := range m.memoryConfigs {
		config := m.memoryConfigs[i]
		if config.Id == conf.Id {
			m.memoryConfigs[i] = conf
			m.store()
			m.sort()
			return nil
		}
	}
	return errors.New("the DNS configuration not exists. Please check the DNS configuration you entered")
}

func (m *Memory) RemoveDnsParser(id uint64) error {
	lock.Lock()
	defer lock.Unlock()
	for i := range m.memoryConfigs {
		config := m.memoryConfigs[i]
		if config.Id == id {
			if i == 0 {
				m.memoryConfigs = m.memoryConfigs[1:]
			} else {
				m.memoryConfigs = append(m.memoryConfigs[0:i], m.memoryConfigs[i+1:]...)
			}
			m.store()
			return nil
		}
	}
	return errors.New(fmt.Sprintf("not found dns parser config,dns id: %v", id))
}

func (m *Memory) AddDnsParser(conf *DnsConfig) error {
	lock.Lock()
	defer lock.Unlock()
	for i := range m.memoryConfigs {
		config := m.memoryConfigs[i]
		if config.Type == conf.Type && strings.EqualFold(config.Pattern, conf.Pattern) {
			return errors.New("the DNS configuration already exists. Please check the DNS configuration you entered")
		}
	}
	conf.Id = m.nextId()
	m.memoryConfigs = append(m.memoryConfigs, conf)
	m.store()
	m.sort()
	return nil
}

func (m *Memory) QueryQuestion(q dns.Question) (error, []dns.RR) {
	for i := range m.memoryConfigs {
		config := m.memoryConfigs[i]
		if match(config, q) {
			rr := make([]dns.RR, len(config.Values))
			for j := range config.Values {
				value := config.Values[j]
				var ip net.IP
				var name string
				if q.Qtype == dns.TypePTR {
					ip = net.ParseIP(util.ResolveIp(q.Name))
					name = value.Val
				} else if q.Qtype == dns.TypeA || q.Qtype == dns.TypeAAAA {
					ip = net.ParseIP(value.Val)
					name = q.Name
				}
				if value.Type == dns.TypeA {
					rr[j] = &dns.A{
						Hdr: dns.RR_Header{
							Name:   name,
							Rrtype: value.Type,
							Class:  dns.ClassINET,
							Ttl:    config.Ttl,
						},
						A: ip,
					}
				} else if value.Type == dns.TypeAAAA {
					rr[j] = &dns.AAAA{
						Hdr: dns.RR_Header{
							Name:   name,
							Rrtype: value.Type,
							Class:  dns.ClassINET,
							Ttl:    config.Ttl,
						},
						AAAA: ip,
					}
				}

			}
			return nil, rr
		}
	}
	//如果查询类型是ipv6
	if q.Qtype == dns.TypeAAAA {
		//然后查找有没有对应的ipv4配置，如果存在，则返回空rr
		for i := range m.memoryConfigs {
			config := m.memoryConfigs[i]
			if config.Type == dns.TypeA && patternMatch(config.Pattern, q.Name) {
				return nil, make([]dns.RR, 0)
			}
		}
	}
	return errors.New("not found memory dns "), nil
}

// match 简单匹配
func match(config *DnsConfig, q dns.Question) bool {
	if config.Type == q.Qtype && config.Type == dns.TypePTR {
		ip := util.ResolveIp(q.Name)
		if strings.EqualFold(config.Pattern, ip) || patternMatch(config.Pattern, ip) {
			return true
		}
	} else if config.Type == q.Qtype {
		return patternMatch(config.Pattern, q.Name)
	}
	return false
}

// patternMatch 正则匹配
func patternMatch(name, val string) bool {
	compile, err := regexp.Compile(name)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return compile.MatchString(val)
}

// 初始化
func (m *Memory) init() {
	var err error
	tokenExpireSeconds, err = util.StringToInt64(util.GetEnv("MEMORY_HTTP_TOKEN_EXPIRE_SECOND", "604800", "内存模式http登录成功之后token有效期,默认七天(604800s)"))
	if err != nil {
		fmt.Println("LOAD ENVIRONMENT(MEMORY_HTTP_TOKEN_EXPIRE_SECOND) CONVERT FAILURE: error ", err)
	}
	tokenKey = util.GetEnv("MEMORY_HTTP_TOKEN_KEY", "ax@!#$148ax){}ax", "内存模式http登录成功token加密key,默认: ax@!#$148ax){}ax")
	authUsername = util.GetEnv("MEMORY_HTTP_AUTH_USERNAME", "admin", "内存模式http认证用户名,默认: admin")
	authPassword = util.GetEnv("MEMORY_HTTP_AUTH_PASSWORD", "123456", "内存模式http认证用户名,默认: 123456")
	persistenceFilename := util.GetEnv("MEMORY_DNS_PARSER_DATA_PERSISTENCE_DIR", "./data", "内存模式数据持久化目录")
	file, err := os.OpenFile(persistenceFilename+"/memory.log", os.O_RDWR, os.ModePerm)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(persistenceFilename, os.ModeDir)
		file, err = os.OpenFile(persistenceFilename+"/memory.log", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
			return
		}
	}
	memoryFile = file
	m.load()
}

// 根据优先级排序
func (m *Memory) sort() {
	util.ArrayBubbleSort(m.memoryConfigs, func(a *DnsConfig, b *DnsConfig) bool { return b.Priority > a.Priority })
}

// 加载本地配置
func (m *Memory) load() {
	bodyLen := util.ReadUInt64(memoryFile, 4)
	if bodyLen == 0 {
		memoryConfigs := make([]*DnsConfig, 0)
		m.memoryConfigs = memoryConfigs
	} else {
		_, _ = memoryFile.Seek(8, 0)
		bytes := make([]byte, bodyLen)
		_, _ = memoryFile.Read(bytes)
		err := json.Unmarshal(bytes, &m.memoryConfigs)
		if err != nil {
			panic(err)
		}
		m.sort()
	}
}

// 持久化
func (m *Memory) store() {
	configs := m.memoryConfigs
	bytes, _ := json.Marshal(configs)
	util.WriteUInt64(memoryFile, 4, uint64(len(bytes)))
	_, _ = memoryFile.Seek(8, 0)
	_, _ = memoryFile.Write(bytes)
}

// 下一个id
func (m *Memory) nextId() uint64 {
	currentId := util.ReadUInt64(memoryFile, 0)
	nextId := currentId + 1
	util.WriteUInt64(memoryFile, 0, nextId)
	return nextId
}

func memoryHandler() *Memory {
	memory := &Memory{}
	memory.init()
	return memory
}
