package main // 定义包名，用关键字：package
import (
	"github.com/farseernet/farseer.go/fsApp"
)

func main() { // main函数，程序执行的入口
	fsApp.Initialize[StartupModule]("fss")
}
