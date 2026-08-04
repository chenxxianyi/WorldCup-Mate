# WorldCup Mate 联赛主题 UI 开发方案

> 方案名称：One App, Six Matchday Worlds  
> 适用范围：世界杯（WC）、英超（PL）、西甲（PD）、德甲（BL1）、意甲（SA）、法甲（FL1）  
> 前端技术栈：Vue 3 + TypeScript + Vite + Pinia + Vue Router + Tailwind CSS 4  
> 文档状态：开发实施稿

---

## 一、项目目标

WorldCup Mate 已经具备多赛事数据、赛事切换 Store 和联赛页面分支。本次开发不复制六套页面，而是在现有组件体系上增加“赛事视觉主题层”：用户切换赛事后，页面的色彩、背景纹理、卡片轮廓、标题表现和局部动效随赛事变化，同时保持导航、信息结构和操作方式一致。

核心体验目标：

1. 用户无需阅读赛事名称，也能通过整体视觉快速判断当前赛事。
2. 赛事切换具有明确反馈，但不引发大面积闪烁、跳动或认知负担。
3. 浅色/深色模式与赛事主题相互独立，可以任意组合。
4. 直播、成功、警告、降级区等业务语义颜色保持稳定，不被赛事主题覆盖。
5. 所有页面继续复用当前业务组件和数据请求逻辑。

---

## 二、现状与主要问题

### 2.1 已有基础

- `frontend/src/stores/useCompetitionStore.ts` 已管理当前赛事，并通过 localStorage 持久化。
- `frontend/src/utils/theme.ts` 已通过 `data-theme` 管理浅色/深色模式。
- `frontend/src/styles/main.css` 已集中定义主要颜色、阴影、圆角和语义变量。
- 首页、赛程、球队、积分榜和球队详情已经存在世界杯/联赛业务分支。
- `TopBar.vue` 已有赛事选择器，但目前使用原生 `select`。

### 2.2 视觉问题

1. 各赛事共用完全相同的色彩和卡片样式，切换后主要只有文案变化。
2. 移动端顶部区域较高，品牌、赛事选择、账户和主题按钮相互竞争。
3. 原生赛事下拉框缺少品牌感和切换仪式感。
4. 深色模式中页面背景与卡片层级接近，内容分区不够明显。
5. 大量空状态卡片高度偏高，首屏有效信息密度不足。
6. 卡片普遍采用同一组大圆角和柔和阴影，无法表达不同赛事气质。

---

## 三、设计原则

### 3.1 稳定层与变化层

保持稳定：

- 主导航结构与路由位置
- 页面信息顺序和主要操作
- 字号层级与正文可读性
- 直播、完成、警告、成功等语义状态
- 表格、表单和无障碍交互逻辑

允许变化：

- 赛事主色、强调色和柔和背景色
- 页面装饰纹理与 Header 背景
- 卡片圆角、切角、边框与装饰线
- 标题的英文字体表现和数字风格
- 激活状态、进度条和按钮外观
- 赛事切换时的局部过渡动效

### 3.2 品牌边界

- 赛事主题是 WorldCup Mate 对赛事气质的产品化演绎，不等同于官方品牌规范。
- 未确认商标和素材授权前，不将官方赛事 Logo 作为大面积背景或核心产品标志。
- 优先使用赛事简称、抽象图形和自主绘制纹理。
- WorldCup Mate 品牌名称和基础图标始终保留。

---

## 四、赛事主题定义

以下颜色为产品 UI 建议值，开发阶段仍需进行 WCAG 对比度检测。

