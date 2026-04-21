package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"syncDevelop/web"
)

func main() {
	addr := flag.String("addr", ":9837", "监听地址")
	dir := flag.String("dir", "./syncfiles", "文件存储目录")
	flag.Parse()

	os.MkdirAll(*dir, 0755)

	localIP := web.GetLocalIP()

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║       SyncDevelop 剪贴板同步         ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Printf("  本机: http://%s%s\n", localIP, *addr)
	fmt.Printf("  局域网其他设备访问上述地址即可同步\n")
	fmt.Println()

	srv := web.NewServer(*addr, *dir)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("[Main] 启动失败: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\n[Main] 正在退出...")
}
