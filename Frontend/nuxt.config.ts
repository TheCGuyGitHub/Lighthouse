// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/color-mode',
    'shadcn-nuxt',
  ],
  colorMode: {
    classSuffix: ''
  },
  css: [
    '@/assets/css/tailwind.css',
  ],
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * @default "./app/components/ui"
     */
    componentDir: './app/components/ui'
  },
  alias: {
    '@/lib': './app/lib'
  },
  vite: {
  plugins: [
    // Use dynamic import for ESM compatibility
    (await import('vite-tsconfig-paths')).default()
  ]
}
})