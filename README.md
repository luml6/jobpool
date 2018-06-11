## 下载说明

- go get远程git库
```
go get git@github.com:luml6/jobpool.git
```


## 使用说明

- import
```
import (
	jobpool "git@github.com:luml6/jobpool.git"
)
```

- 使用示例
```
type YourJob struct {
	Name string
	Age  int
}

func (c *YourJob) Do() error {
	fmt.Println("name is %v,age is %v", c.Name, c.Age)
	return nil
}

func main() {
	dispath := jobpool.NewDispatcher(conf.Conf.MaxWorker, conf.Conf.MaxWorker)
	dispath.Run()
	youjob := &YourJob{Name: "albert", Age: 12}
	dispath.Add(youjob)
}
```
