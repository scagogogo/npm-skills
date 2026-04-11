import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/npm-crawler/',
  
  locales: {
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'NPM Crawler',
      description: '高性能 NPM Registry 客户端，支持多镜像源和代理配置',
      
      themeConfig: {
        logo: 'https://cdn.worldvectorlogo.com/logos/npm-2.svg',
        
        nav: [
          { text: '首页', link: '/' },
          { text: '快速开始', link: '/getting-started' },
          { text: 'API文档', link: '/api/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/npm-crawler' }
        ],

        sidebar: {
          '/': [
            {
              text: '介绍',
              items: [
                { text: '什么是 NPM Crawler', link: '/' },
                { text: '快速开始', link: '/getting-started' },
                { text: '安装指南', link: '/installation' }
              ]
            },
            {
              text: 'API参考',
              items: [
                { text: 'API 概述', link: '/api/' },
                { text: 'Registry 客户端', link: '/api/registry' },
                { text: '数据模型', link: '/api/models' },
                { text: '配置选项', link: '/api/configuration' }
              ]
            },
            {
              text: '示例',
              items: [
                { text: '基本用法', link: '/examples/basic' },
                { text: '高级用法', link: '/examples/advanced' },
                { text: '镜像源配置', link: '/examples/mirrors' },
                { text: '下载 Tarball', link: '/examples/download' }
              ]
            }
          ]
        },

        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/npm-crawler' }
        ],

        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2023-present NPM Crawler Team'
        },

        editLink: {
          pattern: 'https://github.com/scagogogo/npm-crawler/edit/main/docs/:path',
          text: '在 GitHub 上编辑此页面'
        },

        lastUpdated: {
          text: '最后更新于',
          formatOptions: {
            dateStyle: 'short',
            timeStyle: 'medium'
          }
        },

        docFooter: {
          prev: '上一页',
          next: '下一页'
        },

        outline: {
          label: '页面导航'
        },

        returnToTopLabel: '回到顶部',
        sidebarMenuLabel: '菜单',
        darkModeSwitchLabel: '主题',
        lightModeSwitchTitle: '切换到浅色模式',
        darkModeSwitchTitle: '切换到深色模式'
      }
    },
    
    en: {
      label: 'English',
      lang: 'en-US',
      title: 'NPM Crawler',
      description: 'High-performance NPM Registry client with multi-mirror source and proxy support',
      
      themeConfig: {
        logo: 'https://cdn.worldvectorlogo.com/logos/npm-2.svg',
        
        nav: [
          { text: 'Home', link: '/en/' },
          { text: 'Getting Started', link: '/en/getting-started' },
          { text: 'API Docs', link: '/en/api/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/npm-crawler' }
        ],

        sidebar: {
          '/en/': [
            {
              text: 'Introduction',
              items: [
                { text: 'What is NPM Crawler', link: '/en/' },
                { text: 'Getting Started', link: '/en/getting-started' },
                { text: 'Installation', link: '/en/installation' }
              ]
            },
            {
              text: 'API Reference',
              items: [
                { text: 'API Overview', link: '/en/api/' },
                { text: 'Registry Client', link: '/en/api/registry' },
                { text: 'Data Models', link: '/en/api/models' },
                { text: 'Configuration', link: '/en/api/configuration' }
              ]
            },
            {
              text: 'Examples',
              items: [
                { text: 'Basic Usage', link: '/en/examples/basic' },
                { text: 'Advanced Usage', link: '/en/examples/advanced' },
                { text: 'Mirror Configuration', link: '/en/examples/mirrors' },
                { text: 'Download Tarball', link: '/en/examples/download' }
              ]
            }
          ]
        },

        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/npm-crawler' }
        ],

        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2023-present NPM Crawler Team'
        },

        editLink: {
          pattern: 'https://github.com/scagogogo/npm-crawler/edit/main/docs/:path',
          text: 'Edit this page on GitHub'
        },

        lastUpdated: {
          text: 'Last updated',
          formatOptions: {
            dateStyle: 'short',
            timeStyle: 'medium'
          }
        },

        docFooter: {
          prev: 'Previous page',
          next: 'Next page'
        },

        outline: {
          label: 'On this page'
        },

        returnToTopLabel: 'Return to top',
        sidebarMenuLabel: 'Menu',
        darkModeSwitchLabel: 'Theme',
        lightModeSwitchTitle: 'Switch to light theme',
        darkModeSwitchTitle: 'Switch to dark theme'
      }
    }
  },

  markdown: {
    theme: {
      light: 'github-light',
      dark: 'github-dark'
    },
    lineNumbers: true
  }
})
