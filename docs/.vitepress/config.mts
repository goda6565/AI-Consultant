import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/AI-Consultant/',
  title: "AI-Consultant",
  description: "AI-Consultant project documentation",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Backend', link: '/backend/' },
      { text: 'Frontend', link: '/frontend/' },
    ],

    sidebar: [
      {
        text: 'Overview',
        items: [
          { text: 'Project Overview', link: '/' },
        ]
      },
      {
        text: 'Backend',
        items: [
          { text: 'Admin', link: '/backend/admin' },
          { text: 'Agent', link: '/backend/agent' },
          { text: 'Proposal Job', link: '/backend/proposal-job' },
          { text: 'Vector', link: '/backend/vector' }
        ]
      },
      {
        text: 'Frontend',
        items: [
          { text: 'Frontend Overview', link: '/frontend/frontend' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/harutin/sakana/problem-set/ai-consultant' }
    ],

    footer: {
      message: 'AI-Consultant Project',
      copyright: 'Copyright © 2025 AI-Consultant'
    },

    search: {
      provider: 'local'
    }
  }
})
