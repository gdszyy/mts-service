# MTS 前端交互逻辑设计 v4 (专家模式增强版)

**版本**: 4.0  
**日期**: 2025-12-01  
**作者**: Manus AI

---

## 1. 概述

本文档是 MTS 前端交互逻辑的第四次迭代，在 v3 的精细化设计基础上，采纳了更深层次的用户行为洞察，引入了 **“新手/老手模式”** 的概念。这一核心增强旨在为不同经验水平的用户提供截然不同的、但都同样流畅的投注体验，从而实现真正的个性化服务。v4 版本专注于解决 **互斥选项处理** 和 **异常选项清理** 这两大场景下新手与老手用户的核心诉求差异。

---

## 2. 核心交互逻辑 (继承 v3)

核心交互依然遵循 v3 版本的设计，包括：

-   **“单关/串关”双模式切换**
-   **模式切换时的“清空并确认”策略**
-   **可配置的快捷金额选项**
-   **全面的设置页面**
-   **精细化的异常情况处理**

---

## 3. 专家级功能增强 (v4 新增)

### 3.1. 新手/老手模式 (User Mode)

这是 v4 版本的核心。在“设置页面”中，我们新增一个 **“用户模式”** 选项，允许用户在“新手模式”和“老手模式”之间选择。此设置将从根本上改变系统处理互斥选项和异常选项的方式。

| 模式 | 目标用户 | 核心特点 |
| :--- | :--- | :--- |
| **新手模式 (默认)** | 新用户、休闲玩家 | **便捷、防错**：自动处理冲突，简化决策流程。 |
| **老手模式** | 专业玩家、高频用户 | **灵活、可控**：允许用户观察和手动处理复杂情况。 |

### 3.2. 互斥选项处理

当用户在同一场比赛中选择了互斥的选项时（例如，同时选择主队胜和平局），系统将根据用户模式采取不同策略。

```mermaid
flowchart TD
    Start([用户选择新的选项]) --> CheckUserMode{检查用户模式}
    
    CheckUserMode -->|新手模式| CheckConflict1{检查是否与<br/>现有选项冲突}
    CheckUserMode -->|老手模式| AllowMultiple[允许添加互斥选项<br/>用于观察和比对]
    
    CheckConflict1 -->|有冲突| AutoReplace[自动替换冲突选项<br/>保留新选择]
    CheckConflict1 -->|无冲突| AddSelection1[添加选项到投注单]
    
    AutoReplace --> ShowNotification[显示通知<br/>已自动替换冲突选项]
    ShowNotification --> AddSelection1
    
    AllowMultiple --> AddSelection2[添加选项到投注单<br/>标记互斥关系]
    AddSelection2 --> MarkConflict[将互斥选项标记为<br/>橙色边框或警告图标]
    MarkConflict --> ShowWarning[在投注单顶部显示警告<br/>存在互斥选项，请手动移除]
    
    AddSelection1 --> UpdateBetSlip1[更新投注单显示]
    ShowWarning --> UpdateBetSlip2[更新投注单显示]
    
    UpdateBetSlip1 --> ValidateForSubmit1[验证投注单有效性]
    UpdateBetSlip2 --> ValidateForSubmit2[验证投注单有效性]
    
    ValidateForSubmit1 --> EnableSubmit[启用下单按钮]
    ValidateForSubmit2 --> CheckConflictStatus{是否存在<br/>互斥选项?}
    
    CheckConflictStatus -->|是| DisableSubmit[禁用下单按钮<br/>显示: 请移除互斥选项]
    CheckConflictStatus -->|否| EnableSubmit
    
    EnableSubmit --> End([等待用户下单])
    DisableSubmit --> WaitUserAction[等待用户手动移除]
    WaitUserAction --> UserRemoves[用户移除互斥选项]
    UserRemoves --> UpdateBetSlip2
    
    style Start fill:#4CAF50,stroke:#2E7D32,color:#fff
    style End fill:#F44336,stroke:#C62828,color:#fff
    style AutoReplace fill:#E91E63,stroke:#880E4F,color:#fff
    style AllowMultiple fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style MarkConflict fill:#FF9800,stroke:#E65100,color:#fff
    style ShowWarning fill:#FF9800,stroke:#E65100,color:#fff
    style ShowNotification fill:#2196F3,stroke:#1565C0,color:#fff
    style EnableSubmit fill:#4CAF50,stroke:#2E7D32,color:#fff
    style DisableSubmit fill:#F44336,stroke:#C62828,color:#fff
    style CheckUserMode fill:#FFC107,stroke:#F57C00,color:#000
    style CheckConflict1 fill:#FFC107,stroke:#F57C00,color:#000
    style CheckConflictStatus fill:#FFC107,stroke:#F57C00,color:#000
```

**图 1：互斥选项处理流程图**

