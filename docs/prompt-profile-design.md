# AI 解读人格与 PromptProfile 设计文档

日期：2026-06-26

## 0. 一句话定义

把 AI 解读从“写死的提示词”升级为“可配置、可版本化、可复现的解读人格系统”：第一版只把梅花易数打磨成「康节先生」人格，并让直接解读与聊天解读共用同一个 Compose 管道。

## 0.1 乔布斯式审查结论

这份设计的方向是对的，但第一版必须更聚焦。

不要一上来迁移所有占卜方式，也不要做人格选择器。那会把一个本该锋利的产品能力，做成一个设置面板。第一版只证明一件事：**同一条梅花易数结果，在直接 AI 解读和聊天模式里，都由同一个「康节先生」人格、同一个 Compose 函数生成稳定 prompt，并且历史记录能追溯 profile 版本。**

做成这一件事后，再迁移高岛、塔罗、八字、星座。否则工程会分散，用户却感知不到更好的体验。

## 0.2 第一版开发范围裁决

第一版必须做：

- 新增 `prompt_profiles.yaml` 与加载方法。
- 新增纯函数 `Compose`。
- 新增梅花易数 `parseMeihuaFacts`。
- 梅花易数直接 AI 解读与聊天初始解读共用 Compose。
- `DivinationRecord` 保存 profile ID、名称、版本。
- 前端展示后端返回的 profile 名称。
- 旧记录无 profile 字段时 fallback 到默认 profile。

第一版不做：

- 不迁移所有占卜类型。
- 不开放用户选择人格。
- 不做后台编辑。
- 不做 A/B 测试。
- 不把 profile 存数据库。
- 不在前端新增人格选择 UI。
- 不保留梅花易数的第二套 prompt 构建真相源。

## 1. 背景

当前 AI 解读提示词主要写在 `server/internal/service/prompts/*.txt`，并由不同路径拼装：

- `AIService.Interpret` 使用 `text/template` 执行模板。
- `ChatService.buildSystemPrompt` 对通用、六爻、高岛、梅花分别用 `strings.ReplaceAll` 手动替换。
- 前端 `client/src/components/chat/divinationPersona.ts` 维护一套展示人格，如 icon、name、subtitle、welcome 文案。

这导致三个问题：

- 提示词的推理规则、表达风格、占卜结果数据混在一起。
- 直接解读和聊天解读走两套 prompt 构建逻辑。
- 前端展示人格和后端实际 prompt 人格可能逐渐不一致。

本设计目标是把 AI 解读人格做成可配置、可版本化、可复现的系统，同时保持第一版产品体验简单：每种占卜方式默认绑定一个大师人格。

## 2. 设计目标

1. **数据与控制分离**
   占卜结果是数据，PromptProfile 是控制指令。卦象、牌面、八字、星座运势等结果不能被人格配置改写。

2. **Compose 纯函数化**
   Prompt 生成必须是纯函数：

   ```go
   func Compose(profile *ProfileConfig, input InterpretationInput) (*Messages, error)
   ```

   它不读数据库、不访问全局状态、不调用 AI、不写日志、不保存记录。指针传递避免大结构体拷贝，返回指针便于调用方判断 nil。

3. **推理框架与表达风格分离**
   `ReasoningFramework` 决定 AI 如何组织推理链。
   `VoiceStyle` 决定 AI 如何表达。
   二者不能混在一个 `interpretation_rules` 字段里。

4. **双路径统一**
   `AIService.Interpret` 和 `ChatService.buildSystemPrompt` 都必须走同一个 Compose 函数。

5. **人格版本可复现**
   每条 AI 解读记录保存 `prompt_profile_id` 和 `prompt_profile_version`，便于以后追溯同一解读当时使用的提示词版本。

6. **前后端人格不漂移**
   第一版采用务实方案：后端在结果或会话响应里返回人格展示字段，前端只展示，不自行推导后端实际人格。

7. **配置体系一致**
   PromptProfile 配置与现有 `config.yaml` 使用相同的 YAML 格式 + viper 加载机制，不引入额外的 JSON 嵌入方案。

## 3. 非目标

第一版不做以下内容：

- 不做用户自由创建人格。
- 不做后台编辑 PromptProfile。
- 不做 A/B 测试。
- 不把 PromptProfile 存数据库。
- 不让用户在聊天页看到复杂的人格设置面板。

这些能力可以在后续版本基于同一 Compose 签名扩展。

## 4. 核心概念

### 4.1 ProfileConfig

`ProfileConfig` 是 AI 解读人格的配置，也是"可存储的程序"。它描述某种占卜类型下，AI 应该以什么推理框架和表达风格解释数据。

建议结构：

