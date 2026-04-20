//go:build darwin

package clipboard

import (
	"os"
	"os/exec"
)

func readText() (string, error) {
	out, err := exec.Command("pbpaste", "-Prefer", "txt").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeText(text string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(text))
	}()
	return cmd.Run()
}

func readImage() ([]byte, error) {
	// 使用 python3 + AppScript 方式读取剪贴板图片
	out, err := exec.Command("python3", "-c", `
import subprocess, sys, tempfile, os
# 尝试用 pngpaste 或直接用 osascript 获取图片
result = subprocess.run(["osascript", "-e",
    'set theType to clipboard info\nset pngData to missing value\nrepeat with i from 1 to number of items in theType\nif (item i of theType) is {{«class PNGf», 0}} then\nset pngData to the clipboard as «class PNGf»\nexit repeat\nend if\nend repeat\nif pngData is not missing value then\nreturn pngData\nend if'], capture_output=True)
if result.returncode == 0 and len(result.stdout) > 0:
    sys.stdout.buffer.write(result.stdout)
`).Output()
	if err != nil || len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func writeImage(data []byte) error {
	tmpFile := "/tmp/syncclip_img.png"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	return exec.Command("osascript", "-e",
		`set the clipboard to (read (POSIX file "`+tmpFile+`") as «class PNGf»)`).Run()
}
