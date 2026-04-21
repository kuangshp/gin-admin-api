package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func getGoctlPath() string {
	cmd := exec.Command("which", "go")
	output, err := cmd.Output()
	if err != nil {
		return "goctl"
	}
	goPath := strings.TrimSpace(string(output))
	goDir := goPath[:strings.LastIndex(goPath, "/")]
	return goDir + "/goctl"
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	// 使用Title后，再把特定的全大写缩写词恢复
	name = strings.Title(name)
	// 处理常见的缩写词，如 IP, ID, URL, API 等
	acronyms := []string{"IP", "ID", "URL", "API", "IOS", "API", "XML", "JSON", "JWT", "SQL", "ORM"}
	for _, acronym := range acronyms {
		name = strings.ReplaceAll(name, strings.Title(strings.ToLower(acronym)), acronym)
	}
	return strings.Replace(name, " ", "", -1)
}

func LowerCamelCase(name string) string {
	name = Case2Camel(name)
	return strings.ToLower(name[:1]) + name[1:]
}

var inputLines []string
var inputIndex int

func readInput(prompt string) string {
	if inputIndex < len(inputLines) {
		line := inputLines[inputIndex]
		inputIndex++
		fmt.Println(line)
		return line
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}

func init() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, _ := io.ReadAll(os.Stdin)
		if len(data) > 0 {
			allInput := string(data)
			inputLines = strings.Split(allInput, "\n")
			for i, line := range inputLines {
				inputLines[i] = strings.TrimSpace(line)
			}
		}
	}
}

func getTableColumns(db *gorm.DB, tableName string) ([]ColumnInfo, error) {
	type Column struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default *string
		Extra   string
		Comment string
	}
	var columns []Column
	err := db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", tableName)).Scan(&columns).Error

	result := make([]ColumnInfo, len(columns))
	for i, col := range columns {
		result[i] = ColumnInfo{
			Name:    col.Field,
			Type:    col.Type,
			CanNull: col.Null == "YES",
			IsKey:   col.Key == "PRI",
			Extra:   col.Extra,
			Comment: col.Comment,
		}
	}
	return result, err
}

type ColumnInfo struct {
	Name       string
	Type       string
	FieldName  string
	FieldType  string
	JsonTag    string
	JsonTagOpt string
	CanNull    bool
	IsKey      bool
	Extra      string
	Comment    string
	Validate   string
}

type ApiTemplateData struct {
	TableName    string
	ModelName    string
	EntityName   string
	TableComment string
	Columns      []ColumnInfo
}

type RepositoryTemplateData struct {
	ModelName       string
	ModelNameLower  string
	EntityName      string
	EntityNameLower string
	Package         string
	DaoPath         string
	ModelPath       string
	ModelPkgName    string // model包的名称，如 "entity"
	Columns         []ColumnInfo
}

func getTableComment(db *gorm.DB, tableName string) (string, error) {
	var comment string
	err := db.Raw(fmt.Sprintf("SELECT TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = '%s'", tableName)).Scan(&comment).Error
	return comment, err
}

func generateValidateRule(col ColumnInfo) string {
	var rules []string
	if !col.CanNull {
		rules = append(rules, "required")
	}
	if col.IsKey {
		rules = append(rules, "uuid")
	}
	if col.FieldType == "string" && strings.Contains(col.Name, "email") {
		rules = append(rules, "email")
	}
	if col.FieldType == "string" && strings.Contains(col.Name, "mobile") {
		rules = append(rules, "mobile")
	}
	if strings.Contains(strings.ToLower(col.Comment), "1是") && strings.Contains(strings.ToLower(col.Comment), "2是") {
		enumRegex := regexp.MustCompile(`(\d+)是([^，,]+)[,，]?`)
		matches := enumRegex.FindAllStringSubmatch(col.Comment, -1)
		if len(matches) > 0 {
			values := make([]string, 0, len(matches))
			for _, match := range matches {
				if len(match) >= 3 {
					values = append(values, match[1])
				}
			}
			if len(values) > 0 {
				rules = append(rules, fmt.Sprintf("oneof=%s", strings.Join(values, " ")))
			}
		}
	}
	if !col.CanNull && col.FieldType == "int64" && (strings.Contains(col.Name, "status") || strings.Contains(col.Name, "type") || strings.Contains(col.Name, "is_")) {
		rules = append(rules, "gte=1")
	}
	return strings.Join(rules, ",")
}

