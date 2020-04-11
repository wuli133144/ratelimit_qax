# ratelimit_qax
golang 最简单的限速工具

# 使用方法


func main() {

	t := CreateRateLimit(2000, 100)
	for {
		pre := time.Now()
		r := t.TakeT()
		fmt.Println(r.Sub(pre))
	}
}
# plateform_frame
平台框架是处理不同平台消息的基础框架
目前平台包括：
* 核心云
* 天御云
* smac
* 云镜

### frame使用说明
```cassandraql
var (
	g_publisher          *Publisher
	g_smac_plateform     *PlateForm
	g_hexinyun_plateform *PlateForm
)

func init() {
	g_publisher = NewPulisher(time.Duration(2)*time.Second, 100)
	g_smac_plateform = NewPlateForm(`smac`, g_publisher)
	g_hexinyun_plateform = NewPlateForm(`hexinyun`, g_publisher)
}

//go tool pprof
func main() {

	defer func() {
		g_publisher.Close()
	}()

	t := &Task{
		Id:   0,
		Desc: "ioc",
		Args: `ioc`,
		Handle: func(v interface{}) interface{} {
			if v, ok := v.(string); ok {
				fmt.Println(v)
			} else {
				panic(fmt.Errorf(`type unvalid`).Error())
			}
			return nil
		},
	}
	g_publisher.Publish(t)

	for i := 0; i < 400; i++ {
		go func() {
			for {
				g_publisher.Publish(t)
			}
		}()
	}

	g_hexinyun_plateform.LoopHanle()
	g_smac_plateform.LoopHanle()
	//fmt.Println(`out select`)
	http.ListenAndServe("0.0.0.0:9999", nil)
}

```
平台可定制不同的消息，也可以提供统一的消息
# 消息发生器
消息发生器是适配NGFW和skydas_hub的统一适配器，确保经过消息适配器的消息都是
plateform_framework可以识别的
