# AI测试用例生成API接口规范

## 1. 接口概述

### 1.1 基础信息
- **基础URL**: `/api/v1/testcase`
- **认证方式**: JWT Token
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`

### 1.2 通用响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2025-09-21T22:00:00Z",
  "request_id": "uuid"
}
```

### 1.3 错误码定义
| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 200 | 成功 | 200 |
| 400 | 请求参数错误 | 400 |
| 401 | 未授权 | 401 |
| 403 | 权限不足 | 403 |
| 404 | 资源不存在 | 404 |
| 429 | 请求频率超限 | 429 |
| 500 | 服务器内部错误 | 500 |
| 1001 | LLM服务不可用 | 503 |
| 1002 | 生成质量不达标 | 422 |
| 1003 | 提示词构建失败 | 400 |
| 1004 | 测试用例解析失败 | 422 |

## 2. 核心接口

### 2.1 生成测试用例

#### 接口信息
- **URL**: `POST /api/v1/testcase/generate`
- **描述**: 根据题目信息生成测试用例
- **权限**: 需要教师或管理员权限

#### 请求参数
```json
{
  "problem_id": "string",           // 题目ID（可选，用于关联）
  "problem_info": {                 // 题目信息
    "title": "string",              // 题目标题
    "description": "string",        // 题目描述
    "input_format": "string",       // 输入格式说明
    "output_format": "string",      // 输出格式说明
    "constraints": "string",        // 约束条件
    "difficulty": "easy|medium|hard", // 难度等级
    "tags": ["string"],             // 题目标签
    "examples": [                   // 示例用例
      {
        "input": "string",
        "output": "string",
        "explanation": "string"     // 可选：解释说明
      }
    ]
  },
  "generation_config": {            // 生成配置
    "count": 10,                    // 生成数量（1-50）
    "types": [                      // 测试用例类型
      "basic",                      // 基础测试用例
      "boundary",                   // 边界测试用例
      "extreme",                    // 极端测试用例
      "corner"                      // 角落测试用例
    ],
    "difficulty_range": {           // 难度范围
      "min": 1,                     // 最小难度（1-10）
      "max": 8                      // 最大难度（1-10）
    },
    "llm_provider": "qianfan",      // LLM提供商
    "model": "ERNIE-Bot-turbo",     // 模型名称（可选）
    "temperature": 0.7,             // 创造性参数（0-1）
    "max_tokens": 4000              // 最大token数
  },
  "validation_config": {            // 验证配置
    "enable_validation": true,      // 是否启用验证
    "enable_answer_check": true,    // 是否检查答案正确性
    "max_retries": 3,               // 最大重试次数
    "quality_threshold": 80         // 质量阈值（0-100）
  },
  "options": {                      // 其他选项
    "save_to_problem": false,       // 是否保存到题目
    "include_explanation": true,    // 是否包含解释
    "timeout": 30                   // 超时时间（秒）
  }
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "测试用例生成成功",
  "data": {
    "generation_id": "gen_123456789",
    "problem_id": "prob_123",
    "test_cases": [
      {
        "id": "tc_001",
        "input": "5 3",
        "expected_output": "8",
        "type": "basic",
        "description": "基础加法测试",
        "difficulty": 3,
        "tags": ["addition", "positive"],
        "explanation": "测试两个正整数的加法运算"
      },
      {
        "id": "tc_002", 
        "input": "0 0",
        "expected_output": "0",
        "type": "boundary",
        "description": "零值边界测试",
        "difficulty": 2,
        "tags": ["boundary", "zero"],
        "explanation": "测试零值输入的边界情况"
      }
    ],
    "quality_assessment": {
      "overall_quality": "high",     // low|medium|high|excellent
      "overall_score": 85,           // 0-100
      "scores": {
        "correctness": 90,           // 正确性评分
        "coverage": 85,              // 覆盖度评分
        "diversity": 80,             // 多样性评分
        "difficulty": 85,            // 难度合理性评分
        "discrimination": 88         // 区分度评分
      },
      "suggestions": [               // 改进建议
        "建议增加更多极端情况测试",
        "可以添加负数测试用例"
      ]
    },
    "coverage_analysis": {
      "has_basic_cases": true,
      "has_boundary_cases": true,
      "has_extreme_cases": false,
      "has_corner_cases": true,
      "coverage_score": 75,
      "missing_types": ["extreme"]
    },
    "statistics": {
      "total_count": 10,
      "type_distribution": {
        "basic": 4,
        "boundary": 3,
        "extreme": 0,
        "corner": 3
      },
      "difficulty_distribution": {
        "1": 1, "2": 2, "3": 3, "4": 2, "5": 2
      },
      "avg_difficulty": 3.2
    },
    "generation_info": {
      "generated_at": "2025-09-21T22:00:00Z",
      "generated_by": "user_123",
      "llm_provider": "qianfan",
      "model_used": "ERNIE-Bot-turbo",
      "prompt_tokens": 1200,
      "completion_tokens": 800,
      "total_tokens": 2000,
      "generation_time": 15.5,       // 生成耗时（秒）
      "retry_count": 0
    }
  }
}
```

