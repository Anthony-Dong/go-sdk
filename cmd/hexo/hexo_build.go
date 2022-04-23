package hexo

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/collections"
	"github.com/anthony-dong/go-sdk/commons/logs"
	git "github.com/sabhiram/go-gitignore"
	"gopkg.in/yaml.v3"
)

const (
	delimiter         = "---"
	abstractDelimiter = "<!-- more -->"
)

type Config struct {
	Title      string   `yaml:"title"`       // 标题(如果没有设置，为源文件的名称)
	TargetFile string   `yaml:"target_file"` // 目标文件，值得是生成的文件
	OriginFile string   `yaml:"origin_file"` // 原文件，指的是我们写的文件
	Date       string   `yaml:"date"`        // 日期(为文件的修改日期)
	Tags       []string `yaml:"tags,omitempty"`
	Categories []string `yaml:"categories,omitempty"`
}

type CheckFileCanHexoResult struct {
	CanHexo       bool
	HasAbstract   bool
	Content       []string
	Config        *Config
	FileInfo      os.FileInfo
	FileName      string // 源文件
	WriteFilePath string // 写入的目录
	NeedWrite     bool   // 是否需要重新写入
	ContentData   []byte
}

// WriteFile
// 1. 生成文件内容
// 2. 同时格式化源文件.
func (c *CheckFileCanHexoResult) WriteFile() error {
	checkFileNeedReWrite := func(data []byte) (bool, error) {
		// 当前文件
		oldData, err := ioutil.ReadFile(c.FileName)
		if err != nil {
			return false, err
		}
		return codec.Md5HexString(oldData) != codec.Md5HexString(data), nil
	}

	// 组装文件
	buffer := &bytes.Buffer{}
	// write hexo header
	buffer.Write([]byte(delimiter))
	buffer.WriteByte('\n')
	cfg, err := yaml.Marshal(c.Config)
	if err != nil {
		return err
	}
	buffer.Write(cfg)
	buffer.Write([]byte(delimiter))
	buffer.WriteByte('\n')
	buffer.WriteByte('\n')
	// write content
	buffer.Write(commons.UnsafeBytes(commons.LinesToString(c.Content)))

	// 写入的文件
	data := buffer.Bytes()
	// 检测文件是否需要重写，通过MD5校验
	needWrite, err := checkFileNeedReWrite(data)
	if err != nil {
		return err
	}
	if needWrite {
		logs.Infof("[Hexo] 发现MD5比较不一致需要重新写入到源文件, 文件: %s", c.FileName)
		if err := ioutil.WriteFile(c.FileName, data, commons.DefaultFileMode); err != nil {
			return err
		}
	}
	c.ContentData = data
	return nil
}

