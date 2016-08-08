   /* YaMoney with Go @bolshaaan */
  $(function() {

      var labels;

      var chart_options = {
            chart: { type : 'line' },
             plotOptions: {
                line: {
                    marker: {
                        enabled: true
                    }
                }
            },
            yAxis: {
                title: { text: "Рубли" },
                  min: -5
            },
            xAxis: {
                type: 'datetime',

                 dateTimeLabelFormats: { // don't display the dummy year
                    month: '%e. %b',
                    year: '%b'
                },

            },
            title: { text: 'YaMoney' },
            series: [
                {
                    marker: { enabled: true, radius: 5, symbol: "diamond" },
                    type: 'line',
                    name: "Потрачено",
                    data : [], //ddata["AggregatedOut"],
                    color: "red",

                    tooltip: {
                        useHTML: true,
                        pointFormatter : function() { return labels[this.category]; }
                    },

                    dataGrouping: {
                        approximation: "sum",
                        enabled: true,
                        forced: true,
                        units: [[ 'day', [1] ]]
                    }
                 }
            ]
       };

      $('#container').highcharts('StockChart', chart_options );

      function update_graph( call_from ) {

        console.log("call_from: " + call_from);

        $.getJSON('/data/', { "from": get_date("#datepicker"),  "till": get_date("#datepicker2") } ).done( function (ddata) {

           labels = ddata["Labels"];
           console.log(labels);
           $('#total').text(ddata["Total"]);

           var chart = $('#container').highcharts();
           labels = ddata["Labels"];
           $('#total').text(ddata["Total"]);
           chart.series[0].setData(ddata["AggregatedOut"] , true);

        });

      }

      function get_date(comp) {
        var date = $(comp).datepicker("getDate");

        var month = date.getMonth() + 1;
        month = month < 10
          ? "0" + month
          : month;

        var day = date.getDate();
        day = day < 10
          ? "0" + day
          : day;

        return date.getFullYear() + "-" + month + "-" + day + "T00:00:00Z";
      }

      var dd = 0;
      var default_dates = [ 0, -30 ];

      // datepicker
      $( "#datepicker,#datepicker2" )
        .datepicker({
          dateFormat: "yy-mm-dd",

          defaultDate: default_dates[dd++],

          onSelect : function( date, inst ) {
             update_graph("from_select");

//            $.ajax({
//              url: "/data/",
//              data: {
//                "from": get_date("#datepicker"),
//                "till": get_date("#datepicker2"),
//              },
//              dataType: "json",
//              success: function (data) {
//                var chart = $('#container').highcharts();
//
//                labels = data["Labels"];
//                $('#total').text(data["Total"]);
//
//                chart.series[0].setData(data["AggregatedOut"] , true);
//              }
//            });

          }
        });

       $("#datepicker").datepicker("setDate", -30);
       $("#datepicker2").datepicker("setDate", 0);
//       $("#datepicker2").datepicker("onSelect");

       update_graph("from_init");


  });
