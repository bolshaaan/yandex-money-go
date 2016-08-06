   /* YaMoney with Go @bolshaaan */
  $(function() {

      var labels;

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

      // datepicker
      $( "#datepicker,#datepicker2" )
        .datepicker({
          dateFormat: "yy-mm-dd",
          onSelect : function( date, inst ) {

            console.log(get_date('#datepicker'));
            console.log(get_date('#datepicker2'));

            $.ajax({
              url: "/data/",
              data: {
                "from": get_date("#datepicker"),
                "till": get_date("#datepicker2"),
              },
              dataType: "json",
              success: function (data) {
                var chart = $('#container').highcharts();

                labels = data["Labels"];
                $('#total').text(data["Total"]);

                chart.series[0].setData(data["AggregatedOut"] , true);
              }
            });

          }
        });

      $.getJSON('/data/', function (ddata) {

         labels = ddata["Labels"];
         console.log(labels);
         $('#total').text(ddata["Total"]);

         $('#container').highcharts('StockChart', {
              chart: { type : 'line' },
               plotOptions: {
                  line: {
                      marker: {
                          enabled: true
                      }
                  }
              },
              yAxis: {
                  title: { text: "Рубли" }
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
                      data : ddata["AggregatedOut"],
                      color: "red",

                      tooltip: {
                          useHTML: true,
                          pointFormatter : function() {

                              return labels[this.category];
                          }
                      },

                      dataGrouping: {
                          approximation: "sum",
                          enabled: true,
                          forced: true,
                          units: [[ 'day', [1] ]]

                      }

                   }
              ]
         });

      });
  });