| 代码 | 赛事 | 主题气质 | 主色 | 强调色 | 深色背景 | 浅色背景 | 图形语言 | 圆角倾向 |
|---|---|---|---|---|---|---|---|---|
| `WC` | 世界杯 | 全球庆典、城市文化、多元融合 | `#8E2145` | `#D8B44A` | `#09090B` | `#F7F2E8` | 经纬线、城市拼贴、数字 26 | 20–24px |
| `PL` | 英超 | 夜场、电能、速度 | `#37003C` | `#DFFF45` | `#100814` | `#F7F2F8` | 声浪、弧形轨迹、聚光灯 | 16–20px + 局部切角 |
| `PD` | 西甲 | 热烈、年轻、数字娱乐 | `#EF3340` | `#20CDBB` | `#120D0E` | `#FFF4F2` | 脉冲圆环、节拍线、锐角 | 14–18px |
| `BL1` | 德甲 | 精准、直接、力量 | `#D20515` | `#FFFFFF` | `#110B0C` | `#FFF5F5` | 斜切线、速度残影、矩形构图 | 10–14px |
| `SA` | 意甲 | 优雅、技术、未来感 | `#0068B5` | `#65D8FF` | `#061521` | `#EFF8FF` | 棱镜、空气动力线、细三色线 | 16–20px |
| `FL1` | 法甲 | 巴黎街头、文化、模块化 | `#16161A` | `#E6FA32` | `#08080A` | `#F6F6F1` | 字母切片、模块拼贴、不对称标签 | 12–18px |

### 4.1 固定语义色

以下颜色不随赛事改变，避免用户对业务状态产生误判：

| 语义 | 建议变量 | 建议值 |
|---|---|---|
| 直播 | `--status-live` | `#FF453A` |
| 成功/完成 | `--status-success` | `#30C76B` |
| 警告 | `--status-warning` | `#F59E0B` |
| 错误/降级 | `--status-danger` | `#E5484D` |
| 信息 | `--status-info` | `#3B82F6` |

---

## 五、主题系统技术设计

### 5.1 HTML 状态模型

浅色/深色模式和赛事模式分别通过两个属性控制：

```html
<html data-theme="dark" data-competition="pl">
```

- `data-theme`：`light` / `dark`
- `data-competition`：`wc` / `pl` / `pd` / `bl1` / `sa` / `fl1`

两个维度互不覆盖。切换赛事不修改用户的浅色/深色偏好，切换明暗模式也不改变当前赛事。

### 5.2 CSS Token 分层

全局 Token 建议拆分为四层：

```css
:root {
  /* 中性基础层 */
  --surface-page: #f5f5f7;
  --surface-card: #ffffff;
  --surface-soft: #f3f4f6;
  --text-primary: #171719;
  --text-secondary: #5f626a;
  --border-subtle: #e4e6eb;

  /* 赛事表达层 */
  --competition-primary: #8e2145;
  --competition-accent: #d8b44a;
  --competition-soft: #f7edf1;
  --competition-on-primary: #ffffff;
  --competition-pattern-opacity: 0.05;
  --competition-card-radius: 20px;
  --competition-card-cut: 0px;

  /* 固定语义层 */
  --status-live: #ff453a;
  --status-success: #30c76b;
  --status-warning: #f59e0b;
  --status-danger: #e5484d;

  /* 动效层 */
  --motion-fast: 120ms;
  --motion-normal: 200ms;
  --motion-theme: 240ms;
  --motion-ease-out: cubic-bezier(0.2, 0.8, 0.2, 1);
}
```

赛事覆盖示例：

```css
[data-competition='pl'] {
  --competition-primary: #37003c;
  --competition-accent: #dfff45;
  --competition-soft: #f3eaf5;
  --competition-on-primary: #ffffff;
  --competition-card-radius: 18px;
  --competition-card-cut: 10px;
}

[data-theme='dark'][data-competition='pl'] {
  --surface-page: #100814;
  --surface-card: #18101d;
  --surface-soft: #211626;
  --competition-soft: #2b1730;
}
```

### 5.3 TypeScript 配置

新增建议文件：

```text
frontend/src/themes/competitionThemes.ts
frontend/src/composables/useCompetitionTheme.ts
```

建议类型：

```ts
export type CompetitionThemeCode = 'WC' | 'PL' | 'PD' | 'BL1' | 'SA' | 'FL1'

export interface CompetitionTheme {
  code: CompetitionThemeCode
  slug: Lowercase<CompetitionThemeCode>
  displayName: string
  displayNameEn: string
  motif: 'global' | 'soundwave' | 'pulse' | 'diagonal' | 'prism' | 'modular'
  cardStyle: 'soft' | 'cut' | 'block' | 'precision' | 'glass' | 'editorial'
}
```

