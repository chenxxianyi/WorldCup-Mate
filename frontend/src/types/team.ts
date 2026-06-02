export type Continent = '亚洲' | '欧洲' | '南美洲' | '北美洲' | '非洲' | '大洋洲'

/** Backend raw team */
export interface ApiTeam {
  id: number
  name: string
  name_en: string
  fifa_code: string
  flag_url: string
  continent: string
  group_id: number
  group: { id: number; name: string } | null
  description: string
  coach: string
}

/** Frontend normalized team */
export interface Team {
  id: number
  name: string
  name_en: string
  code: string
  flag: string
  group_name: string
  continent: Continent
  is_followed?: boolean
}

export function normalizeTeam(t: ApiTeam): Team {
  return {
    id: t.id,
    name: t.name,
    name_en: t.name_en,
    code: t.fifa_code,
    flag: t.flag_url,
    group_name: t.group?.name || '',
    continent: t.continent as Continent,
  }
}
