package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/88250/lute"
	"github.com/Wsine/feishu2md/core"
	"github.com/Wsine/feishu2md/utils"
	"github.com/pkg/errors"
)

type HugoOpts struct {
	outputDir string
	dump      bool
	batch     bool
	wiki      bool
	tags      string
}

var hugoOpts = HugoOpts{}

func hugoDocument(ctx context.Context, client *core.Client, url string, opts *HugoOpts) error {
	// 获取当前工作目录
	here, err := os.Getwd()
	if err != nil {
		return err
	}

	// 向上查找hugo根目录
	// 有子文件夹config、content、static、themes
	// 有子文件hugo.yaml或hugo.toml或hugo.json或hugo.yml
	sub_dirs := []string{"config", "content", "static", "themes"}
	sub_files := []string{"hugo.yaml", "hugo.toml", "hugo.json", "hugo.yml"}
	// 需要有所有的子目录，文件有一个即可
	for {
		// 如果here是根目录了，/ 或者 Windows中的某个盘符
		if here == "/" || (len(here) == 2 && here[1] == ':') {
			return errors.New("hugo root directory not found")
		}
		all_dirs_exist := true
		for _, sub_dir := range sub_dirs {
			if _, err := os.Stat(filepath.Join(here, sub_dir)); os.IsNotExist(err) {
				all_dirs_exist = false
				break
			}
		}
		if !all_dirs_exist {
			// 向上一级目录查找
			here = filepath.Dir(here)
			continue
		}

		// 检查文件是否存在
		file_exists := false
		for _, sub_file := range sub_files {
			if _, err := os.Stat(filepath.Join(here, sub_file)); !os.IsNotExist(err) {
				file_exists = true
				break
			}
		}
		if file_exists {
			break // 找到hugo根目录了
		}

		// 向上一级目录查找
		here = filepath.Dir(here)
	}
	// Validate the url to download
	docType, docToken, err := utils.ValidateDocumentURL(url)
	if err != nil {
		return err
	}
	fmt.Println("Captured document token:", docToken)

	// for a wiki page, we need to renew docType and docToken first
	if docType == "wiki" {
		node, err := client.GetWikiNodeInfo(ctx, docToken)
		if err != nil {
			err = fmt.Errorf("GetWikiNodeInfo err: %v for %v", err, url)
		}
		utils.CheckErr(err)
		docType = node.ObjType
		docToken = node.ObjToken
	}
	if docType == "docs" {
		return errors.Errorf(
			`Feishu Docs is no longer supported. ` +
				`Please refer to the Readme/Release for v1_support.`)
	}

	// Process the download
	docx, blocks, err := client.GetDocxContent(ctx, docToken)
	utils.CheckErr(err)

	parser := core.NewParser(dlConfig.Output)
	parser.ForHugo = true

	title := docx.Title
	markdown := parser.ParseDocxContent(docx, blocks)

	sanitizeFileName := utils.SanitizeFileName(title)
	img_dir := filepath.Join(here, "static", "post_imgs", sanitizeFileName)
	img_dir = strings.ReplaceAll(img_dir, " ", "_")
	if _, err := os.Stat(img_dir); os.IsNotExist(err) {
		if err := os.MkdirAll(img_dir, 0o755); err != nil {
			return err
		}
	}

	first_img := ""

	if !dlConfig.Output.SkipImgDownload {
		for _, imgToken := range parser.ImgTokens {
			localLink, err := client.DownloadImage(
				ctx, imgToken, img_dir,
			)
			if err != nil {
				return err
			}
			if first_img == "" {
				first_img = localLink
			}
			localLink = strings.Replace(localLink, "static/", "", 1)
			markdown = strings.Replace(markdown, imgToken, localLink[2:], 1)
		}
	}

	// Format the markdown document
	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})
	result := engine.FormatStr("md", markdown)

	// 所有标题降一级
	re := regexp.MustCompile(`(?m)^(#{1,6})(\s+.*)$`)
	result = re.ReplaceAllStringFunc(result, func(line string) string {
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			return line
		}
		hashes := matches[1]
		content := matches[2]

		// 如果已经是 6 级标题，就不再降级
		if len(hashes) >= 6 {
			return "######" + content
		}
		return strings.Repeat("#", len(hashes)+1) + content
	})
	// // Handle the output directory and name
	// if _, err := os.Stat(opts.outputDir); os.IsNotExist(err) {
	// 	if err := os.MkdirAll(opts.outputDir, 0o755); err != nil {
	// 		return err
	// 	}
	// }

	tags := strings.Split(hugoOpts.tags, ",")

	// Write to markdown file
	mdName := fmt.Sprintf("%s.md", docToken)
	if dlConfig.Output.TitleAsFilename {
		mdName = fmt.Sprintf("%s.md", sanitizeFileName)
	}
	now := time.Now().Format("2006-01-02T15:04:05-07:00")
	metaContent := "---\n"
	metaContent += "title: \"" + title + "\"\n"
	metaContent += "date: " + now + "\n"
	metaContent += "categories:\n"
	metaContent += "  - 算法\n"
	metaContent += "  - 论文\n"
	metaContent += "tags:\n"
	for _, tag := range tags {
		metaContent += "  - " + tag + "\n"
	}
	metaContent += "---\n\n"
	outputDir := filepath.Join(here, "content", "posts", mdName[:len(mdName)-3])
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return err
		}
	}
	outputPath := filepath.Join(outputDir, "index.md")
	if err = os.WriteFile(outputPath, []byte(metaContent+result), 0o644); err != nil {
		return err
	}
	// 复制第一张图片作为封面
	if first_img != "" {
		coverPath := filepath.Join(outputDir, "featured"+filepath.Ext(first_img))
		input, err := os.ReadFile(first_img)
		if err != nil {
			return err
		}
		if err = os.WriteFile(coverPath, input, 0o644); err != nil {
			return err
		}
	}
	fmt.Printf("Downloaded markdown file to %s\n", outputPath)

	return nil
}

