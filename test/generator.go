package main

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
	name = strings.Title(name)
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
	ModelName  string
	EntityName string
	Columns    []ColumnInfo
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

func generateApiFile(tableName string, columns []ColumnInfo, modelName string, db *gorm.DB) (string, error) {
	templatePath := "./template/api_template.txt"
	tmpl, err := loadTemplate(templatePath)
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
		EntityName:   Case2Camel(strings.ToUpper(tableName[:1])+tableName[1:]) + "Entity",
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

func generateRepositoryFile(columns []ColumnInfo, modelName string) (string, error) {
	templatePath := "./template/repository_template.txt"
	tmpl, err := loadTemplate(templatePath)
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
		ModelName:  modelName,
		EntityName: modelName + "Entity",
		Columns:    columnData,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
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

func main() {
	username := "esim-api"
	password := "R8S2xr2iR6eYLaXe"
	hostname := "43.138.220.19"
	port := 3306
	database := "esim-api"
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	tableName := readInput("请输入表名（直接回车同步所有表）: ")
	generateApi := readInput("是否生成go-zero api文件? (y/N): ")
	generateRepo := readInput("是否生成repository代码? (y/N): ")

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./query/dao",
		ModelPkgPath:      "./query/model",
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
			fmt.Printf("获取表结构失败: %v\n", err)
			return
		}
		modelName = Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:])
	}

	if strings.ToLower(generateRepo) == "y" && tableName != "" {
		repoDir := "./query/repository"
		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			os.MkdirAll(repoDir, 0755)
		}

		repoContent, err := generateRepositoryFile(columns, modelName)
		if err != nil {
			fmt.Printf("生成repository内容失败: %v\n", err)
		} else {
			repoFileName := fmt.Sprintf("%s/%s.go", repoDir, strings.ToLower(modelName))
			// 检查文件是否已存在
			if _, err := os.Stat(repoFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoFileName, []byte(repoContent), 0644)
				if err != nil {
					fmt.Printf("写入repository文件失败: %v\n", err)
				} else {
					fmt.Printf("repository文件已生成: %s\n", repoFileName)
				}
			} else {
				fmt.Printf("repository文件已存在，不覆盖更新: %s\n", repoFileName)
			}
		}
	} else if generateRepo == "y" {
		fmt.Println("批量模式需要在每个表中分别生成repository文件，请指定具体表名")
	}

	if strings.ToLower(generateApi) == "y" {
		apiDir := "./apps/admin/desc"
		if _, err := os.Stat(apiDir); os.IsNotExist(err) {
			os.MkdirAll(apiDir, 0755)
		}

		if tableName != "" {
			columns, err := getTableColumns(db, tableName)
			if err != nil {
				fmt.Printf("获取表结构失败: %v\n", err)
				return
			}
			modelName := Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:])
			apiContent, err := generateApiFile(tableName, columns, modelName, db)
			if err != nil {
				fmt.Printf("生成api内容失败: %v\n", err)
				return
			}
			apiFileName := fmt.Sprintf("%s/%s.api", apiDir, tableName)
			// 检查文件是否已存在
			if _, err := os.Stat(apiFileName); os.IsNotExist(err) {
				err = os.WriteFile(apiFileName, []byte(apiContent), 0644)
				if err != nil {
					fmt.Printf("写入api文件失败: %v\n", err)
					return
				}
				fmt.Printf("api文件已生成: %s\n", apiFileName)

				goctlPath := getGoctlPath()
				cmd := exec.Command(goctlPath, "api", "go", "-api", apiFileName, "--dir", "./apps/admin", "--style=goZero")
				cmd.Dir = "."
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("执行goctl失败: %v\n%s\n", err, output)
					return
				}
				fmt.Printf("go-zero代码生成成功\n%s\n", output)
			} else {
				fmt.Printf("api文件已存在，不覆盖更新: %s\n", apiFileName)
			}
		} else {
			fmt.Println("批量模式需要在每个表中分别生成api文件，请指定具体表名")
		}
	}

	if strings.ToLower(generateRepo) == "y" {
		repoDir := "./query/repository"
		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			os.MkdirAll(repoDir, 0755)
		}

		if tableName != "" {
			columns, err := getTableColumns(db, tableName)
			if err != nil {
				fmt.Printf("获取表结构失败: %v\n", err)
				return
			}
			modelName := Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:])
			repoContent, err := generateRepositoryFile(columns, modelName)
			if err != nil {
				fmt.Printf("生成repository内容失败: %v\n", err)
				return
			}
			repoFileName := fmt.Sprintf("%s/%s.go", repoDir, strings.ToLower(modelName))
			// 检查文件是否已存在
			if _, err := os.Stat(repoFileName); os.IsNotExist(err) {
				err = os.WriteFile(repoFileName, []byte(repoContent), 0644)
				if err != nil {
					fmt.Printf("写入repository文件失败: %v\n", err)
					return
				}
				fmt.Printf("repository文件已生成: %s\n", repoFileName)
			} else {
				fmt.Printf("repository文件已存在，不覆盖更新: %s\n", repoFileName)
			}
		} else {
			fmt.Println("批量模式需要在每个表中分别生成repository文件，请指定具体表名")
		}
	}

	fmt.Println("生成完成!")
}
