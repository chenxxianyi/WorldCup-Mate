#!/usr/bin/env bash
set -euo pipefail

echo "==> Scanning for prohibited files in git index..."

errors=0

# node_modules
count=$(git ls-files frontend/node_modules 2>/dev/null | wc -l)
if [ "$count" -gt 0 ]; then
  echo "ERROR: frontend/node_modules/ is still tracked by git ($count files)"
  errors=$((errors + 1))
fi

# dist
count=$(git ls-files frontend/dist 2>/dev/null | wc -l)
if [ "$count" -gt 0 ]; then
  echo "ERROR: frontend/dist/ is still tracked by git ($count files)"
  errors=$((errors + 1))
fi

# uploads
count=$(git ls-files backend/uploads 2>/dev/null | wc -l)
if [ "$count" -gt 0 ]; then
  echo "ERROR: backend/uploads/ is still tracked by git ($count files)"
  errors=$((errors + 1))
fi

# .env files (except .env.example)
count=$(git ls-files '.env' 2>/dev/null | wc -l)
if [ "$count" -gt 0 ]; then
  echo "ERROR: .env files are tracked by git ($count files)"
  errors=$((errors + 1))
fi

# private keys / secrets in tracked files.
# Flag hardcoded values, not env-var interpolations like ${JWT_SECRET}. Values to
# catch: `=password`, `=change_this`, `default_secret`, actual JWT/API keys.
if git grep -lE 'JWT_SECRET=[^$]|MYSQL_ROOT_PASSWORD=[^$]|SMTP_PASSWORD=[^$]|FOOTBALL_DATA_API_KEY=[^$=]' -- '*.yml' '*.yaml' '*.env' '*.json' '*.go' 2>/dev/null | grep -v '\.env\.example' | grep -q .; then
  echo "ERROR: potential hardcoded secrets found in tracked config files"
  echo "Files:"
  git grep -lE 'JWT_SECRET=[^$]|MYSQL_ROOT_PASSWORD=[^$]|SMTP_PASSWORD=[^$]|FOOTBALL_DATA_API_KEY=[^$=]' -- '*.yml' '*.yaml' '*.env' '*.json' '*.go' 2>/dev/null | grep -v '\.env\.example'
  errors=$((errors + 1))
fi

if [ "$errors" -gt 0 ]; then
  echo ""
  echo "Scan FAILED with $errors error(s). Fix before committing."
  exit 1
fi

echo "Scan passed: no prohibited files in git index."
