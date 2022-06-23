# Practice2022_MPEI

## Содержание
1. [Задание](#Task)
2. [Подготовительный этап](#Implement)
 + [Работа со статическими данными](#Static)
 + [Работа с динамическими данными](#Dynamic)
 + [Разработка простенького REST API](#Rest_api)
 + [Текущие результаты](#Results)
3. [Ход работы над заданием](#Deal)
4. [Дополнительные теоретические материалы](#Article)

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
В главную функцию добавим обработчики необходимых страниц.

```golang
func main() {
  ...
  mux.HandleFunc("/", home)
  mux.HandleFunc("/postform", postform)
  mux.HandleFunc("/test1", test1)
  mux.HandleFunc("/test2", test2)
  mux.HandleFunc("/test3", test3)
  ...
}
```
Для рендеринга статического содержимого создадим отдельную функцию.
```golang
func (aps *application) render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
	files := []string{
		name,
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	rs, err := template.ParseFiles(files...)
	if err != nil {
		aps.infoLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = rs.Execute(w, td)
	if err != nil {
		aps.infoLog.Println(err.Error())
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
В другой html-странице будем выводить переданный текст.
```golang
func postform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error! Locked.", 423)
		return
	}
	name := r.FormValue("username")
	...
	err = ts.Execute(w, name)
	...
}
```

#### <a name="Dynamic"></a> Работа с динамическими данными
Для работы с динамическими данными можно использовать файлы формата json. Для их обработки в Golang представлена стандартная библиотека "encoding/json".

Для кодирования данных JSON используется Marshal функция.
```golang
func Marshal(v interface{}) ([]byte, error)
```
Для декодирования данных JSON используется Unmarshal функция.
```golang
func Unmarshal(data []byte, v interface{}) error
```
Пакет json обращается только к экспортированным полям struct типов (те, которые начинаются с заглавной буквы). Поэтому в выводе JSON будут присутствовать только экспортированные поля структуры.

Использование тегов в структуре кодируемой в JSON позволяет получить названия полей в результирующем JSON, отличающиеся от названия полей в структуре. В следующем примере в результирующем JSON поле ID_user будет выглядеть как id:
```golang
type Item struct {
 ID_user      uint   `json:"id"`
 Title   string `json:"title"`
}
```
Добавим в проект динамики. Свяжем html-страницы с серверной api. Для этого будем использовать кросс-доменные ajax-запросы.
```js
function Action1Message() {
    $.ajax({
        async: true,
        type: 'get',
        url: 'http://127.0.0.1:4001/api/records',
        crossDomain: true,
        cache:false,
        dataType: 'json',
        success: function (data, textStatus, jqXHR ){
            var obj = JSON.parse(jqXHR.responseText);
            document.getElementById("test").innerHTML = obj.title + " " + obj.text;
        },
        error: function () {
            alert('Failed...');
        }
    });
}
```
При этом разметка html-страницы.
```html
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.0/jquery.min.js"></script>
<script src="/static/js/ajax.js" type="text/javascript"></script>

<form name="f1">
    <label>Title</label>
    <input id="title" name="title" type="text">
    <label>Message</label>
    <textarea placeholder="Введите ваше сообщение" name="text" id="text"></textarea>
</form>

<input onclick="Action1Message()" type="submit" value="Action1" />
<input onclick="Action2Message()" type="submit" value="Action2" /> <br> <br>
```
В данном случае, обработка кросс-доменных запросов (в том числе запросов другим портам) требует либо включения на сервере-ответчике специальных хидеров, либо других методов (например, использование прокси-сервера). В исследовательской мини-задаче пойдем по пути наименьшего сопротивления, тогда обработчик запроса примет следующий вид.
```golang
func getrecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		{
			newq, err := json.Marshal(Quote{2, m.Title, m.Text + time.Now().Format(" 2006-01-02 15:04:05")})
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Write(newq)
			w.WriteHeader(200)
		}
	case http.MethodPost:
		{
			body, _ := ioutil.ReadAll(r.Body)
			_ = json.Unmarshal(body, &m)
			w.WriteHeader(200)
		}
	default: ...
	}
}
```

#### <a name="Rest_api"></a> Разработка простенького REST API


#### <a name="Results"></a> Текущие результаты
Вид главной страницы.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/75f4c14f32e6c7a21feaf36ec264f2e97eb789b9/current_results_pict/%D0%9F%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%BE%D1%81%D0%BD%D0%BE%D0%B2%D0%BD%D0%BE%D0%B9%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD.png)

Вид второй страницы.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/75f4c14f32e6c7a21feaf36ec264f2e97eb789b9/current_results_pict/%D0%9F%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%B2%D1%82%D0%BE%D1%80%D0%BE%D0%B9%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD.png)

Вид третьей страницы.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/75f4c14f32e6c7a21feaf36ec264f2e97eb789b9/current_results_pict/%D0%9F%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B5%20%D1%82%D1%80%D0%B5%D1%82%D0%B8%D0%B9%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD.png)

Вид четвертой страницы.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/75f4c14f32e6c7a21feaf36ec264f2e97eb789b9/current_results_pict/%D0%9F%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B5%20%D1%87%D0%B5%D1%82%D0%B2%D0%B5%D1%80%D1%82%D1%8B%D0%B9%20%D1%8D%D0%BA%D1%80%D0%B0%D0%BD.png)

При попытке перейти по URL.

![](https://github.com/MariaKlm0519/Practice2022_MPEI/blob/961700bfbfc113d1c2a9f8be0cbf8aeba0bddf2e/current_results_pict/URL_%D0%B7%D0%B0%D0%BF%D1%80%D0%BE%D1%81.png)

Удалось:
1. Создать хранилище данных
2. Усовершенствовать обработчики запросов
3. Начать процесс разработки простенького Rest API

Не удалось:
1. Добавить в проект работу с динамическими данными json
2. Создать базу данных, вместо локального хранилища

Необходимо подобрать простую библиотеку по REST API. Возможные варианты:
1. [Gin](https://github.com/gin-gonic/gin)
2. [resty](https://github.com/go-resty/resty)
3. [echo](https://github.com/labstack/echo)
4. Использование стандартной библиотеки

Примерные планы на день:
1. Добавить работу с динамическими данными
2. Продолжить работу над REST API

### <a name="Deal"></a> Ход работы над заданием
...

### <a name="Article"></a> Дополнительные теоретические материалы
Работа со статическими данными в net/http:
1. [Статические файлы](https://metanit.com/go/web/1.3.php)
2. [Обработка статических файлов](https://golangify.com/serving-static-files)

Работа с динамическими данными. Rest API:
1. [Работа с JSON в Golang](https://golang-blog.blogspot.com/2019/11/json-golang.html)
2. [Принципы Rest](https://habr.com/ru/post/590679/)
3. [Разработка Rest-серверов на Golang](https://habr.com/ru/company/ruvds/blog/559816/)
4. [Обработка запросов в Golang, пример](https://uproger.com/vazhnye-konczepczii-obrabotchikov-veb-serverov-v-golang/)

Другое:
1. [Про протокол http](https://habr.com/ru/post/215117/)
2. [Основы HTML](https://html5book.ru/html-html5/)
3. [Основы CSS](https://html5book.ru/css-css3/)
