export interface ApiUser {
  id: number
  username: string
  email: string
  avatar: string
  timezone: string
  language: string
  role: string
  notification_email: string
}

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  timezone: string
  language: string
  notificationEmail: string
  followed_teams?: number[]
  favorite_matches?: number[]
}

export function normalizeUser(u: ApiUser): User {
  return {
    id: u.id,
    username: u.username,
    nickname: u.username,
    avatar: u.avatar || u.username.charAt(0).toUpperCase(),
    timezone: u.timezone,
    language: u.language,
    notificationEmail: u.notification_email || '',
  }
}
