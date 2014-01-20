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
      # _.bindAll @, 'render'

    render: ->
      $(@el).html @template (@model.toJSON())
      $(@el).find('.date').text moment(@model.get('Date')).format("MMM D")
      @

  app.AccomplishmentsView = Backbone.View.extend
    el: $("#accomplishments")

    initialize: ->
      @collection = new app.Accomplishments app.accomplishmentData
      @render()

    render: ->
      @collection.each (model) =>
        accomplishmentView = new app.AccomplishmentView {model: model}
        $(@el).append accomplishmentView.render().el 

  app.accomplishmentsView = new app.AccomplishmentsView

  container = document.querySelector('#accomplishments');
  imagesLoaded container, ->
    msnry = new Masonry container, 
      columnWidth: 75,
      itemSelector: '.item'
