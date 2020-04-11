package main

import (
	"bytes"
	"sync"

	//"bytes"
	//"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	//"sync"

	//"strings"

	//"sync/atomic"
	"time"
)

var (
	httpClient *http.Client
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

const (
	MaxIdleConns        int = 4
	MaxIdleConnsPerHost int = 4
	IdleConnTimeout     int = 60
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 20 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	return client
}

func init() {
	httpClient = createHTTPClient()
}

func If(cond bool, left, right interface{}) interface{} {
	if cond {
		return left
	} else {
		return right
	}
}

func test() {

	//ctx := make(chan int)
	go func() {
		defer func() {
			//ctx <- 1
		}()

		resp, err := httpClient.Get(`http://www.baidu.com`)
		if err != nil && resp == nil {
			fmt.Println(`error ` + err.Error())
		} else {
			defer resp.Body.Close()

			buff := new(bytes.Buffer)
			buff.ReadFrom(resp.Body)
			fmt.Println(buff.String())
		}
	}()

	for {

		t := time.NewTimer(time.Duration(30000) * time.Millisecond)
		select {
		//case <-ctx:
		//	fmt.Println(`success`)
		//	//break
		//	//return
		//	break
		case <-t.C:
			fmt.Println(`time.out`)
			break
		//returnd
		default:
			fmt.Println(`default`)
			t.Stop()
			break
		}
	}

}

//一种发布订阅模式golang

type (
	subscriber chan interface{}
	topicFunc  func(v interface{}) bool
)

type Publisher struct {
	m           sync.RWMutex
	buffer      int
	expire      time.Duration
	subscribers map[subscriber]topicFunc
}

func NewPulisher(expire time.Duration, num int) *Publisher {
	return &Publisher{
		m:           sync.RWMutex{},
		buffer:      num,
		expire:      expire,
		subscribers: make(map[subscriber]topicFunc),
	}
}

func (s *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(subscriber, s.buffer)
	s.m.Lock()
	defer s.m.Unlock()
	s.subscribers[ch] = topic
	return ch
}

func (s *Publisher) Publish(v interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	var wg sync.WaitGroup
	for sub, top := range s.subscribers {
		wg.Add(1)
		//go s.sendTopic(sub,)
		go s.sendTopic(sub, top, v, &wg)
	}
	//fmt.Printf(`%v`,v)
	wg.Wait()
}

func (s *Publisher) sendTopic(sub subscriber, top topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if top != nil && !top(v) {
		return
	}
	//control message number
	if len(sub) < s.buffer {
		fmt.Printf("%d\n", len(sub))
		select {
		case sub <- v:
			fmt.Println(`sub <-v `)
		case <-time.After(s.expire):
		}
	} else {
		fmt.Println(`overload handler power`)
	}
}

func (s *Publisher) Evict(sub chan interface{}) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.subscribers, sub)
	close(sub)
}

func (s *Publisher) Close() {
	s.m.Lock()
	defer s.m.Unlock()
	for key, _ := range s.subscribers {
		delete(s.subscribers, key)
		close(key)
	}
}

type TaskFunc func(v interface{}) interface{}

type Task struct {
	Id          int64
	Desc        string
	MessageType int
	Args        interface{}
	Handle      TaskFunc
}

const (
	PLATE_FORM_YUNJING   = 0
	PLATE_FORM_SMAC      = 1
	PLATE_FORM_TIANYUYUN = 2
	PLATE_FORM_HEXINYUN  = 3
)

const (
	MESSAGE_URL_TYPE    = 0
	MESSAGE_THREAT_TYPE = 1
	MESSAGE_IOC_TYPE    = 2
	MESSAGE_SAMPLE_TYPE = 3
)

type PlateForm struct {
	Name                  string
	Pub                   *Publisher
	PlateFormHanle        TaskFunc
	Switch                bool
	m                     sync.Mutex
	Url_subscriber        subscriber
	Ioc_subscriber        subscriber
	Threat_log_subscriber subscriber
	Sample_md5_subscriber subscriber
	Subscribers           []subscriber
}

func NewPlateForm(name string, pub *Publisher) *PlateForm {
	p := &PlateForm{
		Name:                  name,
		Pub:                   pub,
		PlateFormHanle:        nil,
		Switch:                true,
		m:                     sync.Mutex{},
		Url_subscriber:        nil,
		Ioc_subscriber:        nil,
		Threat_log_subscriber: nil,
		Sample_md5_subscriber: nil,
		Subscribers:           []subscriber{},
	}
	p.PlateFormHanle = func(v interface{}) interface{} {
		fmt.Println(`plateform: ` + p.Name)
		return nil
	}
	return p
}