```go
// ProfileConfig 定义单个 AI 解读人格。
type ProfileConfig struct {
    Version            string   `mapstructure:"version"`
    DivinationType     string   `mapstructure:"divination_type"`
    Name               string   `mapstructure:"name"`
    Title              string   `mapstructure:"title"`
    Subtitle           string   `mapstructure:"subtitle"`
    Icon               string   `mapstructure:"icon"`
    Description        string   `mapstructure:"description"`
    Enabled            bool     `mapstructure:"enabled"`
    SystemIdentity     string   `mapstructure:"system_identity"`
    ReasoningFramework []string `mapstructure:"reasoning_framework"`
    VoiceStyle         []string `mapstructure:"voice_style"`
    OutputStructure    []string `mapstructure:"output_structure"`
    Guardrails         []string `mapstructure:"guardrails"`
}
```

注意：`ID` 不放在结构体内部，而是作为外层 map 的 key（见第 5.2 节 YAML 结构），与 viper 的 `mapstructure` 反序列化模式一致。

字段说明：

- `ID`：稳定标识，如 `shaoyong_meihua`。
- `Version`：版本号，如 `v1`。配置变化必须升级版本。
- `DivinationType`：适用占卜类型，如 `meihua`。
- `Name/Title/Subtitle/Icon/Description`：给前端展示使用。
- `SystemIdentity`：身份设定，不包含具体卦象数据。
- `ReasoningFramework`：推理步骤，是控制流核心。
- `VoiceStyle`：表达风格，只影响语言，不改变推理顺序。
- `OutputStructure`：要求 AI 输出的段落结构。
- `Guardrails`：禁止编造、必须基于结果、偏题引导等规则。
- `Enabled`：是否启用。

### 4.2 InterpretationInput

`InterpretationInput` 是一次解读的输入数据。不使用 `Context` 命名，避免与标准库 `context.Context` 混淆。

```go
// InterpretationInput 是单次占卜解读的输入数据。
type InterpretationInput struct {
    DivinationType string            `json:"divination_type"`
    Question       string            `json:"question"`
    ResultJSON     string            `json:"result_json"`
    ResultFacts    map[string]string `json:"result_facts"`
    Mode           string            `json:"mode"` // "direct" or "chat"
}
```

`ResultFacts` 由各占卜类型的解析器生成。例如梅花易数：

```go
map[string]string{
    "Method":          "时间起卦（丙午年五月初九子时）",
    "BenGua":          "泽火革（兑上离下）",
    "HuGua":           "天风姤（乾上巽下）",
    "BianGua":         "泽山咸（兑上艮下）",
    "MovingLine":      "初爻动",
    "Ti":              "兑（金）",
    "Yong":            "离（火）",
    "TiYongRelation":  "用克体",
}
```

### 4.3 Messages

Compose 输出统一消息结构，避免有些路径把全部内容塞进 user prompt，有些路径塞进 system prompt。类型名不加 `Prompt` 前缀，因为所在包已经提供上下文。

```go
// Messages 是 Compose 的输出，包含 system 和 user 两条消息。
type Messages struct {
    System string `json:"system"`
    User   string `json:"user"`
}
```

AI 请求统一使用：

```go
msgs := []map[string]string{
    {"role": "system", "content": messages.System},
    {"role": "user", "content": messages.User},
}
```

## 5. 配置文件设计

### 5.1 文件位置与格式

参照 `server/config/config.yaml` 的设计模式，新增一个 YAML 配置文件：

```
server/config/prompt_profiles.yaml
```

与 `config.yaml` 同级、同格式、同加载机制。

### 5.2 配置文件结构