`useCompetitionTheme` 负责：

1. 监听 `useCompetitionStore().currentCode`。
2. 将赛事代码转换为可用主题，未知代码回退到 `WC`。
3. 更新 `document.documentElement.dataset.competition`。
4. 暴露当前主题元数据给 Header、赛事选择器和装饰组件。
5. 在应用卸载或服务端渲染场景下进行安全判断。

### 5.4 初始化时机

主题必须在应用首次绘制前尽量完成初始化，避免先显示世界杯颜色、随后再跳到其他联赛颜色。

建议顺序：

1. 从 `wm-competition` 读取本地赛事代码。
2. 在 `main.ts` 挂载 Vue 应用之前设置 `data-competition`。
3. API 返回赛事列表后再校验代码是否合法。
4. 非法或已停用赛事回退到 `WC`。

---

## 六、组件改造方案

### 6.1 CompetitionSwitcher 新赛事选择器

新增：

```text
frontend/src/components/common/CompetitionSwitcher.vue
frontend/src/components/common/CompetitionMark.vue
```

移动端：

- Header 中仅展示当前赛事胶囊按钮：赛事图形、名称、赛季、展开箭头。
- 点击后打开底部面板，使用 `2 × 3` 网格展示六项赛事。
- 每个赛事项显示缩写、自定义抽象图形、中文名称和赛季。
- 当前赛事同时使用边框、背景和勾选图标表达，不能只依赖颜色。
- 面板支持点击遮罩关闭、Escape 关闭、焦点锁定和焦点归还。

桌面端：

- 使用横向分段赛事导航。
- 宽度不足时降级为 Popover，而不是让整个页面产生横向滚动。

替换位置：

- 将 `TopBar.vue` 中现有原生 `select` 替换为 `CompetitionSwitcher`。

### 6.2 TopBar 顶部导航

移动端目标高度控制在约 104–124px：

```text
┌ WorldCup Mate                         头像 主题 ┐
│ [赛事图形] 英超 · PREMIER LEAGUE  25/26  ▾ │
└──────────────────────────────────────┘
```

改造要点：

- 品牌标志缩小，保留 WorldCup Mate 主品牌。
- 登录、退出和主题切换改为统一尺寸图标按钮。
- 退出登录放入用户菜单，避免长期占据 Header 主视觉位置。
- Header 背景使用赛事渐变和低透明度 Motif。
- 删除会导致布局位移的 Hover 上浮效果，改用边框、背景和亮度变化。

### 6.3 CompetitionMasthead 赛事首屏

新增首页核心组件：

```text
frontend/src/components/competition/CompetitionMasthead.vue
```

合并当前首页的“下一场比赛”和“赛事进度”区域，包含：

- 当前赛事、赛季、当前轮次或阶段
- 下一场比赛双方、队徽、开球时间
- 倒计时或直播比分
- 已完成场次和总场次进度
- 推荐、焦点或直播标签

无下一场比赛时使用紧凑空状态，不保留大型空白卡片。

### 6.4 MatchCard 比赛卡片

业务结构不变，增加赛事主题表达槽位：

- 左侧或顶部赛事强调条。
- Featured 状态使用赛事强调色，直播仍使用固定直播红。
- 轮次、日期和状态标签重新排列，提升扫描效率。
- 比分数字使用运动型窄体字体或 `font-stretch` 可用字体。
- 赛事主题仅改变边框、强调线和背景纹理，不改变收藏、提醒等操作位置。
- 移动端点击热区不小于 44×44px。

### 6.5 BottomTabBar / DesktopRail

- 导航背景保持中性，避免大面积高饱和颜色影响阅读。
- 激活项使用赛事主色和强调色组合。
- 激活状态同时使用图标填充、文字加粗和短指示条。
- 支持移动设备 `safe-area-inset-bottom`。
- 页面切换不重新播放赛事主题切换动画。

