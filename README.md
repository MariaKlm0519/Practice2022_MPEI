# Practice2022_MPEI

## Содержание
1. [Задание](#Task)
2. [Исследовательская часть](#Implement)
 + [Работа со статическими данными](#Static)
 + [Работа с динамическими данными](#Dynamic)
 + [Текущие результаты](#Results)
3. [Дополнительные теоретические материалы](#Article)

### <a name="Task"></a> Задание
Разработать небольшую службу для удаленного администрирования. Подзадачи: 
1. Найти и исследовать библиотеку раздачи статических данных; 
2. Найти и исследовать библиотеку для разработки Rest API (не навороченную);

### <a name="Implement"></a> Исследовательская часть
#### <a name="Static"></a> Работа со статическими данными
Для решения данной задачи оптимально выбрать пакет net/http из стандартной библиотеки Golang. Для прямой отправки статических файлов в пакете http определена функция FileServer, которая возващает объект Handler:
```golang
func FileServer(root FileSystem) Handler
```
Для нашего приложения создадим папку static, куда поместим все статические файлы, с которыми будем работать. Затем, все запросы, начинающиеся со /static/ будем обрабатывать с помощью FileServer.
    
```golang
func main() {
  mux := http.NewServeMux()

  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static", http.NotFoundHandler())
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))
  
  log.Println("Запуск сервера на http://127.0.0.1:4000")
  err := http.ListenAndServe(":4000", mux)
  log.Fatal(err)
}
```
В главную функцию добавим обработчик домашней страницы.

```golang
func main() {
  ...
  mux.HandleFunc("/", home)
  ...
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
```
В документе html домашней страницы создадим поле ввода и кнопку, при нажатии на которую пользователь будет переброшен на другую страницу.
```html
<form method="POST" action="postform"> <br>
            <label>Введите ваше имя</label><br>
            <input type="text" name="username" /><br><br>
            <input type="submit" value="Отправить" />
        </form>
```
```golang
func main() {
  ...
  mux.HandleFunc("/postform", postform)
  ...
}
func postform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Ошибка!", 405)
		return
	}
	name := r.FormValue("username")
  fmt.Fprintf(w, "У тебя всё получится,  %s !", name)
}
```
#### <a name="Dynamic"></a> Работа с динамическими данными
Для отправки запросов в пакете net/http определен ряд функций:

```golang
func Get(url string) (resp *Response, err error)
func Head(url string) (resp *Response, err error)
func Post(url string, contentType string, body io.Reader) (resp *Response, err error)
func PostForm(url string, data url.Values) (resp *Response, err error)
```
+ Get(): отправляет запрос GET
+ Head(): отправляет запрос HEAD
+ Post(): отправляет запрос POST
+ PostForm(): отправляет форму в запросе POST

#### <a name="Results"></a> Текущие результаты
Вид главной страницы.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/961700bfbfc113d1c2a9f8be0cbf8aeba0bddf2e/current_results_pict/%D0%9F%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%BE%D1%81%D0%BD%D0%BE%D0%B2%D0%BD%D0%BE%D0%B9%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD.png)

При нажатии на кнопку.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/961700bfbfc113d1c2a9f8be0cbf8aeba0bddf2e/current_results_pict/POST_%D0%B7%D0%B0%D0%BF%D1%80%D0%BE%D1%81.png)

При попытке перейти по URL.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/961700bfbfc113d1c2a9f8be0cbf8aeba0bddf2e/current_results_pict/URL_%D0%B7%D0%B0%D0%BF%D1%80%D0%BE%D1%81.png)

Удалось:
1. Создать репозиторий с проектом на гитхаб
2. Обработчики пары запросов
3. Теоретические основы Rest API

Неудалось (посмотрю на выходных):
1. Подключить html-разметку к странице /postform
2. Общая функция обработки запросов? Больше вариантов запросов?

### <a name="Article"></a> Дополнительные теоретические материалы
Работа со статическими данными в net/http:
1. [Статические файлы](https://metanit.com/go/web/1.3.php)
2. [Обработка статических файлов](https://golangify.com/serving-static-files)

Работа с динамическими данными. Rest API:
1. [Принципы Rest](https://habr.com/ru/post/590679/)
2. [Разработка Rest-серверов на Golang](https://habr.com/ru/company/ruvds/blog/559816/)
3. [Обработка запросов в Golang, пример](https://uproger.com/vazhnye-konczepczii-obrabotchikov-veb-serverov-v-golang/)

Другое:
1. [Про протокол http](https://habr.com/ru/post/215117/)
2. [Основы HTML](https://html5book.ru/html-html5/)
3. [Основы CSS](https://html5book.ru/css-css3/)
