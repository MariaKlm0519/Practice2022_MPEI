let host = "http://127.0.0.1" + ":4001";

function Action1Message() {
    $.ajax({
        async: true,
        type: 'get',
        url: host + "/api/records", //'http://127.0.0.1:4001/api/records',
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

function Action2Message(){
    jQuery.support.cors = true;
    let title = document.forms.f1.elements.title.value;
    let txt = document.forms.f1.elements.text.value;
    var data = JSON.stringify({"title": title, "text" : txt});
    $.ajax({
        async: true,
        type: 'post',
        url: host + "/api/records", //'http://127.0.0.1:4001/api/records',
        crossDomain: true,
        cache: false,
        data: data,
        processData: false,
        success: function (){
            alert("Данные отправлены на сервер");
        },
        error: function () {
            alert('Failed...');
        }
    });
    document.forms.f1.elements.title.value = ""
    document.forms.f1.elements.text.value = ""
    return false;
}

function Action3Message() {
    jQuery.support.cors = true;
    let st = document.forms.f2.elements.status.value;
    var data = JSON.stringify({"status": st});
    $.ajax({
        async: true,
        type: 'post',
        url: 'http://127.0.0.1:4001/api/list',
        crossDomain: true,
        cache:false,
        data: data,
        processData: false,
        success: function (data, textStatus, jqXHR){
            let object = JSON.parse(data);
            CreateTable(object);
        },
        error: function () {
            alert('Failed...');
        }
    });
    return false;
}

function Action4Message() {
    jQuery.support.cors = true;
    let name = document.getElementById("file_name").value
    var data = JSON.stringify({"name": name});
    $.ajax({
        async: true,
        type: 'post',
        url: 'http://127.0.0.1:4001/api/ini',
        crossDomain: true,
        cache:false,
        data: data,
        processData: false,
        success: function () {
            alert('Файл изменен.');
        },
        error: function () {
            alert('Failed...');
        }
    });
}

function Action5Message() {
    jQuery.support.cors = true;
    let data1 = document.forms.f3.elements.date_start.value;
    let time1 = document.forms.f3.elements.time_start.value;
    let data2 = document.forms.f3.elements.date_end.value;
    let time2 = document.forms.f3.elements.time_end.value;
    var data = JSON.stringify({"date_start": data1, "time_start" : time1, "date_end": data2, "time_end" : time2});

    $.ajax({
        async: true,
        type: 'post',
        url: 'http://127.0.0.1:4001/api/searchLog',
        crossDomain: true,
        cache: false,
        processData: false,
        data: data,
        success: function (data) {
            var blob = new Blob([data], {type: "application/zip"});
            window.location.href = window.URL.createObjectURL(blob);
        },
        error: function () {
            alert('Failed...');
        }
    });
}

function Action6Message() {
    $.ajax({
        async: true,
        type: 'get',
        url: 'http://127.0.0.1:4001/api/searchLog',
        crossDomain: true,
        cache:false,
        dataType: 'json',
        success: function (data, textStatus, jqXHR){
            var obj = JSON.parse(jqXHR.responseText);
            document.getElementById("test").innerHTML = "Hostname: " + obj.hostname + "\n"
            + "CPU: " + obj.cpu + "\n" + "Platform: " + obj.platform + "\n" + "RAM: " + obj.ram + "\n"
            + "Disk: " + obj.disk + "\n";
        },
        error: function () {
            alert('Failed...');
        }
    });
}

function CreateTable(object) {
    let table = document.getElementById('table');
    ClearTable(table)
    if (table.rows.length === 0) {
        let tr = document.createElement('tr');
        let td1 = document.createElement('td');
        td1.innerHTML = "Название";
        tr.appendChild(td1);

        let td2 = document.createElement('td');
        td2.innerHTML = "Описание";
        tr.appendChild(td2);

        let td3 = document.createElement('td');
        td3.innerHTML = "Зависимости";
        tr.appendChild(td3);

        let td4 = document.createElement('td');
        td4.innerHTML = "Путь";
        tr.appendChild(td4);

        let td5 = document.createElement('td');
        td5.innerHTML = "Статус";
        tr.appendChild(td5);

        table.appendChild(tr);
    }

    for (let obj of object) {
        let tr = document.createElement('tr');
        let td1 = document.createElement('td');
        td1.innerHTML = obj.name;
        tr.appendChild(td1);

        let td2 = document.createElement('td');
        td2.innerHTML = obj.config.Description;
        tr.appendChild(td2);

        let td3 = document.createElement('td');
        td3.innerHTML = obj.config.Dependencies;
        tr.appendChild(td3);

        let td4 = document.createElement('td');
        td4.innerHTML = obj.config.BinaryPathName;
        tr.appendChild(td4);

        let td5 = document.createElement('td');
        td5.innerHTML = obj.status;
        tr.appendChild(td5);
        table.appendChild(tr);
    }
}

function ClearTable(table) {
    for (var i = table.rows.length-1; i > 0; i--) {
        table.deleteRow(i);
    }
}