func hugoDocuments(ctx context.Context, client *core.Client, url string) error {
	// Validate the url to download
	folderToken, err := utils.ValidateFolderURL(url)
	if err != nil {
		return err
	}
	fmt.Println("Captured folder token:", folderToken)

	// Error channel and wait group
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	// Recursively go through the folder and download the documents
	var processFolder func(ctx context.Context, folderPath, folderToken string) error
	processFolder = func(ctx context.Context, folderPath, folderToken string) error {
		files, err := client.GetDriveFolderFileList(ctx, nil, &folderToken)
		if err != nil {
			return err
		}
		opts := HugoOpts{outputDir: folderPath, dump: hugoOpts.dump, batch: false}
		for _, file := range files {
			if file.Type == "folder" {
				_folderPath := filepath.Join(folderPath, file.Name)
				if err := processFolder(ctx, _folderPath, file.Token); err != nil {
					return err
				}
			} else if file.Type == "docx" {
				// concurrently download the document
				wg.Add(1)
				go func(_url string) {
					if err := hugoDocument(ctx, client, _url, &opts); err != nil {
						errChan <- err
					}
					wg.Done()
				}(file.URL)
			}
		}
		return nil
	}
	if err := processFolder(ctx, hugoOpts.outputDir, folderToken); err != nil {
		return err
	}

	// Wait for all the downloads to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		return err
	}
	return nil
}

func hugoWiki(ctx context.Context, client *core.Client, url string) error {
	prefixURL, spaceID, err := utils.ValidateWikiURL(url)
	if err != nil {
		return err
	}

	folderPath, err := client.GetWikiName(ctx, spaceID)
	if err != nil {
		return err
	}
	if folderPath == "" {
		return fmt.Errorf("failed to GetWikiName")
	}

	errChan := make(chan error)

	var maxConcurrency = 10 // Set the maximum concurrency level
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, maxConcurrency) // Create a semaphore with the maximum concurrency level

	var downloadWikiNode func(ctx context.Context,
		client *core.Client,
		spaceID string,
		parentPath string,
		parentNodeToken *string) error

	downloadWikiNode = func(ctx context.Context,
		client *core.Client,
		spaceID string,
		folderPath string,
		parentNodeToken *string) error {
		nodes, err := client.GetWikiNodeList(ctx, spaceID, parentNodeToken)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			if n.HasChild {
				_folderPath := filepath.Join(folderPath, n.Title)
				if err := downloadWikiNode(ctx, client,
					spaceID, _folderPath, &n.NodeToken); err != nil {
					return err
				}
			}
			if n.ObjType == "docx" {
				opts := HugoOpts{outputDir: folderPath, dump: hugoOpts.dump, batch: false}
				wg.Add(1)
				semaphore <- struct{}{}
				go func(_url string) {
					if err := hugoDocument(ctx, client, _url, &opts); err != nil {
						errChan <- err
					}
					wg.Done()
					<-semaphore
				}(prefixURL + "/wiki/" + n.NodeToken)
				// hugoDocument(ctx, client, prefixURL+"/wiki/"+n.NodeToken, &opts)
			}
		}
		return nil
	}

	if err = downloadWikiNode(ctx, client, spaceID, folderPath, nil); err != nil {
		return err
	}

	// Wait for all the downloads to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		return err
	}
	return nil
}

func handleHugoCommand(url string) error {
	// Load config
	configPath, err := core.GetConfigFilePath()
	if err != nil {
		return err
	}
	config, err := core.ReadConfigFromFile(configPath)
	if err != nil {
		return err
	}
	dlConfig = *config

	// Instantiate the client
	client := core.NewClient(
		dlConfig.Feishu.AppId, dlConfig.Feishu.AppSecret,
	)
	ctx := context.Background()

	if hugoOpts.batch {
		return hugoDocuments(ctx, client, url)
	}

	if hugoOpts.wiki {
		return hugoWiki(ctx, client, url)
	}

	return hugoDocument(ctx, client, url, &hugoOpts)
}
