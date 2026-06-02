import type { User } from '@/types/user'

export const mockUser: User = {
  id: 1,
  username: 'leochen',
  nickname: 'Leo Chen',
  avatar: 'L',
  timezone: 'Asia/Shanghai',
  language: 'zh-CN',
  followed_teams: [5, 11],
  favorite_matches: [1, 2, 3],
}
