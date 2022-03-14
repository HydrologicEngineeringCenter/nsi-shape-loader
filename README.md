go run . -u admin -p notPassword -n gis -t host.docker.internal -o 25432 -a nsi_v2022 -d /workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp

./upload -s "host=host.docker.internal port=25432 user=admin password=notPassword dbname=gis" -d . -c nsi -t table

ogrinfo -so -al test/nsi/NSI_V2_Archives/V2022/15001.shp

THE GO CLI TOOL IS DEPRECATED IN FAVOR OF THE MIX GO/OGR2OGR BASH SCRIPTS