### 6.6 EmptyState 统一空状态

新增通用组件：

```text
frontend/src/components/common/EmptyState.vue
```

支持：

- `compact` / `default` 两种尺寸
- 自定义标题、说明、CTA 和抽象 SVG 图形
- 图形颜色自动引用赛事 Token
- 空状态不使用 Emoji

首页默认使用 `compact`，列表页可以使用 `default`。

---

## 七、页面适配

### 7.1 首页

推荐顺序：

1. CompetitionMasthead
2. 赛季进度与关注球队两张紧凑统计卡
3. 今日比赛横向卡片区
4. 我的关注赛程
5. 热门推荐
6. 积分速览

桌面端保持主内容 + 右侧栏布局；移动端将积分速览放到热门推荐之后，避免首屏过长。

### 7.2 赛程页

- 世界杯模式继续使用阶段和小组筛选。
- 联赛模式突出轮次选择，并在页头显示当前轮次进度。
- 当前筛选项使用赛事强调色，日期和状态仍保持统一语义。
- 切换赛事后重置不适用于新赛事的筛选状态。

### 7.3 球队页

- 世界杯使用国家队卡片和大洲/小组信息。
- 联赛使用俱乐部队徽、城市和主场信息。
- 卡片背景可以提取队徽主色做极低透明度装饰，但文本区域必须保持可读。
- 队徽加载失败统一进入缩写占位状态。

### 7.4 积分榜页

- 表格本身保持中性，赛事主题主要应用于表头、Tab 和当前球队高亮。
- 欧冠、欧联、附加赛和降级区必须同时使用色条与文字/图标，不只依赖颜色。
- 移动端固定排名、球队和积分三列，其余数据允许受控横向滚动。

### 7.5 比赛详情页

- 顶部 Score Hero 使用当前赛事主题。
- 世界杯展示阶段、小组和比赛编号。
- 联赛展示轮次、双方联赛排名和近期状态。
- 数据区、提醒区和收藏操作保持组件结构一致。

### 7.6 个人中心

- 页面主体保持 WorldCup Mate 中性品牌风格。
- 当前赛事只影响选中状态和少量强调色。
- 跨赛事关注列表使用赛事徽标/缩写区分来源，避免整张卡片同时使用多个主题背景。

---

## 八、动效设计

### 8.1 赛事切换时序

```text
用户选择赛事
    ↓
立即更新选中状态
    ↓
旧主题内容淡出 120ms
    ↓
更新 data-competition 与主题 Token
    ↓
显示已有缓存或 Skeleton
    ↓
新内容淡入 180ms
```

总体视觉过渡控制在约 240ms，最长不超过 300ms。

### 8.2 动画范围

允许动画：

- Header 背景和 Motif 的透明度
- 赛事选择器选中指示器
- 页面主内容的轻微淡入和 4–8px 位移
- 进度条首次载入
- 直播指示点

禁止或限制：

- 所有卡片同时弹跳或旋转
- 持续移动的装饰背景
- 通过动画修改 width、height、top、left 等高重绘属性
- 超过 500ms 的常规 UI 过渡
- 切换赛事时重播所有页面组件动画

### 8.3 减少动态效果

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    scroll-behavior: auto !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 九、字体与图标

### 9.1 字体

- 中文正文：`PingFang SC`, `Microsoft YaHei`, `Noto Sans SC`, system-ui
- 英文标题和赛事名称：`Barlow Condensed`
- 英文正文和数据：`Barlow`
- 比分和倒计时启用 tabular numbers

字体文件建议自托管并使用 `font-display: swap`，避免依赖不可控的外部字体服务。字体未加载时必须使用系统字体回退。

### 9.2 图标

- 继续使用统一 SVG 或 Material Symbols 图标体系。
- 不使用 Emoji 作为功能图标。
- 所有图标统一使用 24×24 viewBox，常规操作尺寸 20–24px。
- 赛事抽象图形使用自主 SVG，不直接猜测或重绘官方 Logo。

---

