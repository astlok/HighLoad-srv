run:
	sudo docker build -t highloadsrv:image .
	sudo docker run -d -p 8080:8080 highloadsrv:image

grafana:
	sudo docker run -d -p 3000:3000 grafana/grafana-enterprise

node:
	sudo docker run -d -p 9100:9100 prom/node-exporter

prom:
	sudo docker run -p 9090:9090 -d --name prometheus -v $(shell pwd)/prometheus:/etc/config prom/prometheus --config.file=/etc/config/prometheus.yml

stop:
	sudo docker rm -f $(shell sudo docker ps -aq)
del:
	sudo docker rmi -f $(shell sudo docker images -a -q)