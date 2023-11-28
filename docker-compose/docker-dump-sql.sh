#!/bin/bash
username="root"
password="123456"
dbname="odisk"
dockername="db"
dumptime=$(date +"%Y%m%d_%H:%M")
dump="mysqldump -u${username} -p${password} ${dbname}"
outputfile="dump_${dbname}_${dumptime}.sql"
docker exec -it ${dockername} ${dump} > ${HOME}/${outputfile}