```yaml
# server/config/prompt_profiles.yaml
# AI 解读人格配置

prompt_profiles:
  # ==================== 梅花易数 ====================
  shaoyong_meihua:
    version: "v1"
    divination_type: "meihua"
    name: "康节先生"
    title: "梅花易数解读者"
    subtitle: "师承邵雍一脉，观物取象"
    icon: "🌸"
    description: "以邵雍梅花易数体系为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「康节先生」人格，一位以邵雍梅花易数体系为核心的 AI 解读者。
      你不自称历史上的邵雍本人，而是以邵雍一脉的观物取象、体用生克、数理成卦方法进行判断。
    reasoning_framework:
      - "先以体用生克定大方向（用生体吉、体生用泄、用克体凶、体克用可控、比和稳）"
      - "再看本卦现状：上下卦象征与卦义"
      - "再看互卦过程：内在因素与隐含趋势"
      - "再看动爻与变卦：变化焦点与最终走向"
      - "如有时间判断，参考卦数推断应期"
    voice_style:
      - "文言与白话交融，沉稳通透"
      - "善用比喻，不空泛玄谈"
      - "语言现代、清楚、直接"
    output_structure:
      - "卦象总览"
      - "本卦解读"
      - "互卦分析"
      - "动爻与变卦"
      - "体用生克"
      - "结论与建议"
    guardrails:
      - "必须基于起卦结果解读，不得编造卦象信息"
      - "不得自称历史上的邵雍本人"
      - "用户问与占卜无关问题时，礼貌引导回占卜话题"

  # ==================== 高岛易断 ====================
  takashima_ekidan:
    version: "v1"
    divination_type: "liuyao_v2"
    name: "高岛吞象"
    title: "高岛易断解读者"
    subtitle: "原典证据，卦辞爻辞"
    icon: "📖"
    description: "以《高岛易断》原典为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「高岛吞象」人格，精通《高岛易断》的卦师。
      你必须通过书中证据解卦，不得把没有检索到的内容说成高岛原文。
    reasoning_framework:
      - "先依据书中证据，再结合用户问题"
      - "引用时标注页码，例如'据第 58 页'"
      - "动爻优先；无动爻时以本卦卦辞、象义和断法原则为主"
      - "证据不足时明确说'书中未见直接说明'，不要补造"
    voice_style:
      - "严谨考据，有理有据"
      - "语言现代、清楚、直接"
    output_structure:
      - "卦象总判断"
      - "本卦卦辞说明"
      - "动爻重点解读"
      - "变卦趋势"
      - "对用户问题的结论"
      - "可执行建议"
    guardrails:
      - "必须基于书中证据解卦"
      - "不得编造高岛原文"
      - "证据不足时必须明说"

  # ==================== 六爻 ====================
  xuanji_liuyao:
    version: "v1"
    divination_type: "liuyao"
    name: "玄机卦师"
    title: "六爻解读者"
    subtitle: "本卦动爻，世应用神"
    icon: "☯️"
    description: "以传统六爻体系为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「玄机卦师」，精通传统六爻占卜的卦师。
      你基于本卦、动爻、变卦、世应用神体系进行判断。
    reasoning_framework:
      - "先看本卦卦象与卦辞"
      - "再看动爻位置与爻辞"
      - "分析世爻、应爻、用神关系"
      - "看变卦趋势"
      - "综合判断吉凶"
    voice_style:
      - "沉稳专业"
      - "语言现代、清楚、直接"
    output_structure:
      - "卦象概述"
      - "动爻分析"
      - "世应用神"
      - "变卦趋势"
      - "结论与建议"
    guardrails:
      - "必须基于起卦结果解读"
      - "不得编造卦象信息"

  # ==================== 八字 ====================
  sizhu_mingjian:
    version: "v1"
    divination_type: "bazi"
    name: "四柱明鉴师"
    title: "八字解读者"
    subtitle: "四柱五行，日主格局"
    icon: "🏛️"
    description: "以四柱八字体系为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「四柱明鉴师」，精通四柱八字命理的解读者。
      你基于四柱、五行、日主、格局体系进行判断。
    reasoning_framework:
      - "先排四柱，确定日主"
      - "分析五行旺衰"
      - "看格局与用神"
      - "结合大运流年"
      - "综合判断"
    voice_style:
      - "温和睿智"
      - "语言现代、清楚、直接"
    output_structure:
      - "四柱排盘"
      - "五行分析"
      - "格局判断"
      - "综合解读"
      - "建议"
    guardrails:
      - "必须基于四柱数据解读"
      - "不得编造命理信息"

  # ==================== 塔罗 ====================
  star_tarot:
    version: "v1"
    divination_type: "tarot"
    name: "星牌解语师"
    title: "塔罗解读者"
    subtitle: "牌面象征，正逆位解读"
    icon: "🔮"
    description: "以塔罗牌体系为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「星牌解语师」，精通塔罗牌解读的占卜师。
      你基于牌面、正逆位、象征关系进行判断。
    reasoning_framework:
      - "先看每张牌的正逆位与核心含义"
      - "分析牌与牌之间的象征关系"
      - "结合牌位含义（如三张牌的过去/现在/未来）"
      - "综合判断整体信息"
    voice_style:
      - "温暖直觉"
      - "善用象征性语言"
    output_structure:
      - "牌面概述"
      - "逐张解读"
      - "牌面关系"
      - "综合信息"
      - "建议"
    guardrails:
      - "必须基于牌面结果解读"
      - "不得编造牌面信息"

  # ==================== 星座 ====================
  xinggui_horoscope:
    version: "v1"
    divination_type: "horoscope"
    name: "星轨知微师"
    title: "星座解读者"
    subtitle: "星座能量，近期节奏"
    icon: "⭐"
    description: "以星座运势体系为核心的 AI 解读者"
    enabled: true
    system_identity: |
      你是「星轨知微师」，精通星座运势解读的占卜师。
      你基于星座能量、近期节奏、状态建议进行判断。
    reasoning_framework:
      - "先看星座整体能量状态"
      - "分析近期节奏与趋势"
      - "给出状态调整建议"
    voice_style:
      - "轻松亲切"
      - "积极引导"
    output_structure:
      - "星座能量"
      - "近期节奏"
      - "状态建议"
    guardrails:
      - "必须基于星座运势数据解读"
      - "不得编造运势信息"

# ==================== 占卜类型与默认人格绑定 ====================
default_bindings:
  meihua: "shaoyong_meihua"
  liuyao_v2: "takashima_ekidan"
  liuyao: "xuanji_liuyao"
  bazi: "sizhu_mingjian"
  tarot: "star_tarot"
  horoscope: "xinggui_horoscope"
```