### 2.2 验证测试用例

#### 接口信息
- **URL**: `POST /api/v1/testcase/validate`
- **描述**: 验证测试用例的正确性和质量
- **权限**: 需要教师或管理员权限

#### 请求参数
```json
{
  "problem_id": "string",           // 题目ID
  "test_cases": [                   // 待验证的测试用例
    {
      "input": "string",
      "expected_output": "string",
      "type": "basic|boundary|extreme|corner"
    }
  ],
  "validation_options": {
    "check_syntax": true,           // 检查语法
    "check_logic": true,            // 检查逻辑
    "check_answer": true,           // 检查答案正确性
    "check_constraints": true       // 检查约束条件
  }
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "验证完成",
  "data": {
    "validation_id": "val_123456789",
    "overall_result": {
      "is_valid": true,
      "overall_score": 92,
      "passed_count": 9,
      "failed_count": 1,
      "warning_count": 2
    },
    "test_case_results": [
      {
        "index": 0,
        "is_valid": true,
        "score": 95,
        "validation_details": {
          "syntax_check": {
            "passed": true,
            "errors": []
          },
          "logic_check": {
            "passed": true,
            "warnings": ["输入数据较简单，建议增加复杂度"]
          },
          "answer_check": {
            "passed": true,
            "execution_time": 0.001,
            "memory_usage": 1024
          },
          "constraint_check": {
            "passed": true,
            "violations": []
          }
        }
      }
    ],
    "suggestions": [
      "测试用例整体质量良好",
      "建议增加更多边界情况测试"
    ]
  }
}
```

### 2.3 查询生成历史

#### 接口信息
- **URL**: `GET /api/v1/testcase/history`
- **描述**: 查询测试用例生成历史记录
- **权限**: 需要登录

#### 请求参数
```
?problem_id=string          // 题目ID（可选）
&user_id=string            // 用户ID（可选，管理员可查看所有）
&start_date=2025-09-01     // 开始日期（可选）
&end_date=2025-09-30       // 结束日期（可选）
&quality=high              // 质量等级过滤（可选）
&llm_provider=qianfan      // LLM提供商过滤（可选）
&page=1                    // 页码（默认1）
&size=20                   // 每页大小（默认20，最大100）
&sort=created_at           // 排序字段
&order=desc                // 排序方向（asc|desc）
```

#### 响应示例
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "total": 156,
    "page": 1,
    "size": 20,
    "total_pages": 8,
    "items": [
      {
        "generation_id": "gen_123456789",
        "problem_id": "prob_123",
        "problem_title": "两数之和",
        "test_case_count": 10,
        "quality": "high",
        "overall_score": 85,
        "llm_provider": "qianfan",
        "generated_at": "2025-09-21T22:00:00Z",
        "generated_by": "user_123",
        "generation_time": 15.5,
        "status": "completed"
      }
    ]
  }
}
```

### 2.4 获取生成详情

#### 接口信息
- **URL**: `GET /api/v1/testcase/generation/{generation_id}`
- **描述**: 获取特定生成记录的详细信息
- **权限**: 需要登录，只能查看自己的记录或管理员权限

#### 响应示例
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "generation_id": "gen_123456789",
    "problem_id": "prob_123",
    "problem_info": {
      "title": "两数之和",
      "description": "给定一个整数数组...",
      "difficulty": "easy"
    },
    "test_cases": [
      // 完整的测试用例列表
    ],
    "quality_assessment": {
      // 质量评估详情
    },
    "generation_config": {
      // 生成时使用的配置
    },
    "generation_info": {
      // 生成信息
    },
    "validation_results": {
      // 验证结果（如果进行了验证）
    }
  }
}
```

### 2.5 导出测试用例

#### 接口信息
- **URL**: `GET /api/v1/testcase/export/{generation_id}`
- **描述**: 导出测试用例为指定格式
- **权限**: 需要登录

