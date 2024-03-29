![emergency-service](assets/banner.png)

## Emergency Service

This system is your trusty and reliable aid designed to keep receivers safe and informed during emergencies. 
It promptly sends messages by SMS, Emails about on-going incidents and safety tips, ensuring everyone stays in the loop with live updates. 
Never feel unprepared with our Emergency Service Notification System at hand!

## Run
```
go run cmd/app/main.go
```

## ENVS:
```                                           
export PORT='rpcuser'
export DATABASE_URL='rpcpass'
export LISTEN_ADDR='http://localhost:18334'
export EMAIL_API_KEY='63ccb9a43cd1b'
export PUBLIC_API_KEY='pubkey-2084d9c'
export EMAIL_DOMAIN='emergency-message.com'
export MOBILE_TWIL_ACCOUNT_SID='3537af2e99b5'
export MOBILE_TWIL_AUTH_TOKEN='3537af2e99b5'
export MOBILE_PHONE_EMERGENCY_SERVICE='+783172873'
```  

## Workflow
![emergency-service](assets/workflow.svg)

## Libriraries
```
route: go-chi
migrate: golang-migrate
postgresql: jackc/pgx
sending message by email: mailgun-go
logs: logrus
sql-queries: uptrace/bun
```

# TODO
- [x] graceful shutdown
- [x] add grpc
- [x] remove lorgus, use slog
- [x] add kafka
- [x] add worker to send messages to clients
- [x] add worker to update a message status
- [ ] add parser from forecasting service (gismeteo)
- [ ] send messages by telegram
- [ ] send messages by what's app