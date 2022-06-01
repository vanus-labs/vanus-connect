package config

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func LoadScalerConfig() map[string][][]string {
	pwd, _ := os.Getwd()
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
	}

	content, err := os.ReadFile("pkg/config/scale_config.json")
	if err != nil {
		fmt.Println("read local file err, start to read from config", err)

		content, err = os.ReadFile("/config/scale_config.json")
		if err != nil {
			fmt.Println("read config file error", err)
			return nil
		}

	}

	retMap := make(map[string][][]string)
	conf, ok := gjson.ParseBytes(content).Value().(map[string]interface{})

	if ok {
		for k, v := range conf {
			arr := v.([]interface{})
			retArr := make([][]string, len(arr))
			for i, vv := range arr {
				arr2 := vv.([]interface{})
				subArr := make([]string, len(arr2))
				for k, vvv := range arr2 {
					subArr[k] = vvv.(string)
				}
				retArr[i] = subArr
			}
			retMap[k] = retArr
		}
	}
	return retMap
}
func PrintScalerConfig() {
	mapa := LoadScalerConfig()
	for k, v := range mapa {
		fmt.Printf("key: %s", k)
		for _, vv := range v {
			fmt.Printf("[")
			for _, vvv := range vv {
				fmt.Printf("%s ,", vvv)
			}
			fmt.Printf("]")
		}
		fmt.Printf("\n")
	}
}

func loadScalerConfig2() map[string][][]string {

	content, err := os.ReadFile("pkg/config/scale_config.json")
	if err != nil {
		fmt.Println("read local file err, start to read from config", err)
		absPath, _ := filepath.Abs("pkg/config/scale_config.json")
		fmt.Println("absPath: " + absPath)
		content, err = os.ReadFile(absPath)
		if err != nil {
			fmt.Println("read config file error", err)
			return nil
		}

	}

	retMap := make(map[string][][]string)
	conf, ok := gjson.ParseBytes(content).Value().(map[string]interface{})

	if ok {
		for k, v := range conf {
			arr := v.([]interface{})
			retArr := make([][]string, len(arr))
			for i, vv := range arr {
				arr2 := vv.([]interface{})
				subArr := make([]string, len(arr2))
				for k, vvv := range arr2 {
					subArr[k] = vvv.(string)
				}
				retArr[i] = subArr
			}
			retMap[k] = retArr
		}
	}
	return retMap
}
