package main // 定义包名，用关键字：package
import (
	"github.com/farseer-go/fs"
)

func main() { // main函数，程序执行的入口
	fs.Initialize[StartupModule]("fss")
}
