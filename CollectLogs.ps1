Start-Process -FilePath ".\collectLog.exe" `
    -ArgumentList '"\\server1\Log, \\server2\Log"' `
    -Wait -RedirectStandardOutput ".\log.txt"