export interface ApiPlayer {
  id: number
  team_id: number
  name: string
  name_en?: string
  shirt_number?: number
  position?: string
  position_label?: string
  photo_url?: string
  club?: string
}

export interface Player {
  id: number
  teamId: number
  name: string
  nameEn: string
  shirtNumber?: number
  position: string
  positionLabel: string
  photoUrl: string
  club: string
}

export function normalizePlayer(value: ApiPlayer): Player {
  const position = (value.position || '').toUpperCase()

  return {
    id: value.id,
    teamId: value.team_id,
    name: value.name,
    nameEn: value.name_en || '',
    shirtNumber: value.shirt_number,
    position,
    positionLabel: value.position_label || '',
    photoUrl: value.photo_url || '',
    club: value.club || '',
  }
}
