import { defineConfig } from 'vitepress'

const REPO = 'https://github.com/scagogogo/npm-skills'

export default defineConfig({
  base: '/npm-skills/',
  lastUpdated: true,
  cleanUrls: true,

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/npm-skills/favicon.svg' }],
    ['meta', { name: 'theme-color', content: '#cb3837' }]
  ],

  locales: {
    // 中文为默认（root）
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/',
      title: 'NPM Skills',
      description: '面向 AI 智能体与开发者的 NPM Registry 客户端 — 查询、发布、审计、镜像、代理一体化',

      themeConfig: {
        logo: '/logo.svg',

        nav: [
          { text: '首页', link: '/' },
          { text: '快速开始', link: '/getting-started' },
          { text: 'CLI 命令', link: '/cli' },
          { text: 'Go SDK', link: '/api/registry' },
          { text: 'MCP 服务器', link: '/mcp-server' },
          { text: 'GitHub', link: REPO }
        ],

        sidebar: {
          '/': [
            {
              text: '开始',
              collapsed: false,
              items: [
                { text: '简介', link: '/' },
                { text: '快速开始', link: '/getting-started' },
                { text: '安装指南', link: '/installation' },
                { text: 'CLI 命令手册', link: '/cli' }
              ]
            },
            {
              text: 'Go SDK',
              collapsed: false,
              items: [
                { text: 'Registry 客户端', link: '/api/registry' },
                { text: '数据模型', link: '/api/models' },
                { text: '配置选项', link: '/api/configuration' }
              ]
            },
            {
              text: '集成方式',
              collapsed: false,
              items: [
                { text: 'MCP 服务器', link: '/mcp-server' },
                { text: '镜像源', link: '/examples/mirrors' },
                { text: '下载 Tarball', link: '/examples/download' },
                { text: '基础用法', link: '/examples/basic' },
                { text: '高级用法', link: '/examples/advanced' }
              ]
            }
          ]
        },

        socialLinks: [{ icon: 'github', link: REPO }],

        footer: {
          message: '基于 MIT 协议发布。',
          copyright: 'Copyright © 2024-present scagogogo'
        },

        editLink: {
          pattern: `${REPO}/edit/main/website/:path`,
          text: '在 GitHub 上编辑此页面'
        },

        lastUpdated: {
          text: '最后更新于',
          formatOptions: { dateStyle: 'short', timeStyle: 'medium' }
        },

        docFooter: { prev: '上一页', next: '下一页' },
        outline: { label: '本页导航' },
        returnToTopLabel: '回到顶部',
        sidebarMenuLabel: '菜单',
        darkModeSwitchLabel: '主题',
        lightModeSwitchTitle: '切换到浅色模式',
        darkModeSwitchTitle: '切换到深色模式',

        search: { provider: 'local' }
      }
    },

    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      title: 'NPM Skills',
      description: 'NPM Registry client for AI agents and developers — query, publish, audit, mirrors, proxy in one',

      themeConfig: {
        logo: '/logo.svg',

        nav: [
          { text: 'Home', link: '/en/' },
          { text: 'Getting Started', link: '/en/getting-started' },
          { text: 'CLI', link: '/en/cli' },
          { text: 'Go SDK', link: '/en/api/registry' },
          { text: 'MCP Server', link: '/en/mcp-server' },
          { text: 'GitHub', link: REPO }
        ],

        sidebar: {
          '/en/': [
            {
              text: 'Getting Started',
              collapsed: false,
              items: [
                { text: 'Introduction', link: '/en/' },
                { text: 'Getting Started', link: '/en/getting-started' },
                { text: 'Installation', link: '/en/installation' },
                { text: 'CLI Reference', link: '/en/cli' }
              ]
            },
            {
              text: 'Go SDK',
              collapsed: false,
              items: [
                { text: 'Registry Client', link: '/en/api/registry' },
                { text: 'Data Models', link: '/en/api/models' },
                { text: 'Configuration', link: '/en/api/configuration' }
              ]
            },
            {
              text: 'Integrations',
              collapsed: false,
              items: [
                { text: 'MCP Server', link: '/en/mcp-server' },
                { text: 'Mirrors', link: '/en/examples/mirrors' },
                { text: 'Download Tarball', link: '/en/examples/download' },
                { text: 'Basic Usage', link: '/en/examples/basic' },
                { text: 'Advanced Usage', link: '/en/examples/advanced' }
              ]
            }
          ]
        },

        socialLinks: [{ icon: 'github', link: REPO }],

        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024-present scagogogo'
        },

        editLink: {
          pattern: `${REPO}/edit/main/website/en/:path`,
          text: 'Edit this page on GitHub'
        },

        lastUpdated: {
          text: 'Last updated',
          formatOptions: { dateStyle: 'short', timeStyle: 'medium' }
        },

        docFooter: { prev: 'Previous', next: 'Next' },
        outline: { label: 'On this page' },
        returnToTopLabel: 'Return to top',
        sidebarMenuLabel: 'Menu',
        darkModeSwitchLabel: 'Theme',
        lightModeSwitchTitle: 'Switch to light theme',
        darkModeSwitchTitle: 'Switch to dark theme',

        search: { provider: 'local' }
      }
    }
  },

  markdown: {
    theme: { light: 'github-light', dark: 'github-dark' },
    lineNumbers: true
  }
})