## 十、响应式与无障碍

### 10.1 目标视口

至少验证：

- 375×667
- 390×844
- 768×1024
- 1024×768
- 1440×900

### 10.2 验收规则

- 页面不产生非预期水平滚动。
- 固定导航不遮挡页面内容。
- 正文颜色对比度不低于 4.5:1。
- 大字号和关键图形对比度不低于 3:1。
- 所有可点击元素具备明确 Hover、Active 和 Focus Visible 状态。
- 主要触摸目标不小于 44×44px。
- 赛事选择器支持键盘方向键、Enter、Space 和 Escape。
- 当前赛事和当前导航不能只靠颜色表达。
- 图像具有有效 alt；纯装饰图形使用 `aria-hidden="true"`。
- 异步切换期间使用 `aria-busy` 或可感知的加载状态。

---

## 十一、实施阶段

### 阶段 1：主题基础设施（0.5–1 天）

- 新增 `competitionThemes.ts`。
- 新增 `useCompetitionTheme.ts`。
- 在应用初始化时设置 `data-competition`。
- 重构 `main.css` Token，分离中性、赛事、语义和动效变量。
- 为六项赛事建立浅色/深色 Token。

交付标准：切换赛事后，全局主色、强调色、背景和圆角可以稳定变化，无首次加载闪色。

### 阶段 2：赛事选择器与顶部区域（1–1.5 天）

- 实现 `CompetitionSwitcher.vue`。
- 实现 `CompetitionMark.vue`。
- 改造 `TopBar.vue`。
- 完成移动端底部面板与桌面端分段导航。
- 补齐键盘操作、焦点管理和响应式状态。

交付标准：六项赛事可可靠切换，移动端顶部高度明显下降，无水平滚动。

### 阶段 3：首页与比赛卡片（1–1.5 天）

- 新增 `CompetitionMasthead.vue`。
- 合并下一场比赛与赛事进度。
- 改造 `MatchCard.vue` 的赛事主题表现。
- 新增通用 `EmptyState.vue`。
- 调整首页模块顺序和空状态高度。

交付标准：首页首屏能够明确显示当前赛事、下一场比赛和今日内容，各赛事具有明显但一致的视觉差异。

### 阶段 4：全页面适配（1–2 天）

- 适配赛程、球队、积分榜、比赛详情和个人中心。
- 适配 `BottomTabBar` 与 `DesktopRail`。
- 统一按钮、标签、进度条和表格选中状态。
- 清理遗留的硬编码主色。

交付标准：所有用户端页面均使用主题 Token，不出现固定红色误用或主题断层。

### 阶段 5：视觉验证与回归（0.5–1 天）

- 在六项赛事 × 浅色/深色组合下截图检查。
- 检查五类目标视口。
- 检查空数据、加载中、加载失败、直播、已结束等状态。
- 执行 TypeScript 构建、Lint 和现有测试。
- 使用浏览器对比度和无障碍检查工具验证。

预计总工作量：4–6 个前端工作日，不含官方素材授权和大规模队徽资源整理。

---

## 十二、建议文件变更清单

### 新增

```text
frontend/src/themes/competitionThemes.ts
frontend/src/composables/useCompetitionTheme.ts
frontend/src/components/common/CompetitionSwitcher.vue
frontend/src/components/common/CompetitionMark.vue
frontend/src/components/common/EmptyState.vue
frontend/src/components/competition/CompetitionMasthead.vue
frontend/src/styles/competition-themes.css
```

### 修改

```text
frontend/src/main.ts
frontend/src/styles/main.css
frontend/src/stores/useCompetitionStore.ts
frontend/src/layouts/UserLayout.vue
frontend/src/components/common/TopBar.vue
frontend/src/components/common/MatchCard.vue
frontend/src/components/common/BottomTabBar.vue
frontend/src/components/common/DesktopRail.vue
frontend/src/pages/user/HomePage.vue
frontend/src/pages/user/SchedulePage.vue
frontend/src/pages/user/TeamsPage.vue
frontend/src/pages/user/StandingsPage.vue
frontend/src/pages/user/MatchDetailPage.vue
frontend/src/pages/user/ProfilePage.vue
```