#### 请求参数
```
?format=json               // 导出格式（json|csv|txt）
&include_metadata=true     // 是否包含元数据
```

#### 响应示例
```json
{
  "code": 200,
  "message": "导出成功",
  "data": {
    "download_url": "https://example.com/download/testcases_gen_123456789.json",
    "file_name": "testcases_gen_123456789.json",
    "file_size": 2048,
    "expires_at": "2025-09-22T22:00:00Z"
  }
}
```

## 3. 管理接口

### 3.1 获取LLM提供商列表

#### 接口信息
- **URL**: `GET /api/v1/testcase/providers`
- **描述**: 获取可用的LLM提供商和模型列表
- **权限**: 需要登录

#### 响应示例
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "providers": [
      {
        "name": "qianfan",
        "display_name": "阿里云千炼",
        "status": "available",
        "models": [
          {
            "name": "ERNIE-Bot-turbo",
            "display_name": "文心一言Turbo",
            "max_tokens": 4096,
            "supports_json": true
          }
        ]
      },
      {
        "name": "openai",
        "display_name": "OpenAI",
        "status": "available",
        "models": [
          {
            "name": "gpt-3.5-turbo",
            "display_name": "GPT-3.5 Turbo",
            "max_tokens": 4096,
            "supports_json": true
          }
        ]
      }
    ]
  }
}
```

### 3.2 获取使用统计

#### 接口信息
- **URL**: `GET /api/v1/testcase/statistics`
- **描述**: 获取测试用例生成的使用统计
- **权限**: 需要管理员权限

#### 请求参数
```
?start_date=2025-09-01     // 开始日期
&end_date=2025-09-30       // 结束日期
&group_by=day              // 分组方式（day|week|month）
```

#### 响应示例
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "summary": {
      "total_generations": 1250,
      "total_test_cases": 12500,
      "avg_quality_score": 82.5,
      "success_rate": 0.95,
      "total_tokens": 1500000,
      "total_cost": 150.00
    },
    "daily_stats": [
      {
        "date": "2025-09-21",
        "generations": 45,
        "test_cases": 450,
        "avg_quality": 85.2,
        "success_rate": 0.96,
        "tokens": 54000,
        "cost": 5.40
      }
    ],
    "provider_stats": {
      "qianfan": {
        "generations": 800,
        "success_rate": 0.96,
        "avg_quality": 83.1
      },
      "openai": {
        "generations": 450,
        "success_rate": 0.94,
        "avg_quality": 81.8
      }
    }
  }
}
```

## 4. 错误处理

### 4.1 错误响应格式
```json
{
  "code": 400,
  "message": "请求参数错误",
  "error": {
    "type": "validation_error",
    "details": [
      {
        "field": "problem_info.description",
        "message": "题目描述不能为空"
      }
    ]
  },
  "timestamp": "2025-09-21T22:00:00Z",
  "request_id": "req_123456789"
}
```

### 4.2 常见错误场景
1. **参数验证错误**: 必填字段缺失、格式不正确
2. **权限错误**: 无权限访问或操作
3. **资源不存在**: 题目ID不存在、生成记录不存在
4. **服务不可用**: LLM服务异常、网络超时
5. **配额限制**: 超出使用限制、频率限制

## 5. 接口调用示例

### 5.1 JavaScript示例
```javascript
// 生成测试用例
const generateTestCases = async (problemInfo) => {
  const response = await fetch('/api/v1/testcase/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      problem_info: problemInfo,
      generation_config: {
        count: 10,
        types: ['basic', 'boundary', 'extreme'],
        llm_provider: 'qianfan'
      },
      validation_config: {
        enable_validation: true,
        quality_threshold: 80
      }
    })
  });
  
  const result = await response.json();
  return result;
};
```

### 5.2 Python示例
```python
import requests

def generate_test_cases(problem_info, token):
    url = '/api/v1/testcase/generate'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }
    data = {
        'problem_info': problem_info,
        'generation_config': {
            'count': 10,
            'types': ['basic', 'boundary', 'extreme'],
            'llm_provider': 'qianfan'
        },
        'validation_config': {
            'enable_validation': True,
            'quality_threshold': 80
        }
    }
    
    response = requests.post(url, json=data, headers=headers)
    return response.json()
```

这个API接口规范提供了完整的接口定义，包括请求参数、响应格式、错误处理和使用示例，为前端开发和接口对接提供了详细的参考。
