# shiny-broccoli
SNMP proxy server SNMPv2c -> SNMPv3


# links

1. Python SNMP server: https://github.com/delimitry/snmp-server
2. Golang SNMP server: https://github.com/slayercat/GoSNMPServer
3. Golang SNMP lib: https://github.com/gosnmp/gosnmp

# Algoritm

1. Слушаем сеть
2. Получили SNMP-запрос
3. Вытащили из полученного SNMP-запроса параметры запроса
4. Из полученных параметров сформировали новый запрос SNMPv3
5. Отправили запрос SNMPv3
6. Ждем
7. Получили ответ SNMPv3
8. Вытащили данные из полученного ответа
9. Упаковали данные в ответ SNMPv2c
10. Продолжаем слушать запросы 

# Algoritm 2

1. Слушаем N портов. На каждый порт приходят SNMPv2c запросы для соотетствующих устройств.
2. ПАралелльно ведется опрос устройств по SNMPv3