### 5.3 配置结构体

在 `server/config/` 下新增 `prompt_profiles.go`，与 `config.go` 分文件组织：

```go
package config

// ProfileConfig 定义单个 AI 解读人格。
type ProfileConfig struct {
    Version            string   `mapstructure:"version"`
    DivinationType     string   `mapstructure:"divination_type"`
    Name               string   `mapstructure:"name"`
    Title              string   `mapstructure:"title"`
    Subtitle           string   `mapstructure:"subtitle"`
    Icon               string   `mapstructure:"icon"`
    Description        string   `mapstructure:"description"`
    Enabled            bool     `mapstructure:"enabled"`
    SystemIdentity     string   `mapstructure:"system_identity"`
    ReasoningFramework []string `mapstructure:"reasoning_framework"`
    VoiceStyle         []string `mapstructure:"voice_style"`
    OutputStructure    []string `mapstructure:"output_structure"`
    Guardrails         []string `mapstructure:"guardrails"`
}

// ProfilesConfig 是 prompt_profiles.yaml 的顶层结构。
type ProfilesConfig struct {
    Profiles        map[string]ProfileConfig `mapstructure:"prompt_profiles"`
    DefaultBindings map[string]string        `mapstructure:"default_bindings"`
}
```

命名说明：
- `ProfileConfig` 而非 `PromptProfileConfig`：包名 `config` 已提供上下文，`Profile` 足够清晰。
- `ProfilesConfig` 而非 `PromptProfilesConfig`：同理，避免冗余前缀。

## 6. 配置读取方法

### 6.1 加载函数

参照 `config.Load()` 的模式，在 `server/config/prompt_profiles.go` 中实现：

```go
package config

import (
    "fmt"

    "github.com/spf13/viper"
)

// LoadProfiles 从 YAML 文件加载所有人格配置。
func LoadProfiles(configPath string) (*ProfilesConfig, error) {
    v := viper.New()

    if configPath != "" {
        v.SetConfigFile(configPath)
    } else {
        v.SetConfigName("prompt_profiles")
        v.SetConfigType("yaml")
        v.AddConfigPath(".")
        v.AddConfigPath("./config")
        v.AddConfigPath("../config")
    }

    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("read prompt profiles config: %w", err)
    }

    var cfg ProfilesConfig
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshal prompt profiles: %w", err)
    }

    return &cfg, nil
}

// DefaultProfile 获取指定占卜类型的默认人格配置。
// 返回指针便于调用方判断 nil；未找到或已禁用时返回非 nil error。
func (c *ProfilesConfig) DefaultProfile(divinationType string) (*ProfileConfig, error) {
    profileID, ok := c.DefaultBindings[divinationType]
    if !ok {
        return nil, fmt.Errorf("no default profile for type: %s", divinationType)
    }

    profile, ok := c.Profiles[profileID]
    if !ok {
        return nil, fmt.Errorf("profile not found: %s", profileID)
    }

    if !profile.Enabled {
        return nil, fmt.Errorf("profile disabled: %s", profileID)
    }

    return &profile, nil
}

// Profile 按 ID 获取人格配置。
func (c *ProfilesConfig) Profile(profileID string) (*ProfileConfig, error) {
    profile, ok := c.Profiles[profileID]
    if !ok {
        return nil, fmt.Errorf("profile not found: %s", profileID)
    }
    return &profile, nil
}
```

命名说明：
- `LoadProfiles` 而非 `LoadPromptProfiles`：包名 `config` 已限定上下文。
- `DefaultProfile` 而非 `ResolveDefaultProfile`：Go 惯用简洁动词，`Resolve` 冗余。
- `Profile` 而非 `GetProfile`：Go 惯例 getter 不加 `Get` 前缀（Effective Go）。
- error message 小写开头、无标点：Go 社区规范（`fmt.Errorf` 不以大写开头、不以句号结尾）。

