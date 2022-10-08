package main

import (
	"errors"
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
	flag.StringVar(&data.LanguageName, "cdk", "java", "The type of cdk. The value can be either java or go")
	flag.StringVar(&data.ConnectorName, "name", "Sns", "The name of your connector")
	flag.StringVar(&data.ConnectorType, "type", "Source", "The type of your connector. The value can be either source or sink")
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

func (data *ConnectorData) validate() (err error) {
	lowerLan := strings.ToLower(data.LanguageName)
	lowerType := strings.ToLower(data.ConnectorType)
	if lowerLan != "java" && lowerLan != "go" {
		err = errors.New("illegal language name")
	} else if lowerType != "source" && lowerType != "sink" {
		err = errors.New("illegal connector type")
	}
	return err
}

func main() {
	data := &ConnectorData{}
	parseFlag(data)
	err := data.validate()
	if err != nil {
		panic(err)
	}
	flag.Parse()
	funcMap = template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}
	if strings.ToLower(data.LanguageName) == "java" {
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
		} else if strings.ToLower(data.ConnectorType) == "sink" {
			processTemplate("sink.tmpl", "./template/java/sink/sink.tmpl",
				outputPath+"/"+data.ConnectorName+"Sink.java", data)
			processTemplate("config.tmpl", "./template/java/sink/config.tmpl",
				outputPath+"/"+data.ConnectorName+"Config.java", data)
		}
	}
}
