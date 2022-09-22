package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type ConnectorData struct {
	LanguageName  string
	ConnectorName string
	ConnectorType string
}

var funcMap template.FuncMap

func parseFlag(data *ConnectorData) {
	flag.StringVar(&data.LanguageName, "cdk", "java", "请输入Connector的语言环境")
	flag.StringVar(&data.ConnectorName, "name", "Sns", "请输入Connector的名字")
	flag.StringVar(&data.ConnectorType, "type", "Source", "请输入Connector的类型")
}

func processTemplate(tmplName string, tmplPath string, outFilePath string, data *ConnectorData) {
	tmpl, err := template.New(tmplName).Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Writing file: ", outFilePath)
	file, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	data := &ConnectorData{}
	parseFlag(data)
	flag.Parse()
	funcMap = template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}
	if data.LanguageName == "java" {
		projectPath := "./" + strings.ToLower(data.ConnectorType) + "-" + strings.ToLower(data.ConnectorName)
		outputPath := projectPath + "/src/main/java/com/vance/" + strings.ToLower(data.ConnectorType) +
			"/" + strings.ToLower(data.ConnectorName)
		err := os.MkdirAll(outputPath, 0777)
		if err != nil {
			fmt.Println(err)
		}
		err = os.MkdirAll(projectPath+"/"+"src/main/resources", 0777)
		if err != nil {
			fmt.Println(err)
		}

		processTemplate("README.tmpl", "./template/README.tmpl",
			projectPath+"/"+"README.md", data)

		processTemplate("pom.tmpl", "./template/java/pom.tmpl", "./"+
			strings.ToLower(data.ConnectorType)+"-"+strings.ToLower(data.ConnectorName)+"/"+"pom.xml", data)
		processTemplate("entrance.tmpl", "./template/java/entrance.tmpl",
			outputPath+"/"+"Entrance.java", data)

		if strings.ToLower(data.ConnectorType) == "source" {
			processTemplate("source.tmpl", "./template/java/source/source.tmpl",
				outputPath+"/"+data.ConnectorName+"Source.java", data)
			processTemplate("config.tmpl", "./template/java/source/config.tmpl",
				outputPath+"/"+data.ConnectorName+"Config.java", data)
			processTemplate("adapter.tmpl", "./template/java/source/adapter.tmpl",
				outputPath+"/"+data.ConnectorName+"Adapter.java", data)
			processTemplate("adapted_sample.tmpl", "./template/java/source/adapted_sample.tmpl",
				outputPath+"/"+data.ConnectorName+"AdaptedSample.java", data)
		} else if data.ConnectorType == "Sink" {
			processTemplate("sink.tmpl", "./template/java/sink/sink.tmpl",
				outputPath+"/"+data.ConnectorName+"Sink.java", data)
			processTemplate("config.tmpl", "./template/java/sink/config.tmpl",
				outputPath+"/"+data.ConnectorName+"Config.java", data)
		}
	}
}