-   **新手模式**：
    -   **自动替换**：当用户选择一个与已有选项互斥的新选项时，系统会自动用新选项替换旧选项。
    -   **用户提示**：同时，界面会显示一个简短的通知（例如，“已自动替换冲突选项”），让用户了解发生了什么。

-   **老手模式**：
    -   **允许多选**：系统允许用户将互斥选项同时保留在投注单中。
    -   **视觉警告**：所有互斥的选项都会被高亮标记（例如，橙色边框或警告图标），并在投注单顶部显示警告信息：“存在互斥选项，无法投注。”
    -   **下单禁用**：只要投注单中存在互斥选项，“下单”按钮就会被禁用，强制要求用户手动移除冲突项后才能继续。

### 3.3. 智能异常选项处理

当投注单中出现锁定的盘口、失效的选项或互斥的串关时，系统将根据用户模式提供不同的清理机制。

```mermaid
flowchart TD
    Start([投注单中存在异常选项]) --> DetectException[检测异常类型<br/>锁定盘口/失效选项/冲突选项]
    
    DetectException --> MarkException[将异常选项标记<br/>灰色/红色边框/警告图标]
    MarkException --> CheckUserMode{检查用户模式}
    
    CheckUserMode -->|新手模式| NewbieMode[新手模式处理]
    CheckUserMode -->|老手模式| ExpertMode[老手模式处理]
    
    NewbieMode --> CheckBtnArea{用户点击区域}
    CheckBtnArea -->|点击下单按钮| AutoCleanup[自动移除所有异常选项]
    CheckBtnArea -->|点击异常选项| ManualRemove1[允许手动移除单个选项]
    
    AutoCleanup --> ShowCleanupMsg[显示通知<br/>已自动移除X个异常选项]
    ShowCleanupMsg --> RevalidateBet1[重新验证投注单]
    
    ManualRemove1 --> UpdateBetSlip1[更新投注单显示]
    UpdateBetSlip1 --> RevalidateBet1
    
    RevalidateBet1 --> CheckValid1{投注单是否有效?}
    CheckValid1 -->|有效| EnableSubmit1[启用下单按钮<br/>允许提交]
    CheckValid1 -->|无效| ShowError1[显示错误提示<br/>请检查投注单]
    
    ExpertMode --> DisableSubmit[禁用下单按钮<br/>显示: 存在异常选项]
    DisableSubmit --> ShowExpertHint[在投注单顶部显示提示<br/>请手动移除异常选项后下单]
    ShowExpertHint --> WaitExpertAction[等待用户手动操作]
    
    WaitExpertAction --> ExpertAction{用户操作}
    ExpertAction -->|手动移除异常选项| ManualRemove2[移除选中的异常选项]
    ExpertAction -->|观察和比对| KeepObserving[保持异常选项<br/>继续观察赔率]
    
    ManualRemove2 --> UpdateBetSlip2[更新投注单显示]
    UpdateBetSlip2 --> RevalidateBet2[重新验证投注单]
    
    RevalidateBet2 --> CheckValid2{投注单是否有效?}
    CheckValid2 -->|有效| EnableSubmit2[启用下单按钮<br/>允许提交]
    CheckValid2 -->|无效| ShowExpertHint
    
    KeepObserving --> WaitExpertAction
    
    EnableSubmit1 --> End([用户可以下单])
    EnableSubmit2 --> End
    ShowError1 --> End
    
    style Start fill:#4CAF50,stroke:#2E7D32,color:#fff
    style End fill:#F44336,stroke:#C62828,color:#fff
    style NewbieMode fill:#2196F3,stroke:#1565C0,color:#fff
    style ExpertMode fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style AutoCleanup fill:#E91E63,stroke:#880E4F,color:#fff
    style ShowCleanupMsg fill:#00BCD4,stroke:#00838F,color:#fff
    style DisableSubmit fill:#F44336,stroke:#C62828,color:#fff
    style ShowExpertHint fill:#FF9800,stroke:#E65100,color:#fff
    style EnableSubmit1 fill:#4CAF50,stroke:#2E7D32,color:#fff
    style EnableSubmit2 fill:#4CAF50,stroke:#2E7D32,color:#fff
    style ShowError1 fill:#F44336,stroke:#C62828,color:#fff
    style MarkException fill:#FF9800,stroke:#E65100,color:#fff
    style CheckUserMode fill:#FFC107,stroke:#F57C00,color:#000
    style CheckBtnArea fill:#FFC107,stroke:#F57C00,color:#000
    style CheckValid1 fill:#FFC107,stroke:#F57C00,color:#000
    style CheckValid2 fill:#FFC107,stroke:#F57C00,color:#000
    style ExpertAction fill:#FFC107,stroke:#F57C00,color:#000
```

**图 2：智能异常选项处理流程图**

