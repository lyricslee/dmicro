module.exports = {
    port: 5678,
    title: 'dmicro',
    description: 'Just playing around',
    markdown: {
        lineNumbers: true // 代码块显示行号
    },
    themeConfig: {
        logo: '/assets/img/logo.png',
        nav: [
            { text: '首页', link: '/' },
            { text: '开发文档', link: '/api/' },
            { text: '指南', link: '/guide/' }
        ],
        sidebar: {
            '/guide/': getGuideSidebar('指南'),
            '/api/': getApiSidebar('开发文档')
        },

        sidebarDepth: 2,
    }
}

function getGuideSidebar (title) {
  return [
    {
      title,
      collapsable: true,
      children: [
        ''
      ]
    }
  ]
}

function getApiSidebar (title) {
  return [
    {
      title,
      collapsable: true,
        children: [
          '',
          'passport'
        ]
    }
  ]
}

