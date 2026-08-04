export type CompetitionCode = 'WC' | 'PL' | 'PD' | 'BL1' | 'SA' | 'FL1'

export type DemoTeam = readonly [name: string, code: string, city: string, color: string, crest?: string]

export interface LeagueTheme {
  code: CompetitionCode
  slug: string
  name: string
  en: string
  season: string
  mark: string
  tagline: string
  description: string
  stage: string
  played: number
  total: number
  venue: string
  teams: DemoTeam[]
}

export interface DemoMatch {
  id: number
  home: DemoTeam
  away: DemoTeam
  time: string
  date: string
  status: 'scheduled' | 'live' | 'finished'
  score: string
  venue: string
  featured: boolean
}

export const leagueThemes: Record<CompetitionCode, LeagueTheme> = {
  WC: {
    code: 'WC', slug: 'wc', name: '世界杯', en: 'WORLD CUP 26', season: '美加墨 2026', mark: '26',
    tagline: '全世界，共赴同一片球场',
    description: '48 支国家队、16 座城市与 104 场比赛，组成这个夏天最盛大的足球庆典。',
    stage: '小组赛 · 第 2 轮', played: 32, total: 104, venue: 'MetLife Stadium · New York',
    teams: [
      ['阿根廷', 'ARG', '布宜诺斯艾利斯', '#69b6e4'], ['西班牙', 'ESP', '马德里', '#c60b1e'],
      ['巴西', 'BRA', '里约热内卢', '#009c3b'], ['法国', 'FRA', '巴黎', '#1b3f8b'],
      ['英格兰', 'ENG', '伦敦', '#d8d8d8'], ['德国', 'GER', '柏林', '#171717'],
      ['葡萄牙', 'POR', '里斯本', '#b51e36'], ['荷兰', 'NED', '阿姆斯特丹', '#f58220'],
      ['日本', 'JPN', '东京', '#173a84'], ['墨西哥', 'MEX', '墨西哥城', '#006847'],
    ],
  },
  PL: {
    code: 'PL', slug: 'pl', name: '英超', en: 'PREMIER LEAGUE', season: '2026–27', mark: 'PL',
    tagline: '每一轮，都有新的主角', description: '高强度、快节奏和永不确定的比赛夜，让每一轮都像决赛一样值得期待。',
    stage: 'MATCHWEEK 12', played: 108, total: 380, venue: 'Emirates Stadium · London',
    teams: [
      ['阿森纳', 'ARS', '伦敦', '#d71920'], ['利物浦', 'LIV', '利物浦', '#c8102e'],
      ['曼城', 'MCI', '曼彻斯特', '#6cabdd'], ['切尔西', 'CHE', '伦敦', '#034694'],
      ['曼联', 'MUN', '曼彻斯特', '#da291c'], ['热刺', 'TOT', '伦敦', '#132257'],
      ['纽卡斯尔', 'NEW', '纽卡斯尔', '#241f20'], ['阿斯顿维拉', 'AVL', '伯明翰', '#670e36'],
      ['布莱顿', 'BHA', '布莱顿', '#0057b8'], ['埃弗顿', 'EVE', '利物浦', '#003399'],
    ],
  },
  PD: {
    code: 'PD', slug: 'pd', name: '西甲', en: 'LALIGA', season: '2026–27', mark: 'LL',
    tagline: '让足球成为一种本能', description: '技术、激情与城市文化在九十分钟内交汇，每一次控球都带着独特的西班牙节拍。',
    stage: 'JORNADA 12', played: 110, total: 380, venue: 'Santiago Bernabéu · Madrid',
    teams: [
      ['皇家马德里', 'RMA', '马德里', '#d8d8d8'], ['巴塞罗那', 'BAR', '巴塞罗那', '#a50044'],
      ['马德里竞技', 'ATM', '马德里', '#c8102e'], ['毕尔巴鄂竞技', 'ATH', '毕尔巴鄂', '#ee2523'],
      ['比利亚雷亚尔', 'VIL', '比利亚雷亚尔', '#f7e600'], ['皇家社会', 'RSO', '圣塞瓦斯蒂安', '#0067b1'],
      ['塞维利亚', 'SEV', '塞维利亚', '#d71920'], ['贝蒂斯', 'BET', '塞维利亚', '#158447'],
      ['瓦伦西亚', 'VAL', '瓦伦西亚', '#f58220'], ['赫罗纳', 'GIR', '赫罗纳', '#cd2534'],
    ],
  },
  BL1: {
    code: 'BL1', slug: 'bl1', name: '德甲', en: 'BUNDESLIGA', season: '2026–27', mark: 'BL',
    tagline: '速度必须足够直接', description: '高压、纵深和看台声浪构成德国比赛日，没有多余动作，只有最快到达球门的路径。',
    stage: 'SPIELTAG 11', played: 90, total: 306, venue: 'Allianz Arena · München',
    teams: [
      ['拜仁慕尼黑', 'FCB', '慕尼黑', '#dc052d'], ['多特蒙德', 'BVB', '多特蒙德', '#fdeb00'],
      ['勒沃库森', 'B04', '勒沃库森', '#e32221'], ['莱比锡', 'RBL', '莱比锡', '#dd0741'],
      ['法兰克福', 'SGE', '法兰克福', '#e1000f'], ['斯图加特', 'VFB', '斯图加特', '#e32219'],
      ['沃尔夫斯堡', 'WOB', '沃尔夫斯堡', '#65b32e'], ['门兴', 'BMG', '门兴格拉德巴赫', '#111111'],
      ['弗赖堡', 'SCF', '弗赖堡', '#e3001b'], ['柏林联合', 'FCU', '柏林', '#eb1923'],
    ],
  },
  SA: {
    code: 'SA', slug: 'sa', name: '意甲', en: 'SERIE A', season: '2026–27', mark: 'A',
    tagline: '精确，也可以很优雅', description: '战术结构、空间控制和瞬间创造力，在意大利式的精确中呈现比赛的高级质感。',
    stage: 'GIORNATA 12', played: 110, total: 380, venue: 'San Siro · Milano',
    teams: [
      ['国际米兰', 'INT', '米兰', '#00529f'], ['AC 米兰', 'MIL', '米兰', '#fb090b'],
      ['尤文图斯', 'JUV', '都灵', '#1d1d1b'], ['那不勒斯', 'NAP', '那不勒斯', '#12a0d7'],
      ['罗马', 'ROM', '罗马', '#8e1f2f'], ['亚特兰大', 'ATA', '贝加莫', '#1e71b8'],
      ['拉齐奥', 'LAZ', '罗马', '#8ec7e8'], ['佛罗伦萨', 'FIO', '佛罗伦萨', '#5b2c83'],
      ['博洛尼亚', 'BOL', '博洛尼亚', '#1a2f48'], ['都灵', 'TOR', '都灵', '#8a1538'],
    ],
  },
  FL1: {
    code: 'FL1', slug: 'fl1', name: '法甲', en: 'LIGUE 1', season: '2026–27', mark: 'L1',
    tagline: '看见法式足球的新锋芒', description: '年轻天赋、街头文化和大胆表达，让法国比赛日成为一张持续更新的视觉海报。',
    stage: 'JOURNÉE 11', played: 88, total: 306, venue: 'Parc des Princes · Paris',
    teams: [
      ['巴黎圣日耳曼', 'PSG', '巴黎', '#004170'], ['马赛', 'OM', '马赛', '#2faee0'],
      ['摩纳哥', 'ASM', '摩纳哥', '#e32219'], ['里尔', 'LOSC', '里尔', '#d71920'],
      ['里昂', 'OL', '里昂', '#1b3f8b'], ['朗斯', 'RCL', '朗斯', '#e30613'],
      ['尼斯', 'OGCN', '尼斯', '#171717'], ['雷恩', 'SRFC', '雷恩', '#d90000'],
      ['斯特拉斯堡', 'RCSA', '斯特拉斯堡', '#009ee0'], ['南特', 'FCN', '南特', '#f8d323'],
    ],
  },
}

