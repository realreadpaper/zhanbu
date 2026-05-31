# 占卜网站 — 技术设计文档

> 版本：v1.0 | 最后更新：2026-05-30

---

## 目录

1. [项目概述](#1-项目概述)
2. [占卜类型设计](#2-占卜类型设计)
3. [功能设计](#3-功能设计)
4. [技术架构](#4-技术架构)
5. [目录结构](#5-目录结构)
6. [开发计划](#6-开发计划)

---

## 1. 项目概述

### 1.1 项目目标与愿景

构建一个现代化、用户体验优秀的在线占卜平台，融合东方（周易六爻、八字）与西方（塔罗牌、星座）占卜体系，通过 AI 大模型为用户提供个性化的解读服务。

**愿景**：让每个人都能便捷地探索命运的奥秘，在传统智慧与现代科技的交汇处获得启发。

### 1.2 目标用户

| 用户群体 | 特征 | 核心需求 |
|---------|------|---------|
| 年轻白领 | 18-35岁，对星座/塔罗感兴趣 | 日常运势、情感指引 |
| 传统文化爱好者 | 对周易、八字有研究兴趣 | 专业的六爻/八字排盘 |
| 好奇用户 | 偶尔尝试，娱乐为主 | 简单易用、结果有趣 |
| 深度用户 | 长期关注个人运势 | 历史记录、趋势分析 |

### 1.3 核心价值

- **专业性**：每种占卜类型均基于经典算法，确保结果的准确性和权威性
- **智能化**：AI 大模型提供深度个性化解读，而非模板化输出
- **美观性**：精心设计的动画效果和 UI，营造沉浸式体验
- **私密性**：用户数据加密存储，占卜记录仅本人可见
- **免费核心**：基础占卜功能永久免费，AI 深度解读可作为增值功能

---

## 2. 占卜类型设计

### 2.1 塔罗牌（Tarot）

#### 2.1.1 牌组结构

塔罗牌共 **78 张**，分为两大类：

**大阿尔卡纳（Major Arcana）— 22 张**

| 编号 | 名称 | 英文 | 关键词 |
|------|------|------|--------|
| 0 | 愚者 | The Fool | 开始、冒险、天真 |
| 1 | 魔术师 | The Manifestation | 创造力、意志力 |
| 2 | 女祭司 | The High Priestess | 直觉、神秘、潜意识 |
| 3 | 皇后 | The Empress | 丰饶、母性、自然 |
| 4 | 皇帝 | The Emperor | 权威、结构、掌控 |
| 5 | 教皇 | The Hierophant | 传统、信仰、教育 |
| 6 | 恋人 | The Lovers | 爱情、选择、和谐 |
| 7 | 战车 | The Chariot | 胜利、意志、行动 |
| 8 | 力量 | Strength | 勇气、耐心、内在力量 |
| 9 | 隐士 | The Hermit | 内省、智慧、孤独 |
| 10 | 命运之轮 | Wheel of Fortune | 转变、命运、循环 |
| 11 | 正义 | Justice | 公平、真相、因果 |
| 12 | 倒吊人 | The Hanged Man | 牺牲、放下、新视角 |
| 13 | 死神 | Death | 结束、转变、重生 |
| 14 | 节制 | Temperance | 平衡、耐心、调和 |
| 15 | 恶魔 | The Devil | 束缚、欲望、物质 |
| 16 | 塔 | The Tower | 突变、破坏、觉醒 |
| 17 | 星星 | The Star | 希望、灵感、宁静 |
| 18 | 月亮 | The Moon | 幻象、恐惧、潜意识 |
| 19 | 太阳 | The Sun | 成功、快乐、活力 |
| 20 | 审判 | Judgement | 觉醒、重生、召唤 |
| 21 | 世界 | The World | 完成、整合、成就 |

**小阿尔卡纳（Minor Arcana）— 56 张**

分为四个花色，每个花色 14 张：

| 花色 | 英文 | 元素 | 象征领域 |
|------|------|------|---------|
| 权杖 | Wands | 火 | 行动、创造力、激情 |
| 圣杯 | Cups | 水 | 情感、关系、直觉 |
| 宝剑 | Swords | 风 | 思想、冲突、真相 |
| 金币 | Pentacles | 土 | 物质、健康、财富 |

每个花色包含：
- 数字牌（Ace-10）：10 张
- 宫廷牌（Page, Knight, Queen, King）：4 张

#### 2.1.2 牌阵设计

**单牌抽取**
- 最简单的占卜方式，适合日常问题
- 抽取 1 张牌，代表当前核心信息

**三牌阵（Past-Present-Future）**
```
┌─────┐ ┌─────┐ ┌─────┐
│ 过去 │ │ 现在 │ │ 未来 │
└─────┘ └─────┘ └─────┘
```
- 牌1：过去的影响力
- 牌2：当前状况
- 牌3：未来趋势

**凯尔特十字阵（Celtic Cross）— 10 张**
```
        ┌─────┐
        │  5  │
        └─────┘
   ┌─────┬─────┬─────┐
   │  6  │  1  │  7  │
   └─────┼─────┼─────┘
      ┌──┤  2  ├──┐
      │3 │     │ 4│
      └──┴─────┴──┘
        ┌─────┐
        │  8  │
        └─────┘
        ┌─────┐
        │  9  │
        └─────┘
        ┌─────┐
        │ 10  │
        └─────┘
```
- 牌1：当前状况/核心问题
- 牌2：挑战/障碍（交叉牌）
- 牌3：潜意识/根基
- 牌4：过去
- 牌5：意识/目标
- 牌6：近期未来
- 牌7：自我认知
- 牌8：外部环境
- 牌9：希望与恐惧
- 牌10：最终结果

**爱情十字阵（Love Cross）— 5 张**
```
      ┌─────┐
      │  3  │
      └─────┘
┌─────┬─────┬─────┐
│  4  │  1  │  5  │
└─────┼─────┴─────┘
   ┌──┤  2  │
   │  │     │
   └──┴─────┘
```
- 牌1：你的现状
- 牌2：对方的现状
- 牌3：关系的挑战
- 牌4：你的期望
- 牌5：关系的走向

#### 2.1.3 抽牌算法

```
算法：Fisher-Yates 洗牌 + 正逆位

输入：
  - deck: 78 张牌数组
  - spread: 牌阵类型（single/three/celtic/love）
  - count: 抽牌数量（由牌阵决定）

步骤：
1. 初始化牌组 deck[0..77]
2. Fisher-Yates 洗牌：
   for i = 77 downto 1:
     j = random(0, i)
     swap(deck[i], deck[j])
3. 从洗好的牌组顶部抽取 count 张
4. 每张牌随机决定正位或逆位：
   orientation = random(0, 1) == 0 ? "upright" : "reversed"
5. 返回抽取结果数组

时间复杂度：O(n)
空间复杂度：O(1)
```

#### 2.1.4 正逆位解读

每张牌存储两套关键词和解读：
- **正位（Upright）**：牌面正常含义
- **逆位（Reversed）**：牌面倒置，通常表示能量受阻、内在化或相反含义

### 2.2 星座运势

#### 2.2.1 十二星座数据

| 星座 | 英文 | 日期范围 | 元素 | 守护星 |
|------|------|---------|------|--------|
| 白羊座 | Aries | 3.21 - 4.19 | 火 | 火星 |
| 金牛座 | Taurus | 4.20 - 5.20 | 土 | 金星 |
| 双子座 | Gemini | 5.21 - 6.21 | 风 | 水星 |
| 巨蟹座 | Cancer | 6.22 - 7.22 | 水 | 月亮 |
| 狮子座 | Leo | 7.23 - 8.22 | 火 | 太阳 |
| 处女座 | Virgo | 8.23 - 9.22 | 土 | 水星 |
| 天秤座 | Libra | 9.23 - 10.23 | 风 | 金星 |
| 天蝎座 | Scorpio | 10.24 - 11.22 | 水 | 冥王星 |
| 射手座 | Sagittarius | 11.23 - 12.21 | 火 | 木星 |
| 摩羯座 | Capricorn | 12.22 - 1.19 | 土 | 土星 |
| 水瓶座 | Aquarius | 1.20 - 2.18 | 风 | 天王星 |
| 双鱼座 | Pisces | 2.19 - 3.20 | 海王星 | 海王星 |

#### 2.2.2 运势维度

每个星座运势包含以下维度评分（1-5星）：

```json
{
  "zodiac": "aries",
  "period": "daily",
  "date": "2026-05-30",
  "overall": 4,
  "love": 3,
  "career": 5,
  "wealth": 3,
  "health": 4,
  "lucky_number": 7,
  "lucky_color": "红色",
  "summary": "今天适合开展新项目...",
  "detail": {
    "love": "感情方面...",
    "career": "事业方面...",
    "wealth": "财运方面...",
    "health": "健康方面..."
  }
}
```

#### 2.2.3 运势生成算法

```
算法：基于日期的确定性运势生成

输入：
  - zodiac: 星座名称
  - date: 目标日期
  - period: daily | weekly | monthly

步骤：
1. 计算种子值：
   seed = hash(zodiac + date + period)
   // 使用 SHA-256 取前 8 字节作为种子

2. 用种子初始化伪随机数生成器（确保相同输入得到相同输出）

3. 生成各维度评分：
   overall = seeded_random(1, 5)
   love = seeded_random(1, 5)
   career = seeded_random(1, 5)
   wealth = seeded_random(1, 5)
   health = seeded_random(1, 5)

4. 生成幸运元素：
   lucky_number = seeded_random(1, 9)
   lucky_color = colors[seeded_random(0, len(colors)-1)]

5. 从预置模板库中选取解读文本：
   - 按评分区间（1-2低/3中/4-5高）选取对应模板
   - 结合星座特质填充具体描述

6. 返回运势数据

注意：
- 每日运势以日期为种子，保证同一天结果一致
- 周运势以 ISO 周数为种子
- 月运势以年月为种子
- 后续可通过 AI 大模型替换模板文本，生成更自然的解读
```

#### 2.2.4 运势模板库结构

```
templates/horoscope/
├── daily/
│   ├── love_high.txt      # 爱情运势（4-5星）
│   ├── love_mid.txt       # 爱情运势（3星）
│   ├── love_low.txt       # 爱情运势（1-2星）
│   ├── career_high.txt
│   ├── career_mid.txt
│   ├── career_low.txt
│   ├── wealth_high.txt
│   ├── wealth_mid.txt
│   ├── wealth_low.txt
│   ├── health_high.txt
│   ├── health_mid.txt
│   └── health_low.txt
├── weekly/
│   └── ... (同上结构)
└── monthly/
    └── ... (同上结构)
```

### 2.3 周易六爻

#### 2.3.1 基础概念

**八卦（Eight Trigrams）**

| 卦名 | 符号 | 二进制 | 自象 | 属性 |
|------|------|--------|------|------|
| 乾 | ☰ | 111 | 天 | 阳金 |
| 兑 | ☱ | 110 | 泽 | 阴金 |
| 离 | ☲ | 101 | 火 | 阴火 |
| 震 | ☳ | 100 | 雷 | 阳木 |
| 巽 | ☴ | 011 | 风 | 阴木 |
| 坎 | ☵ | 010 | 水 | 阳水 |
| 艮 | ☶ | 001 | 山 | 阳土 |
| 坤 | ☷ | 000 | 地 | 阴土 |

**六十四卦（64 Hexagrams）**

由两个八卦上下组合而成（上卦 + 下卦），共 64 种组合。每卦包含：
- 卦名（如"乾为天"、"坤为地"）
- 卦辞：对卦的整体描述
- 爻辞：6 条爻的各自解释
- 象辞：卦象的象征意义

#### 2.3.2 掷铜钱算法

传统六爻使用三枚铜钱，投掷六次生成六爻。

```
算法：掷铜钱法

工具：三枚铜钱（有字面为阳，无字面为阴）

投掷一次的规则：
  三枚铜钱落地后，计算有字面朝上的数量：
  - 3 个字面 → 老阳（阳爻，动爻）→ 记为 6（⚋变⚊）
  - 2 个字面 → 少阴（阴爻，静爻）→ 记为 8
  - 1 个字面 → 少阳（阳爻，静爻）→ 记为 7
  - 0 个字面 → 老阴（阴爻，动爻）→ 记为 9（⚊变⚋）

程序实现：
  function throw_once():
    coins = [random(0,1), random(0,1), random(0,1)]
    yang_count = sum(coins)  // 1 = 字面
    switch yang_count:
      case 3: return { value: 6, type: "old_yang", mutable: true }   // ⚋→⚊
      case 2: return { value: 8, type: "young_yin", mutable: false }
      case 1: return { value: 7, type: "young_yang", mutable: false }
      case 0: return { value: 9, type: "old_yin", mutable: true }    // ⚊→⚋

生成本卦：
  lines = []
  for i in 0..5:
    lines[i] = throw_once()  // 初爻到上爻
  本卦 = 组合6爻得到的卦

生成变卦（如有动爻）：
  变爻规则：老阳变阴，老阴变阳
  变卦 = 将本卦中的动爻取反后得到的卦

概率分析：
  老阳(6): 1/8 = 12.5%
  少阴(8): 3/8 = 37.5%
  少阳(7): 3/8 = 37.5%
  老阴(9): 1/8 = 12.5%
  动爻概率: 25%（每次投掷有动爻的概率）
```

#### 2.3.3 数据结构

```go
type Trigram struct {
    Name     string   // 卦名：乾、兑、离、震、巽、坎、艮、坤
    Symbol   string   // 符号：☰、☱、☲、☳、☴、☵、☶、☷
    Binary   int      // 二进制：7,6,5,4,3,2,1,0
    Nature   string   // 自象：天、泽、火、雷、风、水、山、地
    Element  string   // 五行：金、金、火、木、木、水、土、土
    YinYang  string   // 阴阳：阳、阴、阴、阳、阴、阳、阳、阴
}

type Hexagram struct {
    Name       string     // 卦名
    Upper      Trigram    // 上卦
    Lower      Trigram    // 下卦
    Judgment   string     // 卦辞
    Image      string     // 象辞
    Lines      [6]string  // 六爻爻辞
    Binary     string     // 6位二进制表示
}

type LineResult struct {
    Value    int     // 6/7/8/9
    Type     string  // old_yang/young_yang/old_yin/young_yin
    Mutable  bool    // 是否为动爻
    Symbol   string  // ⚊ 或 ⚋
}

type LiuYaoResult struct {
    Question   string       // 用户问题
    Lines      [6]LineResult // 六爻结果
    BenGua     Hexagram     // 本卦
    BianGua    *Hexagram    // 变卦（有动爻时）
    Mutables   []int        // 动爻位置索引
    Timestamp  time.Time    // 占卜时间
}
```

#### 2.3.4 解读逻辑

```
1. 确定本卦：根据6爻组合找到对应的64卦
2. 分析动爻：
   - 无动爻：以本卦卦辞为主
   - 有动爻：以本卦卦辞为主，动爻爻辞为辅，变卦为参考
   - 多个动爻：取主要动爻（上爻优先）
3. 综合解读：
   - 卦象分析（上下卦关系）
   - 五行生克（上下卦五行关系）
   - 爻位分析（每爻的位置意义）
4. 结合用户问题，给出针对性解读
```

### 2.4 八字排盘

#### 2.4.1 基础概念

**天干（10个）**

| 序号 | 天干 | 阴阳 | 五行 |
|------|------|------|------|
| 1 | 甲 | 阳 | 木 |
| 2 | 乙 | 阴 | 木 |
| 3 | 丙 | 阳 | 火 |
| 4 | 丁 | 阴 | 火 |
| 5 | 戊 | 阳 | 土 |
| 6 | 己 | 阴 | 土 |
| 7 | 庚 | 阳 | 金 |
| 8 | 辛 | 阴 | 金 |
| 9 | 壬 | 阳 | 水 |
| 10 | 癸 | 阴 | 水 |

**地支（12个）**

| 序号 | 地支 | 生肖 | 阴阳 | 五行 |
|------|------|------|------|------|
| 1 | 子 | 鼠 | 阳 | 水 |
| 2 | 丑 | 牛 | 阴 | 土 |
| 3 | 寅 | 虎 | 阳 | 木 |
| 4 | 卯 | 兔 | 阴 | 木 |
| 5 | 辰 | 龙 | 阳 | 土 |
| 6 | 巳 | 蛇 | 阴 | 火 |
| 7 | 午 | 马 | 阳 | 火 |
| 8 | 未 | 羊 | 阴 | 土 |
| 9 | 申 | 猴 | 阳 | 金 |
| 10 | 酉 | 鸡 | 阴 | 金 |
| 11 | 戌 | 狗 | 阳 | 土 |
| 12 | 亥 | 猪 | 阴 | 水 |

**六十甲子**：天干（10）× 地支（12）的最小公倍数 = 60 种组合

#### 2.4.2 排盘算法

```
算法：八字排盘

输入：
  - birthDate: 出生日期（公历）
  - birthTime: 出生时间（24小时制）
  - gender: 性别

步骤：

1. 公历转农历：
   - 使用寿星万年历算法或查表法
   - 考虑闰月情况
   - 输出：农历年、月、日

2. 计算年柱：
   - 天干 = (农历年 - 4) % 10 → 取天干数组索引
   - 地支 = (农历年 - 4) % 12 → 取地支数组索引
   - 示例：2024年 → (2024-4)%10=0 → 甲, (2024-4)%12=0 → 子 → 甲子年

3. 计算月柱：
   - 年干确定月干起点（年上起月法）：
     甲己之年丙作首（甲/己年，正月天干从丙开始）
     乙庚之岁戊为头
     丙辛之年寻庚上
     丁壬壬寅顺水流
     戊癸之年何处起，甲寅之上好追求
   - 月支固定：正月=寅, 二月=卯, ..., 十二月=丑
   - 注意：月柱以节气为界，非农历初一

4. 计算日柱：
   - 使用公式计算从已知基准日到目标日的天数差
   - 基准日：1900年1月1日 = 甲戌日（序号10）
   - 天干序号 = (天数差 + 基准天干序号) % 10
   - 地支序号 = (天数差 + 基准地支序号) % 12

5. 计算时柱：
   - 时辰对照：
     23:00-01:00 → 子时
     01:00-03:00 → 丑时
     03:00-05:00 → 寅时
     ...
     21:00-23:00 → 亥时
   - 日干确定时干起点（日上起时法）：
     甲己还加甲（甲/己日，子时天干从甲开始）
     乙庚丙作初
     丙辛从戊起
     丁壬庚子居
     戊癸何方发，壬子是真途

6. 分析五行：
   - 统计八字中各五行的数量
   - 分析五行旺衰（得令、得地、得生、得助）
   - 判断日主强弱
   - 确定用神和忌神

7. 十神分析：
   - 以日干为"我"，分析其他天干与日干的关系：
     比肩、劫财、食神、伤官、偏财、
     正财、七杀、正官、偏印、正印

8. 排大运：
   - 阳年男/阴年女：顺排
   - 阴年男/阳年女：逆排
   - 起运时间：从出生日到下一个/上一个节气的天数，3天折1年
```

#### 2.4.3 数据结构

```go
type BaZi struct {
    BirthDate   time.Time    // 公历出生日期
    BirthTime   string       // 出生时间
    LunarDate   LunarInfo    // 农历信息
    Gender      string       // 性别

    YearPillar  Pillar       // 年柱
    MonthPillar Pillar       // 月柱
    DayPillar   Pillar       // 日柱
    HourPillar  Pillar       // 时柱

    FiveElements FiveElementAnalysis  // 五行分析
    TenGods      []TenGod             // 十神
    DaYun        []DaYunPeriod        // 大运
}

type Pillar struct {
    TianGan    string  // 天干
    DiZhi      string  // 地支
    TianGanWuXing string // 天干五行
    DiZhiWuXing   string // 地支五行
    NaYin      string  // 纳音
    HiddenGan  []string // 地支藏干
}

type FiveElementAnalysis struct {
    Metal   int  // 金
    Wood    int  // 木
    Water   int  // 水
    Fire    int  // 火
    Earth   int  // 土
    Strongest string  // 最旺五行
    Weakest   string  // 最弱五行
    DayMaster string  // 日主强弱
    YongShen  string  // 用神
    JiShen    string  // 忌神
}
```

#### 2.4.4 节气计算

```
八字排盘中，年柱和月柱的分界点是节气，而非农历初一。

节气对照表（简化）：
- 立春（2月4日前后）→ 年柱切换点
- 惊蛰、清明、立夏、芒种、小暑、
  立秋、白露、寒露、立冬、大雪、小寒

程序实现：使用天文算法或查表法计算精确节气时间
```

---

## 3. 功能设计

### 3.1 页面结构

```
┌─────────────────────────────────────────────────────────┐
│  Header: Logo | 导航栏 | 登录/用户头像                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  [占卜类型卡片]                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │  🔮 塔罗  │  │  ♈ 星座  │  │  ☰ 六爻  │  │  📅 八字  │ │
│  │  牌占卜   │  │  运势     │  │  占卜    │  │  排盘    │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
│                                                         │
│  [精选占卜]                                              │
│  [历史记录]                                              │
│  [AI 解读]                                              │
│                                                         │
├─────────────────────────────────────────────────────────┤
│  Footer                                                │
└─────────────────────────────────────────────────────────┘
```

### 3.2 各页面详细设计

#### 3.2.1 首页（Home）

**路由**：`/`

**功能**：
- 展示四种占卜类型的卡片入口（带图标和简短描述）
- 热门推荐占卜
- 今日运势概览（需登录）
- 最近占卜记录（需登录）

**交互**：
- 卡片悬停动画（放大、阴影加深）
- 点击卡片进入对应占卜页面

#### 3.2.2 塔罗牌占卜页（Tarot）

**路由**：`/tarot`

**流程**：
```
用户输入问题 → 选择牌阵 → 抽牌动画 → 展示牌面 → 逐牌解读 → 综合解读 → AI深度解读（可选）
```

**交互设计**：
1. **输入阶段**：
   - 输入框：输入想要占卜的问题（可选）
   - 牌阵选择：展示可用牌阵的可视化预览
   - 开始按钮

2. **抽牌动画**：
   - 牌组展示：78张牌堆叠的3D效果
   - 洗牌动画：牌组分裂、交错、合并
   - 翻牌动画：牌面翻转，先显示牌背再显示正面
   - 正/逆位动画：逆位牌旋转180度

3. **结果展示**：
   - 牌阵布局：按选定牌阵排列已抽取的牌
   - 单牌解读：点击每张牌显示详细含义
   - 综合解读：所有牌的整体分析
   - AI解读按钮：调用大模型生成个性化解读

#### 3.2.3 星座运势页（Horoscope）

**路由**：`/horoscope`

**流程**：
```
选择星座 → 选择运势周期（日/周/月） → 展示运势 → AI深度解读（可选）
```

**交互设计**：
1. **星座选择**：
   - 12宫圆形布局，每个星座一个图标
   - 可通过生日自动匹配星座

2. **运势展示**：
   - 雷达图：展示各维度评分
   - 运势卡片：图文并茂展示各维度详情
   - 幸运元素：数字和颜色的视觉展示

3. **运势周期切换**：
   - Tab切换：日/周/月
   - 日期选择器（日运势）

#### 3.2.4 周易六爻页（LiuYao）

**路由**：`/liuyao`

**流程**：
```
输入问题 → 掷铜钱动画（6次） → 展示卦象 → 卦辞解读 → 变卦分析 → AI深度解读（可选）
```

**交互设计**：
1. **掷铜钱动画**：
   - 三枚铜钱的3D模型
   - 点击"掷"按钮，铜钱抛起、旋转、落下
   - 每次结果显示（老阳/少阴/少阳/老阴）
   - 已掷次数进度指示（1/6, 2/6, ...）

2. **卦象展示**：
   - 本卦的六爻从下到上依次展示
   - 动爻特殊标记（闪烁或颜色区分）
   - 变卦对比展示

3. **解读展示**：
   - 卦名和卦辞
   - 逐爻解读
   - 五行分析
   - 综合建议

#### 3.2.5 八字排盘页（BaZi）

**路由**：`/bazi`

**流程**：
```
输入出生信息 → 排盘计算 → 展示八字 → 五行分析 → 大运流年 → AI深度解读（可选）
```

**交互设计**：
1. **信息输入**：
   - 日期选择器（公历/农历切换）
   - 时间选择器（24小时制）
   - 性别选择

2. **八字展示**：
   - 四柱表格布局：
     ```
     年柱    月柱    日柱    时柱
     甲      丙      庚      壬
     子      寅      午      午
     ```
   - 纳音显示
   - 地支藏干

3. **五行分析**：
   - 五行数量柱状图
   - 日主强弱判断
   - 用神/忌神说明

4. **大运展示**：
   - 时间轴布局
   - 每步大运的天干地支
   - 当前大运高亮

#### 3.2.6 历史记录页（History）

**路由**：`/history`

**功能**：
- 按时间倒序展示占卜记录
- 按占卜类型筛选
- 点击查看详情
- 删除记录

#### 3.2.7 用户中心（Profile）

**路由**：`/profile`

**功能**：
- 个人信息编辑
- 星座/生肖展示
- 占卜统计（各类型次数）
- 设置（主题、语言等）

### 3.3 AI 解读功能

**触发方式**：占卜结果页面的"AI 解读"按钮

**Prompt 构造**：

```
你是一位经验丰富的占卜师，精通东方和西方占卜术。请根据以下占卜结果，为用户提供个性化的深度解读。

【占卜类型】{type}
【用户问题】{question}
【占卜结果】{result_data}

请从以下角度进行解读：
1. 整体解读：综合分析占卜结果的核心信息
2. 具体建议：基于结果给出可操作的建议
3. 注意事项：需要警惕或关注的方面
4. 积极引导：用温暖鼓励的语言给予用户力量

要求：
- 语言温暖亲切，避免过于玄学的表述
- 结合用户的具体问题进行解读
- 给出实际可操作的建议
- 字数控制在 300-500 字
```

**接入方式**：
- 支持多种大模型 API（OpenAI、文心一言、通义千问等）
- 通过配置文件切换模型
- 流式输出（SSE），逐字显示解读内容

### 3.4 错误处理与边界情况

| 场景 | 处理方式 |
|------|---------|
| 未登录用户查看历史记录 | 跳转登录页，登录后返回 |
| AI 服务不可用 | 显示基础解读，提示"AI解读暂不可用" |
| 八字输入农历闰月 | 自动转换，提示用户确认 |
| 塔罗牌重复抽取 | 严格洗牌算法，确保不重复 |
| 网络请求超时 | 加载状态 → 超时提示 → 重试按钮 |
| 大量并发请求 | 限流 + 队列，AI 解读设置超时 |

---

## 4. 技术架构

### 4.1 整体架构图

```
┌─────────────────────────────────────────────────────────┐
│                      用户浏览器                          │
│            React + TypeScript + TailwindCSS              │
└───────────────────────┬─────────────────────────────────┘
                        │ HTTP/HTTPS
                        ▼
┌─────────────────────────────────────────────────────────┐
│                    Nginx / 静态服务                       │
│              （前端静态文件 + 反向代理）                    │
└───────────────────────┬─────────────────────────────────┘
                        │ /api/*
                        ▼
┌─────────────────────────────────────────────────────────┐
│                    Go 后端服务                            │
│                 Gin + GORM + SQLite                      │
│                                                         │
│  ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │ 用户模块 │  │ 塔罗模块  │  │ 星座模块  │  │ 六爻模块  │ │
│  └─────────┘  └──────────┘  └──────────┘  └──────────┘ │
│  ┌─────────┐  ┌──────────┐  ┌──────────┐              │
│  │ 八字模块 │  │ AI解读    │  │ 历史记录  │              │
│  └─────────┘  └──────────┘  └──────────┘              │
└───────────────────────┬─────────────────────────────────┘
                        │
           ┌────────────┼────────────┐
           ▼            ▼            ▼
    ┌────────────┐ ┌─────────┐ ┌──────────┐
    │   SQLite   │ │ AI API  │ │ 静态资源  │
    │  数据库     │ │ 大模型   │ │ 牌面图片  │
    └────────────┘ └─────────┘ └──────────┘
```

### 4.2 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端语言 |
| Gin | v1.9+ | HTTP 框架 |
| GORM | v1.25+ | ORM 框架 |
| SQLite | 3.x | 数据库 |
| JWT | - | 身份认证 |
| golang-jwt | v5 | JWT 库 |
| go-playground/validator | v10 | 请求验证 |
| zerolog | v1 | 日志库 |

### 4.3 前端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| React | 18+ | UI 框架 |
| TypeScript | 5+ | 类型安全 |
| Vite | 5+ | 构建工具 |
| TailwindCSS | 3+ | 样式框架 |
| React Router | 6+ | 路由管理 |
| Zustand | - | 状态管理 |
| Axios | - | HTTP 客户端 |
| Framer Motion | - | 动画库 |
| Recharts | - | 图表库（五行分析） |
| React Query | - | 服务端状态管理 |

### 4.4 API 设计

#### 4.4.1 认证相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 否 |
| POST | `/api/auth/login` | 用户登录 | 否 |
| POST | `/api/auth/refresh` | 刷新 Token | 是 |
| GET | `/api/auth/profile` | 获取用户信息 | 是 |
| PUT | `/api/auth/profile` | 更新用户信息 | 是 |

**POST /api/auth/register**

请求：
```json
{
  "username": "string",      // 3-20字符
  "email": "string",         // 邮箱格式
  "password": "string"       // 6-50字符
}
```

响应（201）：
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "created_at": "2026-05-30T16:00:00Z"
  }
}
```

**POST /api/auth/login**

请求：
```json
{
  "email": "string",
  "password": "string"
}
```

响应（200）：
```json
{
  "code": 0,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "zhangsan",
      "email": "zhangsan@example.com"
    }
  }
}
```

#### 4.4.2 塔罗牌相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/tarot/cards` | 获取所有牌数据 | 否 |
| GET | `/api/tarot/cards/:id` | 获取单张牌详情 | 否 |
| GET | `/api/tarot/spreads` | 获取可用牌阵列表 | 否 |
| POST | `/api/tarot/draw` | 抽牌 | 否 |

**POST /api/tarot/draw**

请求：
```json
{
  "spread": "celtic",       // single | three | celtic | love
  "question": "string"      // 可选，用户问题
}
```

响应（200）：
```json
{
  "code": 0,
  "data": {
    "spread": "celtic",
    "question": "我的事业发展如何？",
    "cards": [
      {
        "position": 1,
        "position_name": "当前状况",
        "card": {
          "id": 7,
          "name": "战车",
          "name_en": "The Chariot",
          "image": "/assets/tarot/major/07.jpg",
          "orientation": "upright",
          "keywords": ["胜利", "意志", "行动"],
          "meaning": "战车牌正位代表..."
        }
      }
    ],
    "timestamp": "2026-05-30T16:00:00Z"
  }
}
```

#### 4.4.3 星座运势相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/horoscope/:zodiac` | 获取星座运势 | 否 |
| GET | `/api/horoscope/:zodiac/history` | 获取运势历史 | 否 |

**GET /api/horoscope/:zodiac?period=daily&date=2026-05-30**

参数：
- `zodiac`: 星座名称（aries, taurus, ...）
- `period`: daily | weekly | monthly
- `date`: 日期（可选，默认今天）

响应（200）：
```json
{
  "code": 0,
  "data": {
    "zodiac": "aries",
    "zodiac_cn": "白羊座",
    "period": "daily",
    "date": "2026-05-30",
    "overall": 4,
    "love": 3,
    "career": 5,
    "wealth": 3,
    "health": 4,
    "lucky_number": 7,
    "lucky_color": "红色",
    "summary": "今天适合开展新项目...",
    "detail": {
      "love": "感情方面运势平稳...",
      "career": "事业运旺盛...",
      "wealth": "财运一般...",
      "health": "注意休息..."
    }
  }
}
```

#### 4.4.4 周易六爻相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/liuyao/hexagrams` | 获取64卦列表 | 否 |
| GET | `/api/liuyao/hexagrams/:id` | 获取单卦详情 | 否 |
| POST | `/api/liuyao/throw` | 掷铜钱占卜 | 否 |

**POST /api/liuyao/throw**

请求：
```json
{
  "question": "string"      // 可选，用户问题
}
```

响应（200）：
```json
{
  "code": 0,
  "data": {
    "question": "这次投资是否可行？",
    "lines": [
      { "value": 7, "type": "young_yang", "mutable": false, "symbol": "⚊" },
      { "value": 9, "type": "old_yin", "mutable": true, "symbol": "⚋" },
      { "value": 8, "type": "young_yin", "mutable": false, "symbol": "⚋" },
      { "value": 7, "type": "young_yang", "mutable": false, "symbol": "⚊" },
      { "value": 8, "type": "young_yin", "mutable": false, "symbol": "⚋" },
      { "value": 7, "type": "young_yang", "mutable": false, "symbol": "⚊" }
    ],
    "ben_gua": {
      "id": 12,
      "name": "天地否",
      "upper": "乾",
      "lower": "坤",
      "judgment": "否之匪人，不利君子贞...",
      "image": "天地不交，否..."
    },
    "bian_gua": {
      "id": 45,
      "name": "泽地萃",
      "upper": "兑",
      "lower": "坤",
      "judgment": "萃，亨..."
    },
    "mutable_lines": [1],
    "timestamp": "2026-05-30T16:00:00Z"
  }
}
```

#### 4.4.5 八字排盘相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/bazi/calculate` | 八字排盘计算 | 否 |

**POST /api/bazi/calculate**

请求：
```json
{
  "birth_date": "1990-05-15",    // 公历日期
  "birth_time": "14:30",         // 24小时制
  "gender": "male"               // male | female
}
```

响应（200）：
```json
{
  "code": 0,
  "data": {
    "birth": {
      "solar": "1990-05-15 14:30",
      "lunar": "庚午年 四月廿一 未时"
    },
    "pillars": {
      "year":  { "tian_gan": "庚", "di_zhi": "午", "wu_xing": "金火", "na_yin": "路旁土" },
      "month": { "tian_gan": "辛", "di_zhi": "巳", "wu_xing": "金火", "na_yin": "白蜡金" },
      "day":   { "tian_gan": "丙", "di_zhi": "寅", "wu_xing": "火木", "na_yin": "炉中火" },
      "hour":  { "tian_gan": "乙", "di_zhi": "未", "wu_xing": "木土", "na_yin": "沙中金" }
    },
    "five_elements": {
      "metal": 2,
      "wood": 2,
      "water": 0,
      "fire": 3,
      "earth": 1,
      "strongest": "火",
      "weakest": "水",
      "day_master": "丙火",
      "strength": "偏弱",
      "yong_shen": "木",
      "ji_shen": "水"
    },
    "ten_gods": [
      { "position": "年干", "tian_gan": "庚", "god": "偏财" },
      { "position": "月干", "tian_gan": "辛", "god": "正财" },
      { "position": "时干", "tian_gan": "乙", "god": "正印" }
    ],
    "da_yun": [
      { "start_age": 5, "end_age": 14, "tian_gan": "壬", "di_zhi": "午", "period": "1995-2004" },
      { "start_age": 15, "end_age": 24, "tian_gan": "癸", "di_zhi": "未", "period": "2005-2014" }
    ]
  }
}
```

#### 4.4.6 AI 解读相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/ai/interpret` | AI 解读占卜结果 | 是（限流） |

**POST /api/ai/interpret**

请求：
```json
{
  "type": "tarot",           // tarot | horoscope | liuyao | bazi
  "result_id": 123,          // 占卜记录 ID
  "question": "string"       // 可选，补充问题
}
```

响应（200，SSE 流式）：
```
data: {"text": "根据"}
data: {"text": "您的塔罗牌"}
data: {"text": "占卜结果..."}
data: [DONE]
```

#### 4.4.7 历史记录相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/history` | 获取占卜历史列表 | 是 |
| GET | `/api/history/:id` | 获取单条记录详情 | 是 |
| DELETE | `/api/history/:id` | 删除占卜记录 | 是 |

**GET /api/history?type=tarot&page=1&page_size=10**

响应（200）：
```json
{
  "code": 0,
  "data": {
    "total": 42,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 123,
        "type": "tarot",
        "type_cn": "塔罗牌",
        "spread": "celtic",
        "question": "我的事业发展如何？",
        "summary": "战车牌正位，事业运势旺盛...",
        "created_at": "2026-05-30T16:00:00Z"
      }
    ]
  }
}
```

### 4.5 数据模型

#### 4.5.1 用户表（users）

```sql
CREATE TABLE users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,               -- bcrypt 加密
    avatar        TEXT DEFAULT '',              -- 头像 URL
    zodiac        TEXT DEFAULT '',              -- 星座
    birth_date    TEXT DEFAULT '',              -- 生日
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

#### 4.5.2 占卜记录表（divination_records）

```sql
CREATE TABLE divination_records (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,            -- 关联用户
    type          TEXT NOT NULL,                -- tarot | horoscope | liuyao | bazi
    question      TEXT DEFAULT '',              -- 用户问题
    result        TEXT NOT NULL,                -- JSON 格式的占卜结果
    ai_reading    TEXT DEFAULT '',              -- AI 解读内容
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_records_user_id ON divination_records(user_id);
CREATE INDEX idx_records_type ON divination_records(type);
CREATE INDEX idx_records_created_at ON divination_records(created_at);
```

#### 4.5.3 塔罗牌数据表（tarot_cards）

```sql
CREATE TABLE tarot_cards (
    id            INTEGER PRIMARY KEY,
    name          TEXT NOT NULL,                -- 中文名
    name_en       TEXT NOT NULL,                -- 英文名
    type          TEXT NOT NULL,                -- major | minor
    suit          TEXT DEFAULT '',              -- 花色：wands/cups/swords/pentacles
    number        INTEGER DEFAULT 0,            -- 牌号
    image         TEXT NOT NULL,                -- 图片路径
    keywords_up   TEXT NOT NULL,                -- 正位关键词（JSON数组）
    keywords_down TEXT NOT NULL,                -- 逆位关键词（JSON数组）
    meaning_up    TEXT NOT NULL,                -- 正位含义
    meaning_down  TEXT NOT NULL,                -- 逆位含义
    description   TEXT DEFAULT ''               -- 牌面描述
);
```

#### 4.5.4 六十四卦数据表（hexagrams）

```sql
CREATE TABLE hexagrams (
    id            INTEGER PRIMARY KEY,          -- 1-64
    name          TEXT NOT NULL,                -- 卦名（如"乾为天"）
    name_short    TEXT NOT NULL,                -- 简称（如"乾"）
    upper_trigram TEXT NOT NULL,                -- 上卦
    lower_trigram TEXT NOT NULL,                -- 下卦
    binary        TEXT NOT NULL,                -- 6位二进制
    judgment      TEXT NOT NULL,                -- 卦辞
    image         TEXT NOT NULL,                -- 象辞
    line_texts    TEXT NOT NULL,                -- 六爻爻辞（JSON数组）
    description   TEXT DEFAULT ''               -- 卦的总体描述
);

CREATE TABLE trigrams (
    id            INTEGER PRIMARY KEY,          -- 0-7
    name          TEXT NOT NULL,                -- 卦名
    symbol        TEXT NOT NULL,                -- 符号
    nature        TEXT NOT NULL,                -- 自象
    element       TEXT NOT NULL,                -- 五行
    yin_yang      TEXT NOT NULL                 -- 阴阳
);
```

#### 4.5.5 星座运势模板表（horoscope_templates）

```sql
CREATE TABLE horoscope_templates (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    zodiac        TEXT NOT NULL,                -- 星座
    period        TEXT NOT NULL,                -- daily | weekly | monthly
    dimension     TEXT NOT NULL,                -- love | career | wealth | health
    score_level   TEXT NOT NULL,                -- high | mid | low
    template      TEXT NOT NULL,                -- 模板文本
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_lookup ON horoscope_templates(zodiac, period, dimension, score_level);
```

#### 4.5.6 ER 关系图

```
┌─────────────┐       ┌─────────────────────┐
│   users     │       │  divination_records │
├─────────────┤       ├─────────────────────┤
│ id (PK)     │──┐    │ id (PK)             │
│ username    │  │    │ user_id (FK)        │──┐
│ email       │  └───>│ type                │  │
│ password    │       │ question            │  │
│ avatar      │       │ result (JSON)       │  │
│ zodiac      │       │ ai_reading          │  │
│ birth_date  │       │ created_at          │  │
│ created_at  │       └─────────────────────┘  │
│ updated_at  │                                │
└─────────────┘                                │
                                               │
┌─────────────┐       ┌─────────────────────┐ │
│ tarot_cards │       │  horoscope_templates│ │
├─────────────┤       ├─────────────────────┤ │
│ id (PK)     │       │ id (PK)             │ │
│ name        │       │ zodiac              │ │
│ name_en     │       │ period              │ │
│ type        │       │ dimension           │ │
│ suit        │       │ score_level         │ │
│ number      │       │ template            │ │
│ image       │       │ created_at          │ │
│ meaning_up  │       └─────────────────────┘ │
│ meaning_down│                                │
└─────────────┘                                │
                                               │
┌─────────────┐       ┌─────────────────────┐ │
│ hexagrams   │       │    trigrams         │ │
├─────────────┤       ├─────────────────────┤ │
│ id (PK)     │       │ id (PK)             │ │
│ name        │       │ name                │ │
│ upper       │       │ symbol              │ │
│ lower       │       │ nature              │ │
│ binary      │       │ element             │ │
│ judgment    │       │ yin_yang            │ │
│ image       │       └─────────────────────┘ │
│ line_texts  │                                │
└─────────────┘                                │
```

### 4.6 安全设计

#### 4.6.1 认证与授权

- **JWT 认证**：Access Token（1小时）+ Refresh Token（7天）
- **密码加密**：bcrypt，cost factor = 12
- **Token 存储**：前端 localStorage（考虑 HttpOnly Cookie 更安全，但复杂度增加）

#### 4.6.2 接口安全

- **CORS**：仅允许前端域名访问
- **限流**：AI 解读接口限流（每用户每分钟 5 次）
- **输入验证**：所有请求参数严格验证
- **SQL 注入**：使用 GORM 参数化查询

#### 4.6.3 数据安全

- **敏感数据**：密码 bcrypt 加密，不在日志中记录
- **SQLite 加密**：可选 SQLCipher 加密数据库文件
- **备份策略**：定期备份 SQLite 文件

### 4.7 配置管理

```yaml
# config.yaml
server:
  port: 8080
  mode: debug  # debug | release

database:
  path: ./data/zhanbu.db

jwt:
  secret: "your-secret-key-change-in-production"
  access_ttl: 1h
  refresh_ttl: 168h  # 7 days

ai:
  provider: openai  # openai | zhipu | qwen
  api_key: "sk-..."
  model: "gpt-4"
  base_url: "https://api.openai.com/v1"
  max_tokens: 1000
  temperature: 0.7

rate_limit:
  ai_per_minute: 5
  api_per_minute: 60

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "https://zhanbu.example.com"
```

---

## 5. 目录结构

### 5.1 后端目录结构

```
zhanbu/
├── server/                          # 后端代码
│   ├── main.go                      # 程序入口
│   ├── go.mod                       # Go 模块定义
│   ├── go.sum                       # 依赖校验
│   │
│   ├── config/                      # 配置
│   │   ├── config.go                # 配置结构体定义
│   │   └── config.yaml              # 配置文件
│   │
│   ├── internal/                    # 内部代码（不对外暴露）
│   │   ├── middleware/              # 中间件
│   │   │   ├── auth.go             # JWT 认证中间件
│   │   │   ├── cors.go             # CORS 中间件
│   │   │   ├── logger.go           # 日志中间件
│   │   │   └── ratelimit.go        # 限流中间件
│   │   │
│   │   ├── handler/                 # HTTP 处理器
│   │   │   ├── auth.go             # 认证处理器
│   │   │   ├── tarot.go            # 塔罗牌处理器
│   │   │   ├── horoscope.go        # 星座运势处理器
│   │   │   ├── liuyao.go           # 六爻处理器
│   │   │   ├── bazi.go             # 八字处理器
│   │   │   ├── ai.go               # AI 解读处理器
│   │   │   └── history.go          # 历史记录处理器
│   │   │
│   │   ├── service/                 # 业务逻辑层
│   │   │   ├── auth.go             # 认证服务
│   │   │   ├── tarot.go            # 塔罗牌服务（洗牌、抽牌逻辑）
│   │   │   ├── horoscope.go        # 星座运势服务（运势生成算法）
│   │   │   ├── liuyao.go           # 六爻服务（掷铜钱、卦象解析）
│   │   │   ├── bazi.go             # 八字服务（排盘计算）
│   │   │   ├── ai.go               # AI 解读服务
│   │   │   └── history.go          # 历史记录服务
│   │   │
│   │   ├── model/                   # 数据模型
│   │   │   ├── user.go             # 用户模型
│   │   │   ├── divination.go       # 占卜记录模型
│   │   │   ├── tarot.go            # 塔罗牌模型
│   │   │   ├── hexagram.go         # 卦象模型
│   │   │   └── horoscope.go        # 星座运势模型
│   │   │
│   │   ├── repository/              # 数据访问层
│   │   │   ├── user.go             # 用户数据访问
│   │   │   ├── divination.go       # 占卜记录数据访问
│   │   │   ├── tarot.go            # 塔罗牌数据访问
│   │   │   ├── hexagram.go         # 卦象数据访问
│   │   │   └── horoscope.go        # 星座运势数据访问
│   │   │
│   │   ├── database/                # 数据库
│   │   │   ├── db.go               # 数据库连接
│   │   │   ├── migrate.go          # 数据库迁移
│   │   │   └── seed.go             # 种子数据（64卦、78张牌等）
│   │   │
│   │   └── router/                  # 路由
│   │       └── router.go           # 路由注册
│   │
│   ├── pkg/                         # 公共包
│   │   ├── response/               # 统一响应格式
│   │   │   └── response.go
│   │   ├── errors/                  # 错误定义
│   │   │   └── errors.go
│   │   └── utils/                   # 工具函数
│   │       ├── jwt.go              # JWT 工具
│   │       ├── hash.go             # 密码哈希工具
│   │       └── validator.go        # 验证工具
│   │
│   └── data/                        # 数据文件
│       ├── tarot/                   # 塔罗牌数据
│       │   ├── major_arcana.json   # 大阿尔卡纳数据
│       │   └── minor_arcana.json   # 小阿尔卡纳数据
│       ├── hexagrams/               # 六十四卦数据
│       │   └── hexagrams.json      # 卦辞、爻辞数据
│       └── horoscope/               # 星座运势模板
│           ├── daily.json
│           ├── weekly.json
│           └── monthly.json
│
├── docs/                            # 文档
│   └── design.md                    # 本文档
│
└── README.md                        # 项目说明
```

### 5.2 前端目录结构

```
zhanbu/
├── client/                          # 前端代码
│   ├── index.html                   # HTML 入口
│   ├── package.json                 # 依赖配置
│   ├── tsconfig.json                # TypeScript 配置
│   ├── vite.config.ts               # Vite 配置
│   ├── tailwind.config.js           # TailwindCSS 配置
│   ├── postcss.config.js            # PostCSS 配置
│   │
│   ├── public/                      # 静态资源
│   │   ├── favicon.ico
│   │   └── assets/
│   │       ├── tarot/              # 塔罗牌图片
│   │       │   ├── major/          # 大阿尔卡纳图片
│   │       │   │   ├── 00_fool.jpg
│   │       │   │   ├── 01_magician.jpg
│   │       │   │   └── ...
│   │       │   └── minor/          # 小阿尔卡纳图片
│   │       │       ├── wands/
│   │       │       ├── cups/
│   │       │       ├── swords/
│   │       │       └── pentacles/
│   │       ├── coins/              # 铜钱图片
│   │       ├── zodiac/             # 星座图标
│   │       └── trigrams/           # 八卦符号
│   │
│   └── src/                         # 源代码
│       ├── main.tsx                 # React 入口
│       ├── App.tsx                  # 根组件
│       │
│       ├── components/              # 可复用组件
│       │   ├── common/             # 通用组件
│       │   │   ├── Header.tsx
│       │   │   ├── Footer.tsx
│       │   │   ├── Loading.tsx
│       │   │   ├── ErrorBoundary.tsx
│       │   │   └── Modal.tsx
│       │   │
│       │   ├── tarot/              # 塔罗牌组件
│       │   │   ├── TarotCard.tsx           # 单张牌组件
│       │   │   ├── TarotDeck.tsx           # 牌组组件
│       │   │   ├── SpreadLayout.tsx        # 牌阵布局
│       │   │   ├── CelticCross.tsx         # 凯尔特十字布局
│       │   │   ├── DrawAnimation.tsx       # 抽牌动画
│       │   │   └── CardReading.tsx         # 牌面解读
│       │   │
│       │   ├── horoscope/          # 星座组件
│       │   │   ├── ZodiacWheel.tsx         # 12宫轮盘
│       │   │   ├── HoroscopeCard.tsx       # 运势卡片
│       │   │   ├── RadarChart.tsx          # 雷达图
│       │   │   └── LuckyElements.tsx       # 幸运元素
│       │   │
│       │   ├── liuyao/             # 六爻组件
│       │   │   ├── CoinAnimation.tsx       # 铜钱动画
│       │   │   ├── HexagramDisplay.tsx     # 卦象展示
│       │   │   ├── LineResult.tsx          # 单爻结果
│       │   │   └── TrigramChart.tsx        # 八卦图表
│       │   │
│       │   └── bazi/               # 八字组件
│       │       ├── PillarTable.tsx         # 四柱表格
│       │       ├── FiveElementChart.tsx    # 五行图表
│       │       ├── DaYunTimeline.tsx       # 大运时间轴
│       │       └── TenGodList.tsx          # 十神列表
│       │
│       ├── pages/                   # 页面组件
│       │   ├── Home.tsx            # 首页
│       │   ├── Tarot.tsx           # 塔罗牌占卜页
│       │   ├── Horoscope.tsx       # 星座运势页
│       │   ├── LiuYao.tsx          # 六爻占卜页
│       │   ├── BaZi.tsx            # 八字排盘页
│       │   ├── History.tsx         # 历史记录页
│       │   ├── Profile.tsx         # 用户中心
│       │   ├── Login.tsx           # 登录页
│       │   └── Register.tsx        # 注册页
│       │
│       ├── hooks/                   # 自定义 Hooks
│       │   ├── useAuth.ts          # 认证 Hook
│       │   ├── useTarot.ts         # 塔罗牌 Hook
│       │   ├── useHoroscope.ts     # 星座运势 Hook
│       │   ├── useLiuYao.ts        # 六爻 Hook
│       │   ├── useBaZi.ts          # 八字 Hook
│       │   └── useAI.ts            # AI 解读 Hook
│       │
│       ├── services/                # API 服务
│       │   ├── api.ts              # Axios 实例配置
│       │   ├── auth.ts             # 认证 API
│       │   ├── tarot.ts            # 塔罗牌 API
│       │   ├── horoscope.ts        # 星座运势 API
│       │   ├── liuyao.ts           # 六爻 API
│       │   ├── bazi.ts             # 八字 API
│       │   ├── ai.ts               # AI 解读 API
│       │   └── history.ts          # 历史记录 API
│       │
│       ├── stores/                  # 状态管理
│       │   ├── authStore.ts        # 认证状态
│       │   └── uiStore.ts          # UI 状态
│       │
│       ├── types/                   # TypeScript 类型
│       │   ├── user.ts             # 用户类型
│       │   ├── tarot.ts            # 塔罗牌类型
│       │   ├── horoscope.ts        # 星座运势类型
│       │   ├── liuyao.ts           # 六爻类型
│       │   ├── bazi.ts             # 八字类型
│       │   └── api.ts              # API 响应类型
│       │
│       ├── utils/                   # 工具函数
│       │   ├── format.ts           # 格式化工具
│       │   ├── date.ts             # 日期工具
│       │   └── storage.ts          # 本地存储工具
│       │
│       └── styles/                  # 样式
│           ├── globals.css         # 全局样式
│           └── animations.css      # 动画样式
│
└── README.md                        # 项目说明
```

---

## 6. 开发计划

### Phase 1：基础框架与核心功能（4周）

**目标**：搭建项目框架，实现塔罗牌占卜基础功能

#### Week 1：项目初始化

| 任务 | 说明 | 产出 |
|------|------|------|
| 1.1 | 初始化 Go 后端项目（Gin + GORM） | 可运行的 HTTP 服务 |
| 1.2 | 初始化 React 前端项目（Vite + TS + Tailwind） | 可运行的前端应用 |
| 1.3 | 设计并创建数据库表结构 | SQLite 数据库 + 迁移脚本 |
| 1.4 | 实现 JWT 认证中间件 | 认证模块可用 |
| 1.5 | 实现用户注册/登录 API | 完整的认证流程 |
| 1.6 | 前端登录/注册页面 | 可用的登录注册 UI |

#### Week 2：塔罗牌后端

| 任务 | 说明 | 产出 |
|------|------|------|
| 2.1 | 准备塔罗牌数据（78张牌的 JSON） | 塔罗牌数据文件 |
| 2.2 | 实现 Fisher-Yates 洗牌算法 | 可靠的随机抽牌 |
| 2.3 | 实现牌阵逻辑（单牌、三牌、凯尔特十字） | 4种牌阵可用 |
| 2.4 | 实现塔罗牌 API（/api/tarot/*） | 完整的塔罗牌 API |
| 2.5 | 编写单元测试 | 测试覆盖率 > 80% |

#### Week 3：塔罗牌前端

| 任务 | 说明 | 产出 |
|------|------|------|
| 3.1 | 设计首页 UI（占卜类型卡片） | 首页完成 |
| 3.2 | 实现塔罗牌占卜页面 | 占卜流程可用 |
| 3.3 | 实现抽牌动画（Framer Motion） | 流畅的动画效果 |
| 3.4 | 实现牌阵布局组件 | 4种牌阵布局 |
| 3.5 | 实现牌面解读展示 | 解读 UI 完成 |

#### Week 4：历史记录与联调

| 任务 | 说明 | 产出 |
|------|------|------|
| 4.1 | 实现占卜记录保存 API | 记录保存可用 |
| 4.2 | 实现历史记录页面 | 历史记录 UI |
| 4.3 | 前后端联调 | 完整的塔罗牌占卜流程 |
| 4.4 | UI 细节优化与 Bug 修复 | 稳定的可用版本 |
| 4.5 | 部署到测试环境 | 可访问的测试版本 |

**Phase 1 交付物**：
- ✅ 可运行的塔罗牌占卜网站
- ✅ 用户注册/登录
- ✅ 4种牌阵（单牌、三牌、凯尔特十字、爱情十字）
- ✅ 抽牌动画
- ✅ 历史记录

---

### Phase 2：更多占卜类型（4周）

**目标**：实现星座运势、周易六爻、八字排盘

#### Week 5-6：星座运势

| 任务 | 说明 | 产出 |
|------|------|------|
| 5.1 | 设计运势生成算法 | 运势生成服务 |
| 5.2 | 准备运势模板库 | 运势模板数据 |
| 5.3 | 实现星座运势 API | 完整的运势 API |
| 5.4 | 实现 12 宫轮盘 UI | 星座选择组件 |
| 5.5 | 实现运势展示页面 | 运势 UI |
| 5.6 | 实现运势雷达图 | 图表可视化 |

#### Week 7-8：周易六爻与八字

| 任务 | 说明 | 产出 |
|------|------|------|
| 7.1 | 准备 64 卦数据（卦辞、爻辞） | 六十四卦数据 |
| 7.2 | 实现掷铜钱算法 | 六爻服务 |
| 7.3 | 实现六爻 API | 完整的六爻 API |
| 7.4 | 实现铜钱动画 | 掷铜钱 UI |
| 7.5 | 实现八字排盘算法 | 八字服务（含节气计算） |
| 7.6 | 实现八字 API | 完整的八字 API |
| 7.7 | 实现八字排盘 UI | 四柱展示、五行图表 |
| 7.8 | 各类型联调与测试 | 稳定的多类型占卜 |

**Phase 2 交付物**：
- ✅ 星座运势（12星座 × 日/周/月）
- ✅ 周易六爻（64卦完整数据）
- ✅ 八字排盘（天干地支、五行分析、大运）
- ✅ 各类型的完整 UI 和动画

---

### Phase 3：AI 集成与优化（3周）

**目标**：接入 AI 大模型，优化用户体验

#### Week 9-10：AI 解读

| 任务 | 说明 | 产出 |
|------|------|------|
| 9.1 | 设计 AI Prompt 模板 | 各类型的 Prompt 模板 |
| 9.2 | 实现 AI 服务抽象层 | 支持多种大模型 |
| 9.3 | 实现流式输出（SSE） | 逐字显示 AI 解读 |
| 9.4 | 前端 AI 解读 UI | 流式文本展示 |
| 9.5 | 实现限流和错误处理 | 稳定的 AI 服务 |

#### Week 11：优化与完善

| 任务 | 说明 | 产出 |
|------|------|------|
| 10.1 | UI/UX 全面优化 | 更精美的界面 |
| 10.2 | 性能优化（缓存、懒加载） | 更快的响应速度 |
| 10.3 | 移动端适配 | 响应式设计 |
| 10.4 | 安全审计 | 安全漏洞修复 |
| 10.5 | 编写用户文档 | 使用说明 |

#### Week 12：上线准备

| 任务 | 说明 | 产出 |
|------|------|------|
| 11.1 | 准备生产环境配置 | 生产配置 |
| 11.2 | 部署到生产服务器 | 可访问的线上版本 |
| 11.3 | 监控与日志 | 监控系统 |
| 11.4 | 压力测试 | 性能报告 |
| 11.5 | 正式上线 | 🎉 上线 |

**Phase 3 交付物**：
- ✅ AI 大模型解读（支持多种模型）
- ✅ 流式输出
- ✅ 移动端适配
- ✅ 生产环境部署
- ✅ 监控与日志

---

## 附录

### A. 本地开发环境搭建

#### 后端

```bash
# 进入后端目录
cd server

# 安装依赖
go mod tidy

# 初始化数据库
go run main.go --migrate

# 导入种子数据
go run main.go --seed

# 启动服务
go run main.go
# 服务启动在 http://localhost:8080
```

#### 前端

```bash
# 进入前端目录
cd client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
# 前端启动在 http://localhost:5173
```

### B. API 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 未授权 |
| 1003 | 禁止访问 |
| 1004 | 资源不存在 |
| 1005 | 请求过于频繁 |
| 2001 | 用户已存在 |
| 2002 | 用户名或密码错误 |
| 2003 | Token 已过期 |
| 3001 | AI 服务不可用 |
| 3002 | AI 服务超时 |

### C. 参考资料

- [塔罗牌完整资料](https://www.tarot.com/)
- [周易六十四卦](https://www.iching.net/)
- [万年历算法](https://github.com/nicholasgasior/gova-chinese-calendar)
- [八字排盘算法](https://github.com/lxfriday/bazi)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [React 文档](https://react.dev/)
- [TailwindCSS 文档](https://tailwindcss.com/docs)

---

> **文档维护说明**：本文档应随项目迭代持续更新。重大架构变更需同步更新此文档。
