# Леонид Дятлов - "Система документации IT-инфраструктуры PVA"

### Группа: 11И2
### Электронная почта: dyatlovleonid@gmail.com
### Telegram: https://t.me/bakane6030


### [ Сценарий 1 - Регистрация пользователя ]

1. Пользователь заходит на страницу приложения, если он не вошел в аккаунт, он будет переброшен на страницу регистрации
2. Пользователь вводит логин, с которым он будет создавать аккаунт
3. Пользователь вводит пароль, с которым он будет создавать аккаунт
4. Пользователь нажимает кнопку "Register"
5. Если выбранный логин уже существует в системе, то пользователю сообщается об этом и поля логина и пароля сбрасываются
6. Если заново выбранный логин существует в системе, повторяется пункт 5, пока он не введет уникальный логин
7. При успешной валидации ввода высвечивается сообщение об успешной регистрации, поля ввода сбрасываются,
пользователь может зайти в аккаунт с помощью полей ввода

### [ Сценарий 2 - Вход пользователя ]

1. Пользователь заходит на страницу приложения, если он не вошел в аккаунт, он будет переброшен на страницу регистрации
2. Пользователь вводит логин, с которым он будет входить
3. Пользователь вводит пароль, с которым он будет входить
4. Пользователь нажимает кнопку "Log in"
5. Если логин не найден, пользователь получает информацию об ошибке; поля сбрасываются (Нельзя сообщать точную ошибку, так как
она может помочь злоумышленнику подобрать вход в аккаунт)
6. Если пароль неправильный, повторяется пункт 5, пока он не введет правильную информацию
7. При успешном входе пользователь перенаправляется на домашнюю страницу. Доступ ко всему контенту вне страниц аутентификации
проверяется с помощью подписанного токена, хранящегося в HttpOnly Secure куки.

### [ Сценарий 3 - Просмотр устройств на локальной сети ]

1. Пользователь переходит на URL "/devices" по ссылке или введя адрес в браузере.
2. Если куки аутентификации недействителен или подписан не сервером, пользователь перенаправляется на страницу аутентификации.
3. При наличии и валидности токена пользователь остается на странице.
4. Клиент отправляет запрос на сервер после загрузки страницы у пользователя
5. На сервере с момента запуска работает поток, каждые 2 минуты проводящий сканирование локальной сети с помощью протокола ICMP.
Сервер считывает результаты последнего сканирования, а так же берет данные об устройствах из базы данных (название, описание).
Если сканер нашел устройство, содержащееся в базе данных, рядом с информацией о нем отображается зеленый индикатор. Если нет, то серый.
Если устройство нашел сканер, но его нет в базе данных, оно отображается с зеленым индикатором и пустыми полями Name и Description.
6. Клиент получает данные об устройствах с сервера (ID, IP, имя, описание).
7. Данные пишутся в таблицу на странице вместе с индикатором и кнопками "Edit" и "Delete". Пользователь может искать нужное ему устройство
и менять имена и описания устройств, а так же удалять имена и описания из базы.

### [ Сценарий 4 - Просмотр паролей ]

1. Пользователь переходит на URL "/vaults" по ссылке или введя адрес в браузере.
2. Если куки аутентификации недействителен или подписан не сервером, пользователь перенаправляется на страницу аутентификации.
3. При наличии и валидности токена пользователь остается на странице.
4. Клиент отправляет серверу запрос о получении паролей, к которым есть доступ у пользователя.
5. С помощью токена расшифровывается private key пользователя
6. С помощью private key пользователя расшифровываются ключи сейфов (папок для паролей), к которым пользователь имеет доступ.
7. С помощью ключей сейфа расшифровываются сами пароли
8. Сервер отправляет ответ с паролями в сейфах, клиент отображает их в таблице

### [ Сценарий 5 - Создание сейфов ]

1. Пользователь переходит на URL "/vaults" по ссылке или введя адрес в браузере.
2. Если куки аутентификации недействителен или подписан не сервером, пользователь перенаправляется на страницу аутентификации.
3. При наличии и валидности токена пользователь остается на странице.
4. Пользователь нажимает на кнопку "+ New Vault", перед ним открывается диалог
5. Пользователь вводит название и описание сейфа
6. При успешной валидации ввода кнопка "Create Vault" становится доступной для нажатия
7. Пользователь нажимает на кнопку "Create Vault". Сервер создает случайный ключ сейфа. Ключ шифруется public key всех админов и пользователя, который создает
сейф. Зашифрованные ключи и информация о сейфе хранятся в базе. 
8. Сервер возвращает статус операции, ошибки отображаются пользователю. Пользователь видит новый сейф.

### [ Сценарий 6 - Поделиться сейфом ]

1. Пользователь переходит на URL "/vaults" по ссылке или введя адрес в браузере.
2. Если куки аутентификации недействителен или подписан не сервером, пользователь перенаправляется на страницу аутентификации.
3. При наличии и валидности токена пользователь остается на странице.
4. Пользователь находит сейф в таблице и жмет справа от него кнопку "Share" (иконка).
5. Пользователь вводит имя пользователя, с которым хочет поделиться и жмет кнопку "Share this vault".
6. Сервер с помощью токена расшифровывает private key пользователя, а с его помощью - ключ сейфа. Ключ сейфа шифруется public key
указанного пользователя и хранится в базе.
7. Сервер возвращает клиенту статус операции, он отображается пользователю.


### [ Сценарий 7 - Редактирование сейфов/паролей ]

1. Пользователь переходит на URL "/vaults" по ссылке или введя адрес в браузере.
2. Если куки аутентификации недействителен или подписан не сервером, пользователь перенаправляется на страницу аутентификации.
3. При наличии и валидности токена пользователь остается на странице.
4. Пользователь находит сейф или пароль и жмет кнопку с карандашом справа от него.
5. Пользователь видит диалоговое окно, где он вводит новое имя/описание для сейфа, либо новое имя/описание/пароль для пароля.
6. Пользователь жмет кнопку "Edit this vault"/"Edit this password". Сервер обновляет данные. Если при изменении пароля меняется само значение пароля,
новый пароль шифруется зашифрованным ключом сейфа из сценария 5.
7. Сервер возвращает статус операции, ошибки отображаются пользователю. Пользователь видит измененный сейф/пароль.

