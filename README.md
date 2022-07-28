# Welcome 
## [Documentation](https://maxym.gitbook.io/crypto-mailer/)
``` ENG ```
### **Hi!** This project was written to send the current rate of cryptocurrencies to email
* Backend was written entirely in **Go**, and thoroughly tested in **Postman** :D 
* By convention, this project does not use a database, so email data is stored in a csv file
* The data in this file is stored in order, which allows you to use **binary search** to find duplicates
* The project is written to work out of the box, so all the credentials are public, and some of them are **trial** \
Therefore, if the API stopped working for some reason, first try to set up the [config](https://maxym.gitbook.io/crypto-mailer/reference/setup-config)

\
``` УКР ```
### **Привіт!** Це проєкт для отримання поточного курсу криптовалют і розсилання його по email
* Бекенд повністю написаний на **Go**, і перевірений у **Postman** :D 
* За умовою завдання, не потрібно було підключати бд, тому емейли зберігаються у csv файлі
* Дані в цьому csv файлі зберігаються впорядковано, тому використовується **бінарний пошук** для пошуку дублікатів
* Проєкт написано так, щоб він працював "з коробки", тому деякі облікові дані з **пробним періодом** \
Якщо API перестав працювати як належне, можливо варто переналаштувати [конфіг](https://maxym.gitbook.io/crypto-mailer/reference/setup-config)