func (p *PlateForm) TurnOn() {
	p.m.Lock()
	defer p.m.Unlock()
	p.Switch = true
}

func (p *PlateForm) TurnOff() {
	p.m.Lock()
	defer p.m.Unlock()
	p.Switch = false
}

func (p *PlateForm) Register(message int) {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.Switch {
		return
	}
	var subscriber_hanle subscriber
	switch message {
	case MESSAGE_IOC_TYPE:
		subscriber_hanle = p.Pub.SubscribeTopic(func(v interface{}) bool {
			vv := v.(*Task)
			if vv.MessageType == MESSAGE_IOC_TYPE {
				return true
			}
			return false
		})

		break
	case MESSAGE_SAMPLE_TYPE:
		subscriber_hanle = p.Pub.SubscribeTopic(func(v interface{}) bool {
			vv := v.(*Task)
			if vv.MessageType == MESSAGE_SAMPLE_TYPE {
				return true
			}
			return false
		})
		break
	case MESSAGE_THREAT_TYPE:
		subscriber_hanle = p.Pub.SubscribeTopic(func(v interface{}) bool {
			vv := v.(*Task)
			//if 0 == strings.Compare(vv.Desc, `threatlog`) {
			//	return true
			//}
			if vv.MessageType == MESSAGE_THREAT_TYPE {
				return true
			}
			return false
		})
		break
	case MESSAGE_URL_TYPE:
		subscriber_hanle = p.Pub.SubscribeTopic(func(v interface{}) bool {
			vv := v.(*Task)
			if vv.MessageType == MESSAGE_URL_TYPE {
				return true
			}
			return false
		})
		break
	default:

	}
	p.Subscribers = append(p.Subscribers, subscriber_hanle)
}

func (p *PlateForm) LoopHanle() {

	for _, ch := range p.Subscribers {
		go func(ch subscriber) {
			for j := range ch {
				fmt.Println(p.Name)
				if v, ok := j.(*Task); ok {
					v := v.Handle(v.Args)
					p.PlateFormHanle(v)
				}
			}
		}(ch)
	}
}

var (
	g_publisher          *Publisher
	g_smac_plateform     *PlateForm
	g_hexinyun_plateform *PlateForm
)

func init() {
	g_publisher = NewPulisher(time.Duration(2)*time.Second, 100)
	g_smac_plateform = NewPlateForm(`smac`, g_publisher)
	g_hexinyun_plateform = NewPlateForm(`hexinyun`, g_publisher)
	g_hexinyun_plateform.TurnOff()
	g_smac_plateform.Register(MESSAGE_IOC_TYPE)
	g_smac_plateform.Register(MESSAGE_URL_TYPE)
	g_smac_plateform.Register(MESSAGE_THREAT_TYPE)
	g_hexinyun_plateform.Register(MESSAGE_IOC_TYPE)
	g_hexinyun_plateform.Register(MESSAGE_URL_TYPE)

}

//go tool pprof
func main() {

	defer func() {
		g_publisher.Close()
	}()

	limit := CreateRateLimit(1000, 100)

	t := &Task{
		Id:          0,
		Desc:        "ioc",
		Args:        `ioc`,
		MessageType: MESSAGE_IOC_TYPE,
		Handle: func(v interface{}) interface{} {
			if v, ok := v.(string); ok {
				fmt.Println(v)
			} else {
				panic(fmt.Errorf(`type unvalid`).Error())
			}
			return nil
		},
	}

	tt := &Task{
		Id:          0,
		Desc:        "threat",
		Args:        `threat`,
		MessageType: MESSAGE_THREAT_TYPE,
		Handle: func(v interface{}) interface{} {
			if v, ok := v.(string); ok {
				fmt.Println(v)
			} else {
				panic(fmt.Errorf(`type unvalid`).Error())
			}
			return nil
		},
	}
	//g_publisher.Publish(t)

	for i := 0; i < 400; i++ {
		go func() {
			for {
				if limit.Take() {
					fmt.Println(`success......................`)
					g_publisher.Publish(t)
					g_publisher.Publish(tt)
				} else {
					fmt.Println(`limitrate......................`)
				}

			}
		}()
	}

	g_hexinyun_plateform.LoopHanle()
	g_smac_plateform.LoopHanle()
	//fmt.Println(`out select`)
	http.ListenAndServe("0.0.0.0:9999", nil)
}
