export type Continent = '亚洲' | '欧洲' | '南美洲' | '北美洲' | '非洲' | '大洋洲'

/** Backend raw team */
export interface ApiTeam {
  id: number
  name: string
  name_en: string
  fifa_code: string | null
  flag_url: string
  continent: string
  group_id?: number | null
  group: { id: number; name: string } | null
  description: string
  coach: string
  external_code?: string | null
  team_type?: string
  country?: string
  venue?: string
}

/** Frontend normalized team */
export interface Team {
  id: number
  name: string
  name_en: string
  code: string
  flag: string
  group_id?: number | null
  group_name: string
  continent: Continent
  team_type?: string
  country?: string
  venue?: string
  is_followed?: boolean
}

export function normalizeTeam(t: ApiTeam): Team {
  return {
    id: t.id,
    name: t.name,
    name_en: t.name_en,
    code: t.fifa_code || t.external_code || '',
    flag: t.flag_url,
    group_id: t.group_id,
    group_name: t.group?.name || '',
    continent: (t.continent as Continent) || '欧洲',
    team_type: t.team_type,
    country: t.country,
    venue: t.venue,
  }
}