### 6.2 在 main.go 中加载

```go
// main.go
func main() {
    // ... 现有配置加载 ...

    // 加载人格配置
    profiles, err := config.LoadProfiles("")
    if err != nil {
        log.Printf("warning: load prompt profiles: %v", err)
        // 降级：空配置，后续 DefaultProfile 会返回明确错误
        profiles = &config.ProfilesConfig{
            Profiles:        make(map[string]config.ProfileConfig),
            DefaultBindings: make(map[string]string),
        }
    }

    // 传递给需要的服务
    // ...
}
```

### 6.3 读取流程图

```
server/config/prompt_profiles.yaml
    ↓
viper.ReadInConfig() + Unmarshal
    ↓
ProfilesConfig { Profiles, DefaultBindings }
    ↓
DefaultProfile(divinationType) → *ProfileConfig
    ↓
Compose(profile, input) → *Messages{System, User}
```

## 7. Compose 函数设计

### 7.1 函数签名

```go
// Compose 根据人格配置和解读输入，组装 AI 请求所需的 system/user 消息。
// 它是纯函数：不访问数据库、全局状态、网络或日志。
func Compose(profile *ProfileConfig, input InterpretationInput) (*Messages, error)
```

签名说明：
- `profile *ProfileConfig`：指针传递，避免拷贝大结构体（含多个 `[]string` 字段）。
- `input InterpretationInput`：值传递，结构体较小（5 个 string + 1 个 map），且调用方通常已持有值。
- 返回 `*Messages`：指针返回便于调用方判断 nil，与 error 配合语义一致。
- 参数名不使用 `context`：避免遮蔽标准库 `context.Context`。

### 7.2 输入约束

Compose 必须校验：

- `profile.DivinationType` 非空。
- `profile.DivinationType == input.DivinationType`。
- `input.Question` 非空。
- `input.ResultFacts` 包含该占卜类型的最低必要字段。

例如梅花易数最低必要字段：

- `Method`
- `BenGua`
- `HuGua`
- `BianGua`
- `MovingLine`
- `Ti`
- `Yong`
- `TiYongRelation`

### 7.3 输出规则

`Messages.System` 只包含身份、推理框架、表达风格、边界规则。

`Messages.User` 只包含本次问题、占卜事实、要求输出的结构。

示意：

```text
System:
你是「康节先生」人格...
推理框架：
1. 先以体用生克定大方向
2. 再看本卦现状
3. 再看互卦过程
4. 再看动爻与变卦
...

User:
用户问题：...
起卦方式：...
本卦：...
互卦：...
...
请按以下结构输出：...
```

这样区分后，模型不会把本次数据误认为身份规则，也不会把身份文本混进占卜事实。

## 8. ResultFacts 解析器

为每种占卜类型提供一个 facts 解析接口：

```go
// FactParser 从占卜结果 JSON 中提取关键事实。
// 单方法接口以 er 结尾命名，符合 Go 惯例。
type FactParser interface {
    Parse(resultJSON, question string) (map[string]string, error)
}
```

先落地为函数映射（无需显式实现接口，函数签名匹配即可）：

```go
var factParsers = map[string]func(resultJSON, question string) (map[string]string, error){
    "meihua":    parseMeihuaFacts,
    "liuyao_v2": parseLiuyaoV2Facts,
    // ...
}
```

第一阶段只迁移梅花易数。高岛易断排在第二阶段。理由：梅花易数刚新增，用户体验正在成型，先把一个闭环做到完整，比同时迁移多个类型更有价值。

## 9. 存储设计

### 9.1 DivinationRecord 新增字段

建议在 `DivinationRecord` 增加：

```go
PromptProfileID      string `gorm:"type:text;default:''" json:"prompt_profile_id"`
PromptProfileName    string `gorm:"type:text;default:''" json:"prompt_profile_name"`
PromptProfileVersion string `gorm:"type:text;default:''" json:"prompt_profile_version"`
```

记录保存时写入默认 profile。

### 9.2 为什么存在 DivinationRecord 而不是 ChatSession

AI 解读属于某条占卜记录，而不是某个聊天会话。

同一条记录可能：

- 直接 AI 解读。
- 进入聊天继续追问。
- 在历史记录中查看。

所以 profile 信息应保存到 `DivinationRecord`。

## 10. 前后端协调

第一版采用方案 B：保持前端 `divinationPersona.ts` 存在，但后端返回实际 prompt 人格字段，前端展示时优先使用后端字段。

返回字段来源：

- 历史详情：`DivinationRecord.prompt_profile_name`
- 聊天会话关联记录：record detail 中带出 profile 字段
- 新建聊天后：session 返回的 record 中带出 profile 字段

