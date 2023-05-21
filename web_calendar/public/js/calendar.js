function gettid() {
    alert();
}

function alerted(i) {
    document.getElementById('i').innetHTML = id
    alert(id)
}

function addEvent(date, time, description, event) {
    html += '<td><label type="textbox" class="teckst">' + event + ' </label></td>'
    html += '<td><label type="textbox" class="teckst">' + date + ' </label></td>'
    html += '<td><label type="textbox" class="teckst">' + time + ' </label></td>'
    html += '<td><label type="textbox" class="teckst">' + description + ' </label></td>';
}
var Cal = function(divId) {

    this.divId = divId;

    this.DaysOfWeek = [
        'Пн',
        'Вт',
        'Ср',
        'Чт',
        'Пт',
        'Сб',
        'Вс'
    ];

    this.Months = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'];

    var d = new Date();

    this.currMonth = d.getMonth('5');
    this.currYear = d.getFullYear('23');
    currDay = d.getDate('3');
};

Cal.prototype.nextMonth = function() {
    if (this.currMonth == 11) {
        this.currMonth = 0;
        this.currYear = this.currYear + 1;
    } else {
        this.currMonth = this.currMonth + 1;
    }
    this.showcurr();
};

Cal.prototype.previousMonth = function() {
    if (this.currMonth == 0) {
        this.currMonth = 11;
        this.currYear = this.currYear - 1;
    } else {
        this.currMonth = this.currMonth - 1;
    }
    this.showcurr();
};

Cal.prototype.showcurr = function() {
    this.showMonth(this.currYear, this.currMonth);
};

Cal.prototype.showMonth = function(y, m) {

    var d = new Date(),
        firstDayOfMonth = new Date(y, m, 7).getDay(),
        lastDateOfMonth = new Date(y, m + 1, 0).getDate(),
        lastDayOfLastMonth = m == 0 ? new Date(y - 1, 11, 0).getDate() : new Date(y, m, 0).getDate();


    var html = '<table>';

    html += '<thead><tr>';
    html += '<td colspan="7">' + this.Months[m] + ' ' + y + '</td>';
    html += '</tr></thead>';
    html += '<tr class="days">';
    for (var i = 0; i < this.DaysOfWeek.length; i++) {
        html += '<td>' + this.DaysOfWeek[i] + '</td>';
    }
    html += '</tr>';
    var i = 1;
    do {

        var dow = new Date(y, m, i).getDay();
        if (dow == 1) {
            html += '<tr>';
        } else if (i == 1) {
            html += '<tr>';
            var k = lastDayOfLastMonth - firstDayOfMonth + 1;
            for (var j = 0; j < firstDayOfMonth; j++) {
                html += '<td class="not-current">' + k + '</td>';
                k++;
            }
        }
        l = 1;
        var chk = new Date();
        var chkY = chk.getFullYear();
        var chkM = chk.getMonth();
        if (chkY == this.currYear && chkM == this.currMonth && i == currDay) {
            html += '<td class="submit"><button class="active" onclick="alerted();">' + i + '</button></td>';
        } else {
            html += '<td class="submit"><button class="notact" onclick="alerted();">' + i + '</button></td>';
        }
        if (dow == 0) {
            html += '</tr>';
        } else if (i == lastDateOfMonth) {
            var k = 1;
            for (dow; dow < 7; dow++) {
                html += '<td class="not-current">' + k + '</td>';
                k++;
            }
        }

        i++;
    } while (i <= lastDateOfMonth);

    html += '</table>';

    document.getElementById(this.divId).innerHTML = html;
};

window.onload = function() {

    var c = new Cal("divCal");
    c.showcurr();

    getId('btnNext').onclick = function() {
        c.nextMonth();
    };
    getId('btnPrev').onclick = function() {
        c.previousMonth();
    };
}

function getId(id) {
    return document.getElementById(id);
}