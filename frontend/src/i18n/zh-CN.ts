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
    inferHint: '已根据浏览器自动选择,可手动覆盖'
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
    topics: {
      clientInstall: '客户端安装',
      serverInstall: '服务端安装',
      serverStart: '服务端启动与连接',
      offlineImport: '离线导入扩展',
      faq: '排错 FAQ'
    },
    scriptNote: '以下命令不含任何下载动作;请在具备网络的客户端下载后,放入内网运行。'
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