-   **新手模式**：
    -   **一键清理**：当投注单中存在异常选项时，“下单”按钮会变为 **“移除异常并下单”**。
    -   **用户操作**：用户点击该按钮后，系统会自动移除所有异常选项，并提交剩余的有效投注。
    -   **用户提示**：同时，界面会显示通知：“已自动移除 X 个异常选项并提交投注。”

-   **老手模式**：
    -   **手动清理**：系统会高亮所有异常选项，但不会提供自动清理功能。
    -   **下单禁用**：“下单”按钮会被禁用，并显示提示：“请手动移除所有异常选项。”
    -   **用户操作**：老手用户可以从容地观察、比对这些异常选项，然后根据自己的判断手动移除它们，直到投注单变为有效状态。

### 3.4. “使用上次金额”功能

为了进一步提升高频用户的投注效率，我们在快捷金额选项中新增了 **“Last” (上次)** 按钮。

```mermaid
flowchart TD
    Start([用户查看金额输入区域]) --> DisplayQuickBtns[显示快捷金额按钮<br/>包括: +5, +10, +20, +50, +100, Last, Max]
    
    DisplayQuickBtns --> WaitClick[等待用户操作]
    
    WaitClick --> UserAction{用户操作}
    
    UserAction -->|点击 Last 按钮| CheckHistory{检查历史记录}
    UserAction -->|点击其他按钮| OtherActions[执行其他快捷操作]
    
    CheckHistory -->|有历史记录| GetLastAmount[获取上一次成功投注的金额]
    CheckHistory -->|无历史记录| ShowNoHistory[显示提示: 暂无历史投注记录]
    
    GetLastAmount --> CheckMode{检查当前投注模式}
    
    CheckMode -->|单关模式| FillSingleLast[将上次金额填充到<br/>当前焦点的输入框]
    CheckMode -->|串关模式| FillMultiLast[将上次金额填充到<br/>串关总金额输入框]
    
    FillSingleLast --> RecalcPayout[重新计算预计返还]
    FillMultiLast --> RecalcPayout
    
    ShowNoHistory --> WaitClick
    OtherActions --> RecalcPayout
    
    RecalcPayout --> DisplayPayout[显示预计返还金额]
    DisplayPayout --> SaveToHistory[将当前金额保存到<br/>临时历史记录]
    SaveToHistory --> WaitSubmit[等待用户提交投注]
    
    WaitSubmit --> BetSubmitted{投注是否成功?}
    BetSubmitted -->|成功| PersistHistory[将临时历史记录<br/>持久化为上次金额]
    BetSubmitted -->|失败| DiscardHistory[丢弃临时历史记录]
    
    PersistHistory --> End([完成])
    DiscardHistory --> End
    
    style Start fill:#4CAF50,stroke:#2E7D32,color:#fff
    style End fill:#F44336,stroke:#C62828,color:#fff
    style DisplayQuickBtns fill:#2196F3,stroke:#1565C0,color:#fff
    style GetLastAmount fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style FillSingleLast fill:#00BCD4,stroke:#00838F,color:#fff
    style FillMultiLast fill:#00BCD4,stroke:#00838F,color:#fff
    style RecalcPayout fill:#3F51B5,stroke:#1A237E,color:#fff
    style ShowNoHistory fill:#FF9800,stroke:#E65100,color:#fff
    style PersistHistory fill:#4CAF50,stroke:#2E7D32,color:#fff
    style DiscardHistory fill:#F44336,stroke:#C62828,color:#fff
    style CheckHistory fill:#FFC107,stroke:#F57C00,color:#000
    style CheckMode fill:#FFC107,stroke:#F57C00,color:#000
    style UserAction fill:#FFC107,stroke:#F57C00,color:#000
    style BetSubmitted fill:#FFC107,stroke:#F57C00,color:#000
```

**图 3：“使用上次金额”功能流程图**

-   **功能**：点击“Last”按钮，系统会自动将上一次**成功提交**的投注金额填充到当前输入框中。
-   **逻辑**：
    -   系统只记录**成功**的投注金额。
    -   在“单关模式”下，会将金额填充到当前有焦点的输入框。
    -   在“串关模式”下，会将金额填充到总金额或单位金额输入框。
    -   如果无历史记录，则提示用户“暂无历史投注记录”。

---

## 4. 结论

v4 版本通过引入“新手/老手模式”，成功地在“便捷性”和“灵活性”之间找到了完美的平衡点。它不仅为新手用户提供了无缝、防错的投注路径，也为专业玩家提供了他们所需要的强大控制力和观察能力。这些专家级的功能增强，将使 MTS 前端在市场竞争中脱颖而出，能够同时满足更广泛用户群体的需求。我们强烈建议开发团队以此 v4 文档作为最终的开发蓝图。
