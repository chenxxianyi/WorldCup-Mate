import type { GroupStanding } from '@/types/standing'

export const mockGroupStandings: GroupStanding[] = [
  {
    group_name: 'Group A',
    standings: [
      { team_id: 1, team_name: '墨西哥', team_code: 'MEX', flag: '🇲🇽', played: 3, won: 2, drawn: 1, lost: 0, goals_for: 5, goals_against: 2, goal_difference: 3, points: 7, status: '晋级' },
      { team_id: 2, team_name: '加拿大', team_code: 'CAN', flag: '🇨🇦', played: 3, won: 1, drawn: 2, lost: 0, goals_for: 3, goals_against: 2, goal_difference: 1, points: 5, status: '晋级' },
      { team_id: 3, team_name: '南非', team_code: 'RSA', flag: '🇿🇦', played: 3, won: 1, drawn: 1, lost: 1, goals_for: 2, goals_against: 2, goal_difference: 0, points: 4, status: '待定' },
      { team_id: 4, team_name: '新西兰', team_code: 'NZL', flag: '🇳🇿', played: 3, won: 0, drawn: 0, lost: 3, goals_for: 1, goals_against: 5, goal_difference: -4, points: 0, status: '淘汰' },
    ],
  },
  {
    group_name: 'Group B',
    standings: [
      { team_id: 9, team_name: '巴西', team_code: 'BRA', flag: '🇧🇷', played: 2, won: 2, drawn: 0, lost: 0, goals_for: 5, goals_against: 1, goal_difference: 4, points: 6, status: '晋级' },
      { team_id: 10, team_name: '西班牙', team_code: 'ESP', flag: '🇪🇸', played: 2, won: 1, drawn: 1, lost: 0, goals_for: 3, goals_against: 1, goal_difference: 2, points: 4, status: '晋级' },
      { team_id: 7, team_name: '美国', team_code: 'USA', flag: '🇺🇸', played: 2, won: 1, drawn: 0, lost: 1, goals_for: 2, goals_against: 2, goal_difference: 0, points: 3, status: '待定' },
      { team_id: 8, team_name: '加纳', team_code: 'GHA', flag: '🇬🇭', played: 2, won: 0, drawn: 0, lost: 2, goals_for: 0, goals_against: 6, goal_difference: -6, points: 0, status: '淘汰' },
    ],
  },
  {
    group_name: 'Group C',
    standings: [
      { team_id: 5, team_name: '阿根廷', team_code: 'ARG', flag: '🇦🇷', played: 2, won: 2, drawn: 0, lost: 0, goals_for: 4, goals_against: 1, goal_difference: 3, points: 6, status: '晋级' },
      { team_id: 6, team_name: '法国', team_code: 'FRA', flag: '🇫🇷', played: 2, won: 1, drawn: 0, lost: 1, goals_for: 3, goals_against: 2, goal_difference: 1, points: 3, status: '待定' },
      { team_id: 13, team_name: '英格兰', team_code: 'ENG', flag: '🏴󠁧󠁢󠁥󠁮󠁧󠁿', played: 2, won: 0, drawn: 1, lost: 1, goals_for: 1, goals_against: 2, goal_difference: -1, points: 1, status: '待定' },
      { team_id: 14, team_name: '葡萄牙', team_code: 'POR', flag: '🇵🇹', played: 2, won: 0, drawn: 1, lost: 1, goals_for: 1, goals_against: 4, goal_difference: -3, points: 1, status: '淘汰' },
    ],
  },
]

export const mockBestThird: { team_name: string; team_code: string; flag: string; group_name: string; points: number; goal_difference: number }[] = [
  { team_name: '南非', team_code: 'RSA', flag: '🇿🇦', group_name: 'Group A', points: 4, goal_difference: 0 },
  { team_name: '美国', team_code: 'USA', flag: '🇺🇸', group_name: 'Group B', points: 3, goal_difference: 0 },
  { team_name: '英格兰', team_code: 'ENG', flag: '🏴󠁧󠁢󠁥󠁮󠁧󠁿', group_name: 'Group C', points: 1, goal_difference: -1 },
  { team_name: '葡萄牙', team_code: 'POR', flag: '🇵🇹', group_name: 'Group C', points: 1, goal_difference: -3 },
]
