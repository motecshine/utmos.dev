# 架构文档

## 查看 Mermaid 图表

### 方法 1：在 Cursor 中预览（推荐）

1. **安装扩展**：
   - 打开扩展面板（`Ctrl+Shift+X` 或 `Cmd+Shift+X`）
   - 搜索并安装：`Markdown Preview Mermaid Support`
   - 或者安装：`Markdown Preview Enhanced`

2. **预览文档**：
   - 打开 `middleware-architecture.md`
   - 按 `Ctrl+Shift+V`（Windows/Linux）或 `Cmd+Shift+V`（Mac）打开预览
   - Mermaid 图表会自动渲染

### 方法 2：在线查看

访问 [Mermaid Live Editor](https://mermaid.live/)，复制文档中的 Mermaid 代码块内容即可查看。

### 方法 3：导出为图片

如果安装了 Mermaid CLI：

```bash
# 安装 Mermaid CLI
npm install -g @mermaid-js/mermaid-cli

# 导出为 PNG
mmdc -i docs/architecture/middleware-architecture.md -o docs/architecture/middleware-architecture.png

# 导出为 SVG
mmdc -i docs/architecture/middleware-architecture.md -o docs/architecture/middleware-architecture.svg
```

### 方法 4：VS Code 内置支持

如果 Cursor 基于较新版本的 VS Code，可能已经内置了 Mermaid 预览支持。直接打开 Markdown 预览即可。

## 文档列表

- [基础中间件架构图](./middleware-architecture.md) - 平台基础中间件架构和组件关系