前端展示优先级：

```text
record.prompt_profile_name
> divinationPersona.ts 默认 name
> "AI 占卜师"
```

这样避免第一版大改前端，同时保证后端实际人格可以被用户看到。

## 11. API 行为

第一版不需要开放人格选择 API。

后端内部通过配置提供：

```go
// 在 service 初始化时注入
type AIService struct {
    profiles *config.ProfilesConfig
    // ...
}

func (s *AIService) Interpret(...) {
    profile, err := s.profiles.DefaultProfile(record.Type)
    // ...
}
```

未来需要用户选择人格时，再增加：

```http
GET /api/prompt-profiles?type=meihua
POST /api/chat/sessions
{
  "type": "meihua",
  "question": "...",
  "prompt_profile_id": "shaoyong_meihua"
}
```

第一版可以先保留 `prompt_profile_id` 入参扩展点，但不在前端暴露。

## 12. 实施步骤

### Step 1：创建配置文件与加载方法

新增：

- `server/config/prompt_profiles.yaml` — 所有人格配置
- `server/config/prompt_profiles.go` — 加载与读取方法

要求：

- `LoadProfiles` 使用 viper 加载，与 `Load()` 模式一致。
- `DefaultProfile` 支持按占卜类型查找默认人格。
- `Profile` 支持按 ID 查找人格。
- 单元测试覆盖：加载成功、类型不匹配、profile 不存在、profile 已禁用。

### Step 2：定义 Compose 函数

新增：

- `server/internal/service/prompt_composer.go`
- `server/internal/service/prompt_composer_test.go`

要求：

- Compose 是纯函数。
- 签名：`func Compose(profile *config.ProfileConfig, input InterpretationInput) (*Messages, error)`。
- 单元测试覆盖 profile/type 不匹配、缺字段、正常输出。

### Step 3：将梅花易数 prompt 拆成 profile + facts

把现有 `meihua_prompt.txt` 拆为：

- `ProfileConfig.SystemIdentity`（对应 `prompt_profiles.yaml` 中的 `system_identity`）
- `ProfileConfig.ReasoningFramework`
- `ProfileConfig.VoiceStyle`
- `ProfileConfig.OutputStructure`
- `ProfileConfig.Guardrails`
- `parseMeihuaFacts`

测试要求：

- 同一条梅花结果能生成包含 `康节先生`、`本卦`、`体用生克` 的 prompt。
- `reasoning_framework` 中必须先出现体用生克，再出现本卦/互卦/变卦细解。
- 迁移完成后，梅花易数不再依赖 `server/internal/service/prompts/meihua_prompt.txt` 作为真相源。该文件应删除，或仅作为已废弃参考文档迁入 `docs/`。

### Step 4：统一梅花易数的 AIService 和 ChatService prompt 路径

修改：

- `OpenAIProvider.Interpret`
- `ChatService.buildSystemPrompt`

当 `record.Type == "meihua"` 时，两者都调用：

```go
profile, err := profiles.DefaultProfile(record.Type)
facts, err := parseFacts(record.Type, record.Result, record.Question)
messages, err := Compose(profile, InterpretationInput{...})
```

聊天路径使用 `messages.System` 作为 system prompt，用户消息仍保留历史对话。

直接解读路径使用 `messages.System + messages.User` 组成请求。

其他占卜类型第一版继续走现有逻辑，避免一次性迁移造成风险。

### Step 5：DivinationRecord 保存 profile 信息

新增字段并在保存 record 时写入：

- `PromptProfileID`
- `PromptProfileName`
- `PromptProfileVersion`

覆盖路径：

- 普通占卜 API 保存记录
- ChatModeStarter 保存记录
- 未来重新生成 AI 解读时，继续使用记录上的 profile 信息；如果旧记录为空，则回退默认 profile。

### Step 6：前端展示后端人格名称

调整：

- `client/src/services/chat.ts`
- `client/src/services/history.ts`
- `ChatMessage`
- `DivinationResultCard`
- `History`

展示策略：

```text
record.prompt_profile_name || persona.name
```

不要新增复杂选择器。

### Step 7：迁移其他占卜类型

迁移顺序建议：

1. 梅花易数
2. 高岛易断
3. 通用聊天 prompt
4. 塔罗 / 八字 / 星座

每迁移一个类型，增加对应 facts builder 测试。

## 13. 文件级开发清单

后端新增文件：

- `server/config/prompt_profiles.yaml`
- `server/config/prompt_profiles.go`
- `server/config/prompt_profiles_test.go`
- `server/internal/service/prompt_composer.go`
- `server/internal/service/prompt_composer_test.go`
- `server/internal/service/prompt_facts.go`
- `server/internal/service/prompt_facts_meihua.go`
- `server/internal/service/prompt_facts_meihua_test.go`