---

## 十三、测试方案

### 13.1 单元测试

- 赛事代码映射正确。
- 未知赛事回退到 `WC`。
- Store 赛事变化能更新 `data-competition`。
- 浅色/深色切换不会覆盖赛事属性。
- localStorage 中非法值能够被清理。

### 13.2 组件测试

- CompetitionSwitcher 六项赛事渲染与选中状态。
- 移动端面板的打开、关闭和焦点归还。
- 键盘选择赛事。
- MatchCard 在 scheduled/live/finished 状态下的主题表现。
- EmptyState 的紧凑和默认模式。

### 13.3 端到端测试

1. 从世界杯切换到英超，确认 URL/路由不变、数据重新加载、主题更新。
2. 刷新页面，确认英超选择和明暗模式均被保留。
3. 依次切换六项赛事，确认无报错、无空白屏和无布局跳动。
4. 在赛程、球队和积分榜页面切换赛事，确认筛选状态正确重置。
5. 模拟 API 失败，确认主题仍切换且出现可恢复错误状态。

---

## 十四、验收清单

### 功能

- [ ] 六项赛事均有独立主题配置。
- [ ] 赛事切换与数据切换保持一致。
- [ ] 刷新页面后保留赛事选择。
- [ ] 浅色/深色与赛事主题可以任意组合。
- [ ] 未知赛事能够安全回退。

### 视觉

- [ ] 不看赛事名称也能大致区分六项赛事。
- [ ] 页面仍能清晰识别为同一个 WorldCup Mate 产品。
- [ ] 直播、完成、警告等状态不因主题变化而混淆。
- [ ] 深色模式下页面、卡片和柔和区域层次清楚。
- [ ] 空状态不会占据过多首屏空间。

### 交互与无障碍

- [ ] 动效处于 150–300ms 范围。
- [ ] 支持 `prefers-reduced-motion`。
- [ ] 键盘可以完成赛事切换。
- [ ] Focus Visible 清晰可见。
- [ ] 正文对比度达到 4.5:1。
- [ ] 375px 宽度下无水平滚动。

### 工程质量

- [ ] 页面和组件不再新增硬编码赛事主色。
- [ ] 主题配置有明确 TypeScript 类型。
- [ ] 构建、Lint 和测试通过。
- [ ] 不影响管理后台和现有世界杯业务路径。

---

## 十五、风险与控制

| 风险 | 影响 | 控制措施 |
|---|---|---|
| 六套主题造成维护成本上升 | 中 | 只允许 Token 和少量 Motif 差异，不复制业务组件 |
| 赛事色与语义色冲突 | 高 | 固定直播、成功、警告、错误等语义颜色 |
| 切换时主题和数据不同步 | 中 | 主题立即切换，内容区使用 Skeleton/加载状态，Store 作为唯一赛事来源 |
| 深浅模式组合数量增加 | 中 | 建立自动截图矩阵，重点验证 12 种主题组合 |
| 官方 Logo 或视觉资产授权不明确 | 中 | 使用自主抽象图形，授权确认前不嵌入受限素材 |
| 动效影响低端设备 | 低 | 只使用 transform/opacity，禁用持续背景动画 |
| 首次加载出现主题闪烁 | 中 | Vue 挂载前从 localStorage 设置赛事属性 |

---

## 十六、后续扩展

主题系统完成后，可以继续支持：

1. 欧冠、欧联、足总杯等新增赛事，仅增加主题配置和 Motif。
2. 俱乐部详情页基于队徽色生成低透明度球队主题。
3. 比赛分享图根据赛事主题自动生成。
4. 用户首页提供“跟随当前赛事”与“WorldCup Mate 中性主题”两种偏好。
5. 管理后台增加主题预览和配置校验工具。

本方案的最终目标不是让六项赛事变成六个产品，而是在统一的 WorldCup Mate 体验中，为每项赛事建立清晰、克制且可持续维护的比赛日世界。
