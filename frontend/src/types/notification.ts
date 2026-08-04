/** Backend notification record (models.Notification JSON). */
export interface ApiNotification {
  id: number
  user_id: number
  title: string
  content: string
  type: string
  is_read: boolean
  created_at: string
}
