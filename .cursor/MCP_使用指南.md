# GVA MCP 使用指南

配置好 `gva-helper` MCP 后，在 Cursor 里可以直接用自然语言让 AI 调用 GVA 的能力，完成需求分析、代码生成、菜单/API/字典创建等。

---

## 在 Cursor 里怎么用

1. **正常对话即可**  
   在 Cursor 的 AI 对话里（Chat / Composer），用中文或英文描述你的需求，例如：
   - “帮我做一个简单的图书管理功能，有书名、作者、分类”
   - “根据需求分析一下，需要哪些模块和字典”
   - “生成一个用户管理的前后端代码”

2. **AI 会自动选工具**  
   当你的需求涉及「需求分析、代码生成、菜单、API、字典」等时，AI 会自动调用 GVA 的 MCP 工具，你不需要手动点选工具。

3. **推荐的说法示例**  
   - “用 requirement_analyzer 分析一下：我要做一个 xxx 管理系统，包含……”  
   - “先 gva_analyze 当前项目，再根据需求生成代码”  
   - “帮我查一下系统里现在有哪些字典”  
   - “给刚生成的 API 在系统里创建对应的 API 权限记录”

---

## MCP 提供的工具一览

| 工具名 | 用途 | 典型用法 |
|--------|------|----------|
| **requirement_analyzer** | 需求分析 + 模块设计（首选入口） | 把业务需求变成「模块、字段、字典」等设计 |
| **gva_analyze** | 分析当前 GVA 的包/模块/字典 | 看现有结构，判断要不要新建包、模块、字典 |
| **gva_execute** | 执行代码生成 | 根据执行计划生成后端 + 前端代码 |
| **gva_review** | 代码审查 | 在 gva_execute 之后做一次生成结果审查 |
| **create_api** | 创建 API 权限记录 | 把新接口登记到系统的 API 列表里 |
| **create_menu** | 创建前端菜单 | 把新页面登记到菜单/路由里 |
| **list_all_apis** | 列出所有 API | 查系统里已有接口，方便对接或去重 |
| **list_all_menus** | 列出所有菜单 | 查菜单树、路由、组件路径 |
| **query_dictionaries** | 查询字典 | 查字典类型、选项，用于表单/下拉等 |
| **generate_dictionary_options** | 生成并创建字典 | 按字段描述生成字典选项并写入系统 |

---

## 推荐工作流（从需求到上线）

```
1. requirement_analyzer   → 用自然语言说需求，得到「模块 + 字段 + 字典」设计
2. gva_analyze            → 看当前项目包/模块/字典，决定新建还是复用
3. gva_execute            → 按执行计划生成后端 + 前端代码
4. gva_review（可选）     → 审查生成代码
5. create_api / create_menu → 把新接口、新页面登记到权限和菜单
```

**示例对话：**

- 你：“我要做一个简单的「会议室预约」功能：可以选会议室、选时间段、填预约人。用 requirement_analyzer 分析一下。”  
- AI 会调用 `requirement_analyzer`，返回模块划分、字段设计、是否需要字典等。

- 你：“根据刚才的分析，用 gva_analyze 看下当前项目，再 gva_execute 生成代码。”  
- AI 会先 `gva_analyze`，再按结果调用 `gva_execute` 生成代码。

- 你：“把刚才生成的接口在系统里登记成 API，再加一个前端菜单。”  
- AI 会调用 `create_api` 和 `create_menu`。

---

## 使用前注意

1. **GVA 后端要已启动**  
   MCP 连的是本机 `http://localhost:8888/sse`，确保 gin-vue-admin 的 server 已跑在 8888 端口。

2. **首次用可先“查现状”**  
   可以说：“先 list_all_menus 和 query_dictionaries，让我看看当前系统菜单和字典。”  
   方便后续生成代码时复用已有菜单和字典。

3. **生成代码后要人工检查**  
   `gva_execute` 生成的是初版代码，建议在项目里跑一遍、看接口和页面是否符合预期，再按需改。

4. **在 Cursor 里确认 MCP 已启用**  
   设置里 MCP 的 `gva-helper` 显示为已连接/可用即可；若断连，检查 GVA 是否在跑、`.cursor/mcp.json` 里的 url 是否正确。

---

按上面方式在对话里直接说需求或指定工具名，即可使用 GVA MCP。
