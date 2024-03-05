---
marp: true
_class: lead
style: @import url('https://unpkg.com/tailwindcss@^2/dist/utilities.min.css');
backgroundImage: url('https://marp.app/assets/hero-background.svg')


---

# Включить Скринкаст

---

# Заполнить на портале telegram и gitlab

---

# Отметиться на лекции

(мы очень хотим увидеть ваши отзывы о лекциях)

---

# Преподаватели

<div class="grid grid-cols-4 gap-4">

<div>

![](../common/static/01_lecture/a.sulaev.png)

<center>

### Сулаев<br>Антон

###### Тимлид разработки backend'a Почты

</center>
</div>
<div>

![](../common/static/01_lecture/d.fedorova.png)

<center>

### Федорова<br>Дарья

###### Ведущий разработчик команды backend'а Почты

</center>

</div>
<div>

![](../common/static/01_lecture/a.kiselev.png)

<center>

### Киселев<br>Андрей

###### Разработчик кросс-командных решений для backend'а Почты

</center>

</center>
</div>
<div>

![](../common/static/01_lecture/a.sazonov.jpg)

<center>

### Персиянова<br>Вероника

###### Разработчик команды backend'a Почты

</div>

</div>

---

# Организационное

- Это факультатив. ходить на него - ваше решение
- Можно не делать домашку, если она вам не нравится
- Можно вообще не делать домашки
- Но если вы делаете домашку - вы её делаете сами
- За списывание отчисляем с курса

---

# Почему golang?

---

![bg left](../common/static/01_lecture/industry_yoy.svg)

# Где больше всего используется golang

(<https://go.dev/blog/survey2021-results>)

---

![bg left](../common/static/01_lecture/app_yoy.svg)

# Для чего больше всего используется golang

(<https://go.dev/blog/survey2021-results>)

---

# Кто использует в VK

- Почта
- Реклама
- Юла
- Delivery Club
- Облако
- Медиапроекты
- Процессинг
- vkontakte
- Внутренняя разработка

---

# Кто использует в России

<div class="grid grid-cols-2 gap-4">
<div>

- VK
- 2GIS
- Acronis
- avito.ru
- Gett
- ITooLabs
- Izvestia
- iSpring
- mc² software

</div>
<div>

- OZON.ru
- Сбермаркет
- Positive Technologies
- PostmanQ - High performance Mail Transfer Agent (MTA)
- ThetaPad
- Tinkoff
- Tochka
- TRY.FIT

</div>
</div>

---

# Краткая история golang

## Разрабатывался ветеранами индустрии

- Кен Томпсон (UNIX, UTF-8, C)
- Роб Пайк (UTF-8, Plan 9, Inferno)
- Роберт Гризмер (Java HotSpot, Sawzall, распределённые системы Google)

## В условиях большой компании

- Много кода
- Много программистов
- Много серверов (а на них много ядер)
- Есть легаси код

---

# Про сервера и ядра

![bg left:55% 105%](../common/static/01_lecture/50-years-processor-trend.png)

(<https://github.com/karlrupp/microprocessor-trend-data>)

---
<!-- https://www.quora.com/Scala-vs-Go-Could-people-help-compare-contrast-these-on-relative-merits-demerits -->
# Реалии разработки*

- Меньше кода - проще для понимания. но есть грань.
- Код читается чаще чем пишется
- Код часто живёт дольше, чем мы предполагаем
- Человек, который тестирует или поддерживает код, чаще всего не его первоначальный автор
- Средний уровень разработчика, который читает, пишет, поддерживает или тестирует код - “не эксперт”

---

# Много кода, много программистов, legacy

- Простой и компактный синтаксис
- Мало магии и синтаксического сахара
- Нет даже тернарного оператора
- Ориентация на простоту и читабельность кода
- Жесткий стиль кода и инструмент для авто-форматирования
- Множество синтаксических анализаторов

---

# Много кода, много программистов, legacy

- Быстрая компиляция
  - Веб-апи почты собирается за 2 минуты
  - Микросервисы собираются за 1 минуту
- Тесты из коробки
  - Вместе с покрытием и отчётом
  - Вместе с бенчмарками
- Профилировщик из коробки

---

# Много серверов

- Статический бинарь
  - Нет dependency hell
- Удобная работа с зависимостями
  - Просто кладём их в репозиторий
  - Просто скачиваем их с git
- Кросс-компиляция

---

# Много ядер в процессоре

![bg left:40% 100%](../common/static/01_lecture/50-years-processor-trend.png)

- Асинхронный i/o на уровне языка (по модели CSP)
  - Нет callback-hell
  - Весь код стандартной библиотеки и внешних либ тоже!

- Приложение масштабируется на все ядра процессора

---

# 10 кубиков программиста

- Бизнесу не нужно самовыражение программиста и его игры в песочнице

<table >
<tr>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#fff2cc">магия</td>
    <td bgcolor="#fff2cc">рантайм</td>
    <td bgcolor="#fff2cc">стиль</td>
    <td bgcolor="#fff2cc">отладка</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
</tr>
</table>

<br>

- Бизнесу нужен продукт
- На go скучно заниматься чем-то кроме продукта
  - Потому что там не так прикольно делать всякие классные штуки с новым синтаксисом

<table>
<tr>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#cfe2f3">язык</td>
    <td bgcolor="#fff2cc">магия</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
    <td bgcolor="#d9ead3">задача</td>
</tr>
</table>
