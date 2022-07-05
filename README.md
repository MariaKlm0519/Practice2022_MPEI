# Practice2022_MPEI

## Содержание
1. [Задание](#Task)
2. [Подготовительный этап](#Implement)
 + [Работа со статическими данными](#Static)
 + [Работа с динамическими данными](#Dynamic)
 + [Мини-задача №1](#Task1)
 + [Мини-задача №2](#Task2)
 + [Мини-задача №3](#Task3)
 + [Мини-задача №4](#Task4)
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

#### <a name="Task1"></a> Мини-задача №1. Вывести список служб
Используется пакет golang.org/x/sys/windows/svc/mgr. Для использования данных функций требуются особые права, поэтому необходимо либо запускать приложение от имени администратора, либо передать необходимые права заданному пользователю.
```golang
func ListServices(get_status uint32) []Service {
	m, _ := mgr.Connect()
	names, _ := m.ListServices()
	var result []Service
	for i := 0; i < len(names); i++ {
		serv, err := m.OpenService(names[i])
		if err != nil {
			continue
		}
		status, err := serv.Query()
		if err != nil {
			continue
		}
		if uint32(status.State) == get_status {
			config, err := serv.Config()
			if err != nil {
				continue
			}
			newserv := Service{names[i], config, uint32(status.State), serv}
			result = append(result, newserv)
		}
	}
	return result
}
```

#### <a name="Task2"></a> Мини-задача №2. Изменить ini-файл
Будем использовать пакет gopkg.in/ini.v1.
```golang
cfg, err := ini.Load(name)
if err != nil {
	w.WriteHeader(500)
	return
	}
status := cfg.Section("Options").Key("Enabled").Value()
var new_status string
if status == "1" {
	new_status = "0"
	} else {
	new_status = "1"
	}
cfg.Section("Options").Key("Enabled").SetValue(new_status)
err = cfg.SaveTo(name)
```

#### <a name="Task3"></a> Мини-задача №3. Найти log-файлы, модифицированные в заданном временном диапазоне
Функция возвращает архив с файлами, дата изменения которых находится внутри диапазона, введенного пользователем.
```golang
func listDirByWalk(file_path string, zip_path string, t1 time.Time, t2 time.Time) (*os.File, error) {

	name := time.Now().Format("02012006150405") + ".zip"
	outFile, err := os.Create(zip_path + "\\" + name)
	if err != nil {
		return nil, errors.New("can't create output file")
	}
	zipW := zip.NewWriter(outFile)

	filepath.Walk(file_path, func(wPath string, info os.FileInfo, err error) error {
		if wPath == file_path {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().After(t1) && info.ModTime().Before(t2) {
			dat, _ := ioutil.ReadFile(wPath)
			f, _ := zipW.Create(info.Name())
			f.Write(dat)
		}
		return nil
	})
	err = zipW.Close()
	if err != nil {
		return nil, errors.New("can't close zip writer")
	}
	err = outFile.Close()
	if err != nil {
		return nil, errors.New("can't close output file")
	}
	return outFile, nil
}
```

#### <a name="Task4"></a> Мини-задача №4. Вывести информацию о системе

Для решения будем использовать библиотеку gopsutil.

```golang
hostStat, _ := host.Info()
cpuStat, _ := cpu.Info()
vmStat, _ := mem.VirtualMemory()
diskStat, _ := disk.Usage("\\")

var info SysInfo
info.Hostname = hostStat.Hostname
info.Platform = hostStat.Platform
info.CPU = cpuStat[0].ModelName
info.RAM = vmStat.Total / 1024 / 1024 // в Мб
info.Disk = diskStat.Free / 1024 / 1024  // в Мб
```

Удалось:
1. Вывести список служб с зависимостями
2. Изменить заданный ini-файл
3. Найти log-файлы в заданном временном диапазоне. Заархивировать
4. Вывести некоторую информацию о системе

Не удалось:
1. Передать архив от сервера пользователю

Примерные планы на день:
1. Добавить передачу архива от сервера к пользователю
2. Добавить возможность выбора порта из нескольких
3. Преобразовать консольное приложение к службе
4. Заставить службу в течение своей работы посылать сигнал другой службе

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
