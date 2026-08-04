/** Competition metadata from GET /api/competitions */
export interface Competition {
  id: number
  code: string
  name: string
  name_en: string
  country: string
  logo_url: string
  format: 'league' | 'cup'
  season: number
  status: string
  sort_order: number
}

/** The World Cup code — the default competition, legacy behavior */
export const WC_CODE = 'WC'

export function seasonLabel(season: number): string {
  if (!season) return ''
  return `${season}-${String((season + 1) % 100).padStart(2, '0')}`
}
