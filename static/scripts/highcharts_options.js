// 全局配置，针对页面上所有图表有效
Highcharts.setOptions({
    // chart: {
    //     borderWidth: 1,
    //     plotShadow: true,
    //     plotBorderWidth: 0
    // },
    credits: {
        enabled: false
    }
});

function drawPolar(data) {
    $('#'+data.domId).highcharts({
        chart: {
            polar: true,
            type: 'line'
        },
        title: {
            text: '',
            x: -80
        },
        pane: {
            size: '85%'
        },
        xAxis: {
            categories: ['A艺术型', 'S社会型', 'E企业型', 'C常规型',
                'R实际型', 'I调研型'],
            tickmarkPlacement: 'on',
            lineWidth: 0
        },
        yAxis: {
            gridLineInterpolation: 'polygon',
            lineWidth: 0,
            min: 0
        },
        tooltip: {
            shared: true,
            pointFormat: '<span style="color:{series.color}">{series.name}: <b>{point.y:,.0f}</b><br/>'
        },
        legend: {
            align: 'center',
            verticalAlign: 'bottom'
            // y: 70,
            // layout: 'vertical'
        },
        series: [{
            name: '优势',
            data: data.value1,
            pointPlacement: 'on'
        }, {
            name: '兴趣',
            data: data.value2,
            pointPlacement: 'on'
        }]
    });
}