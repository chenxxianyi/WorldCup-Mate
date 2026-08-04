(function () {
  'use strict'

  const competitions = {
    WC: {
      slug: 'wc',
      name: '世界杯',
      en: 'WORLD CUP 26',
      season: '美加墨 2026',
      mark: '26',
      tagline: '全世界，共赴同一片球场',
      description: '48 支国家队、16 座城市与 104 场比赛，组成这个夏天最盛大的足球庆典。',
      stage: '小组赛 · 第 2 轮',
      played: 32,
      total: 104,
      venue: 'MetLife Stadium · New York',
      teams: [
        ['阿根廷', 'ARG', '布宜诺斯艾利斯', '#69b6e4'],
        ['西班牙', 'ESP', '马德里', '#c60b1e'],
        ['巴西', 'BRA', '里约热内卢', '#009c3b'],
        ['法国', 'FRA', '巴黎', '#1b3f8b'],
        ['英格兰', 'ENG', '伦敦', '#d8d8d8'],
        ['德国', 'GER', '柏林', '#171717'],
        ['葡萄牙', 'POR', '里斯本', '#b51e36'],
        ['荷兰', 'NED', '阿姆斯特丹', '#f58220'],
        ['日本', 'JPN', '东京', '#173a84'],
        ['墨西哥', 'MEX', '墨西哥城', '#006847'],
      ],
    },
    PL: {
      slug: 'pl',
      name: '英超',
      en: 'PREMIER LEAGUE',
      season: '2026–27',
      mark: 'PL',
      tagline: '每一轮，都有新的主角',
      description: '高强度、快节奏和永不确定的比赛夜，让每一轮都像决赛一样值得期待。',
      stage: 'MATCHWEEK 12',
      played: 108,
      total: 380,
      venue: 'Emirates Stadium · London',
      teams: [
        ['阿森纳', 'ARS', '伦敦', '#d71920'],
        ['利物浦', 'LIV', '利物浦', '#c8102e'],
        ['曼城', 'MCI', '曼彻斯特', '#6cabdd'],
        ['切尔西', 'CHE', '伦敦', '#034694'],
        ['曼联', 'MUN', '曼彻斯特', '#da291c'],
        ['热刺', 'TOT', '伦敦', '#132257'],
        ['纽卡斯尔', 'NEW', '纽卡斯尔', '#241f20'],
        ['阿斯顿维拉', 'AVL', '伯明翰', '#670e36'],
        ['布莱顿', 'BHA', '布莱顿', '#0057b8'],
        ['埃弗顿', 'EVE', '利物浦', '#003399'],
      ],
    },
    PD: {
      slug: 'pd',
      name: '西甲',
      en: 'LALIGA',
      season: '2026–27',
      mark: 'LL',
      tagline: '让足球成为一种本能',
      description: '技术、激情与城市文化在九十分钟内交汇，每一次控球都带着独特的西班牙节拍。',
      stage: 'JORNADA 12',
      played: 110,
      total: 380,
      venue: 'Santiago Bernabéu · Madrid',
      teams: [
        ['皇家马德里', 'RMA', '马德里', '#d8d8d8'],
        ['巴塞罗那', 'BAR', '巴塞罗那', '#a50044'],
        ['马德里竞技', 'ATM', '马德里', '#c8102e'],
        ['毕尔巴鄂竞技', 'ATH', '毕尔巴鄂', '#ee2523'],
        ['比利亚雷亚尔', 'VIL', '比利亚雷亚尔', '#f7e600'],
        ['皇家社会', 'RSO', '圣塞瓦斯蒂安', '#0067b1'],
        ['塞维利亚', 'SEV', '塞维利亚', '#d71920'],
        ['贝蒂斯', 'BET', '塞维利亚', '#158447'],
        ['瓦伦西亚', 'VAL', '瓦伦西亚', '#f58220'],
        ['赫罗纳', 'GIR', '赫罗纳', '#cd2534'],
      ],
    },
    BL1: {
      slug: 'bl1',
      name: '德甲',
      en: 'BUNDESLIGA',
      season: '2026–27',
      mark: 'BL',
      tagline: '速度必须足够直接',
      description: '高压、纵深和看台声浪构成德国比赛日，没有多余动作，只有最快到达球门的路径。',
      stage: 'SPIELTAG 11',
      played: 90,
      total: 306,
      venue: 'Allianz Arena · München',
      teams: [
        ['拜仁慕尼黑', 'FCB', '慕尼黑', '#dc052d'],
        ['多特蒙德', 'BVB', '多特蒙德', '#fdeb00'],
        ['勒沃库森', 'B04', '勒沃库森', '#e32221'],
        ['莱比锡', 'RBL', '莱比锡', '#dd0741'],
        ['法兰克福', 'SGE', '法兰克福', '#e1000f'],
        ['斯图加特', 'VFB', '斯图加特', '#e32219'],
        ['沃尔夫斯堡', 'WOB', '沃尔夫斯堡', '#65b32e'],
        ['门兴', 'BMG', '门兴格拉德巴赫', '#111111'],
        ['弗赖堡', 'SCF', '弗赖堡', '#e3001b'],
        ['柏林联合', 'FCU', '柏林', '#eb1923'],
      ],
    },
    SA: {
      slug: 'sa',
      name: '意甲',
      en: 'SERIE A',
      season: '2026–27',
      mark: 'A',
      tagline: '精确，也可以很优雅',
      description: '战术结构、空间控制和瞬间创造力，在意大利式的精确中呈现比赛的高级质感。',
      stage: 'GIORNATA 12',
      played: 110,
      total: 380,
      venue: 'San Siro · Milano',
      teams: [
        ['国际米兰', 'INT', '米兰', '#00529f'],
        ['AC 米兰', 'MIL', '米兰', '#fb090b'],
        ['尤文图斯', 'JUV', '都灵', '#1d1d1b'],
        ['那不勒斯', 'NAP', '那不勒斯', '#12a0d7'],
        ['罗马', 'ROM', '罗马', '#8e1f2f'],
        ['亚特兰大', 'ATA', '贝加莫', '#1e71b8'],
        ['拉齐奥', 'LAZ', '罗马', '#8ec7e8'],
        ['佛罗伦萨', 'FIO', '佛罗伦萨', '#5b2c83'],
        ['博洛尼亚', 'BOL', '博洛尼亚', '#1a2f48'],
        ['都灵', 'TOR', '都灵', '#8a1538'],
      ],
    },
    FL1: {
      slug: 'fl1',
      name: '法甲',
      en: 'LIGUE 1',
      season: '2026–27',
      mark: 'L1',
      tagline: '看见法式足球的新锋芒',
      description: '年轻天赋、街头文化和大胆表达，让法国比赛日成为一张持续更新的视觉海报。',
      stage: 'JOURNÉE 11',
      played: 88,
      total: 306,
      venue: 'Parc des Princes · Paris',
      teams: [
        ['巴黎圣日耳曼', 'PSG', '巴黎', '#004170'],
        ['马赛', 'OM', '马赛', '#2faee0'],
        ['摩纳哥', 'ASM', '摩纳哥', '#e32219'],
        ['里尔', 'LOSC', '里尔', '#d71920'],
        ['里昂', 'OL', '里昂', '#1b3f8b'],
        ['朗斯', 'RCL', '朗斯', '#e30613'],
        ['尼斯', 'OGCN', '尼斯', '#171717'],
        ['雷恩', 'SRFC', '雷恩', '#d90000'],
        ['斯特拉斯堡', 'RCSA', '斯特拉斯堡', '#009ee0'],
        ['南特', 'FCN', '南特', '#f8d323'],
      ],
    },
  }

  const navItems = [
    ['home', '首页', 'i-home'],
    ['schedule', '赛程', 'i-calendar'],
    ['teams', '球队', 'i-shield'],
    ['standings', '积分榜', 'i-table'],
    ['profile', '我的', 'i-user'],
  ]

  const queryParams = new URLSearchParams(location.search)

  const state = {
    competition: queryParams.get('competition') || safeStorageGet('wm-demo-competition') || 'WC',
    theme: queryParams.get('theme') || safeStorageGet('wm-demo-theme') || 'dark',
    route: parseHash().route,
    detailIndex: parseHash().index,
    scheduleFilter: '全部',
    standingsMode: '总榜',
  }

  let toastTimer = null
  let dialogLastFocus = null

  const elements = {
    main: document.querySelector('#main-content'),
    desktopLeagues: document.querySelector('#desktop-leagues'),
    desktopRail: document.querySelector('#desktop-rail'),
    bottomNav: document.querySelector('#bottom-nav'),
    themeToggle: document.querySelector('#theme-toggle'),
    mobileTrigger: document.querySelector('#mobile-competition-trigger'),
    dialog: document.querySelector('#competition-dialog'),
    competitionGrid: document.querySelector('#competition-grid'),
    toast: document.querySelector('#toast'),
  }

  function safeStorageGet(key) {
    try {
      return localStorage.getItem(key)
    } catch {
      return null
    }
  }

  function safeStorageSet(key, value) {
    try {
      localStorage.setItem(key, value)
    } catch {
      // Local files can disable storage in hardened browser configurations.
    }
  }

  function parseHash() {
    const raw = location.hash.replace(/^#\/?/, '') || 'home'
    const [route, index] = raw.split('/')
    const valid = ['home', 'schedule', 'teams', 'standings', 'profile', 'match', 'team', 'login']
    return {
      route: valid.includes(route) ? route : 'home',
      index: Number.isFinite(Number(index)) ? Number(index) : 0,
    }
  }

  function icon(id, className = '') {
    return `<svg class="${className}" aria-hidden="true"><use href="#${id}"></use></svg>`
  }

  function activeMainRoute() {
    if (state.route === 'match') return 'schedule'
    if (state.route === 'team') return 'teams'
    if (state.route === 'login') return 'profile'
    return state.route
  }

  function competition() {
    return competitions[state.competition] || competitions.WC
  }

  function badge(team, size = '') {
    const [, code, , color] = team
    return `<span class="team-badge ${size}" style="--badge-a:${color};--badge-b:${color}">${code}</span>`
  }

  function getMatches(comp = competition()) {
    const t = comp.teams
    return [
      { home: t[0], away: t[1], time: '21:00', date: '08月08日', status: 'scheduled', score: 'VS', venue: comp.venue, featured: true },
      { home: t[2], away: t[3], time: "67'", date: '08月08日', status: 'live', score: '2–1', venue: `${t[2][2]} · Main Stadium`, featured: false },
      { home: t[4], away: t[5], time: '完场', date: '08月07日', status: 'finished', score: '1–1', venue: `${t[4][2]} · National Arena`, featured: false },
      { home: t[6], away: t[7], time: '23:30', date: '08月09日', status: 'scheduled', score: 'VS', venue: `${t[6][2]} · City Stadium`, featured: true },
      { home: t[8], away: t[9], time: '02:45', date: '08月09日', status: 'scheduled', score: 'VS', venue: `${t[8][2]} · Football Park`, featured: false },
    ]
  }

  function statusLabel(status) {
    if (status === 'live') return '直播中'
    if (status === 'finished') return '已结束'
    return '未开始'
  }

  function matchCard(match, index, compact = false) {
    return `
      <article class="card match-card clickable-card ${match.featured ? 'featured' : ''}" data-match-card="${index}">
        <div class="match-card-top">
          <span class="label">${competition().stage}</span>
          <span class="status-pill ${match.status}">${match.status === 'live' ? '<i class="live-dot"></i>&nbsp;' : ''}${statusLabel(match.status)}</span>
        </div>
        <div class="match-teams">
          <div class="match-team">
            ${badge(match.home, 'small')}
            <span class="match-team-copy"><strong>${match.home[0]}</strong><small>${match.home[1]}</small></span>
          </div>
          <button class="match-score text-button" type="button" data-route="match" data-index="${index}" aria-label="查看 ${match.home[0]} 对 ${match.away[0]} 比赛详情">${match.score}</button>
          <div class="match-team away">
            <span class="match-team-copy"><strong>${match.away[0]}</strong><small>${match.away[1]}</small></span>
            ${badge(match.away, 'small')}
          </div>
        </div>
        <div class="match-card-bottom">
          <span class="match-meta"><strong>${match.date} · ${match.time}</strong><span>${compact ? competition().name : match.venue}</span></span>
          <span class="card-actions">
            <button class="icon-action" type="button" data-demo-action="reminder" aria-label="设置比赛提醒">${icon('i-bell')}</button>
            <button class="icon-action" type="button" data-demo-action="favorite" aria-label="收藏比赛">${icon('i-star')}</button>
          </span>
        </div>
      </article>`
  }

  function renderChrome() {
    const comp = competition()
    document.documentElement.dataset.theme = state.theme
    document.documentElement.dataset.competition = comp.slug
    document.title = `WorldCup Mate · ${comp.name}主题 Demo`

    elements.desktopLeagues.innerHTML = Object.entries(competitions)
      .map(([code, item]) => `
        <button class="league-tab ${code === state.competition ? 'active' : ''}" type="button" data-competition-code="${code}" aria-pressed="${code === state.competition}">
          <span class="league-tab-mark">${item.mark}</span><span class="league-tab-name">${item.name}</span>
        </button>`)
      .join('')

    elements.mobileTrigger.innerHTML = `
      <span class="mobile-trigger-copy"><span class="league-tab-mark">${comp.mark}</span><strong>${comp.name} · ${comp.en}</strong><span>${comp.season}</span></span>
      ${icon('i-chevron')}`

    elements.themeToggle.innerHTML = state.theme === 'dark' ? icon('i-sun') : icon('i-moon')
    elements.themeToggle.setAttribute('aria-label', state.theme === 'dark' ? '切换浅色模式' : '切换深色模式')

    const current = activeMainRoute()
    const nav = navItems
      .map(([route, label, iconId]) => `
        <button class="rail-link ${route === current ? 'active' : ''}" type="button" data-route="${route}" aria-current="${route === current ? 'page' : 'false'}">
          ${icon(iconId)}<span>${label}</span>
        </button>`)
      .join('')
    elements.desktopRail.innerHTML = nav
    elements.bottomNav.innerHTML = nav.replaceAll('rail-link', 'bottom-link')

    elements.competitionGrid.innerHTML = Object.entries(competitions)
      .map(([code, item]) => `
        <button class="competition-option ${code === state.competition ? 'active' : ''}" type="button" data-competition-code="${code}" data-preview="${item.slug}" aria-pressed="${code === state.competition}">
          ${code === state.competition ? `<span class="option-check">${icon('i-check')}</span>` : ''}
          <span class="option-mark">${item.mark}</span>
          <span><strong>${item.name}</strong><small>${item.en} · ${item.season}</small></span>
        </button>`)
      .join('')
  }

  function renderHome() {
    const comp = competition()
    const matches = getMatches(comp)
    const progress = Math.round((comp.played / comp.total) * 100)
    const top = comp.teams.slice(0, 5)
    return `
      <div class="page-view">
        <section class="competition-masthead" aria-labelledby="hero-title">
          <div class="hero-copy">
            <div>
              <div class="hero-league-line"><span class="hero-mark">${comp.mark}</span><span><b>${comp.name}</b><br><small>${comp.en} · ${comp.season}</small></span></div>
              <h1 id="hero-title">${comp.tagline}</h1>
              <p class="hero-description">${comp.description}</p>
            </div>
            <div class="hero-progress">
              <div class="hero-progress-meta"><span>${comp.stage}</span><span>${comp.played} / ${comp.total} 场 · ${progress}%</span></div>
              <div class="progress-track"><div class="progress-fill" style="width:${progress}%"></div></div>
            </div>
          </div>
          <div class="next-match-panel">
            <div class="next-match-head"><span class="status-line"><i class="live-dot"></i>下一场焦点</span><span class="next-kickoff">TODAY · ${matches[0].time}</span></div>
            <div class="team-versus">
              <span class="versus-team">${badge(matches[0].home)}<strong>${matches[0].home[0]}</strong></span>
              <span class="versus-center"><strong>VS</strong><small>${comp.stage}</small></span>
              <span class="versus-team">${badge(matches[0].away)}<strong>${matches[0].away[0]}</strong></span>
            </div>
            <button class="primary-button full-button" type="button" data-route="match" data-index="0">查看比赛中心 ${icon('i-arrow')}</button>
          </div>
        </section>

        <section class="quick-grid" aria-label="赛事速览">
          <article class="card quick-card"><span class="quick-card-top"><span>赛季进度</span>${icon('i-trophy')}</span><strong>${progress}%</strong></article>
          <article class="card quick-card"><span class="quick-card-top"><span>关注球队</span>${icon('i-star')}</span><strong>3 支</strong></article>
          <article class="card quick-card"><span class="quick-card-top"><span>今日比赛</span>${icon('i-calendar')}</span><strong>5 场</strong></article>
        </section>

        <div class="content-grid">
          <div>
            <section class="section">
              <div class="section-heading"><div><p class="eyebrow">MATCHDAY</p><h2>今日比赛</h2></div><button class="text-link" type="button" data-route="schedule">完整赛程 ${icon('i-arrow')}</button></div>
              <div class="match-strip">${matches.slice(0, 4).map((match, index) => matchCard(match, index, true)).join('')}</div>
            </section>

            <section class="section">
              <div class="section-heading"><div><p class="eyebrow">FOLLOWING</p><h2>我的关注赛程</h2></div><span>3 支球队</span></div>
              <div class="match-list">${matchCard(matches[3], 3)}${matchCard(matches[4], 4)}</div>
            </section>

            <section class="section">
              <div class="section-heading"><div><p class="eyebrow">PERSONAL</p><h2>比赛提醒</h2></div></div>
              <article class="card empty-compact">
                <span class="empty-art">${icon('i-bell')}</span>
                <span class="empty-copy"><h3>不错过真正重要的比赛</h3><p>开球前 15 分钟，通过站内消息和邮件提醒你。</p></span>
                <button class="primary-button" type="button" data-demo-action="reminder">设置提醒</button>
              </article>
            </section>
          </div>

          <aside>
            <section class="section">
              <div class="section-heading"><div><p class="eyebrow">TABLE</p><h2>积分速览</h2></div><button class="text-link" type="button" data-route="standings">查看全部</button></div>
              <div class="card standings-preview">
                ${top.map((team, index) => `<div class="standing-row"><span class="standing-position ${index < 4 ? 'zone-top' : ''}">${index + 1}</span><span class="standing-team">${badge(team, 'small')}<span>${team[0]}</span></span><strong>${30 - index * 2}</strong></div>`).join('')}
              </div>
            </section>

            <section class="section">
              <div class="section-heading"><div><p class="eyebrow">FORM</p><h2>领头羊近况</h2></div></div>
              <article class="card quick-card">
                <span class="standing-team">${badge(comp.teams[0], 'small')}<span>${comp.teams[0][0]}</span></span>
                <span class="form-dots"><i class="form-dot">胜</i><i class="form-dot">胜</i><i class="form-dot draw">平</i><i class="form-dot">胜</i><i class="form-dot">胜</i></span>
              </article>
            </section>
          </aside>
        </div>
      </div>`
  }

  function renderSchedule() {
    const matches = getMatches()
    const filters = competition().slug === 'wc' ? ['全部', '小组赛', '淘汰赛', '收藏'] : ['全部', '第 11 轮', '第 12 轮', '直播', '收藏']
    return `
      <div class="page-view">
        <header class="page-heading"><div><p class="eyebrow">${competition().en}</p><h1>赛程中心</h1></div><p class="muted">当地时间 · 自动转换为你的时区</p></header>
        <section class="card filter-card">
          <div class="filter-row">${filters.map((item) => `<button class="chip ${item === state.scheduleFilter ? 'active' : ''}" type="button" data-schedule-filter="${item}">${item}</button>`).join('')}</div>
        </section>
        <section class="schedule-day">
          <div class="day-marker"><span>周六</span><strong>08</strong><span>AUG</span></div>
          <div class="match-list">${matches.slice(0, 3).map((match, index) => matchCard(match, index)).join('')}</div>
        </section>
        <section class="schedule-day">
          <div class="day-marker"><span>周日</span><strong>09</strong><span>AUG</span></div>
          <div class="match-list">${matches.slice(3).map((match, index) => matchCard(match, index + 3)).join('')}</div>
        </section>
      </div>`
  }

  function renderTeams() {
    const comp = competition()
    return `
      <div class="page-view">
        <header class="page-heading"><div><p class="eyebrow">${comp.en}</p><h1>${comp.slug === 'wc' ? '国家队' : '俱乐部'}</h1></div><p class="muted">${comp.teams.length} 支示例球队 · 点击查看详情</p></header>
        <label class="search-box">${icon('i-search')}<input id="team-search" type="search" placeholder="搜索球队名称、城市或缩写" aria-label="搜索球队" /></label>
        <section class="section">
          <div class="teams-grid" id="teams-grid">
            ${comp.teams.map((team, index) => `
              <button class="team-card" type="button" data-route="team" data-index="${index}" data-search="${team.join(' ')}" data-code="${team[1]}">
                ${badge(team)}<h3>${team[0]}</h3><p>${team[2]} · ${team[1]}</p>
                <span class="team-card-footer"><span>${index < 4 ? '争冠区' : '查看球队资料'}</span>${icon('i-arrow')}</span>
              </button>`).join('')}
          </div>
        </section>
      </div>`
  }

  function renderStandings() {
    const comp = competition()
    const modes = comp.slug === 'wc' ? ['小组 A', '小组 B', '最佳第三'] : ['总榜', '主场', '客场']
    const standingTeams = comp.slug === 'wc' ? comp.teams.slice(0, 4) : comp.teams
    if (!modes.includes(state.standingsMode)) state.standingsMode = modes[0]
    return `
      <div class="page-view">
        <header class="page-heading">
          <div><p class="eyebrow">${comp.en} · ${comp.season}</p><h1>积分榜</h1></div>
          <div class="standings-tabs">${modes.map((mode) => `<button class="chip ${mode === state.standingsMode ? 'active' : ''}" type="button" data-standing-mode="${mode}">${mode}</button>`).join('')}</div>
        </header>
        <section class="card table-card">
          <table class="standings-table">
            <thead><tr><th>排名</th><th>球队</th><th>赛</th><th>胜</th><th>平</th><th>负</th><th>净胜</th><th>积分</th></tr></thead>
            <tbody>
              ${standingTeams.map((team, index) => {
                const played = comp.slug === 'wc' ? 3 : 11
                const points = comp.slug === 'wc' ? Math.max(1, 7 - index * 2) : 31 - index * 2
                const wins = Math.floor(points / 3)
                const draws = points % 3
                const losses = Math.max(0, played - wins - draws)
                const zoneClass = index < (comp.slug === 'wc' ? 2 : 4) ? 'zone-top' : index > (comp.slug === 'wc' ? 2 : 7) ? 'zone-danger' : ''
                return `<tr><td><span class="standing-position ${zoneClass}">${index + 1}</span></td><td><span class="table-team">${badge(team, 'small')}<span>${team[0]}</span></span></td><td>${played}</td><td>${wins}</td><td>${draws}</td><td>${losses}</td><td>${14 - index * 2 > 0 ? '+' : ''}${14 - index * 2}</td><td><strong>${points}</strong></td></tr>`
              }).join('')}
            </tbody>
          </table>
          <div class="zone-legend"><span class="zone-item"><i class="zone-swatch"></i>晋级 / 欧冠区</span><span class="zone-item"><i class="zone-swatch europa"></i>次级欧战区</span><span class="zone-item"><i class="zone-swatch danger"></i>淘汰 / 降级区</span></div>
        </section>
      </div>`
  }

  function renderProfile() {
    const comp = competition()
    const matches = getMatches()
    return `
      <div class="page-view">
        <header class="page-heading"><div><p class="eyebrow">MY MATCHDAY</p><h1>我的</h1></div><button class="secondary-button" type="button" data-route="login">切换账号</button></header>
        <div class="profile-layout">
          <article class="card profile-card">
            <span class="profile-avatar">M</span><h2>Matchday Fan</h2><p>和 WorldCup Mate 一起看球的第 126 天</p>
            <div class="profile-stats"><span class="profile-stat"><strong>3</strong><span>关注球队</span></span><span class="profile-stat"><strong>12</strong><span>收藏比赛</span></span><span class="profile-stat"><strong>8</strong><span>比赛提醒</span></span></div>
          </article>
          <div class="stack">
            <section class="card settings-card">
              <div class="setting-line"><span class="setting-label"><strong>当前赛事</strong><span>首页与数据范围</span></span><span class="setting-value">${comp.name} · ${comp.season}</span></div>
              <div class="setting-line"><span class="setting-label"><strong>开球提醒</strong><span>比赛开始前 15 分钟</span></span><button class="text-button" type="button" data-demo-action="reminder">已开启</button></div>
              <div class="setting-line"><span class="setting-label"><strong>外观模式</strong><span>可以与赛事主题自由组合</span></span><button class="text-button" type="button" data-demo-action="theme">${state.theme === 'dark' ? '深色' : '浅色'}</button></div>
            </section>
            <section>
              <div class="section-heading"><div><p class="eyebrow">SAVED</p><h2>最近收藏</h2></div></div>
              ${matchCard(matches[0], 0)}
            </section>
          </div>
        </div>
      </div>`
  }

  function renderMatchDetail() {
    const matches = getMatches()
    const match = matches[Math.min(state.detailIndex, matches.length - 1)] || matches[0]
    const isLive = match.status === 'live'
    return `
      <div class="page-view">
        <div class="back-row"><button class="back-button" type="button" data-route="schedule">${icon('i-back')} 返回赛程</button></div>
        <section class="score-hero">
          <div class="match-card-top" style="position:relative"><span class="status-line">${isLive ? '<i class="live-dot"></i> LIVE MATCH' : competition().stage}</span><span class="next-kickoff">${match.date} · ${match.time}</span></div>
          <div class="detail-score-line">
            <div class="detail-team">${badge(match.home, 'large')}<h2>${match.home[0]}</h2><p>${match.home[1]}</p></div>
            <div class="detail-score">${match.score}<small>${isLive ? "67' · 比赛进行中" : match.venue}</small></div>
            <div class="detail-team">${badge(match.away, 'large')}<h2>${match.away[0]}</h2><p>${match.away[1]}</p></div>
          </div>
        </section>
        <div class="detail-grid section">
          <article class="card detail-panel">
            <h3>比赛数据</h3>
            ${[['控球率', 56, 44], ['射门', 14, 9], ['射正', 6, 4], ['角球', 7, 3], ['犯规', 10, 12]].map(([label, home, away]) => `<div class="stat-row"><span>${home}</span><div class="stat-bar" title="${label}"><span style="--value:${home}%"></span></div><small>${label}</small><div class="stat-bar" title="${label}"><span style="--value:${away}%"></span></div><span>${away}</span></div>`).join('')}
          </article>
          <article class="card timeline-card">
            <h3>比赛事件</h3>
            <div class="timeline-item"><span class="timeline-minute">67'</span><i class="timeline-node"></i><span class="timeline-copy"><strong>${match.home[0]} 换人调整</strong><span>中场位置发生变化</span></span></div>
            <div class="timeline-item"><span class="timeline-minute">54'</span><i class="timeline-node"></i><span class="timeline-copy"><strong>${match.away[0]} 扳回一球</strong><span>禁区内右脚低射</span></span></div>
            <div class="timeline-item"><span class="timeline-minute">31'</span><i class="timeline-node"></i><span class="timeline-copy"><strong>${match.home[0]} 扩大领先</strong><span>快速反击完成破门</span></span></div>
            <div class="timeline-item"><span class="timeline-minute">12'</span><i class="timeline-node"></i><span class="timeline-copy"><strong>${match.home[0]} 首开纪录</strong><span>定位球机会转化为进球</span></span></div>
          </article>
        </div>
      </div>`
  }

  function renderTeamDetail() {
    const comp = competition()
    const team = comp.teams[Math.min(state.detailIndex, comp.teams.length - 1)] || comp.teams[0]
    const opponent = comp.teams[(state.detailIndex + 1) % comp.teams.length]
    const match = { home: team, away: opponent, time: '21:00', date: '08月12日', status: 'scheduled', score: 'VS', venue: `${team[2]} · Main Stadium`, featured: true }
    return `
      <div class="page-view">
        <div class="back-row"><button class="back-button" type="button" data-route="teams">${icon('i-back')} 返回球队</button></div>
        <section class="team-hero">
          <div class="team-hero-content">
            ${badge(team, 'large')}
            <div class="team-hero-copy"><p class="eyebrow" style="color:var(--competition-accent)">${comp.en}</p><h1>${team[0]}</h1><p>${team[2]} · ${team[1]}</p><div class="team-detail-meta"><span class="hero-meta-pill">当前排名 1</span><span class="hero-meta-pill">主场 ${team[2]}</span><span class="hero-meta-pill">关注 128.4K</span></div></div>
          </div>
        </section>
        <div class="detail-grid section">
          <div class="stack">
            <article class="card detail-panel"><h3>赛季表现</h3><div class="quick-grid" style="margin-top:0"><span class="quick-card"><small class="muted">积分</small><strong>31</strong></span><span class="quick-card"><small class="muted">进球</small><strong>28</strong></span><span class="quick-card"><small class="muted">失球</small><strong>11</strong></span></div></article>
            <article class="card detail-panel"><h3>近期状态</h3><div class="form-dots"><i class="form-dot">胜</i><i class="form-dot">胜</i><i class="form-dot draw">平</i><i class="form-dot">胜</i><i class="form-dot loss">负</i></div></article>
          </div>
          <section><div class="section-heading"><div><p class="eyebrow">NEXT MATCH</p><h2>下一场比赛</h2></div></div>${matchCard(match, 0)}</section>
        </div>
      </div>`
  }

  function renderLogin() {
    const comp = competition()
    return `
      <div class="page-view login-wrap">
        <section class="login-panel">
          <div class="login-visual"><span class="hero-mark">${comp.mark}</span><div><h1>你的比赛日，应该更简单。</h1><p>登录后同步关注球队、收藏比赛与开球提醒，在六个赛事世界之间无缝切换。</p></div><small>${comp.name} · ${comp.season}</small></div>
          <form class="login-form" id="demo-login-form"><p class="eyebrow">WELCOME BACK</p><h2>登录 WorldCup Mate</h2><p>这是视觉演示，填写任意内容即可体验。</p>
            <div class="form-stack">
              <label class="field-label">邮箱<div class="field-control">${icon('i-mail')}<input type="email" value="demo@worldcupmate.app" required /></div></label>
              <label class="field-label">密码<div class="field-control">${icon('i-lock')}<input type="password" value="worldcupmate" required /></div></label>
              <div class="form-row"><label><input type="checkbox" checked /> 记住我</label><button class="text-button" type="button" data-demo-action="password">忘记密码？</button></div>
              <button class="primary-button full-button" type="submit">登录并继续 ${icon('i-arrow')}</button>
              <button class="secondary-button full-button" type="button" data-route="profile">返回个人中心</button>
            </div>
          </form>
        </section>
      </div>`
  }

  function renderPage() {
    const renderers = {
      home: renderHome,
      schedule: renderSchedule,
      teams: renderTeams,
      standings: renderStandings,
      profile: renderProfile,
      match: renderMatchDetail,
      team: renderTeamDetail,
      login: renderLogin,
    }
    elements.main.innerHTML = (renderers[state.route] || renderHome)()
    elements.main.setAttribute('aria-busy', 'false')
  }

  function renderAll() {
    renderChrome()
    renderPage()
  }

  function navigate(route, index = 0) {
    state.route = route
    state.detailIndex = Number(index) || 0
    const nextHash = `#/${route}${route === 'match' || route === 'team' ? `/${state.detailIndex}` : ''}`
    if (location.hash !== nextHash) history.pushState(null, '', nextHash)
    renderChrome()
    renderPage()
    window.scrollTo({ top: 0, behavior: matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth' })
    requestAnimationFrame(() => elements.main.focus({ preventScroll: true }))
  }

  function setCompetition(code) {
    if (!competitions[code]) return
    state.competition = code
    safeStorageSet('wm-demo-competition', code)
    closeCompetitionDialog()
    elements.main.setAttribute('aria-busy', 'true')
    renderAll()
    showToast(`已切换到${competition().name} · ${competition().season}`)
  }

  function toggleTheme() {
    state.theme = state.theme === 'dark' ? 'light' : 'dark'
    safeStorageSet('wm-demo-theme', state.theme)
    renderChrome()
    renderPage()
    showToast(state.theme === 'dark' ? '已切换深色模式' : '已切换浅色模式')
  }

  function showToast(message) {
    clearTimeout(toastTimer)
    elements.toast.textContent = message
    elements.toast.classList.add('show')
    toastTimer = setTimeout(() => elements.toast.classList.remove('show'), 1900)
  }

  function openCompetitionDialog() {
    dialogLastFocus = document.activeElement
    elements.dialog.hidden = false
    document.body.style.overflow = 'hidden'
    requestAnimationFrame(() => elements.dialog.querySelector('.competition-option')?.focus())
  }

  function closeCompetitionDialog() {
    if (elements.dialog.hidden) return
    elements.dialog.hidden = true
    document.body.style.overflow = ''
    dialogLastFocus?.focus?.()
  }

  function trapDialogFocus(event) {
    if (elements.dialog.hidden || event.key !== 'Tab') return
    const focusable = [...elements.dialog.querySelectorAll('button:not([disabled])')]
    if (!focusable.length) return
    const first = focusable[0]
    const last = focusable[focusable.length - 1]
    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault()
      last.focus()
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault()
      first.focus()
    }
  }

  document.addEventListener('click', (event) => {
    const routeButton = event.target.closest('[data-route]')
    if (routeButton) {
      navigate(routeButton.dataset.route, routeButton.dataset.index || 0)
      return
    }

    const competitionButton = event.target.closest('[data-competition-code]')
    if (competitionButton) {
      setCompetition(competitionButton.dataset.competitionCode)
      return
    }

    if (event.target.closest('#mobile-competition-trigger')) {
      openCompetitionDialog()
      return
    }

    if (event.target.closest('[data-close-dialog]')) {
      closeCompetitionDialog()
      return
    }

    if (event.target.closest('#theme-toggle')) {
      toggleTheme()
      return
    }

    const filter = event.target.closest('[data-schedule-filter]')
    if (filter) {
      state.scheduleFilter = filter.dataset.scheduleFilter
      renderPage()
      showToast(`赛程筛选：${state.scheduleFilter}`)
      return
    }

    const mode = event.target.closest('[data-standing-mode]')
    if (mode) {
      state.standingsMode = mode.dataset.standingMode
      renderPage()
      return
    }

    const action = event.target.closest('[data-demo-action]')
    if (action) {
      const type = action.dataset.demoAction
      if (type === 'theme') toggleTheme()
      else {
        action.classList.toggle('active')
        const copy = { reminder: '比赛提醒状态已更新', favorite: '收藏状态已更新', password: '已发送演示重置邮件' }
        showToast(copy[type] || '操作已完成')
      }
    }
  })

  document.addEventListener('submit', (event) => {
    if (event.target.matches('#demo-login-form')) {
      event.preventDefault()
      showToast('登录成功，欢迎回来')
      setTimeout(() => navigate('profile'), 500)
    }
  })

  document.addEventListener('input', (event) => {
    if (!event.target.matches('#team-search')) return
    const keyword = event.target.value.trim().toLowerCase()
    document.querySelectorAll('#teams-grid .team-card').forEach((card) => {
      card.hidden = keyword && !card.dataset.search.toLowerCase().includes(keyword)
    })
  })

  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') closeCompetitionDialog()
    trapDialogFocus(event)
  })

  window.addEventListener('popstate', () => {
    const parsed = parseHash()
    state.route = parsed.route
    state.detailIndex = parsed.index
    renderChrome()
    renderPage()
  })

  if (!competitions[state.competition]) state.competition = 'WC'
  if (!['light', 'dark'].includes(state.theme)) state.theme = 'dark'
  renderAll()
})()
