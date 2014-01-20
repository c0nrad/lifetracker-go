// Generated by CoffeeScript 1.6.3
(function() {
  $(document).ready(function() {
    var app;
    app = window || {};
    app.Accomplishment = Backbone.Model.extend({
      defaults: function() {
        return {
          Name: "Anonymous",
          Body: "I worked on calendar view....",
          Date: new Date(),
          ImagePath: ""
        };
      }
    });
    app.Accomplishments = Backbone.Collection.extend({
      model: app.Accomplishment
    });
    app.AccomplishmentView = Backbone.View.extend({
      tagName: "div",
      className: "accomplishment",
      template: _.template($('#accomplishment-template').html()),
      initialize: function() {},
      render: function() {
        $(this.el).html(this.template(this.model.toJSON()));
        $(this.el).find('.date').text(moment(this.model.get('Date')).format("MMM D"));
        return this;
      }
    });
    app.AccomplishmentsView = Backbone.View.extend({
      el: $("#accomplishments"),
      initialize: function() {
        this.collection = new app.Accomplishments(app.accomplishmentData);
        return this.render();
      },
      render: function() {
        var _this = this;
        return this.collection.each(function(model) {
          var accomplishmentView;
          accomplishmentView = new app.AccomplishmentView({
            model: model
          });
          return $(_this.el).append(accomplishmentView.render().el);
        });
      }
    });
    return app.accomplishmentsView = new app.AccomplishmentsView;
  });

}).call(this);
