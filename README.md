 # GoTrackDb
 GoTrackDb is an application written in Go that receives location data from clients via HTTP POST requests (values of "application/x-www-form-urlencoded") and stores it in a time-series database (InfluxDB).

Edit the config.file according to your needs. The XML file data.xml contains markers for the current location data and can be used to display it on a map.

**config.file**
 ```
dataxmlpathlinux "/var/www/html/data.xml"
dataxmlpathother "data.xml" 
pattern "/req"
port "22222"
dbhost "127.0.0.1"
dbport "8086"
dbuser "admin"
dbpass "admin"
dbname "locdb"
```


 # Future improvements

- Implement communication over HTTPS.


 # Contributing
 Feedbacks and recommendations are welcomed.

 # Licensing
 This project is licensed under the GNU GPLv3 License - see the LICENSE file for details.
