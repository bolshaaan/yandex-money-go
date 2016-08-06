   /* YaMoney with Go @bolshaaan */
  $(function() {


      function get_date(comp) {
        var date = $(comp).datepicker("getDate");

        var month = date.getMonth();
        month = month < 10
          ? "0" + month
          : month;

        var day = date.getDate();
        day = day < 10
          ? "0" + day
          : day;

        return date.getFullYear() + "-" + month + "-" + day;
      }


      // datepicker
      $( "#datepicker,#datepicker2" )
        .datepicker({
          dateFormat: "yy-mm-dd",
          onSelect : function( date, inst ) {
            console.log(date);

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

                chart.series[0].setData(data["AggregatedOut"] , true);
              }
            });

          }
        });

      $.getJSON('/data/', function (ddata) {

         var labels = ddata["Labels"];
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
