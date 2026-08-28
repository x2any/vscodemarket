export const zhCN = {
  app: {
    title: 'VS Code 下载中心',
    intro: '辅助内网用户获取与本地 VSCode 客户端配套的官方下载物。'
  },
  nav: {
    home: '首页',
    search: '扩展',
    releases: '版本',
    trending: '热榜',
    docs: '文档'
  },
  home: {
    channel: '频道',
    version: '版本号',
    platform: '平台',
    architecture: '架构',
    submit: '查询',
    clientCard: '客户端下载',
    serverCard: '服务端下载(vscode-server)',
    commitHash: '对应 commit',
    inferHint: '已根据浏览器自动选择,可手动覆盖',
    matrixHint: '不选平台/架构时返回全平台 × 全架构的下载链接矩阵'
  },
  extension: {
    searchPlaceholder: '搜索扩展(支持 publisher.name 或显示名片段)',
    noResults: '无匹配扩展',
    engines: '适配 VSCode'
  },
  trending: {
    client: '客户端热榜',
    server: '服务端热榜',
    extension: '扩展热榜',
    window24h: '24 小时',
    window7d: '7 天',
    window30d: '30 天',
    empty: '该时段暂无数据'
  },
  docs: {
    title: '文档',
    intro: '所有脚本仅含运行 / 安装命令,不含任何下载动作。请在有网络的客户端下载二进制,再拷入内网执行。',
    topics: {
      clientInstall: {
        title: '客户端安装',
        body: '从首页或版本列表拿到与本地架构一致的客户端直链;在 Windows 上运行安装程序,在 macOS 上拖入 Applications,在 Linux 上解压 .tar.gz 至 /opt。',
        snippet: '# Linux 示例(将下载好的文件预先放入 /tmp)\nsudo tar -xzf /tmp/code-stable-x64-*.tar.gz -C /opt\nsudo ln -sf /opt/VSCode-linux-x64/code /usr/local/bin/code'
      },
      serverInstall: {
        title: '服务端安装',
        body: '使用首页主路径同时拿到与 Stable 客户端 commit hash 严格对齐的 vscode-server tar 包,上传到内网远端,解压到 ~/.vscode-server 或指定目录。',
        snippet: 'mkdir -p ~/.vscode-server\ntar -xzf vscode-server-linux-x64.tar.gz -C ~/.vscode-server\nls ~/.vscode-server/bin'
      },
      serverStart: {
        title: '服务端启动与连接',
        body: '在客户端中通过 Remote-SSH 连接到内网主机,首次连接会自动调用 code-server 二进制;若内网无网络,请提前放置二进制。',
        snippet: '# 在 Remote-SSH 设置中覆盖默认路径\n# settings.json\n{\n  "remote.SSH.serverInstallPath": {\n    "host": "~/.vscode-server"\n  }\n}'
      },
      offlineImport: {
        title: '离线导入扩展',
        body: '在有网机器下载 .vsix 后,通过 VSIX 扩展命令安装。批量导入可将 .vsix 放入 extensions 目录后重启客户端。',
        snippet: 'code --install-extension /path/to/extension.vsix\n# 批量(将所有 .vsix 放入目录后)\nfor f in /opt/extensions/*.vsix; do code --install-extension "$f"; done'
      },
      faq: {
        title: '排错 FAQ',
        body: '常见问题:commit hash 不匹配请回到首页重新查 Stable 版本;vscode-server 找不到请确认 Remote-SSH 设置中的 serverInstallPath;扩展安装失败请检查 engines.vscode 是否覆盖当前客户端。',
        snippet: '# 查看客户端实际 commit\ncode --version\n# 查看 server 实际 commit\ncat ~/.vscode-server/bin/*/commit.txt 2>/dev/null || true'
      }
    }
  },
  errors: {
    VERSION_NOT_FOUND: '版本不存在',
    EXTENSION_VERSION_NOT_FOUND: '扩展版本不存在',
    INVALID_PLATFORM_ARCH: '平台 / 架构组合无效',
    UPSTREAM_FAILURE: '上游服务暂时不可用,请稍后重试',
    INVALID_REQUEST: '请求参数非法',
    INTERNAL: '服务器内部错误',
    NON_OFFICIAL_URL_BLOCKED: '检测到非官方域名,已拒绝'
  }
} as const

export type Dict = typeof zhCN
