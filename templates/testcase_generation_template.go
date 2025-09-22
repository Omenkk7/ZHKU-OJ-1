package templates

import (
	"fmt"
	"strings"
)

// ProblemForAI AI生成测试用例需要的题目信息
type ProblemForAI struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	InputExample  string   `json:"input_example"`
	OutputExample string   `json:"output_example"`
	Difficulty    string   `json:"difficulty"`
	Labels        []string `json:"labels"`
}

// TestCaseConfig 测试用例生成配置
type TestCaseConfig struct {
	TotalCount      int `json:"total_count"`      // 总用例数
	BasicCount      int `json:"basic_count"`      // 基础功能用例数
	BoundaryCount   int `json:"boundary_count"`   // 边界用例数
	PerformanceCount int `json:"performance_count"` // 性能用例数
	TargetedCount   int `json:"targeted_count"`   // 针对性用例数
}

// DefaultTestCaseConfig 默认测试用例配置
var DefaultTestCaseConfig = TestCaseConfig{
	TotalCount:       10,
	BasicCount:       3,
	BoundaryCount:    4,
	PerformanceCount: 2,
	TargetedCount:    1,
}

// TestCaseGenerationTemplate 测试用例生成的完整模板
const TestCaseGenerationTemplate = `请为以下编程题目生成测试用例：

【题目标题】
%s

【难度等级】
%s

【题目标签】
%s

【题目描述】
%s

【输入示例】
%s

【输出示例】
%s

## 📌 测试用例设计策略（基于6大原则）：

### 1. 基础功能用例 (2-3个)
- **目的**：验证代码能解决最基本的问题
- **特点**：小规模、典型输入、预期结果明确
- **包含**：样例用例的变体、标准正常输入

### 2. 边界用例 (3-4个)
- **数值边界**：最大值、最小值、0、负数、溢出边界
- **数据结构边界**：空输入、单元素、两元素、全相同元素
- **逻辑边界**：无解情况、多解情况、特殊条件
- **约束边界**：题目约束条件的临界值

### 3. 性能用例 (2-3个)
- **目的**：淘汰时间/空间复杂度过高的解法
- **特点**：最大规模数据、最坏情况构造
- **考察**：算法效率的关键用例

### 4. 针对性用例 (2-3个)
- **目的**：针对常见错误解法设计"陷阱"
- **特点**：专门构造的反例数据
- **针对**：可能的错误算法和边界处理错误

## 📋 输入输出格式要求：
请严格按照题目描述的输入输出格式生成测试用例：
- 输入格式必须完全符合题目要求
- 输出格式必须精确匹配题目规范
- 注意空格、换行符、特殊字符的处理
- 确保数据范围在题目约束内

## 🧩 输出格式要求：
请以 **纯JSON格式** 返回结果，包含一个字段 "test_cases"，它是一个列表，列表中每个元素是一个测试用例对象，格式如下：

{
  "test_cases": [
    {
      "name": "用例类型-具体描述（如：基础功能用例-标准输入）",
      "category": "用例分类（basic/boundary/performance/targeted）",
      "is_hidden": false,
      "input": "输入数据（严格按照题目输入格式，使用\\n表示换行）",
      "expected": "期望输出（严格按照题目输出格式，使用\\n表示换行）",
      "time_complexity": "时间复杂度（如：O(n)）",
      "space_complexity": "空间复杂度（如：O(1)）",
      "cpu_limit": CPU时间限制纳秒数,
      "memory_limit": 内存限制字节数,
      "description": "用例设计意图和考察点"
    }
  ]
}

## 📊 复杂度和资源限制参考：

### 时间复杂度常见格式：
- O(1) - 常数时间
- O(log n) - 对数时间
- O(n) - 线性时间
- O(n log n) - 线性对数时间
- O(n²) - 平方时间
- O(2^n) - 指数时间

### 空间复杂度常见格式：
- O(1) - 常数空间
- O(log n) - 对数空间
- O(n) - 线性空间
- O(n²) - 平方空间

### 资源限制参考值：
- **基础用例CPU限制**：1000000000 纳秒 (1秒)
- **性能用例CPU限制**：2000000000-5000000000 纳秒 (2-5秒)
- **基础用例内存限制**：134217728 字节 (128MB)
- **性能用例内存限制**：268435456-536870912 字节 (256-512MB)

## ⚠️ 重要约束：
1. **严禁编程语法**：绝对不能使用任何编程语言的语法，包括但不限于：
   - 字符串拼接：+、concat()
   - 循环语法：.repeat()、for、while
   - 数组方法：Array.from()、range()、join()
   - 变量赋值：=、let、var
   - 函数调用：任何形式的函数调用
2. **纯JSON输出**：不要包含markdown标记（json）或解释文字
3. **数据直接书写**：所有测试数据必须完整直接写出，不能使用任何代码生成
4. **数据一致性**：确保输入输出格式严格符合题目要求
5. **复杂度准确性**：时间/空间复杂度必须准确反映算法特性
6. **资源合理性**：CPU和内存限制必须是正整数，符合实际判题需求
7. **用例完整性**：每个测试用例必须包含所有必需字段
8. **隐藏策略**：性能用例和针对性用例建议设为隐藏(is_hidden: true)
9. **换行符处理**：所有字符串中的换行符使用 \\n 表示
10. **大规模数据**：如需大规模测试数据，必须手动完整写出，不得使用任何简化语法

## 🎯 质量标准：
- **正确性**：所有预期输出必须100%%正确
- **全面性**：覆盖所有可能的输入场景和边界情况
- **针对性**：能够有效区分正确解法和常见错误解法
- **渐进性**：从简单到复杂，从小规模到大规模
- **有效性**：资源限制合理，不影响判题效率
- **隐蔽性**：隐藏用例逻辑对用户保密，防止硬编码

## 📝 特别注意：
- 根据题目特点调整测试用例的具体内容
- 考虑题目的算法类型（排序、搜索、动态规划、图论等）
- 针对题目的约束条件设计边界用例
- 确保测试用例能够验证算法的正确性和效率

【生成要求】
请严格按照以上要求生成%d个高质量测试用例，确保能够全面、准确地评估算法实现的正确性和效率。
生成的测试用例分布：
- 基础功能用例：%d个
- 边界用例：%d个  
- 性能用例：%d个
- 针对性用例：%d个

请直接返回纯JSON格式，不要包含任何markdown标记或解释文字和函数代码。`

