export interface PostMatchSummary {
  summary: string
  score_line: string
  key_takeaways: string[]
  qualification_impact: string
  worth_watching: string
  spoiler_level: string
  data_note: string
  generated_at: string
}

export function hasPostMatchSummary(value: any): value is PostMatchSummary {
  return Boolean(value?.summary || value?.score_line)
}