// buildHexo
// targetDir: /Users/bytedance/note/note/hexo-home/source/_posts
// dir: /Users/bytedance/note/note
func buildHexo(ctx context.Context, dir string, targetDir string, firmCode []string, ignore []string) error {
	// 1. 获取 dir 全部需要处理的文件
	// 2. 获取 targetDir 已经有的文件
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	logs.Debugf("[Hexo] 开始全部的Markdown的文件, 目录: %s", dir)
	allPage, err := GetAllPage(dir, ignore)
	if err != nil {
		return err
	}
	logs.Debugf("[Hexo] 获取全部的Markdown文件成功, 目录: %s, 总数: %d", dir, len(allPage))

	logs.Debugf("[Hexo] 开始全部的Target-Markdown的文件, 目录: %s", targetDir)
	targetPage, err := GetAllMarkDownPage(targetDir)
	if err != nil {
		return err
	}
	logs.Debugf("[Hexo] 获取全部的Target-Markdown文件成功, 目录: %s, 总数: %d", targetDir, len(targetPage))
	targetPageSet := collections.NewSet(targetPage)
	newTargetPage := collections.NewSetInitSize(targetPageSet.Size())
	needWriteNewTargetPage := collections.NewSetInit()

	wg := sync.WaitGroup{}
	for _, file := range allPage {
		wg.Add(1)
		func(fileName string) {
			needEnd := false
			targetFile := ""
			getRelativePath := func(filename string) string {
				relativePath, err := filepath.Rel(dir, filename)
				if err != nil {
					return filename
				}
				return relativePath
			}
			defer func() {
				if err := recover(); err != nil {
					logs.Errorf("[Hexo] 运行期间发现 panic, 文件: %s, 异常: %v", fileName, err)
				}
				wg.Done()
				if needEnd {
					logs.Debugf("[Hexo] 结束操作文件, 文件: %s, 目标文件: %s", getRelativePath(fileName), getRelativePath(targetFile))
				}
			}()
			// 检测是否是hexo
			result, err := CheckFileCanHexo(fileName, dir)
			if err != nil {
				logs.Errorf("[Hexo] 检测文件是否是hexo文件发现异常, 文件: %s, 异常: %v", fileName, err)
				return
			}
			if result == nil {
				return
			}
			logs.Debugf("[Hexo] 开始操作文件, 文件: %s", getRelativePath(fileName))
			if !result.HasAbstract {
				logs.Warnf("[Hexo] 警告, 发现没有摘要, 文件: %s", fileName)
			}
			needEnd = true
			// 开始检测是否有公司代码
			if err := CheckFileHasFirmCode(fileName, result.Content, firmCode); err != nil {
				logs.Errorf("[Hexo] 检测公司代码失败, 异常: %s, 文件: %s", err, fileName)
				return
			}

			// 检测文件格式写入（格式化原文件）
			if err := result.WriteFile(); err != nil {
				logs.Errorf("[Hexo] 检测原文件格式失败, 异常: %s, 文件: %s", err, fileName)
				return
			}

			// copy文件
			targetFile = filepath.Join(targetDir, result.Config.TargetFile)
			writeTargetFile := func() {
				if err := ioutil.WriteFile(targetFile, result.ContentData, commons.DefaultFileMode); err != nil {
					logs.Errorf("[Hexo] 发现需要写入到Hexo的post目录文件发现了异常, 文件: %s, 异常: %v", targetFile, err)
					return
				}
				needWriteNewTargetPage.Put(targetFile)
			}

			// 如果原来post目录不存在
			if !targetPageSet.Contains(targetFile) {
				logs.Infof("[Hexo] 发现Hexo的post目录文件不存在需要写入的文件, 目标文件: %s, 源文件: %s", targetFile, fileName)
				// 直接写入
				writeTargetFile()
			} else {
				// 否则读出来比较MD5是否相同，再写入！
				readBody, err := ioutil.ReadFile(targetFile)
				if err != nil {
					logs.Errorf("[Hexo] 发现读取Hexo的post目录文件发现了异常, 文件: %s, 异常: %v", targetFile, err)
					return
				}
				if codec.Md5HexString(result.ContentData) != codec.Md5HexString(readBody) {
					logs.Infof("[Hexo] 发现读取Hexo的post目录文件和原文件MD5值不一样, 需要重写, 文件: %s,  源文件: %s", getRelativePath(targetFile), getRelativePath(fileName))
					writeTargetFile()
				}
			}

			// 操作完成
			newTargetPage.Put(targetFile)
		}(file)
	}
	wg.Wait()

	// delete
	logs.Infof("[Hexo] 操作脚本完成, 一共写入: %d, 总页数: %d", needWriteNewTargetPage.Size(), newTargetPage.Size())

	slice := targetPageSet.ToSlice()
	for _, elem := range slice {
		if newTargetPage.Contains(elem) {
			targetPageSet.Delete(elem)
		}
	}

	logs.Infof("[Hexo] 操作脚本完成需要删除文件: %d", targetPageSet.Size())

	if targetPageSet.Size() == 0 {
		return nil
	}

	for _, elem := range targetPageSet.ToSlice() {
		logs.Infof("[Hexo] 删除文件, 文件: %s", elem)
		if err := os.Remove(elem); err != nil {
			logs.Errorf("[Hexo] 删除文件失败, 文件: %s, 异常: %s", elem, err)
			return err
		}
	}

	return nil
}

// CheckFileHasFirmCode 检测是否有公司代码.
func CheckFileHasFirmCode(fileName string, content []string, firmCode []string) error {
	if len(firmCode) == 0 || len(content) == 0 {
		return nil
	}
	for index, line := range content {
		for _, elem := range firmCode {
			if strings.Contains(line, elem) {
				logs.Warnf("[Hexo] 发现公司代码, 文件: %s, 公司代码: %s, 原文: %s", fileName, elem, line)
				newElem := commons.NewString('x', len(elem))
				line = strings.ReplaceAll(line, elem, newElem)
			}
		}
		content[index] = line
	}
	return nil
}

// CheckFileCanHexoPre 检测是否可能是 hexo文件，防止遍历全文.
func CheckFileCanHexoPre(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer file.Close()
	trueError := errors.New("true")
	falseError := errors.New("false")
	count := 0
	if err := commons.ReadLineByFunc(file, func(line string) error {
		if line == "" {
			return nil
		}
		count++
		if count == 1 && line == delimiter { // 如果第一行不为空的数据是 delimiter则是，否则不是
			return trueError
		} else {
			return falseError
		}
	}); err != nil {
		return err == trueError
	}
	return false
}

