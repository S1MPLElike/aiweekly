package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Deepseek API 的响应结构
type DeepseekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepseekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Service struct {
	APIKey string
	Client *http.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		APIKey: apiKey,
		Client: &http.Client{
			Timeout: 60 * time.Second, // 设置超时 60秒
		},
	}
}

// GenerateDailyReport 生成日报
func (s *Service) GenerateDailyReport(records []WorkRecordInput) (*DailyReportOutput, error) {
	// 构建提示词 (Prompt)
	prompt := fmt.Sprintf(`你是一个专业的日报助手。请根据用户提供的工作记录，生成一份简洁专业的日报。

【重要要求】
1. 必须在日报中保留所有工作记录中的关键信息（如项目名称、具体任务、关键细节等）
2. 不要过度概括或省略具体内容
3. 确保每个工作记录的核心内容都体现在日报中

【工作记录】
`)

	for i, r := range records {
		prompt += fmt.Sprintf(`
%d. 标题：%s
   内容：%s
   时间：%d:00 - %d:00
`, i+1, r.Title, r.Content, r.StartHour, r.EndHour)
	}

	prompt += `

【输出要求】
- 请严格按照以下 JSON 格式输出，不要包含任何其他文字（不要有 Markdown 标记或说明）
- JSON 必须是有效的，不要有语法错误

{
	"title": "日报标题",
	"content": "日报内容"
}

【内容要求】
1. 标题要概括当天工作主题，不超过 30 字，要包含关键信息
2. 内容要对所有工作进行总结，清晰有条理，格式美观
3. 必须保留每个工作记录的关键细节和核心内容
4. 保持专业、正式的风格
`
	messages := []Message{
		{
			Role:    "system",
			Content: "你是一个专业的助手，擅长写工作总结和报告。你只输出 JSON，不输出其他任何内容。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.CallDeepseek(messages, 0.7)
	if err != nil {
		return nil, err
	}

	// 解析输出
	var output DailyReportOutput
	err = json.Unmarshal([]byte(response), &output)
	if err != nil {
		// 如果解析失败，把整个响应当作内容
		return &DailyReportOutput{
			Title:   "今日工作总结",
			Content: response,
		}, nil
	}

	return &output, nil
}

// GenerateWeeklyReport 生成周报
func (s *Service) GenerateWeeklyReport(dailyReports []DailyReportInput, style string) (*WeeklyReportOutput, error) {
	// 构建提示词 - 给领导看的周报格式
	prompt := fmt.Sprintf(`你是一个专业的周报助手。请根据用户提供的日报，生成一份给领导看的精简周报。

【重要要求】
1. 必须在周报中保留所有日报中的关键信息（如项目名称、具体任务、关键细节等）
2. 不要过度概括或省略具体内容
3. 确保每个日报的核心内容都体现在周报中
4. 重点突出成果和亮点

【日报内容】
`)

	for i, r := range dailyReports {
		prompt += fmt.Sprintf(`
%d. 日期：%s
   标题：%s
   内容：%s
`, i+1, r.Date, r.Title, r.Content)
	}

	prompt += `

【输出要求】
- 请严格按照以下 JSON 格式输出，不要包含任何其他文字（不要有 Markdown 标记或说明）
- JSON 必须是有效的，不要有语法错误

{
	"title": "周报标题",
	"summary": "本周工作概述",
	"content": "周报详细内容"
}

【内容格式要求】
请严格按照以下格式生成 content：

一、本周工作概述
（2-3句话，简要说明本周的主要工作方向和重点）

二、主要工作成果
- 具体项目/任务1 + 成果/进度
- 具体项目/任务2 + 成果/进度
- 具体项目/任务3 + 成果/进度
...

三、工作亮点
🌟 本周最有价值/最重要的突破（1-2个）

【其他要求】
1. 标题要概括本周工作主题，不超过 30 字，要包含关键信息
2. summary 就是"一、本周工作概述"的内容
3. content 要包含完整的三个部分（概述、主要成果、工作亮点）
4. 保持正式、严谨的风格，适合向领导汇报
5. 重点突出，精简但不丢关键细节
`

	messages := []Message{
		{
			Role:    "system",
			Content: "你是一个专业的助手，擅长写工作总结和报告。你只输出 JSON，不输出其他任何内容。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.CallDeepseek(messages, 0.7)
	if err != nil {
		return nil, err
	}

	// 解析输出
	var output WeeklyReportOutput
	err = json.Unmarshal([]byte(response), &output)
	if err != nil {
		// 如果解析失败，把整个响应当作内容
		return &WeeklyReportOutput{
			Title:   "本周工作总结",
			Summary: "本周工作已完成",
			Content: response,
		}, nil
	}

	return &output, nil
}

// CallDeepseek 调用 Deepseek API
func (s *Service) CallDeepseek(messages []Message, temperature float32) (string, error) {
	if s.APIKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY 未设置")
	}

	log.Printf("[LLM] 开始调用 Deepseek API，消息数量: %d", len(messages))

	reqBody := DeepseekRequest{
		Model:       "deepseek-v4-flash", // 使用 deepseek-v4-flash 模型
		Messages:    messages,
		Temperature: temperature,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %w", err)
	}

	log.Printf("[LLM] 请求 URL: https://api.deepseek.com/chat/completions")
	log.Printf("[LLM] 请求 Body: %s", string(reqBytes))

	httpReq, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.APIKey)

	httpResp, err := s.Client.Do(httpReq)
	if err != nil {
		log.Printf("[LLM] API 请求失败: %v", err)
		return "", fmt.Errorf("API 请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	log.Printf("[LLM] 响应状态码: %d", httpResp.StatusCode)

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("[LLM] 读取响应失败: %v", err)
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	log.Printf("[LLM] 响应内容: %s", string(respBytes))

	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("API 错误 (状态码 %d): %s", httpResp.StatusCode, string(respBytes))
	}

	var resp DeepseekResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		log.Printf("[LLM] 解析响应失败: %v", err)
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("API 响应无内容")
	}

	log.Printf("[LLM] API 调用成功，返回内容长度: %d", len(resp.Choices[0].Message.Content))
	return resp.Choices[0].Message.Content, nil
}

// ========= 输入输出类型 =========

type WorkRecordInput struct {
	Title     string
	Content   string
	StartHour int
	EndHour   int
}

type DailyReportInput struct {
	Date    string
	Title   string
	Content string
}

type DailyReportOutput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type WeeklyReportOutput struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}
