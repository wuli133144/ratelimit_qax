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

