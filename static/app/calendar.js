// Generated by CoffeeScript 1.6.3
(function() {
  $(document).ready(function() {
    var app;
    app = window || {};
    app.CalendarModel = Backbone.Model.extend({
      defaults: function() {
        return {
          currMonth: moment().format('MMMM'),
          currYear: moment().format('YYYY'),
          currDay: moment().format('DD')
        };
      },
      firstCellBlock: function() {
        var dayName, days, firstDay;
        days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
        firstDay = "" + (this.get('currMonth')) + " 1 " + (this.get('currYear'));
        dayName = moment(firstDay).format('dddd');
        return days.indexOf(dayName);
      }
    });
    app.CalendarCellView = Backbone.View.extend({
      template: _.template($('#calendarCell-template').html()),
      initialize: function() {},
      render: function() {
        $(this.el).html(this.template({
          day: this.model.day,
          body: this.model.body
        }));
        return this;
      }
    });
    app.CalendarView = Backbone.View.extend({
      el: $('#calendarView'),
      initialize: function() {
        return this.model = new app.CalendarModel();
      },
      hide: function() {
        return $(this.el).hide();
      },
      unhide: function() {
        return $(this.el).show();
      },
      render: function() {
        var cellView, currAccomplishment, currDay, firstCellNumber, firstDay, x, _i, _ref;
        $(this.el).find("#calendarTitle").text("" + (this.model.get('currMonth')) + " " + (this.model.get('currYear')));
        firstDay = moment("" + (this.model.get('currMonth')) + " 1 " + (this.model.get('currYear')));
        currDay = firstDay;
        firstCellNumber = this.model.firstCellBlock();
        for (x = _i = 0, _ref = firstDay.daysInMonth(); 0 <= _ref ? _i <= _ref : _i >= _ref; x = 0 <= _ref ? ++_i : --_i) {
          cellView = new app.CalendarCellView({
            model: {
              day: x
            }
          });
          currAccomplishment = app.accomplishments.getAccomplishment(currDay);
          if (currAccomplishment != null) {
            cellView = new app.CalendarCellView({
              model: {
                day: x,
                body: currAccomplishment.get('Body')
              }
            });
          } else {
            cellView = new app.CalendarCellView({
              model: {
                day: x,
                body: ""
              }
            });
          }
          $(this.el).find("#cell" + (x + firstCellNumber)).html(cellView.render().el);
          currDay.add('days', 1);
        }
        return this.unhide();
      }
    });
    return app.calendarView = new app.CalendarView;
  });

}).call(this);
