# ipcam-alarmserver

Motion alerts to shinobi
```
curl http://<host>:8080/<api_key>/motion/<group>/<monitor_id>?data={"plug":"<monitor_id>","name":"fuCdJ","reason":"motion","confidence":200}
```

# Docker
```
docker run -it -p 15002:15002 -v $(pwd)/config.json:/app/config.json shreddedbacon/ipcam-alarmserver:latest
```