// CheckFileCanHexo 检测文件是否可以转换为hexo文件.
func CheckFileCanHexo(fileName string, fileParentPath string) (*CheckFileCanHexoResult, error) {
	if !CheckFileCanHexoPre(fileName) {
		return nil, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 是否有摘要分隔符
	hasAbstract := false

	// hexo 头部信息分隔符出现的次数
	delimiterCount := 0
	delimiterErr := errors.New("delimiter")

	// yaml-config
	yamlConfig := make([]string, 0)
	// body
	content := make([]string, 0)

	// 是否有空格
	// 1、必须携带有分隔符
	// 2、必须是 isspace=true 如果出现了非空格则设置为false
	if err := commons.ReadLineByFunc(file, func(line string) error {
		// 如果分隔符
		if !hasAbstract && line == delimiter {
			delimiterCount++
			return nil
		}

		// 不是
		if !hasAbstract && strings.Contains(line, abstractDelimiter) {
			hasAbstract = true
		}

		// 如果是刚开始则为 yaml
		if delimiterCount == 1 {
			yamlConfig = append(yamlConfig, line)
			return nil
		}

		// 为正文
		if delimiterCount >= 2 {
			if line == "" && len(content) == 0 {
				return nil
			}
			content = append(content, line)
			return nil
		}
		return nil
	}); err != nil {
		if err != delimiterErr {
			return nil, err
		}
	}
	canHexo := delimiterCount == 2
	if !canHexo {
		return nil, nil
	}

	yamlConfigContent := commons.LinesToString(yamlConfig)
	fileConfig := new(Config)
	err = yaml.Unmarshal(commons.UnsafeBytes(yamlConfigContent), fileConfig)
	if err != nil {
		return nil, err
	}
	// title为空
	if fileConfig.Title == "" {
		return nil, errors.New("the hexo title can not null")
	}
	// 源文件(相对路径)
	fileConfig.OriginFile, err = filepath.Rel(fileParentPath, fileName)
	if err != nil {
		return nil, err
	}

	// 如果目标文件已经定义过了，这里就不需要再创建了
	if fileConfig.TargetFile == "" {
		fileConfig.TargetFile = codec.Md5HexString([]byte(fileConfig.Title)) + ".md"
	}

	// 文件修改时间
	if fileConfig.Date == "" {
		fileConfig.Date = fileInfo.ModTime().Format(commons.FormatTimeV1)
	}
	return &CheckFileCanHexoResult{
		CanHexo:     canHexo,
		HasAbstract: hasAbstract,
		Content:     content,
		FileInfo:    fileInfo,
		Config:      fileConfig,
		FileName:    fileName,
	}, nil
}
func GetAllPage(dir string, ignoreLine []string) ([]string, error) {
	var (
		ignore *git.GitIgnore
	)

	ignoreFileName := filepath.Join(dir, ".gitignore")
	if commons.Exist(ignoreFileName) {
		var err error
		ignore, err = git.CompileIgnoreFileAndLines(ignoreFileName, ignoreLine...)
		if err != nil {
			logs.Errorf("[GetAllPage] load ignore file: %s find err: %v", ignoreFileName, err)
			return nil, err
		}
		logs.Infof("[GetAllPage] load ignore file: %s, config: %s, success!", ignoreFileName, commons.ToString(ignoreLine))
	} else {
		ignore = git.CompileIgnoreLines(ignoreLine...)
		logs.Infof("[GetAllPage] load ignore config: %s success!", commons.ToString(ignoreLine))
	}

	return commons.GetAllFiles(dir, func(filePath string) bool {
		fileSuffix := filepath.Ext(filePath)
		fileName := filepath.Base(filePath)
		if !(fileSuffix == ".md" || fileSuffix == ".markdown") {
			return false
		}
		baseDir, err := filepath.Rel(dir, filePath)
		if err != nil {
			return false
		}
		if ignore.MatchesPath(baseDir) { // 如果命中ignore 舍弃
			return false
		}
		if strings.Contains(fileName, "README") {
			return false
		}
		return true
	})
}

func GetAllMarkDownPage(dir string) ([]string, error) {
	return commons.GetAllFiles(dir, func(fileName string) bool {
		suffix := filepath.Ext(fileName)
		return suffix == ".md" || suffix == ".markdown"
	})
}
