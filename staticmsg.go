package main

const (
	HelpMessageEn string = "*Commands:*\n" +
		"/helpen - This message\n" +
		"/help - руководство на русском\n" +
		"/ck <website or IP>\n/check <website or IP> - check something - IP-address or Domain or URL\n" +
		"/donate - Information about donation methods\n\n" +
		"*Simple usage*\n\n" +
		"Send an IP-address or Domain or URL to the bot for restriction checking \n\n" +
		"Send a record ID (ussualy started with #) to the bot for details\n\n" +
		"Or use /ck or /check commands in group chats\n\n" +
		"*Legend*\n\n" +
		"\U000026d4 URL blocking type. It's default blockig type\n" +
		"\U0001f4db HTTPS blocking type. It's not standart blocking type. It's URL blocking type but with HTTPS scheme\n" +
		"\U0001f6ab Domain blockig type\n" +
		"\u274c IP blocking type\n"

	HelpMessage string = "*Команды:*\n" +
		"/helpen - manual in english\n" +
		"/help - это сообщение\n" +
		"/ck <САЙТ или IP-адрес>\n/check <САЙТ или IP-адрес> - проверить IP адрес, домен или URL\n" +
		"/donate - получить информацию о способах пожертвований\n\n" +
		"*Основные возможности*\n\n" +
		"Отправьте IP-адрес или домен, или URL боту для проверки наличия их в списке блокировок\n\n" +
		"Отправьте ID записи (обозначена '#' в ответах) боту для получения подробностей\n\n" +
		"Или воспользуйтесь командами /ck или /check в группах\n\n" +
		"/n\\_<НОМЕР> - показывает подробную информацию о данном номере записи в реестре\n\n" +
		"*Обозначения*\n\n" +
		"\U000026d4 блокировка по URL. Тип блокировки по умолчанию\n" +
		"\U0001f4db блокировка HTTPS. Неформальный тип блокировки\n" +
		"\U0001f6ab блокировка по домену\n" +
		"\u274c блокировка по IP адресу\n\n" +
		"\u26a0\ufe0f *ВНИМАНИЕ* \u26a0\ufe0f\n\n" +
		"Список запрещенных ресурсов в Беларуси предоставляется без электронной подписи. " +
		"Мы доверяем нашему источнику, но с правовой точки зрения информация бота предназначена" +
		"только для ознакомления. Также в Беларуси не существует рекомендаций по способам блокировок. " +
		"Мы позволили себе интерпретировать виды блокировок самостоятельно, опираясь на российский опыт. " +
		"Бот проверки росcийских блокировок Роскомнадзора: https://t.me/u2ckbot"

	DonateMessage string = "\U00002b50 *Традиционные способы:*\n" +
		"Paypal: https://www.paypal.me/schorsx\n" +
		"ЮMoney (бывшие ЯДеньги): https://sobe.ru/na/m2i2s077M0g2\n" +
		"\U0001f3f5 *Сделать меня криптомагнатом:*\n" +
		"BTC: `18YFeAV12ktBxv9hy4wSiSCUXXAh5VR7gE`\n" +
		"LTC: `LVXP51M8MrzaEQi6eBEGWpTSwckybqHU5s`\n" +
		"ETH: `0xba53cebd99157bf412a6bb91165e7dff29abd0a2`\n" +
		"ZEC: `t1McmUhzdsauoXpiu2yCjNpnLKGGH225aAW`\n" +
		"BCH: `1FiXmPZ6eecHVaZbgdadAuzQLU9kqdSzVN`\n" +
		"ETC: `0xeb990a29d4f870b5fdbe331db90d9849ce3dae77`\n" +
		"TON: `EQBrl8BNLWNVvmSCZDNexzoGQLIojnp4xNDT6Wf4AFX4S\\_57`\n\n"
		//"\U0001f36d Мой вишлист: http://mywishlist.ru/me/schors\n\n"

	Footer string = "" /*"--- \n" +
	"Часть проекта @usher2\n\n" /*+
	"\U000026a0 Я хочу тонко намекнуть на толстые обстоятельства. " +
	"Сейчас весь краудфайндинг в 0 уходит на оплату " +
	"хостинга проектов Эшер II - сайт, сбор выгрузок, бот этот. " +
	"Даже писать посты стало решительно некогда. /donate \n" */

	DonateFooter string = "Хочу новый ноут и мониторы: /donate\n"
)