后端修改文件：

- `server/main.go`：加载 `ProfilesConfig` 并传入 router/service 初始化链路。
- `server/internal/router/router.go`：把 profiles 注入 `AIService` 与 `ChatService`。
- `server/internal/model/divination.go`：新增 profile 字段。
- `server/internal/service/ai.go`：梅花易数直接解读改走 Compose。
- `server/internal/service/chat.go`：梅花易数聊天 prompt 改走 Compose。
- `server/internal/service/chat_starter.go`：保存梅花记录时写入默认 profile 信息。
- `server/internal/handler/history.go`：无需特殊逻辑，但确认返回结构包含新增 JSON 字段。

前端修改文件：

- `client/src/services/chat.ts`：`DivinationRecord` 类型增加 profile 字段。
- `client/src/services/history.ts`：`DivinationRecord` 类型增加 profile 字段。
- `client/src/components/chat/ChatMessage.tsx`：优先显示后端 profile 名称。
- `client/src/components/chat/DivinationResultCard.tsx`：结果卡底部优先显示后端 profile 名称。
- `client/src/pages/History.tsx`：详情中展示 profile 名称。

删除或迁移文件：

- `server/internal/service/prompts/meihua_prompt.txt`：迁移后不再作为运行时真相源。

## 14. Compose 输出格式

Compose 不直接返回 `[]map[string]string`，而是返回 `Messages`。调用方负责转为不同 AI Provider 需要的消息格式。

```go
func Compose(profile *config.ProfileConfig, input InterpretationInput) (*Messages, error) {
    if profile == nil {
        return nil, fmt.Errorf("profile is nil")
    }
    if profile.DivinationType == "" {
        return nil, fmt.Errorf("profile divination type is empty")
    }
    if profile.DivinationType != input.DivinationType {
        return nil, fmt.Errorf("profile type %s does not match input type %s", profile.DivinationType, input.DivinationType)
    }
    if strings.TrimSpace(input.Question) == "" {
        return nil, fmt.Errorf("question is empty")
    }

    system := buildSystem(profile)
    user := buildUser(profile, input)
    return &Messages{System: system, User: user}, nil
}
```

`buildSystem(profile)` 包含：

- `SystemIdentity`
- `ReasoningFramework`
- `VoiceStyle`
- `Guardrails`

`buildUser(profile, input)` 包含：

- 用户问题
- `ResultFacts`
- `OutputStructure`

## 15. 梅花易数开发样例

输入：

```go
profile := shaoyongMeihuaProfile
input := InterpretationInput{
    DivinationType: "meihua",
    Question: "最近事业如何？",
    ResultFacts: map[string]string{
        "Method": "时间起卦（丙午年五月初九子时）",
        "BenGua": "泽火革（兑上离下）",
        "HuGua": "天风姤（乾上巽下）",
        "BianGua": "泽山咸（兑上艮下）",
        "MovingLine": "初爻动",
        "Ti": "兑（金）",
        "Yong": "离（火）",
        "TiYongRelation": "用克体",
    },
    Mode: "chat",
}
```

验收输出必须包含：

- `System` 中包含 `康节先生`。
- `System` 中包含“不得自称历史上的邵雍本人”。
- `System` 中先出现 `体用生克`，再出现 `本卦`、`互卦`、`变卦`。
- `User` 中包含用户问题和全部梅花 facts。
- `User` 中包含输出结构。
- 不出现原始 JSON 大段转储，除非 facts 解析失败进入 fallback。

## 16. 测试计划

### 16.1 单元测试

- `LoadProfiles` 能正确加载 YAML 文件。
- `DefaultProfile("meihua")` 返回 `shaoyong_meihua`。
- `DefaultProfile("unknown_type")` 返回错误。
- `Profile("shaoyong_meihua")` 返回正确配置。
- `Profile("not_exist")` 返回错误。
- `Compose` 对同一输入输出稳定。
- `Compose` 不读取全局 profile。
- `Compose` 在 profile/type 不匹配时报错。
- `parseMeihuaFacts` 能正确提取本卦、互卦、变卦、动爻、体用。
- `parseMeihuaFacts` 对缺失字段返回明确错误。

### 16.2 集成测试

- Chat 模式创建梅花易数记录后，record 保存 `shaoyong_meihua/v1`。
- 初始 AI 解读使用 `康节先生` profile。
- 历史详情返回 profile 字段。
- 旧记录没有 profile 字段时仍可解读，自动回退默认 profile。

### 16.3 前端测试

