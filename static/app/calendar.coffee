$(document).ready ->
  app = window || {}

  app.Accomplishment = Backbone.Model.extend
    defaults: ->
      Name: "Anonymous"
      Body: "I worked on calendar view...."
      Date: new Date()
      ImagePath: ""

  app.Accomplishments = Backbone.Collection.extend
    model: app.Accomplishment

  app.AccomplishmentView = Backbone.View.extend
    tagName: "div"
    className: "accomplishment"
    template: _.template($('#accomplishment-template').html())

    initialize: ->
      # _.bindAll @
      console.log "LOL NO WAI I AM AN ACCOMPLISHMENT VIEW!"

  app.accomplishmentView = new app.AccomplishmentView()
  app.accomplishments = new app.Accomplishments app.accomplishmentData