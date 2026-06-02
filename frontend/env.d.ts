/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ImportMeta {
  readonly url: string
}

declare class URL {
  constructor(input: string, base?: string | URL)
}

declare module 'node:url' {
  export function fileURLToPath(url: string | URL): string
}
