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
    inferHint: 'Auto-detected from your browser. Override anytime.',
    matrixHint: 'Leave platform/architecture empty to get the full matrix of download links'
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
    intro: 'All snippets contain run/install commands only — no download actions. Fetch binaries on an internet-enabled host, then transfer and run offline.',
    topics: {
      clientInstall: {
        title: 'Client install',
        body: 'Grab the architecture-matched client direct link from the home page or release list. Run the installer on Windows, drag to Applications on macOS, extract the .tar.gz under /opt on Linux.',
        snippet: '# Linux example (place the downloaded archive in /tmp first)\nsudo tar -xzf /tmp/code-stable-x64-*.tar.gz -C /opt\nsudo ln -sf /opt/VSCode-linux-x64/code /usr/local/bin/code'
      },
      serverInstall: {
        title: 'Server install',
        body: 'Use the home page main path to fetch the vscode-server tarball strictly aligned to the Stable client commit hash. Upload it to the intranet host and extract to ~/.vscode-server (or a custom path).',
        snippet: 'mkdir -p ~/.vscode-server\ntar -xzf vscode-server-linux-x64.tar.gz -C ~/.vscode-server\nls ~/.vscode-server/bin'
      },
      serverStart: {
        title: 'Server start & connect',
        body: 'Connect via Remote-SSH from the client. The first connection invokes the code-server binary automatically. In offline environments, place the binary beforehand.',
        snippet: '# settings.json override\n{\n  "remote.SSH.serverInstallPath": {\n    "host": "~/.vscode-server"\n  }\n}'
      },
      offlineImport: {
        title: 'Offline extension import',
        body: 'Download .vsix files on an internet-enabled host, then install via the VSIX extension command. Bulk import by copying .vsix files into the extensions directory and restarting the client.',
        snippet: 'code --install-extension /path/to/extension.vsix\n# bulk (place all .vsix in the directory first)\nfor f in /opt/extensions/*.vsix; do code --install-extension "$f"; done'
      },
      faq: {
        title: 'Troubleshooting FAQ',
        body: 'Common issues: commit hash mismatch — re-check the Stable version on the home page; vscode-server not found — verify serverInstallPath in Remote-SSH settings; extension install fails — check engines.vscode against the current client.',
        snippet: '# client commit\ncode --version\n# server commit\ncat ~/.vscode-server/bin/*/commit.txt 2>/dev/null || true'
      }
    }
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
