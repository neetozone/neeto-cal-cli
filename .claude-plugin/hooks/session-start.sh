#!/bin/sh
# NeetoCal CLI — session-start hook for Claude Code
# Lightweight auth liveness check. Always exits 0 (informational).

if ! command -v neetocal >/dev/null 2>&1; then
  echo "NeetoCal CLI is not installed or not on PATH."
  exit 0
fi

if neetocal whoami >/dev/null 2>&1; then
  echo "NeetoCal plugin active."
else
  echo "NeetoCal CLI installed but not authenticated. Run 'neetocal login' to authenticate."
fi

exit 0
