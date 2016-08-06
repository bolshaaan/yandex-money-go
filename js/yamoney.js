   /* YaMoney with Go @bolshaaan */
  $(function() {


      // datepicker

      $( "#datepicker,#datepicker2" )
        .datepicker({
          onSelect : function( date, inst ) {
            console.log(date)
          }
          }
        );


      $.getJSON('/data/out', function (ddata) {

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

                              var lll = labels[this.category];
                              return lll;
//                              console.log(lll);

//                              var res = "";
//                              for (var i in lll) {
//                                  res += "<br/> " + lll[i];
//                              }
//                              console.log(lll);
//
//                              return lll;
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
