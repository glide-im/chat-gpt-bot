package chat_gpt

import "time"

type ChatContextMessage struct {
	Speaker string `json:"speaker,omitempty"`
	Text    string `json:"text,omitempty"`
}

type ChatContext struct {
	History []*ChatContextMessage `json:"history"`
}

// ChatCompletionsV2Request 表示向/chat/completions端点发送请求时需要提供的参数
type ChatCompletionsV2Request struct {
	Engine           string      `json:"engine"`                      // 选择使用的模型引擎，可以是OpenAI之前发布的任何引擎，或者是用户在OpenAI平台上训练的私有引擎。
	Prompt           string      `json:"prompt"`                      // 要补全的文本字符串。
	Temperature      float64     `json:"temperature"`                 // 生成文本的创新度。值越高，生成结果越随机和创新。
	MaxTokens        int         `json:"max_tokens"`                  // 生成文本的最大长度（以标记数计算）。这个值必须介于1和2048之间。
	TopP             float64     `json:"top_p,omitempty"`             // 一种限制机制，可以控制生成结果的多样性。
	N                int         `json:"n,omitempty"`                 // 指定要生成的完整响应的数量。
	Stop             []string    `json:"stop,omitempty"`              // 一个可选的字符串数组，用于指定生成文本应该在哪个位置停止。
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`  // 对出现在先前文本中的词汇的惩罚因子。
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"` // 对出现频率较高的词汇惩罚的因子。
	BestOf           int         `json:"best_of,omitempty"`           // 指定要返回的最佳响应数量。
	Logprobs         int         `json:"logprobs,omitempty"`          // 一个布尔值，指定是否应该为每个生成的标记返回log-probability分数。
	Echo             bool        `json:"echo,omitempty"`              // 一个可选的布尔值，指示模型是否应该在响应中包含输入文本。
	Stream           bool        `json:"stream,omitempty"`            // 一个可选的布尔值，用于指示API是否应该使用流式传输返回数据。
	StopSequence     string      `json:"stop_sequence,omitempty"`     // 一个可选的字符串，指示模型应该在哪个序列处停止生成文本。
	MaxCompletions   int         `json:"max_completions,omitempty"`   // 一个可选的整数，指定API应该为每个请求生成的最大响应数。默认值为1。
	Model            string      `json:"model,omitempty"`             // 选择使用的模型（如果您未指定引擎，则可以将其用作选择模型的快捷方式）。
	Context          interface{} `json:"context,omitempty"`           // 包含先前对话消息列表以及与当前会话相关的任何其他信息的JSON对象。
}

// OpenAIResponse 表示/chat/completions端点返回的响应
type OpenAIResponse struct {
	ID      string         `json:"id"`      // 生成的响应标识符。如果您指定了`n`或`best_of`参数，则会返回多个ID。
	Created time.Time      `json:"created"` // 响应生成的时间戳。
	Model   string         `json:"model"`   // 使用的模型名称。
	Choices []OpenAIChoice `json:"choices"` // 包含生成响应的列表。
}

// OpenAIChoice 表示生成的响应的一个选择
type OpenAIChoice struct {
	Text         string          `json:"text"`               // 生成的文本字符串。
	Index        int             `json:"index"`              // 生成结果的索引值。
	Logprobs     *OpenAILogprobs `json:"logprobs,omitempty"` // （可选）包含每个标记的对数概率（如果请求中指定了`logprobs`参数）。
	FinishReason string          `json:"finish_reason"`      // 标识生成过程结束的原因。可能的值包括`stop`、`max_tokens`和`restart`。
	Prompt       *OpenAIPrompt   `json:"prompt,omitempty"`   // （可选）包含原始提示文本的对象（如果请求中指定了`echo`参数）。
}

// OpenAILogprobs 表示每个标记的对数概率
type OpenAILogprobs struct {
	Tokens []OpenAILogprobsToken `json:"tokens"`
}

// OpenAILogprobsToken 表示单个标记的对数概率信息
type OpenAILogprobsToken struct {
	Token         string      `json:"token"`                    // 标记字符串
	Logprob       float64     `json:"logprob"`                  // 标记的对数概率
	TokenLogprobs interface{} `json:"token_logprobs,omitempty"` // （可选）包含每个上下文中此标记的对数概率
	TopLogprobs   interface{} `json:"top_logprobs,omitempty"`   // （可选）包含最高对数概率的其他标记和它们的分数
}

// OpenAIPrompt 包含原始提示的对象
type OpenAIPrompt struct {
	Text string `json:"text"` // 原始提示文本
}

type ChatGPTRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
}

type ChatGPTReply struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs"`
}

type ChatGPTUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type ChatGPTResponse struct {
	Id      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Choices []*ChatGPTReply `json:"choices"`
	Usage   *ChatGPTUsage   `json:"usage"`
}