// BuildTestCasePrompt 构建测试用例生成的prompt
func BuildTestCasePrompt(problem *ProblemForAI, config *TestCaseConfig) string {
	// 使用默认配置如果没有提供
	if config == nil {
		config = &DefaultTestCaseConfig
	}

	// 处理标签
	labels := "无"
	if len(problem.Labels) > 0 {
		labels = strings.Join(problem.Labels, ", ")
	}

	// 处理难度
	difficulty := problem.Difficulty
	if difficulty == "" {
		difficulty = "未指定"
	}

	return fmt.Sprintf(TestCaseGenerationTemplate,
		problem.Title,
		difficulty,
		labels,
		problem.Description,
		problem.InputExample,
		problem.OutputExample,
		config.TotalCount,
		config.BasicCount,
		config.BoundaryCount,
		config.PerformanceCount,
		config.TargetedCount,
	)
}

// BuildTestCasePromptWithDefaults 使用默认配置构建prompt
func BuildTestCasePromptWithDefaults(problem *ProblemForAI) string {
	return BuildTestCasePrompt(problem, &DefaultTestCaseConfig)
}

// ValidateTestCaseConfig 验证测试用例配置
func ValidateTestCaseConfig(config *TestCaseConfig) error {
	if config.TotalCount <= 0 {
		return fmt.Errorf("总用例数必须大于0")
	}

	sum := config.BasicCount + config.BoundaryCount + config.PerformanceCount + config.TargetedCount
	if sum != config.TotalCount {
		return fmt.Errorf("各类型用例数之和(%d)必须等于总用例数(%d)", sum, config.TotalCount)
	}

	if config.BasicCount < 1 {
		return fmt.Errorf("基础功能用例数至少为1")
	}

	if config.BoundaryCount < 1 {
		return fmt.Errorf("边界用例数至少为1")
	}

	return nil
}
