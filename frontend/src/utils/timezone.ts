import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'

dayjs.extend(utc)
dayjs.extend(timezone)

export function toLocalTime(utcString: string, tz?: string): string {
  const userTz = tz || Intl.DateTimeFormat().resolvedOptions().timeZone
  return dayjs.utc(utcString).tz(userTz).format('M月D日 HH:mm')
}

export function toUTCTime(utcString: string): string {
  return dayjs.utc(utcString).format('YYYY-MM-DD HH:mm')
}

export function formatKickoffTime(utcString: string, tz?: string): string {
  const userTz = tz || Intl.DateTimeFormat().resolvedOptions().timeZone
  return dayjs.utc(utcString).tz(userTz).format('HH:mm')
}

export function formatKickoffDate(utcString: string, tz?: string): string {
  const userTz = tz || Intl.DateTimeFormat().resolvedOptions().timeZone
  return dayjs.utc(utcString).tz(userTz).format('M月D日')
}
