/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/yizhixiaokong/BullshitGeneratorGo/utils"

	"github.com/spf13/cobra"
)

// letterBytes 包含英文字母数字换行符和空格的字符列表
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 \n"

// threshold 文件写入阈值
const threshold int64 = 1024 * 1024 * 10 // 10MB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "BullshitGeneratorGo",
	Short: "A bullshit generator written in Go",
	Long: `A bullshit generator written in Go.
It can generate a file with random string.`,

	PreRunE: PreRunE,
	RunE:    RunE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// flag 'n',用于指定生成的名字
	rootCmd.Flags().StringP("name", "n", "bullshit", "name of the file")
	// flag 'a',当指定时，如果文件存在，以追加的方式写入文件
	rootCmd.Flags().BoolP("append", "a", false, "append to the file")
	// flag 's',用于指定生成的字符串大小，单位为字节，优先级最低
	rootCmd.Flags().Int64P("size", "s", 1000, "size of the string, unit is byte, priority is lowest,range is [0,2147483648]")
	// flag 'k',用于指定生成的字符串大小，单位为kb，优先级高于size
	rootCmd.Flags().IntP("kbs", "k", 0, "size of the string, unit is kb, priority is higher than length,range is [0,2097152]")
	// flag 'm',用于指定生成的字符串大小，单位为mb，优先级高于kb和size
	rootCmd.Flags().IntP("mbs", "m", 0, "size of the string, unit is mb, priority is higher than kb and length,range is [0,2048]")
	// flag 'g',用于指定生成的字符串大小，单位为gb，优先级高于mb、kb和size
	rootCmd.Flags().IntP("gbs", "g", 0, "size of the string, unit is gb, priority is higher than mb、kb and length,range is [0,2]")
}

// PreRunE 用于检查flag的值是否合法
func PreRunE(cmd *cobra.Command, args []string) error {

	// 随机数种子
	rand.Seed(time.Now().UnixNano())

	// 检查flag 'gbs'是否在0-2之间
	gbs, err := cmd.Flags().GetInt("gbs")
	if err != nil {
		fmt.Println("get flag 'gbs' failed, err:", err)
		return err
	}
	if gbs < 0 || gbs > 2 {
		err := errors.New("flag 'gbs' must be in 0-2(2GB)")
		return err
	}
	// 使用flag 'g'的值修改flag 'm'的值
	if gbs > 0 {
		// m=1024*g
		err := cmd.Flags().Set("mbs", fmt.Sprintf("%d", 1024*gbs))
		if err != nil {
			fmt.Println("set flag 'mbs' failed, err:", err)
			return err
		}
	}

	// 检查flag 'mbs'是否在0-2*1024之间
	mbs, err := cmd.Flags().GetInt("mbs")
	if err != nil {
		fmt.Println("get flag 'mbs' failed, err:", err)
		return err
	}
	if mbs < 0 || mbs > 2*1024 {
		err := errors.New("flag 'mbs' must be in 0-2048(2GB)")
		return err
	}
	// 使用flag 'm'的值修改flag 'k'的值
	if mbs > 0 {
		// k=1024*m
		err := cmd.Flags().Set("kbs", fmt.Sprintf("%d", 1024*mbs))
		if err != nil {
			fmt.Println("set flag 'kbs' failed, err:", err)
			return err
		}
	}

	// 检查flag 'kbs'是否在0-2*1024*1024之间
	kbs, err := cmd.Flags().GetInt("kbs")
	if err != nil {
		fmt.Println("get flag 'kbs' failed, err:", err)
		return err
	}
	if kbs < 0 || kbs > 2*1024*1024 {
		err := errors.New("flag 'kbs' must be in 0-2097152(2GB)")
		return err
	}
	// 使用flag 'k'的值修改flag 's'的值
	if kbs > 0 {
		// s=1024*k
		err := cmd.Flags().Set("size", fmt.Sprintf("%d", 1024*kbs))
		if err != nil {
			fmt.Println("set flag 'size' failed, err:", err)
			return err
		}
	}

	// 检查flag 'size'是否在0-2*1024*1024*1024之间
	size, err := cmd.Flags().GetInt64("size")
	if err != nil {
		fmt.Println("get flag 'size' failed, err:", err)
		return err
	}
	if size < 0 || size > 2*1024*1024*1024 {
		err := errors.New("flag 'size' must be in 0-2147483648(2GB)")
		return err
	}

	return nil
}

// RunE 用于执行命令
func RunE(cmd *cobra.Command, args []string) error {

	// 获取flag 's'的值
	// size可能大于1024*1024*1024，所以使用int64
	size, err := cmd.Flags().GetInt64("size")
	if err != nil {
		fmt.Println("get flag 'size' failed, err:", err)
		return err
	}

	// 获取flag 'n'的值
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		fmt.Println("get flag 'name' failed, err:", err)
		return err
	}

	// flag 'a' 是否指定
	append := cmd.Flags().Changed("append")

	// 生成随机字符串并写入文件
	writeRandomStringToFile(size, name, append)
	return nil

}

// 将指定大小的随机字符串写入文件
func writeRandomStringToFile(size int64, name string, append bool) error {

	// 打开或创建文件
	// 判断文件是否有任何后缀
	if !strings.Contains(name, ".") {
		// 如果没有后缀，添加后缀.log
		name += ".log"
	}

	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return err
	}
	defer file.Close()

	// 如果指定了追加模式，将文件指针移动到文件末尾
	if append {
		_, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			fmt.Println("seek file failed, err:", err)
			return err
		}
	} else {
		// 如果不是追加模式，清空文件
		err := file.Truncate(0)
		if err != nil {
			fmt.Println("truncate file failed, err:", err)
			return err
		}
	}

	// 生成随机字符串流
	bullshit := utils.GetRandomStringStream(size, letterBytes, threshold)

	// 记录已经写入文件的字节数
	var written int64
	// 进度
	var progress int64

	for s := range bullshit {
		len := len(s)
		// 写入文件
		if _, err := file.WriteString(s); err != nil {
			fmt.Println("write file failed, err:", err)
			return err
		}
		written += int64(len)
		nowProgress := written * 100 / size
		// 只有当写入大小大于阈值，才打印进度
		if size > threshold {
			// 当进度每增加1%时，打印进度
			if nowProgress >= progress+1 {
				// 进度条带箭头
				fmt.Printf("\rwrite file progress: [%-50s] %d%%", strings.Repeat("=", int(nowProgress/2))+">", nowProgress)
				progress = nowProgress
			}
		}
	}
	if size > threshold {
		fmt.Println()
	}
	fmt.Println("write file success, all size:", size, "bytes")
	return nil
}