func loadTemplate(templatePath string) (*template.Template, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	templateName := filepath.Base(templatePath)
	return template.New(templateName).Parse(string(content))
}

func generateApiFile(tableName string, columns []ColumnInfo, modelName string, db *gorm.DB, tmplPath string) (string, error) {
	tmpl, err := loadTemplate(tmplPath)
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	columnData := make([]ColumnInfo, len(columns))
	for i, col := range columns {
		jsonTagOpt := ""
		if col.CanNull {
			jsonTagOpt = ",optional"
		}
		columnData[i] = ColumnInfo{
			Name:       col.Name,
			Type:       col.Type,
			FieldName:  Case2Camel(col.Name),
			FieldType:  getGoType(col.Type),
			JsonTag:    LowerCamelCase(col.Name),
			JsonTagOpt: jsonTagOpt,
			CanNull:    col.CanNull,
			IsKey:      col.IsKey,
			Extra:      col.Extra,
			Comment:    col.Comment,
			Validate:   generateValidateRule(ColumnInfo{CanNull: col.CanNull, IsKey: col.IsKey, FieldType: getGoType(col.Type), Name: col.Name, Comment: col.Comment}),
		}
	}

	tableComment, _ := getTableComment(db, tableName)
	if tableComment == "" {
		tableComment = modelName
	}

	data := ApiTemplateData{
		TableName:    tableName,
		ModelName:    modelName,
		EntityName:   Case2Camel(strings.ToUpper(tableName[:1]+tableName[1:])) + "Entity",
		TableComment: tableComment,
		Columns:      columnData,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	return buf.String(), nil
}

func generateRepositoryFile(columns []ColumnInfo, modelName string, pkg string, daoPath string, modelPath string, tmplPath string) (string, error) {
	tmpl, err := loadTemplate(tmplPath)
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	columnData := make([]ColumnInfo, len(columns))
	for i, col := range columns {
		columnData[i] = ColumnInfo{
			Name:      col.Name,
			Type:      col.Type,
			FieldName: Case2Camel(col.Name),
			FieldType: getGoType(col.Type),
			CanNull:   col.CanNull,
			IsKey:     col.IsKey,
			Comment:   col.Comment,
		}
	}

	data := RepositoryTemplateData{
		ModelName:       modelName,
		ModelNameLower:  LowerCamelCase(modelName),
		EntityName:      modelName + "Entity",
		EntityNameLower: LowerCamelCase(modelName + "Entity"),
		Package:         pkg,
		DaoPath:         pkg + "/" + daoPath,
		ModelPath:       pkg + "/" + modelPath,
		ModelPkgName:    getLastPathSegment(modelPath),
		Columns:         columnData,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	return buf.String(), nil
}

func generateRepositoryInterfaceFile(columns []ColumnInfo, modelName string, pkg string, daoPath string, modelPath string, tmplPath string) (string, error) {
	tmpl, err := loadTemplate(tmplPath)
	if err != nil {
		return "", fmt.Errorf("加载接口模板失败: %w", err)
	}

	columnData := make([]ColumnInfo, len(columns))
	for i, col := range columns {
		columnData[i] = ColumnInfo{
			Name:      col.Name,
			Type:      col.Type,
			FieldName: Case2Camel(col.Name),
			FieldType: getGoType(col.Type),
			CanNull:   col.CanNull,
			IsKey:     col.IsKey,
			Comment:   col.Comment,
		}
	}

	data := RepositoryTemplateData{
		ModelName:       modelName,
		ModelNameLower:  LowerCamelCase(modelName),
		EntityName:      modelName + "Entity",
		EntityNameLower: LowerCamelCase(modelName + "Entity"),
		Package:         pkg,
		DaoPath:         pkg + "/" + daoPath,
		ModelPath:       pkg + "/" + modelPath,
		ModelPkgName:    getLastPathSegment(modelPath),
		Columns:         columnData,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染接口模板失败: %w", err)
	}

	return buf.String(), nil
}

func generateRepositoryExtFile(columns []ColumnInfo, modelName string, pkg string, daoPath string, modelPath string, tmplPath string) (string, error) {
	tmpl, err := loadTemplate(tmplPath)
	if err != nil {
		return "", fmt.Errorf("加载扩展模板失败: %w", err)
	}

	columnData := make([]ColumnInfo, len(columns))
	for i, col := range columns {
		columnData[i] = ColumnInfo{
			Name:      col.Name,
			Type:      col.Type,
			FieldName: Case2Camel(col.Name),
			FieldType: getGoType(col.Type),
			CanNull:   col.CanNull,
			IsKey:     col.IsKey,
			Comment:   col.Comment,
		}
	}

	data := RepositoryTemplateData{
		ModelName:       modelName,
		ModelNameLower:  LowerCamelCase(modelName),
		EntityName:      modelName + "Entity",
		EntityNameLower: LowerCamelCase(modelName + "Entity"),
		Package:         pkg,
		DaoPath:         pkg + "/" + daoPath,
		ModelPath:       pkg + "/" + modelPath,
		ModelPkgName:    getLastPathSegment(modelPath),
		Columns:         columnData,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染扩展模板失败: %w", err)
	}

	return buf.String(), nil
}

func getGoType(sqlType string) string {
	sqlType = strings.ToLower(sqlType)
	if strings.Contains(sqlType, "varchar") || strings.Contains(sqlType, "text") || strings.Contains(sqlType, "char") {
		return "string"
	}
	if strings.Contains(sqlType, "int") {
		return "int64"
	}
	if strings.Contains(sqlType, "decimal") || strings.Contains(sqlType, "float") || strings.Contains(sqlType, "double") {
		return "float64"
	}
	if strings.Contains(sqlType, "datetime") || strings.Contains(sqlType, "timestamp") {
		return "int64"
	}
	if strings.Contains(sqlType, "date") {
		return "int64"
	}
	if strings.Contains(sqlType, "json") {
		return "string"
	}
	if strings.Contains(sqlType, "bool") {
		return "bool"
	}
	return "string"
}

// pathToPkg 将路径转换为包路径，如 "./query/dao" -> "query/dao"
func pathToPkg(path string) string {
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimSuffix(path, "/")
	return path
}

// getLastPathSegment 获取路径的最后一个段，如 "dal/model/entity" -> "entity"
func getLastPathSegment(path string) string {
	path = strings.TrimPrefix(path, "./")
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// Generate 代码生成器主函数
func Generate(cfg *Config) error {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	tableName := readInput("请输入表名（直接回车同步所有表）: ")

	// 获取模板路径 - 优先使用环境变量或当前工作目录
	templateDir := os.Getenv("GENERATOR_TEMPLATE_DIR")
	if templateDir == "" {
		cwd, _ := os.Getwd()
		// 尝试相对于当前工作目录的路径
		templateDir = filepath.Join(cwd, "generator", "template")
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			// 尝试相对于可执行文件的路径
			exePath, _ := os.Executable()
			templateDir = filepath.Join(filepath.Dir(exePath), "generator", "template")
		}
	}
	apiTmplPath := filepath.Join(templateDir, "api_template.txt")
	repoInterfaceTmplPath := filepath.Join(templateDir, "repository_interface_template.txt")
	repoBaseTmplPath := filepath.Join(templateDir, "repository_base_template.txt")
	repoExtTmplPath := filepath.Join(templateDir, "repository_ext_template.txt")

	// 确保路径有"."前缀
	if !strings.HasPrefix(cfg.OutPath, ".") {
		cfg.OutPath = "." + cfg.OutPath
	}
	if !strings.HasPrefix(cfg.ModelPkgPath, ".") {
		cfg.ModelPkgPath = "." + cfg.ModelPkgPath
	}
	if !strings.HasPrefix(cfg.RepoPath, ".") {
		cfg.RepoPath = "." + cfg.RepoPath
	}
	if !strings.HasPrefix(cfg.ApiPath, ".") {
		cfg.ApiPath = "." + cfg.ApiPath
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           cfg.OutPath,
		ModelPkgPath:      cfg.ModelPkgPath,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"json":      func(detailType gorm.ColumnType) (dataType string) { return "JSON" },
		"decimal":   func(detailType gorm.ColumnType) (dataType string) { return "Decimal" },
	}
	g.WithDataTypeMap(dataMap)

	g.WithModelNameStrategy(func(tableName string) (modelName string) {
		return Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:] + "Entity")
	})

	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		if strings.Contains(`deleted_at`, columnName) {
			return "-"
		}
		return LowerCamelCase(columnName)
	})

	autoUpdateTimeField := gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
		return map[string][]string{
			"column":  {"updated_at"},
			"comment": {"更新时间"},
		}
	})
	autoCreateTimeField := gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
		return map[string][]string{
			"column":  {"created_at"},
			"comment": {"创建时间"},
		}
	})
	softDeleteField := gen.FieldType("deleted_at", "gorm.DeletedAt")
	fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	var allModel []interface{}
	if tableName != "" {
		fmt.Printf("生成表 %s 的模型...\n", tableName)
		model := g.GenerateModel(tableName, fieldOpts...)
		allModel = []interface{}{model}
	} else {
		fmt.Println("生成所有表的模型...")
		allModel = g.GenerateAllTable(fieldOpts...)
	}

	g.ApplyBasic(allModel...)
	g.Execute()

	var columns []ColumnInfo
	var modelName string
	if tableName != "" {
		var err error
		columns, err = getTableColumns(db, tableName)
		if err != nil {
			return fmt.Errorf("获取表结构失败: %w", err)
		}
		modelName = Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:])
	}

	if cfg.RepoPath != "" && tableName != "" {
		repoDir := cfg.RepoPath
		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			os.MkdirAll(repoDir, 0755)
		}

		// 1. 生成基础 repository (xxx_base.go) - 始终重新生成
		repoBaseContent, err := generateRepositoryFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoBaseTmplPath)
		if err != nil {
			fmt.Printf("生成repository基础内容失败: %v\n", err)
		} else {
			repoBaseFileName := fmt.Sprintf("%s/%s_base.go", repoDir, strings.ToLower(modelName))
			err = os.WriteFile(repoBaseFileName, []byte(repoBaseContent), 0644)
			if err != nil {
				fmt.Printf("写入repository基础文件失败: %v\n", err)
			} else {
				fmt.Printf("repository基础文件已生成: %s\n", repoBaseFileName)
			}
		}

		// 2. 生成接口定义 (xxx_interface.go) - 如果已存在则跳过
		repoInterfaceContent, err := generateRepositoryInterfaceFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoInterfaceTmplPath)
		if err != nil {
			fmt.Printf("生成repository接口内容失败: %v\n", err)
		} else {
			repoInterfaceFileName := fmt.Sprintf("%s/%s_interface.go", repoDir, strings.ToLower(modelName))
			if _, err := os.Stat(repoInterfaceFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoInterfaceFileName, []byte(repoInterfaceContent), 0644)
				if err != nil {
					fmt.Printf("写入repository接口文件失败: %v\n", err)
				} else {
					fmt.Printf("repository接口文件已生成: %s\n", repoInterfaceFileName)
				}
			} else {
				fmt.Printf("repository接口文件已存在，不覆盖更新: %s\n", repoInterfaceFileName)
			}
		}

		// 3. 生成扩展 repository (xxx.go) - 如果已存在则跳过
		repoExtContent, err := generateRepositoryExtFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoExtTmplPath)
		if err != nil {
			fmt.Printf("生成repository扩展内容失败: %v\n", err)
		} else {
			repoExtFileName := fmt.Sprintf("%s/%s.go", repoDir, strings.ToLower(modelName))
			if _, err := os.Stat(repoExtFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoExtFileName, []byte(repoExtContent), 0644)
				if err != nil {
					fmt.Printf("写入repository扩展文件失败: %v\n", err)
				} else {
					fmt.Printf("repository扩展文件已生成: %s\n", repoExtFileName)
				}
			} else {
				fmt.Printf("repository扩展文件已存在，不覆盖更新: %s\n", repoExtFileName)
			}
		}

	if cfg.ApiPath != "" {
		apiDir := cfg.ApiPath
		if _, err := os.Stat(apiDir); os.IsNotExist(err) {
			os.MkdirAll(apiDir, 0755)
		}

		if tableName != "" {
			columns, err := getTableColumns(db, tableName)
			if err != nil {
				return fmt.Errorf("获取表结构失败: %w", err)
			}
			modelName := Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:])
			apiContent, err := generateApiFile(tableName, columns, modelName, db, apiTmplPath)
			if err != nil {
				return fmt.Errorf("生成api内容失败: %w", err)
			}
			apiFileName := fmt.Sprintf("%s/%s.api", apiDir, tableName)
			// 检查文件是否已存在
			if _, err := os.Stat(apiFileName); os.IsNotExist(err) {
				err = os.WriteFile(apiFileName, []byte(apiContent), 0644)
				if err != nil {
					return fmt.Errorf("写入api文件失败: %w", err)
				}
				fmt.Printf("api文件已生成: %s\n", apiFileName)

				goctlPath := getGoctlPath()
				cmd := exec.Command(goctlPath, "api", "go", "-api", apiFileName, "--dir", filepath.Dir(cfg.ApiPath), "--style=goZero")
				cmd.Dir = "."
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("执行goctl失败: %v\n%s\n", err, output)
				} else {
					fmt.Printf("go-zero代码生成成功\n%s\n", output)
				}
			} else {
				fmt.Printf("api文件已存在，不覆盖更新: %s\n", apiFileName)
			}
		}
	}

	if cfg.RepoPath != "" && tableName != "" {
		repoDir := cfg.RepoPath
		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			os.MkdirAll(repoDir, 0755)
		}

		// 1. 生成基础 repository (xxx_base.go) - 始终重新生成
			repoBaseContent, err := generateRepositoryFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoBaseTmplPath)
			if err != nil {
				return fmt.Errorf("生成repository基础内容失败: %w", err)
			}
			repoBaseFileName := fmt.Sprintf("%s/%s_base.go", repoDir, strings.ToLower(modelName))
			err = os.WriteFile(repoBaseFileName, []byte(repoBaseContent), 0644)
			if err != nil {
				return fmt.Errorf("写入repository基础文件失败: %w", err)
			}
			fmt.Printf("repository基础文件已生成: %s\n", repoBaseFileName)

			// 2. 生成接口定义 (xxx_interface.go) - 如果已存在则跳过
			repoInterfaceContent, err := generateRepositoryInterfaceFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoInterfaceTmplPath)
			if err != nil {
				return fmt.Errorf("生成repository接口内容失败: %w", err)
			}
			repoInterfaceFileName := fmt.Sprintf("%s/%s_interface.go", repoDir, strings.ToLower(modelName))
			if _, err := os.Stat(repoInterfaceFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoInterfaceFileName, []byte(repoInterfaceContent), 0644)
				if err != nil {
					return fmt.Errorf("写入repository接口文件失败: %w", err)
				}
				fmt.Printf("repository接口文件已生成: %s\n", repoInterfaceFileName)
			} else {
				fmt.Printf("repository接口文件已存在，不覆盖更新: %s\n", repoInterfaceFileName)
			}

			// 3. 生成扩展 repository (xxx.go) - 如果已存在则跳过
			repoExtContent, err := generateRepositoryExtFile(columns, modelName, cfg.Package, pathToPkg(cfg.OutPath), pathToPkg(cfg.ModelPkgPath), repoExtTmplPath)
			if err != nil {
				return fmt.Errorf("生成repository扩展内容失败: %w", err)
			}
			repoExtFileName := fmt.Sprintf("%s/%s.go", repoDir, strings.ToLower(modelName))
			if _, err := os.Stat(repoExtFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoExtFileName, []byte(repoExtContent), 0644)
				if err != nil {
					return fmt.Errorf("写入repository扩展文件失败: %w", err)
				}
				fmt.Printf("repository扩展文件已生成: %s\n", repoExtFileName)
			} else {
				fmt.Printf("repository扩展文件已存在，不覆盖更新: %s\n", repoExtFileName)
			}
		}
	}

	fmt.Println("生成完成!")
	return nil
}
