freebsd:
	mkdir -p tmp/freebsd/usr/local/etc/rc.d
	mkdir -p tmp/freebsd/usr/local/sbin/
	cp init/freebsd/gpd tmp/freebsd/usr/local/etc/rc.d/
	GOOS=freebsd GOARCH=amd64 go build -o tmp/freebsd/usr/local/sbin/gpd ./cmd/gpd
	cd tmp/freebsd && tar -jcf freebsd.tar.bz2 usr 