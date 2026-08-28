import type { Dict } from './zh-CN'

export const enUS: Dict = {
  app: {
    title: 'VS Code Download Hub',
    intro: 'Help intranet users fetch official downloads that match their local VSCode client.'
  },
  nav: {
    home: 'Home',
    search: 'Extensions',
    releases: 'Releases',
    trending: 'Trending',
    docs: 'Docs'
  },
  home: {
    channel: 'Channel',
    version: 'Version',
    platform: 'Platform',
    architecture: 'Architecture',
    submit: 'Lookup',
    clientCard: 'Client download',
    serverCard: 'Server download (vscode-server)',
    commitHash: 'Commit hash',
    inferHint: 'Auto-detected from your browser. Override anytime.'
  },
  extension: {
    searchPlaceholder: 'Search extensions (publisher.name or display name fragment)',
    noResults: 'No matching extensions',
    engines: 'Engines.vscode'
  },
  trending: {
    client: 'Client trending',
    server: 'Server trending',
    extension: 'Extension trending',
    window24h: '24 hours',
    window7d: '7 days',
    window30d: '30 days',
    empty: 'No data in this window'
  },
  docs: {
    title: 'Documentation',
    topics: {
      clientInstall: 'Client install',
      serverInstall: 'Server install',
      serverStart: 'Server start & connect',
      offlineImport: 'Offline extension import',
      faq: 'Troubleshooting FAQ'
    },
    scriptNote: 'Commands below contain no download actions. Download on an internet-enabled client, then run offline.'
  },
  errors: {
    VERSION_NOT_FOUND: 'Version not found',
    EXTENSION_VERSION_NOT_FOUND: 'Extension version not found',
    INVALID_PLATFORM_ARCH: 'Invalid platform/architecture combination',
    UPSTREAM_FAILURE: 'Upstream service unavailable, please retry',
    INVALID_REQUEST: 'Invalid request parameters',
    INTERNAL: 'Internal server error',
    NON_OFFICIAL_URL_BLOCKED: 'Non-official host blocked'
  }
}
