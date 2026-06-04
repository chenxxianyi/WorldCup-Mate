import type { User } from '@/types/user'

export const mockUser: User = {
  id: 1,
  username: 'leochen',
  email: 'leochen@example.com',
  nickname: 'Leo Chen',
  avatar: 'L',
  timezone: 'Asia/Shanghai',
  language: 'zh-CN',
  notificationEmail: 'leochen@example.com',
  followed_teams: [5, 11],
  favorite_matches: [1, 2, 3],
}
