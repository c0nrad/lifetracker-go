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
      initialize: function() {
        return console.log("LOL NO WAI I AM AN ACCOMPLISHMENT VIEW!");
      }
    });
    app.accomplishmentView = new app.AccomplishmentView();
    return app.accomplishments = new app.Accomplishments(app.accomplishmentData);
  });

}).call(this);
