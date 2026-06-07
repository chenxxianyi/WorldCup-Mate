const TOKEN_KEY = 'wm-token'
export const LOGOUT_LOGIN_PATH = '/login?logout=1'
export const LOGOUT_QUERY_KEY = 'logout'
export const LOGOUT_QUERY_VALUE = '1'

export function clearAuthStorage() {
  localStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(TOKEN_KEY)
}