- 梅花易数聊天消息显示后端返回的 `康节先生`。
- 没有 `prompt_profile_name` 时仍显示本地 persona name。
- 历史详情显示 profile name。

### 16.4 回归验证命令

完成开发后必须通过：

```bash
cd server && go test ./...
cd client && npm run lint
cd client && npm run build
cd client && npm run test -- History.test.tsx
```

如果新增前端测试文件，补跑对应测试文件。

## 17. 兼容与迁移

旧记录没有 profile 字段：

- 读取时允许为空。
- 构建 prompt 时如果为空，则 `DefaultProfile(record.Type)`。
- 不批量回填旧数据，避免不必要的数据迁移风险。

PromptProfile 第一版使用 YAML 配置文件，与 `config.yaml` 同级。

后续迁移路径：

- YAML 配置 → 数据库存储（如果需要后台编辑）
- 通过同一 `DefaultProfile` 接口，调用方无需感知数据来源变化

## 18. 风险与约束

### 风险 1：人格写得太像角色扮演

控制：

- 不写"你就是历史上的邵雍本人"。
- 写"你是康节先生人格，以邵雍一脉方法解读"。

### 风险 2：语气修改影响推理

控制：

- `reasoning_framework` 和 `voice_style` 分字段。
- 修改 `voice_style` 不允许改变 `reasoning_framework` 测试快照。

### 风险 3：前后端展示不一致

控制：

- 后端返回 `prompt_profile_name`。
- 前端优先显示后端字段。

### 风险 4：prompt 越来越长

控制：

- Compose 输出前统计字符数。
- Profile 中控制输出结构，不堆叠大量背景故事。
- 对高岛易断这类带原典证据的类型保留证据裁剪逻辑。

### 风险 5：YAML 配置文件加载失败

控制：

- `main.go` 中加载失败时打印 Warning，不阻塞启动。
- 降级为空配置，后续 `DefaultProfile` 会返回明确错误。
- 与现有 `config.Load()` 的降级策略一致。

## 19. 开发就绪检查表

这份设计达到可开发状态，需要满足以下检查项：

- 第一版范围只包含梅花易数人格系统闭环。
- Compose 是纯函数，且单测能证明它不依赖全局配置。
- 梅花易数的直接 AI 解读和聊天 AI 解读都走 Compose。
- `meihua_prompt.txt` 不再作为运行时真相源。
- `DivinationRecord` 保存 profile ID、name、version。
- 前端只展示 profile，不允许用户选择 profile。
- 旧记录 profile 为空时能 fallback 到默认 profile。
- 所有新增错误信息小写开头、无句号。
- 所有新增 Go getter 不使用 `Get` 前缀。
- 所有 profile 配置走 YAML + viper。

## 20. Definition of Done

开发完成的定义：

1. 创建一条梅花易数聊天记录，数据库记录包含：
   - `prompt_profile_id = "shaoyong_meihua"`
   - `prompt_profile_name = "康节先生"`
   - `prompt_profile_version = "v1"`
2. 初始 AI 解读的 system prompt 来自 Compose，而不是 `meihua_prompt.txt`。
3. 直接 AI 解读梅花易数时，也使用同一 Compose 输出。
4. 聊天界面显示「康节先生」。
5. 历史详情显示「康节先生」。
6. 删除或停用 `meihua_prompt.txt` 后，梅花易数 AI 解读仍正常。
7. 旧记录缺少 profile 字段时仍能进入聊天并生成解读。
8. `go test ./...`、`npm run lint`、`npm run build` 全部通过。

## 21. 审查重点

请重点审查以下决策：

1. 第一版是否只做默认人格绑定，不做用户选择人格。
2. PromptProfile 配置使用 `server/config/prompt_profiles.yaml` + viper 加载，而不是 `prompt_profiles/*.json` + `//go:embed`。
3. 配置结构体 `ProfileConfig` / `ProfilesConfig` 放在 `server/config/` 包中，与现有配置体系一致。
4. `Compose` 签名：`func Compose(profile *ProfileConfig, input InterpretationInput) (*Messages, error)`，指针传参、避免 `context` 命名。
5. 单方法接口 `FactParser`（方法名 `Parse`），函数映射 `parseMeihuaFacts` / `parseLiuyaoV2Facts`。
6. getter 方法不加 `Get` 前缀：`DefaultProfile()` / `Profile()`。
7. error message 小写开头无标点。
8. profile 字段是否保存到 `DivinationRecord`。
9. 前端是否采用方案 B：优先展示后端 profile 字段，但暂不移除本地 `divinationPersona.ts`。
10. 梅花易数默认人格名称是否定为「康节先生」，而不是「观梅心易师」。
11. YAML 中使用 `system_identity` 多行文本存储身份设定，而不是拆成更细的字段。
