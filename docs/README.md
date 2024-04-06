# Github's Event Files List
GitHub provides 15+ event types, which range from new commits and fork events, to opening new tickets, commenting, and adding members to a project. These events are aggregated into hourly archives, which you can access with any HTTP client:GitHub provides 15+ event types, which range from new commits and fork events, to opening new tickets, commenting, and adding members to a project. These events are aggregated into hourly archives, which you can access with any HTTP client:

|Query|Command|
|---|---|
|Activity for 1/1/2015 @ 3PM UTC|wget https://data.gharchive.org/2015-01-01-15.json.gz|
|Activity for 1/1/2015|wget https://data.gharchive.org/2015-01-01-{0..23}.json.gz|
|Activity for all of January 2015|wget https://data.gharchive.org/2015-01-{01..31}-{0..23}.json.gz|

Use the following bash commmand to download the entire event files:
```bash
for i in {01..31}; do
	for j in {0..23}; do
		wget wget https://data.gharchive.org/2015-01-${i}-${j}.json.gz
    done
done
```