export const competitionCodes = Object.keys(leagueThemes) as CompetitionCode[]

export function isCompetitionCode(value: string | null): value is CompetitionCode {
  return Boolean(value && Object.prototype.hasOwnProperty.call(leagueThemes, value))
}

export function getDemoMatches(theme: LeagueTheme): DemoMatch[] {
  const t = theme.teams
  return [
    { id: 0, home: t[0], away: t[1], time: '21:00', date: '08月08日', status: 'scheduled', score: 'VS', venue: theme.venue, featured: true },
    { id: 1, home: t[2], away: t[3], time: "67'", date: '08月08日', status: 'live', score: '2–1', venue: `${t[2][2]} · Main Stadium`, featured: false },
    { id: 2, home: t[4], away: t[5], time: '完场', date: '08月07日', status: 'finished', score: '1–1', venue: `${t[4][2]} · National Arena`, featured: false },
    { id: 3, home: t[6], away: t[7], time: '23:30', date: '08月09日', status: 'scheduled', score: 'VS', venue: `${t[6][2]} · City Stadium`, featured: true },
    { id: 4, home: t[8], away: t[9], time: '02:45', date: '08月09日', status: 'scheduled', score: 'VS', venue: `${t[8][2]} · Football Park`, featured: false },
  ]
}

export function statusLabel(status: DemoMatch['status']) {
  if (status === 'live') return '直播中'
  if (status === 'finished') return '已结束'
  return '未开始'
}
