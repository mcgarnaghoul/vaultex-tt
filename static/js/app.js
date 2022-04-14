const tableSelect = $('#frm-options_table-select');
const resultCountSelect = $('#frm-options_result-count');

const prevButton = $('#btn-previous');
const nextButton = $('#btn-next');
const refreshButton = $('#btn-refresh');

const currentPageLabel = $('#lbl-page-number');
const resultsContent = $('#results-content_container');

// Build API call URL.
function buildUrl(table, page, count) {
    return 'http://' + window.location.host + '/api/' + table + '?p=' + page + "&c=" + count;
}

// Make a get request to the API and refresh the UI.
//TODO move this dirty inline table generation into a proper framework
function getData(url) {
    const req = new XMLHttpRequest();
    req.overrideMimeType("application/json");
    req.open('GET', url, true);
    req.onload = function() {
        // Clear the content.
        resultsContent.html("");

        const records = JSON.parse(req.responseText);

        // Records should be an array of JSON objects.
        if (records.length < 1) {
            resultsContent.html("No results!");
            return;
        }
        const headers = Object.keys(records[0]);

        let tableHtml = '<table class="uk-table"><tr>';
        for (let i = 0; i < headers.length; i++) {
            tableHtml += '<th>' + headers[i] + '</th>';
        }
        tableHtml += '</tr>';

        for (let i = 0; i < records.length; i++) {
            tableHtml += '<tr>';
            for (let j = 0; j < headers.length; j++) {
                tableHtml += '<td>' + records[i][headers[j]] + '</td>';
            }
            tableHtml += '</tr>';
        }

        tableHtml += '</table>';
        resultsContent.html(tableHtml);
    };
    req.send(null);
}

// Set up click/change handlers.
prevButton.click(function () {
    var nextPage = parseInt(currentPageLabel.text()) - 1;
    if (nextPage < 1) {
        nextPage = 1;
    }
    getData(buildUrl(tableSelect.val(), nextPage, resultCountSelect.val()));

    currentPageLabel.text(nextPage);
});

nextButton.click(function () {
    var nextPage = parseInt(currentPageLabel.text()) + 1;
    if (nextPage < 1) {
        nextPage = 1;
    }
    getData(buildUrl(tableSelect.val(), nextPage, resultCountSelect.val()));

    currentPageLabel.text(nextPage);
});

refreshButton.click(function () {
    getData(buildUrl(tableSelect.val(), currentPageLabel.text(), resultCountSelect.val()));
});

tableSelect.change(function () {
    getData(buildUrl(tableSelect.val(), currentPageLabel.text(), resultCountSelect.val()));
});

resultCountSelect.change(function () {
    getData(buildUrl(tableSelect.val(), currentPageLabel.text(), resultCountSelect.val()));
});

// Load initial data.
getData(buildUrl(tableSelect.val(), currentPageLabel.text(), resultCountSelect.val()));