# Practice2022_MPEI

## Содержание
1. [Постановка задачи](#Task)
2. [Описание объекта задания](#Chapter2)
3. [Описание метода решения задачи](#Chapter3)
4. [Экспериментальная часть](#Chapter4)
   + [Работа со статическими данными](#Chapter4.1)
   + [Работа со динамическими данными](#Chapter4.2)
   + [Вывод служб](#Chapter4.3)
   + [Изменение статуса ini-файла](#Chapter4.4)
   + [Поиск log-файлов. Архивирование](#Chapter4.5)
   + [Информация о системе](#Chapter4.1)
5. [Выводы](#End)
6. [Список литературы](#Articles)

**Результаты выполнения индивидуального задания** 

<a name="Task"></a>**1. Постановка задачи**

Разработать  небольшую  службу  для  удаленного  администрирования  со следующими функциями: 
- Вывод  списка  служб  и  их  параметров  с  выбранным  статусом выполнения. 
- Изменение статуса выбранного ini-файла. 
- Поиск  log-файлов,  модифицированных  в  заданный  промежуток времени и возврат архива с файлами клиенту. 
- Вывод информации о системе.

<a name="Chapter2"></a>**2. Описание объекта задания** 

Стандартный  пакет  net/http  предоставляет  возможности  реализации клиент-серверной архитектуры и обработки REST-запросов и является наиболее оптимальным  для  решения  простых  задач.  REST-запросы  позволяют эффективно  обмениваться  данными  через  веб-приложение;  все  необходимые данные передаются в качестве параметров запроса. Ajax-запросы организуют кроссдоменное взаимодействие между клиентом и сервером. 

<a name="Chapter3"></a>**3. Описание метода решения задачи** 

Создаётся  веб-приложение  с  помощью  стандартного  пакета  net/http. Раздача динамических данных производится с помощью данных в формате json, с  обработкой  через  пакет  encoding/json.  Раздача  статических  данных производится на одном выбранном порту, работа с динамическими данными производится на другом.  

Работа  со  службами  реализуется  с  помощью  пакета  mgr.  Работа  с  ini- файлами  осуществляется  через  пакет  gopkg.in/ini.v1.  Поиск  log-файлов производится с помощью функций библиотеки path/filepath, а их архивирование и отправка – с помощью archive/zip. Вывод информации о системе реализуется через библиотеку gopsutil. 

На клиентской стороне используются шаблоны страниц html и таблицы стилей css. Для реализации кроссдоменных запросов на javascript используются ajax-запросы. 

<a name="Chapter4"></a>**4. Экспериментальная часть** 
<a name="Chapter4.1"></a>*1) Работа со статическими данными*

Для решения данной задачи оптимально выбрать пакет net/http из стандартной библиотеки Golang. Для прямой отправки статических файлов в пакете http определена функция FileServer, которая возващает объект Handler:
```golang
func FileServer(root FileSystem) Handler
```
Для нашего приложения создадим папку static, куда поместим все статические файлы, с которыми будем работать. Затем, все запросы, начинающиеся со /static/ будем обрабатывать с помощью FileServer.

```golang
func routes() *http.ServeMux {
   mux := http.NewServeMux()
   fileServer := http.FileServer(http.Dir("./ui/static/"))
   mux.Handle("/static", http.NotFoundHandler())
   mux.Handle("/static/", http.StripPrefix("/static", fileServer))
   return mux
}
```
Для рендеринга статического содержимого создадим отдельную функцию. 

```golang
func render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
   files := []string{
      name,
      "./ui/html/base.layout.tmpl",
      "./ui/html/footer.partial.tmpl",
   }

   rs, err := template.ParseFiles(files...)
   if err != nil {
      log.Println(err.Error())
      http.Error(w, "Internal Server Error", 500)
      return
   }

   err = rs.Execute(w, td)
   if err != nil {
      log.Println(err.Error())
      http.Error(w, "Internal Server Error", 500)
   }
}
```

<a name="Chapter4.2"></a>*2) Работа с динамическими данными*

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

```javascript
function Action1Message() {
    $.ajax({
        async: true,
        type: 'get',
        url: host + "/api/records",
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

В общем случае, система безопасности браузера предотвращает запросы веб-страницы к другому домену, отличному от того, который обслуживает веб- страницу. Поэтому  обработка  кроссдоменных  запросов  (в  числе  которых запросы  к  другим  портам)  требует  или  включения  на  сервере-ответчике специальных  заголовков  ответа,  или  применения  других  методов  (например, использования  прокси-сервера).  Т.к.  приложение  используется  на  локальном хосте, ограничимся первым вариантом и установим необходимые заголовки: 

```golang
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, 
OPTIONS")
w.Header().Set("Access-Control-Max-Age", "1000")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, XRequested-With")
```

<a name="Chapter4.3"></a>*3) Вывод служб*

Для  начала,  создадим  тип,  который  будет  предоставлять  набор необходимой информации о службе. 

```golang
type Service struct {
  Name string `json:"name"`
  Config mgr.Config `json:"config"`
  Status uint32 `json:"status"`
  srv *mgr.Service `json:"service"`
}
```

Для  поиска  служб  будем  использовать  пакет golang.org/x/sys/windows/svc/mgr.  Для  этого  установим  соединение  с диспетчером управления службами. Затем из списка служб, представленных на устройстве,  выберем  нужные,  с  заданным  статусом  и  запросим  у  них необходимые  данные.  Следует  заметить,  что  не  все  службы  могут  быть опрошены, их мы пропускаем. 

```golang
func ListServices(get_status uint32) ([]Service, error) {
   m, err := mgr.Connect()
   defer m.Disconnect()
   if err != nil {
      return nil, errors.New("can't connect to service control manager")
   }
   names, err := m.ListServices()
   if err != nil {
      return nil, errors.New("can't get service list")
   }
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
   return result, nil
}
```

Тогда обработчик запроса примет вид: 
```golang
func getlist(w http.ResponseWriter, r *http.Request) {
   switch r.Method {
   case http.MethodPost:
      {
         body, _ := ioutil.ReadAll(r.Body)
         var s struct {
            Stat string `json:"status"`
         }
         err := json.Unmarshal(body, &s)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         get_status, err := strconv.Atoi(s.Stat)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         data, err := ListServices(uint32(get_status))
         if err != nil {
            w.WriteHeader(500)
            return
         }
         newq, err := json.Marshal(data)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         w.Header().Set("Access-Control-Allow-Origin", "*")
         w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
         w.Header().Set("Access-Control-Max-Age", "1000")
         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
         w.Header().Set("Content-Type", "text/html; charset=utf-8")
         w.Write(newq)
         w.WriteHeader(200)
      }
   default:
      {
         http.Error(w, "Error! Locked.", 423)
         w.WriteHeader(423)
         return
      }
   }
}
```

<a name="Chapter4.4"></a>*4) Изменение статуса ini-файла* 

Ini-файлы  –  текстовые  файлы  особой  структуры,  содержащие конфигурационные параметры некоторых компонентов ОС Windows. 

INI файл может содержать: 

- пустые строки; 
- комментарии — от символа «;» (точка с запятой), стоящего в начале строки, до конца строки; 
- заголовки  разделов —  строки,  состоящие  из  названия  раздела, заключённого в квадратные скобки «[ ]»; 
- значения параметров — строки вида «ключ=значение». 

Пример содержимого ini-файла: 

[Options]
Enabled   = 1
Brandover = 0
Language  = 1033
SkipUAC   = 1

[Settings]
Display   = 1
Root      = /

По договоренности, в рамках задачи за статус ini-файла отвечает параметр Enabled  секции  Options.  Функция  обработки  файла  загружает  файл  по переданному имени. Затем извлекает текущий статус файла, заменяет его на противоположный и сохраняет ini-файл. 

```golang
func ChangeIni(name string) error {
   cfg, err := ini.Load(name)
   if err != nil {
      return errors.New("can't load ini-file")
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
   if err != nil {
      return errors.New("can't save ini-file")
   }
   return nil
}
```

Тогда обработчик примет вид: 

```golang
func setini(w http.ResponseWriter, r *http.Request) {
   switch r.Method {
   case http.MethodPost:
      {
         body, err := ioutil.ReadAll(r.Body)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         var n struct {
            Name string `json:"name"`
         }
         err = json.Unmarshal(body, &n)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         err = ChangeIni(n.Name)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         w.Header().Set("Access-Control-Allow-Origin", "*")
         w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
         w.Header().Set("Access-Control-Max-Age", "1000")
         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
         w.Header().Set("Content-Type", "text/html; charset=utf-8")
         w.WriteHeader(200)
      }
   default:
      {
         http.Error(w, "Error! Locked.", 423)
         w.WriteHeader(423)
         return
      }
   }
}
```

<a name="Chapter4.5"></a>*5) Поиск log-файлов. Архивирование.*

Сначала создадим архив, в который будем копировать подходящие log- файлы. В качестве имени архива будем использовать текущие дату и время. 

Поиск файлов организуем рекурсивно, с обходом всех файлов и папок внутри заданной директории. Подходящие по времени файлы будем добавлять в архив. 

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

Тогда обработчик примет вид: 

```golang
func searchLog(w http.ResponseWriter, r *http.Request) {
   switch r.Method {
   case http.MethodPost:
      {
         body, _ := ioutil.ReadAll(r.Body)
         var dat struct {
            Date_start string `json:"date_start"`
            Time_start string `json:"time_start"`
            Date_end   string `json:"date_end"`
            Time_end   string `json:"time_end"`
         }
         _ = json.Unmarshal(body, &dat)
         time_before, _ := time.ParseInLocation("2006-01-02 15:04", dat.Date_start+" "+dat.Time_start, time.Local)
         time_after, _ := time.ParseInLocation("2006-01-02 15:04", dat.Date_end+" "+dat.Time_end, time.Local)
         outfile, err := listDirByWalk("\\Документы\\Project_goland\\logs", "\\Документы\\Project_goland\\archive", time_before, time_after)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         fileBytes, _ := ioutil.ReadFile(outfile.Name())
         w.Write(fileBytes)
         w.Header().Set("Access-Control-Allow-Origin", "*")
         w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
         w.Header().Set("Access-Control-Max-Age", "1000")
         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
         w.Header().Set("Content-Type", "application/zip")
         w.WriteHeader(200)
      }
   default:
      {
         http.Error(w, "Error! Locked.", 423)
         w.WriteHeader(423)
         return
      }
   }
}
```

<a name="Chapter4.6"></a>*6) Информация о системе.* 

Выведем некоторую информацию о системе. Для начала, создадим тип, который будет предоставлять набор необходимой информации о системе. 

```golang
type SysInfo struct {
   Hostname string `json:"hostname"`
   Platform string `json:"platform"`
   CPU      string `json:"cpu"`
   RAM      uint64 `json:"ram"`
   Disk     uint64 `json:"disk"`
}
```

Затем используем библиотеку gopsutil для поиска необходимых данных. 
```golang
func GetSystemInfo() (SysInfo, error) {
   var info SysInfo
   hostStat, err := host.Info()
   if err != nil {
      return info, errors.New("can't get system host info")
   }
   cpuStat, err := cpu.Info()
   if err != nil {
      return info, errors.New("can't get system cpu info")
   }
   vmStat, err := mem.VirtualMemory()
   if err != nil {
      return info, errors.New("can't get system memory info")
   }
   diskStat, err := disk.Usage("\\")
   if err != nil {
      return info, errors.New("can't get system disk info")
   }

   info.Hostname = hostStat.Hostname
   info.Platform = hostStat.Platform
   info.CPU = cpuStat[0].ModelName
   info.RAM = vmStat.Total / 1024 / 1024
   info.Disk = diskStat.Free / 1024 / 1024
   return info, nil
}
```

Тогда обработчик запросов примет вид: 

```golang
func systemInfo(w http.ResponseWriter, r *http.Request) {
   switch r.Method {
   case http.MethodGet:
      {
         info, err := GetSystemInfo()
         if err != nil {
            w.WriteHeader(500)
            return
         }
         data, err := json.Marshal(info)
         if err != nil {
            w.WriteHeader(500)
            return
         }
         w.Write(data)
         w.Header().Set("Access-Control-Allow-Origin", "*")
         w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
         w.Header().Set("Access-Control-Max-Age", "1000")
         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
         w.WriteHeader(200)
      }
   default:
      {
         http.Error(w, "Error! Locked.", 423)
         w.WriteHeader(423)
         return
      }
   }
}
```

<a name="End"></a>**5. Выводы** 

Разработано  веб-приложение  в  соответствии  с  поставленной  задачей. Обработка запросов осуществляется на сервере на отдельном порте. Обработка статических  данных  и  обмен  данными  с  пользователем  –  на  другом  порте. Обеспечено кроссдоменное взаимодействие между клиентом и сервером. Для лучшего  визуального  восприятия  на  клиентской  стороне  использованы  html- шаблоны страниц и таблицы стилей  css. Для динамичности добавлен код на javascript. Приложение успешно протестировано. 

<a name="Articles"></a>**Список литературы** 

1. Обработка статических данных // Golangify URL: [https://golangify.com/serving-static-files ](https://golangify.com/serving-static-files)
1. Хабр: [сайт] URL:[ https://habr.com/ru/company/ruvds/blog/559816/ ](https://habr.com/ru/company/ruvds/blog/559816/)- Перевод статьи:[ Разработка REST-серверов на Go. Часть 1: стандартная библиотека ](https://habr.com/ru/company/ruvds/blog/559816/)
1. jQuery [сайт] URL:[ https://api.jquery.com/jquery.ajax/ ](https://api.jquery.com/jquery.ajax/)- jQuery API ajax Documentation 
1. Блог о языке программирования Go [сайт] URL:[ https://golang- blog.blogspot.com/2019/11/json-golang.html ](https://golang-blog.blogspot.com/2019/11/json-golang.html)- Работа с json в Golang 
1. HTML5BOOK [сайт] URL:[ https://html5book.ru/ ](https://html5book.ru/)- Основы HTML, основы CSS 
1. Максим Жашкевич Язык Go для начинающих. - 1-е изд. - 2020. - 109 с. 


