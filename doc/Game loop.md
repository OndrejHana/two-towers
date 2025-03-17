Hráči budou klikat na možnosti v UI, tím se budou posílat zprávy na websocket se serverem

Na serveru bude běžet game loop, která má:
- Player Actions channel: chan, na který se budou agregovat veškeré akce připojených hráčů. Tenhle kanál bude odebírat a akce zapisovat do stavu hry
- Effect timer: Channel na který se agregují timery jako spawn vojáků věžemi. 

## State
stav hry obsahuje:
- mapa: 2D pole
- hráči: pole hráčů, kde každý má `PlayerID` a odkaz na hráče pro kontrolu sessiony při akcích
- věže: pole věží, kde každá věž má `ID`, `PlayerID` a odkaz na channel, na kterém běží její spawner. Ten se resetuje když věž změní hráče. 