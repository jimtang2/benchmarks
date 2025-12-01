# Tuning Notes


### Watch connections

```
watch -n 1 "psql -U postgres -d demo -c \"SELECT count(*) FROM pg_stat_activity WHERE state='active'; SELECT * FROM pg_stat_bgwriter;\""
```